package user

import "github.com/jinzhu/gorm"


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

func (usrService *UserService) CreateUser(user User) (*User, error) {
	return usrService.userRepository.Create(user)
}