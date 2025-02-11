package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
)

func (d *DashboardService) GetGeneratedReportById(generated_report_id, user_id int64) (models.GeneratedReport, error) {
	ctx := context.Background()
	d.log.Log("info", "get generated_report by id", logger.WithInt64Attr("generated_report_id", generated_report_id))
	generated_report, err := d.generated_report.GetGeneratedReportById(ctx, generated_report_id, user_id)
	if err != nil {
		d.log.Log("error", "could not get generated_report by id", logger.WithInt64Attr("generated_report_id", generated_report_id))
		return generated_report, err
	}
	d.log.Log("info", "generated_report found by id", logger.WithInt64Attr("generated_report_id", generated_report_id))
	return generated_report, nil
}

func (d *DashboardService) GetAllGeneratedReports(user_id int64) ([]models.GeneratedReport, error) {
	ctx := context.Background()
	d.log.Log("info", "get all generated_reports by user_id", logger.WithInt64Attr("user_id", user_id))
	generated_reports, err := d.generated_report.GetAllGeneratedReports(ctx, user_id)
	if err != nil {
		d.log.Log("error", "could not get all generated_reports by user_id", logger.WithInt64Attr("user_id", user_id))
		return nil, err
	}
	d.log.Log("info", "all generated_reports found by user_id", logger.WithInt64Attr("user_id", user_id))
	return generated_reports, nil
}

func (d *DashboardService) DeleteGeneratedReportById(generated_report_id, user_id int64) error {
	ctx := context.Background()
	d.log.Log("info", "delete generated_report by id", logger.WithInt64Attr("generated_report_id", generated_report_id))
	if err := d.generated_report.DeleteGeneratedReport(ctx, generated_report_id, user_id); err != nil {
		d.log.Log("error", "could not delete generated_report by id", logger.WithInt64Attr("generated_report_id", generated_report_id))
		return err
	}
	d.log.Log("info", "generated_report deleted by id", logger.WithInt64Attr("generated_report_id", generated_report_id))
	return nil
}
