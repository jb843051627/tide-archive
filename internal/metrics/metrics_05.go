package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MetricsRule05 keeps one small, reusable archive rule close to its domain package.
type MetricsRule05 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewMetricsRule05(id, sessionID, label string, score float64) MetricsRule05 {
	return MetricsRule05{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v MetricsRule05) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("metrics rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("metrics rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("metrics rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("metrics rule score is negative")
	}
	return nil
}

func (v MetricsRule05) Clone() MetricsRule05 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v MetricsRule05) WithTag(tag string) MetricsRule05 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v MetricsRule05) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v MetricsRule05) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v MetricsRule05) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v MetricsRule05) Escalate(delta float64) MetricsRule05 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v MetricsRule05) Deactivate() MetricsRule05 { v.Active = false; return v }

func MetricsRule05Batch(ctx context.Context, values []MetricsRule05) ([]MetricsRule05, error) {
	result := make([]MetricsRule05, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func MetricsRule05Labels(values []MetricsRule05) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func MetricsRule05Average(values []MetricsRule05) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
