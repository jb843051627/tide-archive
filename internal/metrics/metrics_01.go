package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MetricsRule01 keeps one small, reusable archive rule close to its domain package.
type MetricsRule01 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewMetricsRule01(id, sessionID, label string, score float64) MetricsRule01 {
	return MetricsRule01{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v MetricsRule01) Valid() error {
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

func (v MetricsRule01) Clone() MetricsRule01 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v MetricsRule01) WithTag(tag string) MetricsRule01 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v MetricsRule01) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v MetricsRule01) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v MetricsRule01) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v MetricsRule01) Escalate(delta float64) MetricsRule01 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v MetricsRule01) Deactivate() MetricsRule01 { v.Active = false; return v }

func MetricsRule01Batch(ctx context.Context, values []MetricsRule01) ([]MetricsRule01, error) {
	result := make([]MetricsRule01, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func MetricsRule01Labels(values []MetricsRule01) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func MetricsRule01Average(values []MetricsRule01) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
