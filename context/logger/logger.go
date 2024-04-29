package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// 以下のようにログを出力する
// 2024/05/01 02:50:00 ERROR error to SampleFunc "err message"=サンプルエラーメッセージ file=../../pkg/sample/sample_file.go line=50
func Error(funcName string, err error) {
	_, file, line, ok := runtime.Caller(1) // エラーが発生したファイル名と行数を取得
	if !ok {
		log.Println("failed to runtime.Caller")
		return
	}
	// カレントディレクトリを取得
	currentDir, errDir := os.Getwd()
	if errDir != nil {
		log.Println("failed to get current directory")
		return
	}

	// ファイルの相対パスを計算
	relativePath, errRel := filepath.Rel(currentDir, file)
	if errRel != nil {
		log.Println("failed to get relative path")
		return
	}

	slog.Error(fmt.Sprintf("error to %s", funcName),
		slog.String("err message", err.Error()),
		slog.String("file", relativePath),
		slog.Int("line", line))
}

//! 一旦は導入見送りのためコメントアウト

// import (
// 	"fmt"
// 	"log/slog"
// 	"os"
// )

// type Logger struct {
// 	slogLogger *slog.Logger
// }

// // NewLogger creates a new Logger
// func NewLogger(debug bool) *Logger {
// 	level := slog.LevelInfo
// 	if debug {
// 		level = slog.LevelDebug
// 	}

// 	var handler slog.Handler
// 	if debug {
// 		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 			AddSource: false,
// 			Level:     level,
// 		})

// 	} else {
// 		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
// 			AddSource: true,
// 			Level:     level,
// 			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
// 				switch {
// 				case a.Key == slog.MessageKey:
// 					return slog.String("message", a.Value.String())
// 				case a.Key == slog.LevelKey && a.Value.String() == slog.LevelWarn.String():
// 					return slog.String("severity", "WARNING")
// 				case a.Key == slog.LevelKey:
// 					return slog.String("severity", a.Value.String())
// 				}

// 				return a
// 			},
// 		})

// 	}

// 	slogLogger := slog.New(handler)
// 	return &Logger{slogLogger: slogLogger}
// }

// func (l *Logger) Slog() *slog.Logger {
// 	return l.slogLogger
// }

// func (l *Logger) With(args ...any) *Logger {
// 	return &Logger{
// 		slogLogger: l.slogLogger.With(args...),
// 	}
// }

// func (l *Logger) Debug(msg string, arg ...any) {
// 	l.slogLogger.Debug(msg, arg...)
// }

// func (l *Logger) Info(msg string, arg ...any) {
// 	l.slogLogger.Info(msg, arg...)
// }

// func (l *Logger) Warning(msg string, arg ...any) {
// 	l.slogLogger.Warn(msg, arg...)
// }

// func (l *Logger) Error(msg string, arg ...any) {
// 	l.slogLogger.Error(msg, arg...)
// }

// func (l *Logger) Debugf(msg string, arg ...any) {
// 	l.slogLogger.Debug(fmt.Sprintf(msg, arg...))
// }

// func (l *Logger) Infof(msg string, arg ...any) {
// 	l.slogLogger.Info(fmt.Sprintf(msg, arg...))
// }

// func (l *Logger) Warningf(msg string, arg ...any) {
// 	l.slogLogger.Warn(fmt.Sprintf(msg, arg...))
// }

// func (l *Logger) Errorf(msg string, arg ...any) {
// 	l.slogLogger.Error(fmt.Sprintf(msg, arg...))
// }
