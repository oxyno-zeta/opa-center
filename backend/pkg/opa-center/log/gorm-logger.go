package log

import (
	"context"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

type gormLogger struct {
	logger Logger
}

func (gl *gormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return gl
}

func (gl *gormLogger) Info(ctx context.Context, v string, rest ...interface{}) {
	val := []interface{}{v}
	val = append(val, rest...)
	gl.logger.Info(val)
}

func (gl *gormLogger) Warn(ctx context.Context, v string, rest ...interface{}) {
	val := []interface{}{v}
	val = append(val, rest...)
	gl.logger.Warn(val)
}

func (gl *gormLogger) Error(ctx context.Context, v string, rest ...interface{}) {
	val := []interface{}{v}
	val = append(val, rest...)
	gl.logger.Error(val)
}

func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	gl.logger.WithField("rows", rows).WithError(err).Debug(sql)
}
