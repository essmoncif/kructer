package user

import (
	"github.com/jinzhu/gorm"
)

type UserStore struct {
	DB *gorm.DB
}

func (us *UserStore) Create(user User) (*User, error) {
	result := us.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (us *UserStore) FindByID(id string) (*User, error) {
	var user User
	result := us.DB.First(&user, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (us *UserStore) GetUsersWithPagination(page, pageSize int) ([]User, int64, error) {
	var users []User
	var totalRows int64

	us.DB.Model(&User{}).Count(&totalRows)

	result := us.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, totalRows, nil
}

func (us *UserStore) Update(id string, updatedData map[string]interface{}) (*User, error) {
	var user User
	result := us.DB.Model(&user).Where("id = ?", id).Updates(updatedData)
	if result.Error != nil {
		return nil, result.Error
	}
	us.DB.First(&user, "id = ?", id)
	return &user, nil
}

func (us *UserStore) Delete(id string) error {
	result := us.DB.Delete(&User{}, "id = ?", id)
	return result.Error
}
