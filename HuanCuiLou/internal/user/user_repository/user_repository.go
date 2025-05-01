package user_repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"huancuilou/internal/user/user_model"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (ur *UserRepository) GetUserByUserID(userID int) (*user_model.User, error) {
	var user user_model.User
	result := ur.DB.Take(&user, "ID = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepository.GetUserByID err:%w", result.Error)
	}
	return &user, nil
}

func (ur *UserRepository) AddUser(user *user_model.User) (*user_model.User, error) {
	result := ur.DB.Create(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("UserRepository.AddUser err:%w", result.Error)
	}
	return user, nil
}

func (ur *UserRepository) GetUserByPhoneNumber(phoneNumber string) (*user_model.User, error) {
	var user user_model.User
	result := ur.DB.Take(&user, "phone_number = ?", phoneNumber)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepository.GetUserByID err:%w", result.Error)
	}
	return &user, nil
}

func (ur *UserRepository) AddAdminByPhoneNumber(phoneNumber string) error {
	result := ur.DB.Model(&user_model.User{}).Where("phone_number = ?", phoneNumber).Update("is_manager", 1)
	if result.Error != nil {
		return fmt.Errorf("UserRepository.AddAdminByPhoneNumber err:%w", result.Error)
	}
	return nil
}

func (ur *UserRepository) UpdateUserInfo(id int, user *user_model.User) error {
	result := ur.DB.Model(&user_model.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("UserRepository.UpdateUserInfo err:%w", result.Error)
	}
	return nil
}

func (ur *UserRepository) AddScore(scoreRecord *user_model.ScoreRecord) error {
	result := ur.DB.Create(&scoreRecord)
	if result.Error != nil {
		return fmt.Errorf("UserRepository.AddScore err:%w", result.Error)
	}
	return nil
}
func (ur *UserRepository) AddPhoneRecord(phoneRecord *user_model.PhoneRecord) error {
	result := ur.DB.Create(&phoneRecord)
	if result.Error != nil {
		return fmt.Errorf("UserRepository.AddPhoneRecord err:%w", result.Error)
	}
	return nil
}

func (ur *UserRepository) GetPhoneRecordByPhone(phoneNumber string) ([]user_model.PhoneRecord, error) {
	var phoneRecords []user_model.PhoneRecord
	result := ur.DB.Where("user_phone = ?", phoneNumber).Find(&phoneRecords)
	if result.Error != nil {
		return nil, fmt.Errorf("UserRepository.GetPhoneRecordByPhone err:%w", result.Error)
	}
	return phoneRecords, nil
}
