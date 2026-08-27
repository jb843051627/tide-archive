package validation

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ValidationRule04 keeps one small, reusable archive rule close to its domain package.
type ValidationRule04 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewValidationRule04(id, sessionID, label string, score float64) ValidationRule04 {
	return ValidationRule04{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v ValidationRule04) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("validation rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("validation rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("validation rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("validation rule score is negative")
	}
	return nil
}

func (v ValidationRule04) Clone() ValidationRule04 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v ValidationRule04) WithTag(tag string) ValidationRule04 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v ValidationRule04) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v ValidationRule04) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v ValidationRule04) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v ValidationRule04) Escalate(delta float64) ValidationRule04 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v ValidationRule04) Deactivate() ValidationRule04 { v.Active = false; return v }

func ValidationRule04Batch(ctx context.Context, values []ValidationRule04) ([]ValidationRule04, error) {
	result := make([]ValidationRule04, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func ValidationRule04Labels(values []ValidationRule04) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func ValidationRule04Average(values []ValidationRule04) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
