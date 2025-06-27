package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct{}

func (r Repository) WrapTakeError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return errors.WithStack(err)
	}
}
