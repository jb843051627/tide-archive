package audit

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AuditRule06 keeps one small, reusable archive rule close to its domain package.
type AuditRule06 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewAuditRule06(id, sessionID, label string, score float64) AuditRule06 {
	return AuditRule06{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v AuditRule06) Valid() error {
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

func (v AuditRule06) Clone() AuditRule06 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v AuditRule06) WithTag(tag string) AuditRule06 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v AuditRule06) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v AuditRule06) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v AuditRule06) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v AuditRule06) Escalate(delta float64) AuditRule06 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v AuditRule06) Deactivate() AuditRule06 { v.Active = false; return v }

func AuditRule06Batch(ctx context.Context, values []AuditRule06) ([]AuditRule06, error) {
	result := make([]AuditRule06, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func AuditRule06Labels(values []AuditRule06) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func AuditRule06Average(values []AuditRule06) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
