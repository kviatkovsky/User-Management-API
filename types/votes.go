package types

import (
	"time"

	"gorm.io/gorm"
)

type Votes struct {
	gorm.Model
	VoterID   uint      `gorm:"text;not null;" json:"voter_id"`
	ProfileID uint      `gorm:"text;not null" json:"profile_id"`
	VoteValue int       `gorm:"text" json:"value"`
	VoteTime  time.Time `gorm:"text" json:"time"`
}
