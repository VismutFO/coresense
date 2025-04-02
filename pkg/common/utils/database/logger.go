package database

import (
	"context"
	"coresense/pkg/core/model/constants"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
	"time"

	"coresense/pkg/common/config"
)

type Logger struct {
	logger zerolog.Logger
	config config.DatabaseLogger
}

func (l Logger) LogMode(_ logger.LogLevel) logger.Interface {
	return &l
}

func (l Logger) Info(ctx context.Context, s string, i ...interface{}) {
	l.logger.Info().Ctx(ctx).Dict(constants.Args, zerolog.Dict().
		Interface("data", i).Caller(1)).Msg(s)
}

func (l Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.logger.Warn().Ctx(ctx).Dict(constants.Args, zerolog.Dict().
		Interface("data", i).Caller(1)).Msg(s)
}

func (l Logger) Error(ctx context.Context, s string, i ...interface{}) {
	l.logger.Error().Ctx(ctx).Dict(constants.Args, zerolog.Dict().
		Interface("data", i).Caller(1)).Msg(s)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	timeDuration := time.Since(begin)

	var (
		event *zerolog.Event
		dict  = zerolog.Dict().Str("time", timeDuration.String()).Caller(3)
	)
	switch {
	case err != nil:
		event = l.logger.Err(err)
	default:
		event = l.logger.Trace()
	}

	sql, rowsAffected := fc()
	if rowsAffected >= 0 {
		dict.Int64("rows", rowsAffected)
	}
	event.Ctx(ctx).Dict(constants.Args, dict).Msg(sql)
}

func NewLogger(logger zerolog.Logger, config config.DatabaseLogger) *Logger {
	return &Logger{
		logger: logger,
		config: config,
	}
}
