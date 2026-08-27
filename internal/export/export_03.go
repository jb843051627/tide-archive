package export

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ExportRule03 keeps one small, reusable archive rule close to its domain package.
type ExportRule03 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewExportRule03(id, sessionID, label string, score float64) ExportRule03 {
	return ExportRule03{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v ExportRule03) Valid() error {
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

func (v ExportRule03) Clone() ExportRule03 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v ExportRule03) WithTag(tag string) ExportRule03 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v ExportRule03) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v ExportRule03) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v ExportRule03) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v ExportRule03) Escalate(delta float64) ExportRule03 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v ExportRule03) Deactivate() ExportRule03 { v.Active = false; return v }

func ExportRule03Batch(ctx context.Context, values []ExportRule03) ([]ExportRule03, error) {
	result := make([]ExportRule03, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func ExportRule03Labels(values []ExportRule03) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func ExportRule03Average(values []ExportRule03) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
