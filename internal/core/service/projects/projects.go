package projects

import (
	"context"
	"dashboard/internal/application/dto"
	"dashboard/internal/core/models"
	"dashboard/internal/core/ports"
)

type ProjectsService struct {
	repo ports.ProjectRepository
}

func NewProjectsService(repo ports.ProjectRepository) *ProjectsService {
	return &ProjectsService{repo: repo}
}

func (r *ProjectsService) GetProjectById(ctx context.Context, project_id int64) (models.ProjectsContent, error) {
	ProjectDbo, err := r.repo.GetProjectById(ctx, project_id)
	if err != nil {
		return models.ProjectsContent{}, err
	}
	project := ProjectDbo.ToValue()
	return project, nil
}

func (r *ProjectsService) GetAllProjects(ctx context.Context) ([]models.ProjectsContent, error) {
	ProjectsDbo, err := r.repo.GetAllProjects(ctx)
	if err != nil {
		return []models.ProjectsContent{}, err
	}
	var projects []models.ProjectsContent
	for _, project := range ProjectsDbo {
		projects = append(projects, project.ToValue())
	}
	return projects, nil
}

func (r *ProjectsService) CreateProject(ctx context.Context, project models.Project) (int64, error) {
	projectDbo := dto.ToProjectDbo(project)
	project_id, err := r.repo.CreateProject(ctx, projectDbo)
	if err != nil {
		return project_id, err
	}
	return project_id, nil
}

func (r *ProjectsService) UpdateProject(ctx context.Context, project models.Project) (models.Project, error) {
	projectDbo := dto.ToProjectDbo(project)
	newProject, err := r.repo.UpdateProject(ctx, projectDbo)
	if err != nil {
		return models.Project{}, err
	}
	return newProject.ToValue(), nil
}

func (r *ProjectsService) DeleteProject(ctx context.Context, project_id int64) error {
	return r.repo.DeleteProject(ctx, project_id)
}
