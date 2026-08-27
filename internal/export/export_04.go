package export

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ExportRule04 keeps one small, reusable archive rule close to its domain package.
type ExportRule04 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewExportRule04(id, sessionID, label string, score float64) ExportRule04 {
	return ExportRule04{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v ExportRule04) Valid() error {
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

func (v ExportRule04) Clone() ExportRule04 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v ExportRule04) WithTag(tag string) ExportRule04 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v ExportRule04) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v ExportRule04) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v ExportRule04) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v ExportRule04) Escalate(delta float64) ExportRule04 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v ExportRule04) Deactivate() ExportRule04 { v.Active = false; return v }

func ExportRule04Batch(ctx context.Context, values []ExportRule04) ([]ExportRule04, error) {
	result := make([]ExportRule04, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func ExportRule04Labels(values []ExportRule04) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func ExportRule04Average(values []ExportRule04) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
