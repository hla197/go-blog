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

type CreateCommentRequest struct {
	*utils.FieldValidate
	PostID  uint64 `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
	*utils.FieldValidate
	ID      int    `json:"id" binding:"required"`
	PostID  uint64 `json:"post_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1"`
}

type QueryCommentsRequest struct {
	*utils.FieldValidate
	utils.Pagination
	PostID uint64 `json:"post_id"`
	UserId int    `json:"user_id"`
}

func GetPageComments(c *gin.Context) {
	// 获取分页
	var req QueryCommentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}

	var posts []models.Comment
	query := config.DB.Model(&models.Comment{})
	if req.UserId > 0 {
		query.Where("user_id = ?", req.UserId)
	}
	if req.PostID > 0 {
		query.Where("post_id = ?", req.PostID)
	}
	paginatedResult, err := utils.GetPaginatedData(query, &posts, &req.Pagination)

	if err != nil {
		utils.Fail(c, errors.COMMENT_ERROR, "查询失败")
		return
	}
	paginatedResult.Data = posts
	utils.Success(c, paginatedResult, "")
	return
}

func GetUserComment(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.COMMENT_ERROR, "用户未登录")
		return
	}
	var posts []models.Comment
	config.DB.Where("user_id", userId).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Order("ID desc").Limit(10) // 限制只加载 10 条关联数据
	}).Find(&posts)

	utils.Success(c, posts, "")
	return
}

func GetComment(c *gin.Context) {
	id := c.Param("id")
	var post models.Comment
	if err := config.DB.Where("id", id).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Limit(10) // 限制只加载 10 条关联数据
	}).First(&post).Error; err != nil {
		utils.Fail(c, errors.COMMENT_ERROR, "评论没找到")
		return
	}
	utils.Success(c, post, "")
	return
}

func AddComment(c *gin.Context) {
	// 新增评论
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.COMMENT_ERROR, "用户未登录")
		return
	}
	var existPost models.Post
	if err := config.DB.First(&existPost, req.PostID).Error; err != nil {
		utils.Fail(c, errors.COMMENT_ERROR, "文章不存在")
		return
	}

	comment := &models.Comment{
		Content: req.Content,
		UserID:  userId.(uint64),
		PostID:  req.PostID,
	}
	if err := config.DB.Create(&comment).Error; err != nil {
		logger.Log.Error(err)
		utils.Fail(c, errors.COMMENT_ERROR, "添加评论失败")
		return
	}
	utils.Success(c, "", "添加成功")
}

func UpdateComment(c *gin.Context) {
	// 新增评论
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validate utils.FieldValidateIF = req
		msg := validate.Validate(err, req)
		utils.Fail(c, errors.INVALID_PARAMETER, msg)
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.COMMENT_ERROR, "用户未登录")
		return
	}

	var existPost models.Post
	if err := config.DB.First(&existPost, req.PostID).Error; err != nil {
		utils.Fail(c, errors.COMMENT_ERROR, "文章不存在")
		return
	}

	var existComment models.Comment

	if err := config.DB.Where("user_id", userId).Where("post_id", req.PostID).First(&existComment, "id = ?", req.ID).Error; err != nil {
		utils.Fail(c, errors.COMMENT_ERROR, "评论不存在")
		return
	}

	existComment.Content = req.Content

	if err := config.DB.Save(&existComment).Error; err != nil {
		logger.Log.Error(err)
		utils.Fail(c, errors.COMMENT_ERROR, "修改评论失败")
		return
	}
	utils.Success(c, "", "修改评论成功")
}

func DeleteComment(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, errors.COMMENT_ERROR, "用户未登录")
		return
	}
	id := c.Param("id")
	var existComment models.Comment
	if err := config.DB.Where("user_id", userId).Where("id", id).First(&existComment).Error; err != nil {
		utils.Fail(c, errors.COMMENT_ERROR, "评论没找到")
		return
	}
	if err := config.DB.Delete(&existComment).Error; err != nil {
		logger.Log.Error(err)
		utils.Fail(c, errors.COMMENT_ERROR, "删除失败")
		return
	}

	utils.Success(c, "", "删除成功")
	return
}
