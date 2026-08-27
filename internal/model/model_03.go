package model

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ModelRule03 keeps one small, reusable archive rule close to its domain package.
type ModelRule03 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewModelRule03(id, sessionID, label string, score float64) ModelRule03 {
	return ModelRule03{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v ModelRule03) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("model rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("model rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("model rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("model rule score is negative")
	}
	return nil
}

func (v ModelRule03) Clone() ModelRule03 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v ModelRule03) WithTag(tag string) ModelRule03 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v ModelRule03) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v ModelRule03) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v ModelRule03) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v ModelRule03) Escalate(delta float64) ModelRule03 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v ModelRule03) Deactivate() ModelRule03 { v.Active = false; return v }

func ModelRule03Batch(ctx context.Context, values []ModelRule03) ([]ModelRule03, error) {
	result := make([]ModelRule03, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func ModelRule03Labels(values []ModelRule03) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func ModelRule03Average(values []ModelRule03) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
