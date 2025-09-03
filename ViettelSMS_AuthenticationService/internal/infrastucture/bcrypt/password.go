package password

import (
	srv "github.com/th1enq/ViettelSMS_AuthenticationService/internal/domain/service"
	"golang.org/x/crypto/bcrypt"
)

type bcryptService struct{}

func NewBcryptService() srv.PasswordService {
	return &bcryptService{}
}

func (b *bcryptService) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func (b *bcryptService) Verify(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
