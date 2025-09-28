package controllers

import (
	"net/http"
	"strconv"

	"compass-backend/internal/models"
	"compass-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type SpecificationController struct {
	specService services.SpecificationService
}

func NewSpecificationController(specService services.SpecificationService) *SpecificationController {
	return &SpecificationController{
		specService: specService,
	}
}

type CreateSpecificationRequest struct {
	Colour          string   `json:"colour"`
	Ironmongery     string   `json:"ironmongery"`
	UValue          *float64 `json:"u_value"`
	GValue          *float64 `json:"g_value"`
	Vents           string   `json:"vents"`
	Acoustics       string   `json:"acoustics"`
	SBD             string   `json:"sbd"`
	PAS24           string   `json:"pas24"`
	Restrictors     string   `json:"restrictors"`
	SpecialComments string   `json:"special_comments"`
	AttachmentURL   string   `json:"attachment_url"`
}

func (c *SpecificationController) CreateSpecification(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req CreateSpecificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	createdBy, _ := ctx.Get("user_id")
	createdByID := createdBy.(uint64)

	spec := &models.ProjectSpecification{
		ProjectID:       projectID,
		Colour:          req.Colour,
		Ironmongery:     req.Ironmongery,
		UValue:          req.UValue,
		GValue:          req.GValue,
		Vents:           req.Vents,
		Acoustics:       req.Acoustics,
		SBD:             req.SBD,
		PAS24:           req.PAS24,
		Restrictors:     req.Restrictors,
		SpecialComments: req.SpecialComments,
		AttachmentURL:   req.AttachmentURL,
		CreatedBy:       createdByID,
	}

	err = c.specService.CreateSpecification(spec)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Specification created successfully",
		"specification": spec,
	})
}

func (c *SpecificationController) GetProjectSpecifications(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	specs, err := c.specService.GetProjectSpecifications(projectID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"specifications": specs})
}