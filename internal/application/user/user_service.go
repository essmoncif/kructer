package user

import (
	"errors"

	"github.com/jinzhu/gorm"
	"kructer.com/internal/core"
)


type UserService struct {
	userRepository *UserStore
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepository: &UserStore{
			DB: db,
		},
	}
}

func (usrService *UserService) Login(userLogin UserLogin, secret string) (*UserLoginInfo, error) {
	usr, err := usrService.userRepository.FindUserByLogin(userLogin.Login)
	if err != nil {
		return nil, err
	}

	if usr.Login == userLogin.Login && usr.Password == userLogin.Password {
		jwt := core.NewJWTClaims(usr.Login, usr.ID)
		token, err := jwt.GenerateToken(secret)
		if err != nil {
			return nil, err
		}

		usrLoginInfo := &UserLoginInfo{
			ID: usr.ID,
			Login: usr.Login,
			Token: token,
		}

		return usrLoginInfo, nil
	}

	return nil, errors.New("Login or password incorrect")
}

func (usrService *UserService) CreateUser(user User) (*User, error) {
	return usrService.userRepository.Create(user)
}