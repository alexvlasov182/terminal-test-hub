package repository

import (
	"context"
	"fmt"

	"github.com/alexvlasov182/terminal-test-hub/internal/domain"
	"github.com/jmoiron/sqlx"
)

// TerminalRepository handles database operations related to Terminals
type TerminalRepository struct {
	db *sqlx.DB
}

func NewTerminalRepository(db *sqlx.DB) *TerminalRepository {
	return &TerminalRepository{db: db}
}

// Create new terminal in the database
func (r *TerminalRepository) Create(ctx context.Context, terminal *domain.Terminal) error {
	query := `INSERT INTO terminals (serial_number, status, last_seen , metadata)
		VALUES ($1, $2, NOW(), $3)
		RETURNING id, created_at, last_seen`

	err := r.db.QueryRowxContext(ctx, query, terminal.SerialNumber, terminal.Status, terminal.Metadata).Scan(&terminal.ID, &terminal.CreatedAT, &terminal.LastSeen)

	if err != nil {
		return fmt.Errorf("can't find terminal %w", err)
	}

	return nil
}
