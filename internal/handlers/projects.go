package handlers

import (
	"net/http"

	"github.com/4planet/backend/internal/config"
	"github.com/4planet/backend/pkg/pagination"
	"github.com/4planet/backend/pkg/projects"
	"github.com/gin-gonic/gin"
)

type ProjectsHandler struct {
	projectsService *projects.Service
	config          *config.Config
}

func NewProjectsHandler(projectsService *projects.Service, config *config.Config) *ProjectsHandler {
	return &ProjectsHandler{
		projectsService: projectsService,
		config:          config,
	}
}

func (h *ProjectsHandler) GetProjects(c *gin.Context) {
	// Extract pagination parameters
	params := pagination.ExtractPagination(c)

	projects, total, err := h.projectsService.GetProjects(params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	response := pagination.NewPaginatedResponse(projects, total, params)
	c.JSON(http.StatusOK, response)
}

func (h *ProjectsHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	project, err := h.projectsService.GetProjectByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		return
	}
	c.JSON(http.StatusOK, project)
}
