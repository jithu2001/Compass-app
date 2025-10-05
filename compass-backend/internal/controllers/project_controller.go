package controllers

import (
	"net/http"
	"strconv"

	"compass-backend/internal/models"
	"compass-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	projectService services.ProjectService
}

func NewProjectController(projectService services.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

type CreateProjectRequest struct {
	ProjectName     string                              `json:"project_name" binding:"required"`
	CompanyName     string                              `json:"company_name"`
	CompanyAddress  string                              `json:"company_address"`
	ProjectType     models.ProjectType                  `json:"project_type" binding:"required,oneof=windows doors"`
	Specifications  []ProjectSpecificationRequest       `json:"specifications,omitempty"`
	RFIs           []ProjectRFIRequest                 `json:"rfis,omitempty"`
}

type ProjectSpecificationRequest struct {
	VersionNo             int     `json:"version_no" binding:"required"`
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

type ProjectRFIRequest struct {
	QuestionText string `json:"question_text" binding:"required"`
}

type UpdateProjectStatusRequest struct {
	Status models.ProjectStatus `json:"status" binding:"required,oneof=not_yet_started progress completed"`
}

func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	createdBy, _ := ctx.Get("user_id")
	createdByID := createdBy.(uint64)

	project := &models.Project{
		ProjectName:    req.ProjectName,
		CompanyName:    req.CompanyName,
		CompanyAddress: req.CompanyAddress,
		ProjectType:    req.ProjectType,
		CreatedBy:      createdByID,
	}

	// Prepare specifications if provided
	var specifications []models.ProjectSpecification
	for _, specReq := range req.Specifications {
		spec := models.ProjectSpecification{
			VersionNo:             specReq.VersionNo,
			Colour:                specReq.Colour,
			ColourAttachment:      specReq.ColourAttachment,
			Ironmongery:           specReq.Ironmongery,
			IronmongeryAttachment: specReq.IronmongeryAttachment,
			UValue:                specReq.UValue,
			UValueAttachment:      specReq.UValueAttachment,
			GValue:                specReq.GValue,
			GValueAttachment:      specReq.GValueAttachment,
			Vents:                 specReq.Vents,
			VentsAttachment:       specReq.VentsAttachment,
			Acoustics:             specReq.Acoustics,
			AcousticsAttachment:   specReq.AcousticsAttachment,
			SBD:                   specReq.SBD,
			SBDAttachment:         specReq.SBDAttachment,
			PAS24:                 specReq.PAS24,
			PAS24Attachment:       specReq.PAS24Attachment,
			Restrictors:           specReq.Restrictors,
			RestrictorsAttachment: specReq.RestrictorsAttachment,
			SpecialComments:       specReq.SpecialComments,
			AttachmentURL:         specReq.AttachmentURL,
			CreatedBy:             createdByID,
		}
		specifications = append(specifications, spec)
	}

	// Prepare RFIs if provided
	var rfis []models.ProjectRFI
	for _, rfiReq := range req.RFIs {
		rfi := models.ProjectRFI{
			QuestionText: rfiReq.QuestionText,
		}
		rfis = append(rfis, rfi)
	}

	err := c.projectService.CreateProjectWithDetails(project, specifications, rfis)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"project": project,
	})
}

func (c *ProjectController) GetProject(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := c.projectService.GetProject(projectID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"project": project})
}

func (c *ProjectController) ListProjects(ctx *gin.Context) {
	projects, err := c.projectService.ListProjects()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (c *ProjectController) UpdateProjectStatus(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req UpdateProjectStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	updatedBy, _ := ctx.Get("user_id")
	updatedByID := updatedBy.(uint64)

	err = c.projectService.UpdateProjectStatus(projectID, req.Status, updatedByID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Project status updated successfully"})
}

func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Get user role from context
	userRole, _ := ctx.Get("user_role")
	role := userRole.(models.UserRole)

	err = c.projectService.DeleteProject(projectID, role)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}