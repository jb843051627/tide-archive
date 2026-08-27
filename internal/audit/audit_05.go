package audit

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AuditRule05 keeps one small, reusable archive rule close to its domain package.
type AuditRule05 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewAuditRule05(id, sessionID, label string, score float64) AuditRule05 {
	return AuditRule05{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v AuditRule05) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("audit rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("audit rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("audit rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("audit rule score is negative")
	}
	return nil
}

func (v AuditRule05) Clone() AuditRule05 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v AuditRule05) WithTag(tag string) AuditRule05 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v AuditRule05) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v AuditRule05) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v AuditRule05) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v AuditRule05) Escalate(delta float64) AuditRule05 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v AuditRule05) Deactivate() AuditRule05 { v.Active = false; return v }

func AuditRule05Batch(ctx context.Context, values []AuditRule05) ([]AuditRule05, error) {
	result := make([]AuditRule05, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func AuditRule05Labels(values []AuditRule05) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func AuditRule05Average(values []AuditRule05) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
