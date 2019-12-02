package factory

import (
	"time"

	"github.com/kimitzu/kimitzu-go/repo"
)

func NewAPITime(t time.Time) *repo.APITime {
	return repo.NewAPITime(t)
}
