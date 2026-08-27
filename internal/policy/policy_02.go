package policy

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// PolicyRule02 keeps one small, reusable archive rule close to its domain package.
type PolicyRule02 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewPolicyRule02(id, sessionID, label string, score float64) PolicyRule02 {
	return PolicyRule02{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v PolicyRule02) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("policy rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("policy rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("policy rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("policy rule score is negative")
	}
	return nil
}

func (v PolicyRule02) Clone() PolicyRule02 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v PolicyRule02) WithTag(tag string) PolicyRule02 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v PolicyRule02) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v PolicyRule02) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v PolicyRule02) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v PolicyRule02) Escalate(delta float64) PolicyRule02 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v PolicyRule02) Deactivate() PolicyRule02 { v.Active = false; return v }

func PolicyRule02Batch(ctx context.Context, values []PolicyRule02) ([]PolicyRule02, error) {
	result := make([]PolicyRule02, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func PolicyRule02Labels(values []PolicyRule02) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func PolicyRule02Average(values []PolicyRule02) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
