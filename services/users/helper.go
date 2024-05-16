package users

import (
	"main/types"
	"time"

	"gorm.io/gorm"
)

type Helper struct {
	db    *gorm.DB
	store types.UserStore
}

func NewHelper(db *gorm.DB, userStore types.UserStore) *Helper {
	return &Helper{db: db, store: userStore}
}

func (h *Helper) IsAlreadyVoted(uid uint, profileId uint) (bool, error) {
	vote, err := h.store.GetVoteByVoterIdAndProfile(uid, profileId)
	if err != nil {
		return false, err
	}

	return vote.ID > 0, nil
}

func (h *Helper) IsWithdrawVote(vote types.Votes) bool {
	return vote.VoteValue == 0
}

func (h *Helper) IsVoteChange(uuid uint, vote *types.Votes) (bool, error) {
	existingVote, err := h.store.GetVoteByVoterIdAndProfile(uuid, vote.ProfileID)
	if err != nil {
		return false, err
	}

	if existingVote.VoteValue == vote.VoteValue {
		return false, err
	}

	vote.ID = existingVote.ID

	return true, nil
}

func (h *Helper) IsVoteActionAvailable(uid uint) (bool, error) {
	votes, err := h.store.GetAllVotesByUserId(uid)
	if err != nil {
		return false, err
	}

	if len(*votes) > 0 {
		currentTime := time.Now()
		vote := (*votes)[0]

		timeDiff := currentTime.Sub(vote.VoteTime)
		return timeDiff > time.Hour, nil
	}

	return true, nil
}
