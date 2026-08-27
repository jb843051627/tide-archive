package audit

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// AuditRule02 keeps one small, reusable archive rule close to its domain package.
type AuditRule02 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewAuditRule02(id, sessionID, label string, score float64) AuditRule02 {
	return AuditRule02{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v AuditRule02) Valid() error {
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

func (v AuditRule02) Clone() AuditRule02 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v AuditRule02) WithTag(tag string) AuditRule02 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v AuditRule02) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v AuditRule02) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v AuditRule02) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v AuditRule02) Escalate(delta float64) AuditRule02 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v AuditRule02) Deactivate() AuditRule02 { v.Active = false; return v }

func AuditRule02Batch(ctx context.Context, values []AuditRule02) ([]AuditRule02, error) {
	result := make([]AuditRule02, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func AuditRule02Labels(values []AuditRule02) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func AuditRule02Average(values []AuditRule02) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
