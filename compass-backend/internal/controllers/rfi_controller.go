package controllers

import (
	"net/http"
	"strconv"

	"compass-backend/internal/models"
	"compass-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type RFIController struct {
	rfiService services.RFIService
}

func NewRFIController(rfiService services.RFIService) *RFIController {
	return &RFIController{
		rfiService: rfiService,
	}
}

type CreateRFIRequest struct {
	QuestionText string `json:"question_text" binding:"required"`
}

type AnswerRFIRequest struct {
	AnswerValue models.AnswerValue `json:"answer_value" binding:"required,oneof=yes no"`
}

func (c *RFIController) CreateRFI(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req CreateRFIRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defaultAnswer := models.AnswerNo
	rfi := &models.ProjectRFI{
		ProjectID:    projectID,
		QuestionText: req.QuestionText,
		AnswerValue:  &defaultAnswer,
	}

	err = c.rfiService.CreateRFI(rfi)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "RFI created successfully",
		"rfi":     rfi,
	})
}

func (c *RFIController) AnswerRFI(ctx *gin.Context) {
	rfiIDStr := ctx.Param("id")
	rfiID, err := strconv.ParseUint(rfiIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid RFI ID"})
		return
	}

	var req AnswerRFIRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	answeredBy, _ := ctx.Get("user_id")
	answeredByID := answeredBy.(uint64)

	err = c.rfiService.AnswerRFI(rfiID, req.AnswerValue, answeredByID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "RFI answered successfully"})
}

func (c *RFIController) GetProjectRFIs(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	rfis, err := c.rfiService.GetProjectRFIs(projectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"rfis": rfis})
}