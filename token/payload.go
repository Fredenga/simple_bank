package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	ID uuid.UUID 		`json:"id"`
	Username string 	`json:"username"`
	IssuedAt time.Time  `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

//create a token payload
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom() //generate token ID
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID: tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	//check expiration time
	if time.Now().After(payload.ExpiredAt){
		return ErrExpiredToken
	}
	return nil
}