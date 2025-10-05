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
	Colour                string  `json:"colour"`
	ColourAttachment      *string `json:"colour_attachment"`
	Ironmongery           string  `json:"ironmongery"`
	IronmongeryAttachment *string `json:"ironmongery_attachment"`
	UValue                *string `json:"u_value"`
	UValueAttachment      *string `json:"u_value_attachment"`
	GValue                *string `json:"g_value"`
	GValueAttachment      *string `json:"g_value_attachment"`
	Vents                 string  `json:"vents"`
	VentsAttachment       *string `json:"vents_attachment"`
	Acoustics             string  `json:"acoustics"`
	AcousticsAttachment   *string `json:"acoustics_attachment"`
	SBD                   string  `json:"sbd"`
	SBDAttachment         *string `json:"sbd_attachment"`
	PAS24                 string  `json:"pas24"`
	PAS24Attachment       *string `json:"pas24_attachment"`
	Restrictors           string  `json:"restrictors"`
	RestrictorsAttachment *string `json:"restrictors_attachment"`
	SpecialComments       string  `json:"special_comments"`
	AttachmentURL         string  `json:"attachment_url"`
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
		ProjectID:             projectID,
		Colour:                req.Colour,
		ColourAttachment:      req.ColourAttachment,
		Ironmongery:           req.Ironmongery,
		IronmongeryAttachment: req.IronmongeryAttachment,
		UValue:                req.UValue,
		UValueAttachment:      req.UValueAttachment,
		GValue:                req.GValue,
		GValueAttachment:      req.GValueAttachment,
		Vents:                 req.Vents,
		VentsAttachment:       req.VentsAttachment,
		Acoustics:             req.Acoustics,
		AcousticsAttachment:   req.AcousticsAttachment,
		SBD:                   req.SBD,
		SBDAttachment:         req.SBDAttachment,
		PAS24:                 req.PAS24,
		PAS24Attachment:       req.PAS24Attachment,
		Restrictors:           req.Restrictors,
		RestrictorsAttachment: req.RestrictorsAttachment,
		SpecialComments:       req.SpecialComments,
		AttachmentURL:         req.AttachmentURL,
		CreatedBy:             createdByID,
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