package server

import (
	"github.com/gavin/blog/config"
	"github.com/gavin/blog/errors"
	"github.com/gavin/blog/logger"
	"github.com/gavin/blog/models"
	"github.com/gavin/blog/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatePostRequest struct {
	*utils.FieldValidate
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdatePostRequest struct {
	*utils.FieldValidate
	ID      int    `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1"`
}

type QueryPostsRequest struct {
	*utils.FieldValidate
	utils.Pagination
	UserId int `json:"user_id"`
}

func GetPagePosts(c *gin.Context) {
	// 获取分页
	var req QueryPostsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}

	var posts []models.Post
	query := config.DB.Model(&models.Post{})
	if req.UserId > 0 {
		query.Where("user_id = ?", req.UserId)
	}
	paginatedResult, err := utils.GetPaginatedData(query, &posts, &req.Pagination)

	if err != nil {
		utils.Fail(c, errors.POST_ERROR, "查询失败")
		return
	}
	paginatedResult.Data = posts
	utils.Success(c, paginatedResult, "")
	return
}

func GetUserPost(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.POST_ERROR, "用户未登录")
		return
	}
	var posts []models.Post
	config.DB.Where("user_id", userId).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("ID desc").Limit(10) // 限制只加载 10 条关联数据
	}).Find(&posts)

	utils.Success(c, posts, "")
	return
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := config.DB.Where("id", id).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Limit(10) // 限制只加载 10 条关联数据
	}).First(&post).Error; err != nil {
		utils.Fail(c, errors.POST_ERROR, "文章没找到")
		return
	}
	utils.Success(c, post, "")
	return
}

func AddPost(c *gin.Context) {
	// 新增文章
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.POST_ERROR, "用户未登录")
		return
	}
	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId.(uint64),
	}
	if err := config.DB.Create(post).Error; err != nil {
		logger.Log.Error(err)
		utils.Fail(c, errors.POST_ERROR, "添加文章失败")
		return
	}
	utils.Success(c, "", "添加成功")
}

func UpdatePost(c *gin.Context) {
	// 新增文章
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.POST_ERROR, "用户未登录")
		return
	}

	var existPost models.Post

	if err := config.DB.Where("user_id", userId).First(&existPost, "id = ?", req.ID).Error; err != nil {
		utils.Fail(c, errors.POST_ERROR, "文章不存在")
		return
	}

	existPost.Title = req.Title
	existPost.Content = req.Content

	if err := config.DB.Save(&existPost).Error; err != nil {
		logger.Log.Error(err)
		utils.Fail(c, errors.POST_ERROR, "修改文章失败")
		return
	}
	utils.Success(c, "", "修改文章成功")
}

func DeletePost(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.POST_ERROR, "用户未登录")
		return
	}
	id := c.Param("id")
	var existPost models.Post
	if err := config.DB.Where("user_id", userId).Where("id", id).First(&existPost).Error; err != nil {
		utils.Fail(c, errors.POST_ERROR, "文章没找到")
		return
	}
	if err := config.DB.Delete(&existPost).Error; err != nil {
		logger.Log.Error(err)
		utils.Fail(c, errors.POST_ERROR, "删除失败")
		return
	}

	utils.Success(c, "", "删除成功")
	return
}
