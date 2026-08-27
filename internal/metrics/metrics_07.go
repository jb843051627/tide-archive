package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MetricsRule07 keeps one small, reusable archive rule close to its domain package.
type MetricsRule07 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewMetricsRule07(id, sessionID, label string, score float64) MetricsRule07 {
	return MetricsRule07{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v MetricsRule07) Valid() error {
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

func (v MetricsRule07) Clone() MetricsRule07 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v MetricsRule07) WithTag(tag string) MetricsRule07 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v MetricsRule07) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v MetricsRule07) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v MetricsRule07) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v MetricsRule07) Escalate(delta float64) MetricsRule07 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v MetricsRule07) Deactivate() MetricsRule07 { v.Active = false; return v }

func MetricsRule07Batch(ctx context.Context, values []MetricsRule07) ([]MetricsRule07, error) {
	result := make([]MetricsRule07, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func MetricsRule07Labels(values []MetricsRule07) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func MetricsRule07Average(values []MetricsRule07) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
