package userstate

import (
	"time"
)

type userState struct {
	UserID      int       `json:"userID"`
	State       string    `json:"state"`
	LastUpdated time.Time `json:"lastUpdated"`
}
