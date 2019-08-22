package factory

import (
	"time"

	"github.com/djali-foundation/djali-go/repo"
)

func NewAPITime(t time.Time) *repo.APITime {
	return repo.NewAPITime(t)
}
