package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
)

func (d *DashboardService) GetProjectById(project_id int64) (models.ProjectsContent, error) {
	ctx := context.Background()
	d.log.Log("info", "getting project by id", logger.WithInt64Attr("project_id", project_id))
	project, err := d.project.GetProjectById(ctx, project_id)
	if err != nil {
		d.log.Log("error", "error while getting project by id", logger.WithErrAttr(err))
		return project, err
	}
	d.log.Log("info", "project successfully get", logger.WithStrAttr("project_id", project.Project.Name))
	return project, nil
}

func (d *DashboardService) GetAllProjects() ([]models.ProjectsContent, error) {
	ctx := context.Background()
	d.log.Log("info", "getting all projects")
	projects, err := d.project.GetAllProjects(ctx)
	if err != nil {
		d.log.Log("error", "error while fetching projects", logger.WithErrAttr(err))
		return projects, err
	}
	d.log.Log(
		"info",
		"projects successfully fetched",
		logger.WithInt64Attr("project_count", int64(len(projects))),
	)
	return projects, nil
}

func (d *DashboardService) CreateProject(project models.Project) (int64, error) {
	ctx := context.Background()
	d.log.Log("info", "creating project")
	project_id, err := d.project.CreateProject(ctx, project)
	if err != nil {
		d.log.Log("error", "error while creating project", logger.WithErrAttr(err))
		return project_id, err
	}
	d.log.Log("info", "project successfully created", logger.WithStrAttr("project_id", project.Name))
	return project_id, nil
}

func (d *DashboardService) UpdateProject(project models.Project) (models.Project, error) {
	ctx := context.Background()
	d.log.Log("info", "updating project", logger.WithInt64Attr("project_id", project.Id))
	project, err := d.project.UpdateProject(ctx, project)
	if err != nil {
		d.log.Log("error", "error while updating project", logger.WithErrAttr(err))
		return project, err
	}
	d.log.Log("info", "project successfully updated", logger.WithStrAttr("project_id", project.Name))
	return project, nil
}

func (d *DashboardService) DeleteProject(project_id int64) error {
	ctx := context.Background()
	d.log.Log("info", "deleting project", logger.WithInt64Attr("project_id", project_id))
	err := d.project.DeleteProject(ctx, project_id)
	if err != nil {
		d.log.Log("error", "error while deleting project", logger.WithErrAttr(err))
		return err
	}
	d.log.Log("info", "project successfully deleted", logger.WithInt64Attr("project_id", project_id))
	return nil
}
