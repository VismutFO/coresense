package connection

import (
	"gorm.io/driver/postgres"

	"coresense/pkg/common/utils"
)

var errCodes = map[string]error{
	"23505": utils.ErrDuplicatedKey,
}

type dialector struct {
	*postgres.Dialector
}

func newDialector(url string) dialector {
	return dialector{Dialector: &postgres.Dialector{Config: &postgres.Config{DSN: url}}}
}

func (d dialector) Translate(err error) error {
	/*
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				if translatedErr, ok := errCodes[pgErr.Code]; ok {
					return translatedErr
				}
			}
		}
	*/
	return err
}
