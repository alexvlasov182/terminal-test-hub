package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alexvlasov182/terminal-test-hub/internal/domain"
	"github.com/jmoiron/sqlx"
)

// JobRepository is responsible for working with jobs in the database
type JobRepository struct {
	db *sqlx.DB
}

// NewJobRepository creates a new reopository for jobs
func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create new job in the database
func (r *JobRepository) Create(ctx context.Context, job *domain.Job) error {
	query := `INSERT INTO jobs (terminal_id, type, payload, status) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, job.TerminalID, job.Type, job.Payload, job.Status).Scan(&job.ID, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return fmt.Errorf("cannot create job: %w", err)
	}
	return nil
}

// GetByID retrieves a job by its ID
func (r *JobRepository) GetByID(ctx context.Context, id string) (*domain.Job, error) {
	var job domain.Job
	query := `SELECT * FROM jobs WHERE id = $1`
	err := r.db.GetContext(ctx, &job, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("job with ID %s cant find", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving job with ID %s: %w", id, err)
	}
	return &job, nil
}

// ListByTerminal retrieves all jobs for a specific terminal
func (r *JobRepository) ListByTerminal(ctx context.Context, terminalID string) ([]domain.Job, error) {
	var jobs []domain.Job
	query := `SELECT * FROM jobs WHERE terminal_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &jobs, query, terminalID)
	if err != nil {
		return nil, fmt.Errorf("cannot retrives job from terminal: %w", err)
	}

	if jobs == nil {
		jobs = []domain.Job{}
	}

	return jobs, nil
}

// ListAll retrieves all jobs from the database (last 100)
func (r *JobRepository) ListAll(ctx context.Context) ([]domain.Job, error) {
	var jobs []domain.Job
	query := `SELECT * FROM jobs ORDER BY created_at DESC LIMIT 100`
	err := r.db.SelectContext(ctx, &jobs, query)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve all jobs: %w", err)
	}

	if jobs == nil {
		jobs = []domain.Job{}
	}

	return jobs, nil
}

// Update updates an existing job in the database (status and result)
func (r *JobRepository) Update(ctx context.Context, job *domain.Job) error {
	query := `UPDATE jobs SET status = $1, result = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query, job.Status, job.Result, job.ID).Scan(&job.UpdatedAt)

	if err != nil {
		return fmt.Errorf("cannot updates jobs: %w", err)
	}

	return nil
}

// UpdateStatus updates only the status of an existing job in the database (quicy operation)
func (r *JobRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE jobs SET status = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("cannot update status of the job: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("job with ID %s not found", id)
	}
	return nil
}
