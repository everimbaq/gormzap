package gormzap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

// Record is gormzap log record.
type Record struct {
	Message string
	Source  string
	Level   zapcore.Level

	Duration     time.Duration
	SQL          string
	RowsAffected int64
}

// RecordToFields func can encode gormzap Record into a slice of zap fields.
type RecordToFields func(r Record) []zapcore.Field

// DefaultRecordToFields is default encoder func for gormzap log records.
func DefaultRecordToFields(r Record) []zapcore.Field {
	// Note that Level field is ignored here, because it is handled outside
	// by zap itself.

	if r.SQL != "" {
		return []zapcore.Field{
			zap.String("source", SimplifyCodeSource(r.Source)),
			zap.Duration("dur", r.Duration),
			zap.String("query", r.SQL),
			zap.Int64("rows_affected", r.RowsAffected),
		}
	}

	return []zapcore.Field{zap.String("source", SimplifyCodeSource(r.Source))}
}

/*
simplify
/usr/local/share/GOHOME/gopath/src/zgcloud.com/backend/goim/store/reminder_store.go:35
to
store/reminder_store.go:35
*/
func SimplifyCodeSource(source string) string {
	idx := strings.LastIndexByte(source, '/')
	if idx == -1 {
		return source
	}
	idx = strings.LastIndexByte(source[:idx], '/')
	if idx == -1 {
		return source
	}
	return source[idx+1:]
}
