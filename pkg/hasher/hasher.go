package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashed, plain string) error
}

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func (h *BcryptHasher) Compare(hashed, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))

	if err != nil {
		return err
	}

	return nil
}
