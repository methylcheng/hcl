package user_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"huancuilou/common"
	"huancuilou/common/error_handler"
	"huancuilou/configs"
	"huancuilou/internal/user/user_model"
	"huancuilou/internal/user/user_service"
	"huancuilou/response"
	"log"
	"net/http"
	"strconv"
)

//处理用户操作的相关接口

//发送验证码接口
//登录接口
//社区管理员认证接口
//修改个人信息接口
//获取用户个人信息

// UserController 处理请求和返回响应
type UserController struct {
	userService *user_service.UserService
	jwtConfig   configs.JwtConfig
}

func NewUserController(userService *user_service.UserService, jwtConfig configs.JwtConfig) *UserController {
	return &UserController{
		userService: userService,
		jwtConfig:   jwtConfig,
	}
}

func (uc *UserController) SendCode(c *gin.Context) {
	phoneNumber := c.Param("phoneNumber")
	if !common.ValidatePhoneNumber(phoneNumber) {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.SendCode err: 400:手机号格式错误"))
	} else {
		if err := uc.userService.SendCode(phoneNumber); err != nil {
			error_handler.HandleUserError(c, fmt.Errorf("UserController.SendCode err: %w", err))
		} else {
			log.Printf("UserController.SendCode 成功发送验证码")
			c.JSON(http.StatusOK, response.SuccessWithoutData())
		}
	}
}

func (uc *UserController) Login(c *gin.Context) {
	var userCode user_model.UserCode
	if err := c.BindJSON(&userCode); err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.Login err: 400:将json数据绑定到结构体失败:%w", err))
		return
	}

	user, err := uc.userService.Login(&userCode)
	if err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.Login err: %w", err))
		return
	}

	accessToken, err := common.GenAccessToken(user.ID, user.IsManager, uc.jwtConfig.AccessTokenExpireDuration, uc.jwtConfig.SecretKey)
	if err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.Login err: 500:生成accessToken失败:%w", err))
		return
	}
	refreshToken, err := common.GenRefreshToken(user.ID, user.IsManager, uc.jwtConfig.RefreshTokenExpireDuration, uc.jwtConfig.SecretKey)
	if err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.Login err: 500:生成refreshToken失败:%w", err))
		return
	}

	token := map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}
	log.Printf("UserController.Login 成功登录")
	c.JSON(http.StatusOK, response.Success(token))
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	//类型断言
	userID := c.MustGet("userID").(int)
	user, err := uc.userService.GetUserByID(userID)
	if err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.GetUserInfo err: %w", err))
	} else {
		log.Printf("UserController.GetUserInfo 成功获取用户信息:%d", userID)
		c.JSON(http.StatusOK, response.Success(user))
	}
}

func (uc *UserController) AddAdministrator(c *gin.Context) {
	phoneNumber := c.Param("phoneNumber")
	if !common.ValidatePhoneNumber(phoneNumber) {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.SendCode err: 400:错误请求:手机号格式错误"))
	} else {
		if err := uc.userService.AddAdminByPhoneNumber(phoneNumber); err != nil {
			error_handler.HandleUserError(c, fmt.Errorf("UserController.AddAdministrator err: %w", err))
		} else {
			log.Printf("UserController.AddAdministrator 成功添加管理员")
			c.JSON(http.StatusOK, response.SuccessWithoutData())
		}
	}
}

func (uc *UserController) UpdateUserInfo(c *gin.Context) {
	//类型断言
	userID := c.MustGet("userID").(int)

	var user user_model.User
	if err := c.BindJSON(&user); err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.UpdateUserInfo err: 400:将json数据绑定到结构体失败:%w", err))
		return
	}
	if err := uc.userService.UpdateUserInfo(userID, &user); err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.UpdateUserInfo err: %w", err))
	} else {
		log.Printf("UserController.UpdateUserInfo 成功更新用户信息:%d", userID)
		c.JSON(http.StatusOK, response.SuccessWithoutData())
	}
}

func (uc *UserController) AddScore(c *gin.Context) {
	var scoreRecord user_model.ScoreRecord
	if err := c.BindJSON(&scoreRecord); err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.AddScore 错误请求:将json数据绑定到结构体失败:%w", err))
		return
	}
	if err := uc.userService.AddScore(&scoreRecord); err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.AddScore err:%w", err))
	} else {
		log.Printf("UserController.AddScore 成功添加评分")
		c.JSON(http.StatusOK, response.SuccessWithoutData())
	}
}

func (uc *UserController) AddPhoneRecord(c *gin.Context) {
	// 从表单数据中获取userID并转换为int
	userIDStr := c.PostForm("userID")
	managerID, err := strconv.Atoi(userIDStr)
	if err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.AddPhoneRecord err:无法将userID转换为int:%w", err))
		return
	}
	//管理员输入内容
	content := c.PostForm("content")
	if err := uc.userService.AddPhoneRecord(managerID, content); err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.AddPhoneRecord err:%w", err))
	} else {
		log.Printf("UserController.AddPhoneRecord 成功添加求助记录")
		c.JSON(http.StatusOK, response.SuccessWithoutData())
	}
}

func (uc *UserController) GetPhoneRecordByPhone(c *gin.Context) {
	phoneNumber := c.Param("phoneNumber")
	if !common.ValidatePhoneNumber(phoneNumber) {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.GetPhoneRecordByPhone err:错误请求:手机号格式错误"))
		return
	}
	phoneRecords, err := uc.userService.GetPhoneRecordByPhone(phoneNumber)
	if err != nil {
		error_handler.HandleUserError(c, fmt.Errorf("UserController.GetPhoneRecordByPhone err:%w", err))
	} else {
		log.Printf("UserController.GetPhoneRecordByPhone 成功获取求助记录:%d", len(phoneRecords))
		c.JSON(http.StatusOK, response.Success(phoneRecords))
	}
}
