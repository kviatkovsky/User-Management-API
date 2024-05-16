package users

import (
	"context"
	"encoding/json"
	"fmt"
	"main/types"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Store struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewStore(db *gorm.DB, r *redis.Client) *Store {
	return &Store{
		db:    db,
		redis: r,
	}
}

func (s *Store) GetList() ([]types.UserWithRating, error) {
	var users []types.UserWithRating
	rKey := "users:list:*"
	ctx := context.Background()

	exist, rErr := s.redis.Exists(ctx, rKey).Result()
	if rErr != nil {
		return nil, fmt.Errorf("error checking redis user key: %s", rErr.Error())
	}

	if exist == 1 {
		usersFromRedis, rErr := s.redis.Get(ctx, "users:list:*").Result()

		if rErr != nil {
			return nil, fmt.Errorf("error getting users from redis: %s", rErr)
		}

		err := json.Unmarshal([]byte(usersFromRedis), &users)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling users: %s", err.Error())
		}

		return users, nil
	} else {
		dbErr := s.db.Table("users").
			Select("users.*, COALESCE(SUM(votes.vote_value), 0) AS total_rating").
			Joins("LEFT JOIN votes ON users.id = votes.profile_id").
			Group("users.id").
			Scan(&users)

		if dbErr.Error != nil {
			return nil, fmt.Errorf("error getting users: %s", dbErr.Error)
		}

		jsonString, err := json.Marshal(users)
		if err != nil {
			return nil, fmt.Errorf("error converting users struct: %s", err.Error())
		}

		rErr = s.redis.Set(ctx, "users:list:*", jsonString, 60*time.Second).Err()
		if rErr != nil {
			return nil, fmt.Errorf("user list to redis: %s", rErr)
		}
	}

	return users, nil
}

func (s *Store) CreateUser(u *types.User) (uint, error) {
	err := s.db.Create(&u)
	if err != nil {
		return 1, fmt.Errorf("error creating user: %s", err.Error)
	}

	return u.ID, nil
}

func (s *Store) UpdateUser(u *types.User) error {
	err := s.db.Updates(&u)
	if err != nil {
		return fmt.Errorf("error updating user: %s", err.Error)
	}

	return nil
}

func (s *Store) GetRoleByLevel(level int) (types.UserRole, error) {
	var userRole types.UserRole

	err := s.db.Find(&userRole, "access_level = ?", level)
	if err != nil {
		return userRole, fmt.Errorf("getting role by level: %s", err.Error)
	}

	return userRole, nil
}

func (s *Store) AssigneeUserRole(uid uint, role uint) (uint, error) {
	acl := types.ACL{
		UserID: uid,
		RoleID: role,
	}
	err := s.db.Create(&acl)
	if err != nil {
		return 1, fmt.Errorf("error assigne user role: %s", err.Error)
	}

	return acl.RoleID, nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user types.User

	s.db.Find(&user, "email = ?", email)
	if user.ID == 0 {
		return nil, fmt.Errorf("user with email: %s not found", email)
	}

	return &user, nil
}

func (s *Store) GetUserById(id int) (*types.User, error) {
	var user types.User

	s.db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return nil, fmt.Errorf("user with id %d not found", id)
	}

	return &user, nil
}

func (s *Store) DeleteUserById(uid string) error {
	var user types.User
	err := s.db.Unscoped().Delete(&user, "id = ?", uid)

	if err.Error != nil {
		return fmt.Errorf("error deleting user: %s", err.Error)
	}

	return nil
}

func (s *Store) GetUserRoleByUserId(uid uint) (*types.UserRole, error) {
	var acl types.ACL
	var userRole types.UserRole

	err := s.db.Find(&acl, "user_id = ?", uid)
	if err.Error != nil {
		return nil, fmt.Errorf("error getting user: %s", err.Error)
	}

	err = s.db.Find(&userRole, acl.RoleID)
	if err.Error != nil {
		return nil, fmt.Errorf("error getting role by user: %s", err.Error)
	}

	return &userRole, nil
}

func (s *Store) UpdateUserAccessLevel(u *types.User, accessLevel uint) error {
	err := s.db.Model(&types.ACL{}).Where("user_id = ?", u.ID).Update("role_id", accessLevel)

	if err.Error != nil {
		return fmt.Errorf("error updating user access level: %s", err.Error)
	}
	return nil
}

func (s *Store) AssigneeVote(v *types.Votes) (uint, error) {
	err := s.db.Create(&v)

	if err.Error != nil {
		return 1, fmt.Errorf("error assigning vote: %s", err.Error)
	}

	return v.ID, nil
}

func (s *Store) UpdateVote(v *types.Votes) error {
	err := s.db.Updates(&v)
	if err.Error != nil {
		return fmt.Errorf("error updating votes: %s", err.Error)
	}

	return nil
}

func (s *Store) GetVoteByVoterIdAndProfile(voterId uint, profileId uint) (*types.Votes, error) {
	var vote types.Votes
	err := s.db.Find(&vote, "voter_id = ? AND profile_id = ?", voterId, profileId)
	if err.Error != nil {
		return nil, fmt.Errorf("error getting vote by voter id and profile id %d", voterId)
	}

	return &vote, nil
}

func (s *Store) DeleteVote(voterId uint, profileId uint) error {
	var vote types.User
	err := s.db.Unscoped().Delete(&vote, "voter_id = ? AND profile_id = ?", voterId, profileId)
	if err.Error != nil {
		return fmt.Errorf("error deleting vote: %s", err.Error)
	}

	return nil
}

func (s *Store) GetAllVotesByUserId(voterId uint) (*[]types.Votes, error) {
	var votes []types.Votes
	err := s.db.Order("vote_time DESC").Find(&votes, "voter_id = ?", voterId)
	if err.Error != nil {
		return nil, fmt.Errorf("error getting all votes: %s", err.Error)
	}

	return &votes, nil
}
