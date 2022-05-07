package token

import "time"

//Maker is an interface for managig creation and verification tokens
type Maker interface {
	//create and sign and sign new token for a username and valid duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	//check the validity of the token
	VerifyToken(token string) (*Payload, error)
}