package constants

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// ExpireTime for token
	ExpireTime = time.Duration(12) * time.Hour
	// SigningMethod represents algorithm token method
	SigningMethod = jwt.SigningMethodHS256
	// SignatureKey represents random string for token signature
	SignatureKey = []byte("CJ59d3WFwke9r75q0bcYm6MDCwBxVqY3")
)
