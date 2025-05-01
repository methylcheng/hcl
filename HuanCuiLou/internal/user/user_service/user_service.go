package user_service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"huancuilou/common"
	"huancuilou/configs"
	"huancuilou/internal/user/user_model"
	"huancuilou/internal/user/user_repository"
	"log"
	"math/rand"
	"strings"
	"time"
)

// UserService 处理业务逻辑
type UserService struct {
	userRepository    *user_repository.UserRepository
	config            *configs.Config
	userMdbRepository *user_repository.UserMemoryDBRepository
}

func NewUserService(userRepository *user_repository.UserRepository, config *configs.Config, userMemoryDBRepository *user_repository.UserMemoryDBRepository) *UserService {
	return &UserService{
		userRepository:    userRepository,
		config:            config,
		userMdbRepository: userMemoryDBRepository,
	}
}

// SendCode 发送验证码
func (us *UserService) SendCode(phoneNumber string) error {
	// 使用当前时间作为随机数种子
	rand.Seed(time.Now().UnixNano())
	var code string
	for i := 0; i < 6; i++ {
		// 生成 0 到 9 的随机数字
		digit := rand.Intn(10)
		code += fmt.Sprintf("%d", digit)
	}
	userCode := &user_model.UserCode{
		Code:        code,
		PhoneNumber: phoneNumber,
	}
	maskPhone := common.MaskPhoneNumber(phoneNumber)
	// 将验证码插入内存或更新验证码
	if err := us.userMdbRepository.AddCode(userCode, us.config.Code.ExpireDuration); err != nil {
		return fmt.Errorf("UserService.SendCode err: 500: 向内存插入或更新验证码错误: 手机号: %s,err: %w", maskPhone, err)
	}

	// 启动 goroutine 异步调用外部接口发送验证码
	go func() {
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			log.Printf(code)
			return
		}
		log.Printf("UserService.SendCode err:500: 达到最大重试次数，发送验证码失败，手机号: %s", maskPhone)
	}()

	return nil
}

// Login 一键登录注册
func (us *UserService) Login(userCode *user_model.UserCode) (*user_model.User, error) {
	//比对验证码
	err := us.userMdbRepository.ValidateCode(userCode.Code, userCode.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("UserService.Login err: 401: 验证码验证失败: %w", err)
	}
	//检查用户是否存在，存在则登录，不存在则注册
	exists, err := us.userRepository.GetUserByPhoneNumber(userCode.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("UserService.Login err: 500: 通过手机号查找用户错误: %w", err)
	}
	if exists == nil {
		maxRetries := 5
		for i := 0; i < maxRetries; i++ {
			userDB := user_model.User{
				PhoneNumber: userCode.PhoneNumber,
				UserName:    common.GenerateRandomUsername(20),
				Biography:   "添加个人简介，让大家更好地认识你~",
			}
			user, err := us.userRepository.AddUser(&userDB)
			if err != nil {
				if strings.Contains(err.Error(), "Error 1062 (23000): Duplicate entry") {
					log.Printf("用户名 %s 已存在，尝试重新生成...\n", userDB.UserName)
					continue
				}
				return nil, fmt.Errorf("UserService.Login err: 500:添加用户错误:%w", err)
			}
			return user, nil
		}
		return nil, errors.New("UserService.Login err: 500:达到最大重试次数，无法插入唯一用户名")
	}

	return exists, nil
}

func (us *UserService) GetUserByID(userID int) (*user_model.User, error) {
	var user *user_model.User
	user, err := us.userRepository.GetUserByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetUserByID err: 500:通过ID查找用户错误:%w", err)
	}
	return user, nil
}

func (us *UserService) AddAdminByPhoneNumber(phoneNumber string) error {
	user, err := us.userRepository.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf("UserService.AddAdminByPhoneNumber err: 500:通过手机号查找用户错误:%w", err)
	}
	if user.IsManager == 1 {
		return fmt.Errorf("UserService.AddAdminByPhoneNumber err: 400:该用户已经是管理员")
	}
	if err = us.userRepository.AddAdminByPhoneNumber(phoneNumber); err != nil {
		return fmt.Errorf("UserService.AddAdminByPhoneNumber err: 500:更新管理员出错:%w", err)
	}
	return nil
}

func (us *UserService) UpdateUserInfo(userID int, user *user_model.User) error {
	if err := us.userRepository.UpdateUserInfo(userID, user); err != nil {
		return fmt.Errorf("UserService.UpdateUserInfo err: 500:更新用户信息出错:%w", err)
	}
	return nil
}

func (us *UserService) AddScore(scoreRecord *user_model.ScoreRecord) error {
	if err := us.userRepository.AddScore(scoreRecord); err != nil {
		return fmt.Errorf("UserService.AddScore 数据库操作错误:添加积分记录错误:%w", err)
	}
	return nil
}

func (us *UserService) AddPhoneRecord(managerID int, content string) error {
	// 查找最新的评分记录
	var latestScoreRecord user_model.ScoreRecord
	result := us.userRepository.DB.Where("manager_id = ?", managerID).Order("created_at desc").First(&latestScoreRecord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("UserService.AddPhoneRecord 错误请求:未找到评分记录")
		}
		return fmt.Errorf("UserService.AddPhoneRecord 数据库操作错误:查找评分记录错误:%w", result.Error)
	}

	// 查找用户电话号码
	user, err := us.userRepository.GetUserByUserID(latestScoreRecord.UserID)
	if err != nil {
		return fmt.Errorf("UserService.AddPhoneRecord 数据库操作错误:查找用户错误:%w", err)
	}
	if user == nil {
		return fmt.Errorf("UserService.AddPhoneRecord 错误请求:用户不存在")
	}

	// 创建电话记录，Content字段让管理员自己输入
	phoneRecord := &user_model.PhoneRecord{
		UserID:       latestScoreRecord.UserID,
		ManagerID:    latestScoreRecord.ManagerID,
		Satisfaction: latestScoreRecord.Score,
		UserPhone:    user.PhoneNumber,
		CreatedAt:    time.Now(),
		Content:      content,
	}

	// 添加电话记录到数据库
	if err := us.userRepository.AddPhoneRecord(phoneRecord); err != nil {
		return fmt.Errorf("UserService.AddPhoneRecord 数据库操作错误:添加电话记录错误:%w", err)
	}
	return nil
}

func (us *UserService) GetPhoneRecordByPhone(phoneNumber string) ([]user_model.PhoneRecord, error) {
	phoneRecords, err := us.userRepository.GetPhoneRecordByPhone(phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("UserService.GetPhoneRecordByPhone err:%w", err)
	}
	return phoneRecords, nil
}
