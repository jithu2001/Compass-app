package controllers

import (
	"net/http"
	"strconv"

	"compass-backend/internal/models"
	"compass-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type CreateUserRequest struct {
	FullName string          `json:"full_name" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password"`
	Role     models.UserRole `json:"role" binding:"required,oneof=admin user"`
}

type UpdateUserStatusRequest struct {
	Status models.AccountStatus `json:"status" binding:"required,oneof=pending active disabled"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the admin user ID from context
	invitedBy, _ := ctx.Get("user_id")
	invitedByID := invitedBy.(uint64)

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Role:     req.Role,
	}

	err := c.userService.CreateUser(user, req.Password, invitedByID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"user_id":   user.UserID,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
			"status":    user.AccountStatus,
		},
	})
}

func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.userService.UpdateUserStatus(userID, req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

func (c *UserController) ListUsers(ctx *gin.Context) {
	users, err := c.userService.ListUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"users": users})
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that new password and confirm password match
	if req.NewPassword != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirm password do not match"})
		return
	}

	// Get user ID from context
	userID, _ := ctx.Get("user_id")
	userIDValue := userID.(uint64)

	err := c.userService.ChangePassword(userIDValue, req.CurrentPassword, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.userService.ResetPassword(userID, req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}