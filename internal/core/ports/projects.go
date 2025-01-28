package ports

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
)

type ProjectsService interface {
	GetProjectById(context.Context, int64) (models.ProjectsContent, error)
	GetAllProjects(context.Context) ([]models.ProjectsContent, error)
	CreateProject(context.Context, models.Project) (int64, error)
	UpdateProject(context.Context, models.Project) (models.Project, error)
	DeleteProject(context.Context, int64) error
}

type ProjectRepository interface {
	GetProjectById(context.Context, int64) (*dto.ProjectsContentDbo, error)
	GetAllProjects(context.Context) ([]dto.ProjectsContentDbo, error)
	CreateProject(context.Context, dto.ProjectDbo) (int64, error)
	UpdateProject(context.Context, dto.ProjectDbo) (*dto.ProjectDbo, error)
	DeleteProject(context.Context, int64) error
}
