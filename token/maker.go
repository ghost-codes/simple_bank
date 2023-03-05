package token

import "time"

// managin token makers
type Maker interface {
	//create and sign new token for user
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
