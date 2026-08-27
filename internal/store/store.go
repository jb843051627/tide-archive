package store

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type Repository struct{ db *sql.DB }

func Open(path string) (*Repository, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`PRAGMA foreign_keys=ON; PRAGMA journal_mode=WAL; CREATE TABLE IF NOT EXISTS sessions(id TEXT PRIMARY KEY, site TEXT NOT NULL, state TEXT NOT NULL, created_at TEXT NOT NULL); CREATE TABLE IF NOT EXISTS fragments(id TEXT PRIMARY KEY, session_id TEXT NOT NULL, label TEXT NOT NULL, score REAL NOT NULL, reviewed INTEGER NOT NULL DEFAULT 0); CREATE TABLE IF NOT EXISTS audit_events(id INTEGER PRIMARY KEY AUTOINCREMENT, session_id TEXT NOT NULL, message TEXT NOT NULL)`); err != nil {
		db.Close()
		return nil, err
	}
	return &Repository{db: db}, nil
}
func (r *Repository) Close() error { return r.db.Close() }
func (r *Repository) CreateSession(ctx context.Context, id, site, created string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO sessions(id,site,state,created_at) VALUES(?,?,?,?)`, id, site, "draft", created)
	return err
}
func (r *Repository) SetState(ctx context.Context, id, state string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE sessions SET state=? WHERE id=?`, state, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("session %s not found", id)
	}
	return nil
}
func (r *Repository) AddFragment(ctx context.Context, id, sessionID, label string, score float64) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO fragments(id,session_id,label,score) VALUES(?,?,?,?)`, id, sessionID, label, score)
	return err
}
func (r *Repository) ReviewFragment(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE fragments SET reviewed=1 WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("fragment %s not found", id)
	}
	return nil
}
func (r *Repository) CountUnreviewed(ctx context.Context, sessionID string) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM fragments WHERE session_id=? AND reviewed=0`, sessionID).Scan(&n)
	return n, err
}
func (r *Repository) AppendAudit(ctx context.Context, sessionID, message string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO audit_events(session_id,message) VALUES(?,?)`, sessionID, message)
	return err
}
