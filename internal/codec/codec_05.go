package codec

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// CodecRule05 keeps one small, reusable archive rule close to its domain package.
type CodecRule05 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewCodecRule05(id, sessionID, label string, score float64) CodecRule05 {
	return CodecRule05{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v CodecRule05) Valid() error {
	if v.ID == "" {
		return fmt.Errorf("codec rule id is empty")
	}
	if v.SessionID == "" {
		return fmt.Errorf("codec rule session is empty")
	}
	if v.Label == "" {
		return fmt.Errorf("codec rule label is empty")
	}
	if v.Score < 0 {
		return fmt.Errorf("codec rule score is negative")
	}
	return nil
}

func (v CodecRule05) Clone() CodecRule05 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v CodecRule05) WithTag(tag string) CodecRule05 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v CodecRule05) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v CodecRule05) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v CodecRule05) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v CodecRule05) Escalate(delta float64) CodecRule05 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v CodecRule05) Deactivate() CodecRule05 { v.Active = false; return v }

func CodecRule05Batch(ctx context.Context, values []CodecRule05) ([]CodecRule05, error) {
	result := make([]CodecRule05, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func CodecRule05Labels(values []CodecRule05) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func CodecRule05Average(values []CodecRule05) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
