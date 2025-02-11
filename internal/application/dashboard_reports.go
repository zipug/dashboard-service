package application

import (
	"context"
	logger "dashboard/internal/common/service/logger/zerolog"
	"dashboard/internal/core/models"
)

func (d *DashboardService) GetAllReports() ([]models.Report, error) {
	ctx := context.Background()
	d.log.Log("info", "get all reports")
	reports, err := d.report.GetAllReports(ctx)
	if err != nil {
		d.log.Log("error", "could not get all reports")
		return nil, err
	}
	d.log.Log("info", "all reports found")
	return reports, nil
}

func (d *DashboardService) ExecuteReport(report_id, user_id int64, date_from, date_to string) error {
	ctx := context.Background()
	d.log.Log("info", "execute report", logger.WithInt64Attr("user_id", user_id), logger.WithInt64Attr("report_id", report_id))
	if err := d.report.ExecuteReport(ctx, report_id, user_id, date_from, date_to); err != nil {
		d.log.Log("error", "could not execute report", logger.WithInt64Attr("user_id", user_id), logger.WithInt64Attr("report_id", report_id))
		return err
	}
	d.log.Log("info", "report executed", logger.WithInt64Attr("user_id", user_id), logger.WithInt64Attr("report_id", report_id))
	return nil
}
