package users

import (
	"app/internal/pkg/repository"
	"app/pkg/gormx"

	"gorm.io/gorm"
)

type UserRepository struct {
	repository.Repository
	db *gorm.DB
}

func NewUserRepository(ds gormx.DataSources) *UserRepository {
	return &UserRepository{
		db: ds["mysql"],
	}
}

func (r UserRepository) ExistsByUsername(username string) bool {
	var count int64
	r.db.Model(&User{}).Where("username = ?", username).Limit(1).Count(&count)
	return count > 0
}

func (r UserRepository) FindById(id uint) (*User, error) {
	var user *User
	if err := r.db.Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, r.WrapTakeError(err)
	}
	return user, nil
}

func (r UserRepository) FindByUsername(username string) (*User, error) {
	var user *User
	if err := r.db.Where("username = ?", username).Take(&user).Error; err != nil {
		return nil, r.WrapTakeError(err)
	}
	return user, nil
}

func (r UserRepository) Create(user User) error {
	return r.db.Create(&user).Error
}
