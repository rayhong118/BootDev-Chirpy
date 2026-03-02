package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {

	hash, hashErr := argon2id.CreateHash(password, argon2id.DefaultParams)

	if hashErr != nil {
		return "", hashErr
	}

	return hash, nil

}

func CheckPasswordHash(password, hash string) (bool, error) {
	result, compareErr := argon2id.ComparePasswordAndHash(password, hash)

	if compareErr != nil {
		return false, compareErr
	}

	return result, nil
}
