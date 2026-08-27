package worker

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// WorkerRule04 keeps one small, reusable archive rule close to its domain package.
type WorkerRule04 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewWorkerRule04(id, sessionID, label string, score float64) WorkerRule04 {
	return WorkerRule04{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v WorkerRule04) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("worker rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("worker rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("worker rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("worker rule score is negative")
	}
	return nil
}

func (v WorkerRule04) Clone() WorkerRule04 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v WorkerRule04) WithTag(tag string) WorkerRule04 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v WorkerRule04) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v WorkerRule04) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v WorkerRule04) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v WorkerRule04) Escalate(delta float64) WorkerRule04 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v WorkerRule04) Deactivate() WorkerRule04 { v.Active = false; return v }

func WorkerRule04Batch(ctx context.Context, values []WorkerRule04) ([]WorkerRule04, error) {
	result := make([]WorkerRule04, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func WorkerRule04Labels(values []WorkerRule04) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func WorkerRule04Average(values []WorkerRule04) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
