package repository

import (
	"context"
	"time"

	"github.com/sysu-ecnc-dev/shift-manager/backend/internal/domain"
)

func (r *Repository) GetAllSchedulePlans() ([]*domain.SchedulePlan, error) {
	query := `
		SELECT 
			id, 
			name, 
			description, 
			submission_start_time, 
			submission_end_time, 
			active_start_time, 
			active_end_time,
			(SELECT name FROM schedule_template_meta WHERE id = schedule_template_id) AS schedule_template_name,
			created_at, 
			version
		FROM schedule_plans
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Database.QueryTimeout)*time.Second)
	defer cancel()

	rows, err := r.dbpool.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []*domain.SchedulePlan{}
	for rows.Next() {
		var plan domain.SchedulePlan
		dst := []any{
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.SubmissionStartTime,
			&plan.SubmissionEndTime,
			&plan.ActiveStartTime,
			&plan.ActiveEndTime,
			&plan.ScheduleTemplateName,
			&plan.CreatedAt,
			&plan.Version,
		}
		if err := rows.Scan(dst...); err != nil {
			return nil, err
		}
		plans = append(plans, &plan)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

func (r *Repository) UpdateSchedulePlan(plan *domain.SchedulePlan) error {
	// 最好不要让用户更新所使用的模板，不然后续会带来很多麻烦
	query := `
		UPDATE schedule_plans 
		SET
			name = $1,
			description = $2,
			submission_start_time = $3,
			submission_end_time = $4,
			active_start_time = $5,
			active_end_time = $6,
			version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Database.QueryTimeout)*time.Second)
	defer cancel()

	params := []any{
		plan.Name,
		plan.Description,
		plan.SubmissionStartTime,
		plan.SubmissionEndTime,
		plan.ActiveStartTime,
		plan.ActiveEndTime,
		plan.ID,
		plan.Version,
	}

	if err := r.dbpool.QueryRowContext(ctx, query, params...).Scan(&plan.Version); err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertSchedulePlan(plan *domain.SchedulePlan) error {
	query := `
		INSERT INTO schedule_plans (
			name,
			description,
			submission_start_time,
			submission_end_time,
			active_start_time,
			active_end_time,
			schedule_template_id
		) VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM schedule_template_meta WHERE name = $7))
		RETURNING id, created_at, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Database.QueryTimeout)*time.Second)
	defer cancel()

	params := []any{
		plan.Name,
		plan.Description,
		plan.SubmissionStartTime,
		plan.SubmissionEndTime,
		plan.ActiveStartTime,
		plan.ActiveEndTime,
		plan.ScheduleTemplateName,
	}
	dst := []any{&plan.ID, &plan.CreatedAt, &plan.Version}
	if err := r.dbpool.QueryRowContext(ctx, query, params...).Scan(dst...); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetSchedulePlanByID(id int64) (*domain.SchedulePlan, error) {
	query := `
		SELECT 
			name, 
			description, 
			submission_start_time, 
			submission_end_time, 
			active_start_time, 
			active_end_time, 
			created_at, 
			version
		FROM schedule_plans
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.cfg.Database.QueryTimeout)*time.Second)
	defer cancel()

	plan := &domain.SchedulePlan{
		ID: id,
	}

	dst := []any{
		&plan.Name,
		&plan.Description,
		&plan.SubmissionStartTime,
		&plan.SubmissionEndTime,
		&plan.ActiveStartTime,
		&plan.ActiveEndTime,
		&plan.CreatedAt,
		&plan.Version,
	}

	if err := r.dbpool.QueryRowContext(ctx, query, id).Scan(dst...); err != nil {
		return nil, err
	}

	return plan, nil
}
