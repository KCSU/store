// JSON logs with a severity on Cloud Run, plain text elsewhere
package logging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var onCloudRun = os.Getenv("K_SERVICE") != "" || os.Getenv("CLOUD_RUN_JOB") != ""

func Infof(format string, args ...interface{})  { write("INFO", format, args...) }
func Warnf(format string, args ...interface{})  { write("WARNING", format, args...) }
func Errorf(format string, args ...interface{}) { write("ERROR", format, args...) }

func Fatalf(format string, args ...interface{}) {
	write("CRITICAL", format, args...)
	os.Exit(1)
}

func write(severity string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	if !onCloudRun {
		log.Printf("%s: %s", severity, message)
		return
	}
	line, _ := json.Marshal(struct {
		Severity string `json:"severity"`
		Message  string `json:"message"`
	}{severity, message})
	os.Stdout.Write(append(line, '\n'))
}

const slowQuery = 200 * time.Millisecond

type Gorm struct {
	Level gormlogger.LogLevel
}

func (g Gorm) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return Gorm{level}
}

func (g Gorm) Info(_ context.Context, format string, args ...interface{}) {
	if g.Level >= gormlogger.Info {
		Infof(format, args...)
	}
}

func (g Gorm) Warn(_ context.Context, format string, args ...interface{}) {
	if g.Level >= gormlogger.Warn {
		Warnf(format, args...)
	}
}

func (g Gorm) Error(_ context.Context, format string, args ...interface{}) {
	if g.Level >= gormlogger.Error {
		Errorf(format, args...)
	}
}

func (g Gorm) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && g.Level >= gormlogger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		Errorf("SQL failed: %v [%s, %d rows] %s", err, elapsed, rows, sql)
	case elapsed > slowQuery && g.Level >= gormlogger.Warn:
		sql, rows := fc()
		Warnf("slow SQL [%s, %d rows] %s", elapsed, rows, sql)
	case g.Level >= gormlogger.Info:
		sql, rows := fc()
		Infof("SQL [%s, %d rows] %s", elapsed, rows, sql)
	}
}
