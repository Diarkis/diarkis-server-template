// © 2019-2024 Diarkis Inc. All rights reserved.

package log

import (
	"fmt"
	"time"

	"github.com/Diarkis/diarkis/log"
	"github.com/Diarkis/diarkis/log/lib/colorize"
	"github.com/Diarkis/diarkis/log/lib/msg"
	"github.com/Diarkis/diarkis/util"
)

type ScenarioUID string

type ScenarioLogger struct {
	*log.Logger
}

func New(name string) *ScenarioLogger {
	logger := log.New(name)
	return &ScenarioLogger{Logger: logger}
}

func appendUIDAsPrefix(userID string, vals ...interface{}) []interface{} {
	vals[0] = fmt.Sprintf("[UID: %s] %v", userID, vals[0])
	return vals
}

// Verboseu outputs log with UID as prefix
func (sl ScenarioLogger) Verboseu(userID string, vals ...interface{}) {
	sl.Verbose(appendUIDAsPrefix(string(userID), vals...)...)
}

// Networkf outputs log with UID as prefix
func (sl ScenarioLogger) Networkf(userID string, vals ...interface{}) {
	sl.Network(appendUIDAsPrefix(string(userID), vals...)...)
}

// Sysu outputs log with UID as prefix
func (sl ScenarioLogger) Sysu(userID string, vals ...interface{}) {
	sl.Sys(appendUIDAsPrefix(string(userID), vals...)...)
}

// Traceu outputs log with UID as prefix
func (sl ScenarioLogger) Traceu(userID string, vals ...interface{}) {
	sl.Trace(appendUIDAsPrefix(string(userID), vals...)...)
}

// Debugu outputs log with UID as prefix
func (sl ScenarioLogger) Debugu(userID string, vals ...interface{}) {
	sl.Debug(appendUIDAsPrefix(string(userID), vals...)...)
}

// Infou outputs log with UID as prefix
func (sl ScenarioLogger) Infou(userID string, vals ...interface{}) {
	sl.Info(appendUIDAsPrefix(string(userID), vals...)...)
}

// Noticeu outputs log with UID as prefix
func (sl ScenarioLogger) Noticeu(userID string, vals ...interface{}) {
	sl.Notice(appendUIDAsPrefix(string(userID), vals...)...)
}

// Warnu outputs log with UID as prefix
func (sl ScenarioLogger) Warnu(userID string, vals ...interface{}) {
	sl.Warn(appendUIDAsPrefix(string(userID), vals...)...)
}

// Erroru outputs log with UID as prefix
func (sl ScenarioLogger) Erroru(userID string, vals ...interface{}) {
	sl.Error(appendUIDAsPrefix(string(userID), vals...)...)
}

// Fatalu outputs log with UID as prefix
func (sl ScenarioLogger) Fatalu(userID string, vals ...interface{}) {
	sl.Fatal(appendUIDAsPrefix(string(userID), vals...)...)
}

func EnableContextLogger() {
	log.UseCustomOutput()
	log.SetCustomOutput(func(formatted bool, name string, prefix string, level string, vals []interface{}) string {
		message := fmt.Sprintf("%v", vals[0])
		if len(vals) > 1 {
			now := time.Now()
			makePrefix := func() string {
				return fmt.Sprintf(
					msg.OutputPrefix,
					now.UTC().Format("2006/01/02 15:04:05.000"),
					colorize.Get(colorize.Blue, name),
					level,
					colorize.Get(colorize.Grey, prefix),
				)
			}
			if uid, ok := vals[0].(ScenarioUID); ok {
				prefix := util.StrConcat(makePrefix(), "[", string(uid), "] ")
				message = util.StrConcat(prefix, fmt.Sprintf(fmt.Sprintf("%v", vals[1]), vals[2:]...))
			} else {
				prefix := makePrefix()
				message = util.StrConcat(prefix, fmt.Sprintf(message, vals[1:]...))
			}
		}
		return message
	})
}
