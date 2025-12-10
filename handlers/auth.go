package handlers

import (
	"github.com/gavin/blog/config"
	"github.com/gavin/blog/errors"
	"github.com/gavin/blog/logger"
	"github.com/gavin/blog/models"
	"github.com/gavin/blog/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

type RegisterRequest struct {
	*utils.FieldValidate
	Username       string `json:"username" binding:"required,min=3,max=20" label:"用户名"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	RepeatPassword string `json:"repeat_password" binding:"required,min=6"`
}

type LoginRequest struct {
	*utils.FieldValidate
	Username string `json:"username" binding:"required,min=3,max=20" label:"用户名"`
	Password string `json:"password" binding:"required,min=6" label:"密码"`
}

type AuthResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Fail(c, errors.AUTH_ERROR, "user not found")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Fail(c, errors.AUTH_ERROR, "password incorrect")
		return
	}

	token, err := utils.GenerateToken(uint64(user.ID), user.Username)

	if err != nil {
		logger.Log.Errorf("generate token err: %v", err)
		utils.Fail(c, errors.AUTH_ERROR, "generate token failed")
		return
	}

	utils.Success(c, &AuthResponse{
		Username: user.Username,
		Token:    token,
	}, "")
	return
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}

	if req.Password != req.Password {
		utils.Fail(c, errors.AUTH_ERROR, "two password not match")
		return
	}

	var existUser models.User
	config.DB.Where("username = ?", req.Username).First(&existUser)
	if existUser.ID != 0 {
		utils.Fail(c, errors.AUTH_ERROR, "username is exist")
		return
	}

	config.DB.Where("email = ?", req.Email).First(&existUser)
	if existUser.ID != 0 {
		utils.Fail(c, errors.AUTH_ERROR, "email is exist")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Infof("hash err: %v", err)
		utils.Error(c, "bcrypt password err")
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		logger.Log.Infof("create user err: %v", err)
		utils.Error(c, "create user fail")
		return
	}

	token, err := utils.GenerateToken(uint64(user.ID), user.Username)

	if err != nil {
		logger.Log.Errorf("generate token err: %v", err)
		utils.Fail(c, errors.AUTH_ERROR, "generate token failed")
		return
	}
	utils.Success(c, &AuthResponse{
		Username: user.Username,
		Token:    token,
	}, "register success")
	return
}
