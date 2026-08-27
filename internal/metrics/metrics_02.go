package metrics

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// MetricsRule02 keeps one small, reusable archive rule close to its domain package.
type MetricsRule02 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewMetricsRule02(id, sessionID, label string, score float64) MetricsRule02 {
	return MetricsRule02{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v MetricsRule02) Valid() error {
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

func (v MetricsRule02) Clone() MetricsRule02 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v MetricsRule02) WithTag(tag string) MetricsRule02 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v MetricsRule02) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v MetricsRule02) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v MetricsRule02) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v MetricsRule02) Escalate(delta float64) MetricsRule02 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v MetricsRule02) Deactivate() MetricsRule02 { v.Active = false; return v }

func MetricsRule02Batch(ctx context.Context, values []MetricsRule02) ([]MetricsRule02, error) {
	result := make([]MetricsRule02, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func MetricsRule02Labels(values []MetricsRule02) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func MetricsRule02Average(values []MetricsRule02) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
