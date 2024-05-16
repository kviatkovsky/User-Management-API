package types

import (
	"gorm.io/gorm"
)

type UserPayload struct {
	User     User      `json:"user"`
	UserRole *UserRole `json:"user_role"`
}

type UserWithRating struct {
	User
	TotalRating int
}

type User struct {
	gorm.Model
	ID        uint     `gorm:"primarykey"`
	FirstName string   `gorm:"text;default:null" json:"first_name"`
	LastName  string   `gorm:"text;default:null" json:"last_name"`
	Email     string   `gorm:"text;not null;" json:"email"`
	Password  string   `gorm:"text;not null;" json:"password"`
	Rating    []*Votes `gorm:"foreignKey:ProfileID"`
}

type UserStore interface {
	GetList() ([]UserWithRating, error)
	CreateUser(user *User) (uint, error)
	UpdateUser(user *User) error
	GetRoleByLevel(level int) (UserRole, error)
	AssigneeUserRole(uid uint, role uint) (uint, error)
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	DeleteUserById(id string) error
	GetUserRoleByUserId(uid uint) (*UserRole, error)
	UpdateUserAccessLevel(u *User, accessLevel uint) error
	AssigneeVote(t *Votes) (uint, error)
	GetVoteByVoterIdAndProfile(voterId uint, profileId uint) (*Votes, error)
	GetAllVotesByUserId(voterId uint) (*[]Votes, error)
	DeleteVote(u uint, id uint) error
	UpdateVote(user *Votes) error
}

type Helper interface {
	IsVoteActionAvailable(uid uint) (bool, error)
	IsAlreadyVoted(uid uint, profileId uint) (bool, error)
	IsWithdrawVote(vote Votes) bool
	IsVoteChange(u uint, t *Votes) (bool, error)
}
