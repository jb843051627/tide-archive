package codec

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// CodecRule07 keeps one small, reusable archive rule close to its domain package.
type CodecRule07 struct {
	ID        string
	SessionID string
	Label     string
	Score     float64
	Active    bool
	CreatedAt time.Time
	Tags      []string
}

func NewCodecRule07(id, sessionID, label string, score float64) CodecRule07 {
	return CodecRule07{ID: strings.TrimSpace(id), SessionID: strings.TrimSpace(sessionID), Label: strings.TrimSpace(label), Score: score, Active: true, CreatedAt: time.Now().UTC()}
}

func (v CodecRule07) Valid() error {
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

func (v CodecRule07) Clone() CodecRule07 {
	v.Tags = append([]string(nil), v.Tags...)
	return v
}

func (v CodecRule07) WithTag(tag string) CodecRule07 {
	tag = strings.TrimSpace(tag)
	if tag != "" {
		v.Tags = append(v.Tags, tag)
	}
	return v
}

func (v CodecRule07) HasTag(tag string) bool {
	for _, current := range v.Tags {
		if current == tag {
			return true
		}
	}
	return false
}

func (v CodecRule07) Summary() string {
	return fmt.Sprintf("%s:%s:%.2f", v.SessionID, v.Label, v.Score)
}

func (v CodecRule07) Apply(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return v.Valid()
}

func (v CodecRule07) Escalate(delta float64) CodecRule07 {
	if delta > 0 {
		v.Score += delta
	}
	return v
}

func (v CodecRule07) Deactivate() CodecRule07 { v.Active = false; return v }

func CodecRule07Batch(ctx context.Context, values []CodecRule07) ([]CodecRule07, error) {
	result := make([]CodecRule07, 0, len(values))
	for _, value := range values {
		if err := value.Apply(ctx); err != nil {
			return nil, err
		}
		result = append(result, value.Clone())
	}
	return result, nil
}

func CodecRule07Labels(values []CodecRule07) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		if value.Label != "" {
			labels = append(labels, value.Label)
		}
	}
	return labels
}

func CodecRule07Average(values []CodecRule07) float64 {
	if len(values) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range values {
		total += value.Score
	}
	return total / float64(len(values))
}
