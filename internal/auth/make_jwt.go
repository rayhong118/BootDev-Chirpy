package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	expiresIn := time.Duration(1) * time.Second
	issuedAt := jwt.NumericDate{Time: time.Now().UTC()}
	expiresAt := jwt.NumericDate{Time: time.Now().Add(expiresIn).UTC()}

	newToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy-access",
			IssuedAt:  &issuedAt,
			ExpiresAt: &expiresAt,
			Subject:   userID.String(),
		},
	)

	signedToken, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
