package export

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ExportRule05 keeps one small, reusable archive rule close to its domain package.
type ExportRule05 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewExportRule05(id, sessionID, label string, score float64) ExportRule05 {
	return ExportRule05{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v ExportRule05) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("export rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("export rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("export rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("export rule score is negative")
	}
	return nil
}

func (v ExportRule05) Clone() ExportRule05 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v ExportRule05) WithTag(tag string) ExportRule05 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v ExportRule05) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v ExportRule05) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v ExportRule05) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v ExportRule05) Escalate(delta float64) ExportRule05 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v ExportRule05) Deactivate() ExportRule05 { v.Active = false; return v }

func ExportRule05Batch(ctx context.Context, values []ExportRule05) ([]ExportRule05, error) {
	result := make([]ExportRule05, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func ExportRule05Labels(values []ExportRule05) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func ExportRule05Average(values []ExportRule05) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
