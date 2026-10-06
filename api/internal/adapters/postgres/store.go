package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonydawson1000/home-archery/api/internal/application"
	"github.com/tonydawson1000/home-archery/api/internal/domain"
)

// Store persists Archer and Session aggregates via pgx.
type Store struct {
	pool *pgxpool.Pool
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("postgres: connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	if s != nil && s.pool != nil {
		s.pool.Close()
	}
}

func (s *Store) Archers() application.ArcherRepository   { return archerRepo{s} }
func (s *Store) Sessions() application.SessionRepository { return sessionRepo{s} }

func (s *Store) truncateScoring(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `TRUNCATE arrows, ends, sessions RESTART IDENTITY CASCADE`)
	return err
}

type archerRepo struct{ store *Store }

func (r archerRepo) Get(ctx context.Context, id string) (domain.Archer, error) {
	var archer domain.Archer
	err := r.store.pool.QueryRow(ctx,
		`SELECT id::text, display_name FROM archers WHERE id = $1::uuid`, id,
	).Scan(&archer.ID, &archer.DisplayName)
	if err != nil {
		return domain.Archer{}, mapNotFound(err, application.ErrArcherNotFound)
	}
	return archer, nil
}

func (r archerRepo) List(ctx context.Context) ([]domain.Archer, error) {
	rows, err := r.store.pool.Query(ctx,
		`SELECT id::text, display_name FROM archers ORDER BY display_name`)
	if err != nil {
		return nil, wrap(err)
	}
	defer rows.Close()
	out := make([]domain.Archer, 0)
	for rows.Next() {
		var archer domain.Archer
		if err := rows.Scan(&archer.ID, &archer.DisplayName); err != nil {
			return nil, wrap(err)
		}
		out = append(out, archer)
	}
	return out, wrap(rows.Err())
}

type sessionRepo struct{ store *Store }

func (r sessionRepo) Save(ctx context.Context, session domain.Session) error {
	tx, err := r.store.pool.Begin(ctx)
	if err != nil {
		return wrap(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sum := session.Summary()
	_, err = tx.Exec(ctx, `
		INSERT INTO sessions (
			id, archer_id, started_at, completed_at, status, arrows_per_end,
			arrow_count, total_score, hit_count, gold_count, x_count
		) VALUES (
			$1::uuid, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
		ON CONFLICT (id) DO UPDATE SET
			completed_at = EXCLUDED.completed_at,
			status = EXCLUDED.status,
			arrows_per_end = EXCLUDED.arrows_per_end,
			arrow_count = EXCLUDED.arrow_count,
			total_score = EXCLUDED.total_score,
			hit_count = EXCLUDED.hit_count,
			gold_count = EXCLUDED.gold_count,
			x_count = EXCLUDED.x_count`,
		session.ID, session.ArcherID, session.StartedAt, session.CompletedAt, string(session.Status),
		session.ArrowsPerEnd, sum.ArrowCount, sum.Total, sum.Hits, sum.Golds, sum.XCount,
	)
	if err != nil {
		return wrap(err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM ends WHERE session_id = $1::uuid`, session.ID); err != nil {
		return wrap(err)
	}

	for _, end := range session.Ends {
		var endID string
		err := tx.QueryRow(ctx,
			`INSERT INTO ends (id, session_id, end_number)
			 VALUES (gen_random_uuid(), $1::uuid, $2) RETURNING id::text`,
			session.ID, end.Number,
		).Scan(&endID)
		if err != nil {
			return wrap(err)
		}
		for _, arrow := range end.Arrows {
			_, err := tx.Exec(ctx, `
				INSERT INTO arrows (id, end_id, arrow_number, score_code, score_value, recorded_at)
				VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6)`,
				arrow.ID, endID, arrow.ArrowNumber, arrow.Score.Code, arrow.Score.NumericValue, arrow.RecordedAt,
			)
			if err != nil {
				return wrap(err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return wrap(err)
	}
	return nil
}

func (r sessionRepo) Get(ctx context.Context, id string) (domain.Session, error) {
	session, err := scanSession(ctx, r.store.pool, `
		SELECT id::text, archer_id::text, started_at, completed_at, status, arrows_per_end
		FROM sessions WHERE id = $1::uuid`, id)
	if err != nil {
		return domain.Session{}, mapNotFound(err, application.ErrSessionNotFound)
	}
	if err := loadEnds(ctx, r.store.pool, &session); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (r sessionRepo) ListByArcher(ctx context.Context, archerID string) ([]domain.Session, error) {
	rows, err := r.store.pool.Query(ctx, `
		SELECT id::text, archer_id::text, started_at, completed_at, status, arrows_per_end
		FROM sessions WHERE archer_id = $1::uuid
		ORDER BY started_at DESC`, archerID)
	if err != nil {
		return nil, mapNotFound(err, application.ErrArcherNotFound)
	}
	defer rows.Close()

	out := make([]domain.Session, 0)
	for rows.Next() {
		session, err := scanSessionRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, session)
	}
	if err := rows.Err(); err != nil {
		return nil, wrap(err)
	}
	for i := range out {
		if err := loadEnds(ctx, r.store.pool, &out[i]); err != nil {
			return nil, err
		}
	}
	return out, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSession(ctx context.Context, db *pgxpool.Pool, query string, args ...any) (domain.Session, error) {
	row := db.QueryRow(ctx, query, args...)
	return scanSessionRow(row)
}

func scanSessionRow(row rowScanner) (domain.Session, error) {
	var session domain.Session
	var status string
	if err := row.Scan(
		&session.ID, &session.ArcherID, &session.StartedAt, &session.CompletedAt, &status, &session.ArrowsPerEnd,
	); err != nil {
		return domain.Session{}, err
	}
	session.Status = domain.SessionStatus(status)
	session.StartedAt = session.StartedAt.UTC()
	if session.CompletedAt != nil {
		t := session.CompletedAt.UTC()
		session.CompletedAt = &t
	}
	return session, nil
}

func loadEnds(ctx context.Context, db *pgxpool.Pool, session *domain.Session) error {
	rows, err := db.Query(ctx, `
		SELECT id::text, end_number FROM ends
		WHERE session_id = $1::uuid ORDER BY end_number`, session.ID)
	if err != nil {
		return wrap(err)
	}
	defer rows.Close()

	type endRow struct {
		id     string
		number int
	}
	var ends []endRow
	for rows.Next() {
		var row endRow
		if err := rows.Scan(&row.id, &row.number); err != nil {
			return wrap(err)
		}
		ends = append(ends, row)
	}
	if err := rows.Err(); err != nil {
		return wrap(err)
	}

	session.Ends = make([]domain.End, 0, len(ends))
	for _, row := range ends {
		end := domain.End{Number: row.number}
		arrows, err := loadArrows(ctx, db, row.id)
		if err != nil {
			return err
		}
		end.Arrows = arrows
		session.Ends = append(session.Ends, end)
	}
	return nil
}

func loadArrows(ctx context.Context, db *pgxpool.Pool, endID string) ([]domain.Arrow, error) {
	rows, err := db.Query(ctx, `
		SELECT id::text, arrow_number, score_code, score_value, recorded_at
		FROM arrows WHERE end_id = $1::uuid ORDER BY arrow_number`, endID)
	if err != nil {
		return nil, wrap(err)
	}
	defer rows.Close()

	out := make([]domain.Arrow, 0)
	for rows.Next() {
		var arrow domain.Arrow
		var code string
		var value int
		if err := rows.Scan(&arrow.ID, &arrow.ArrowNumber, &code, &value, &arrow.RecordedAt); err != nil {
			return nil, wrap(err)
		}
		arrow.Score = domain.ArrowScore{Code: code, NumericValue: value}
		arrow.RecordedAt = arrow.RecordedAt.UTC()
		out = append(out, arrow)
	}
	return out, wrap(rows.Err())
}

func mapNotFound(err error, sentinel error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return sentinel
	}
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) && pgErr.SQLState() == "22P02" {
		return sentinel
	}
	return wrap(err)
}

func wrap(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("postgres: %w", err)
}
