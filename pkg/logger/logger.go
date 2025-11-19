package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type LogStyle int

const (
	JsonStyle LogStyle = iota + 1
	TextStyle
)

// Logger 定义一个Logger，方便扩展和自定义一些功能
type Logger struct {
	slog.Handler
}

// Handle 如果ctx有 request_id 附加到日志输出
func (log *Logger) Handle(ctx context.Context, r slog.Record) error {
	requestID, ok := ctx.Value("request_id").(string)

	// 附加到slog属性上
	if ok && requestID != "" {
		r.AddAttrs(slog.String("request_id", requestID))
	}

	// if r.Level == slog.LevelInfo {
	// // 做些什么
	// }

	return log.Handler.Handle(ctx, r)
}

// NewLogger 创建一个slog日志实例 level=(DEBUG,INFO,WARN,ERROR); output 日志输出位置
func NewLogger(output io.Writer, level string, logStyle LogStyle) *slog.Logger {
	logLevel := toLevel(level)

	var handler slog.Handler

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {

			// if a.Key == slog.LevelKey {
			// 	// 会被转义，无效果
			// 	a.Value = slog.StringValue(fmt.Sprintf("\033[31m%s\033[0m", a.Value))
			// }

			// 取相对路径，输出更简短的路径
			if a.Key == slog.SourceKey {
				// 此时的a.Value 是 slog.Source 指针
				ss, ok := a.Value.Any().(*slog.Source)
				if !ok || ss.File == "" {
					return a
				}
				dir, _ := os.Getwd()
				var sep = filepath.Base(dir) + "/"
				var parts = strings.Split(ss.File, sep)
				if len(parts) > 1 {
					dir, _ := os.Getwd()
					filepath.Base(dir)
					relativePath := sep + parts[1]
					a.Value = slog.StringValue(fmt.Sprintf("%s %d", relativePath, ss.Line))
				}
			}

			if a.Key == slog.TimeKey {
				datetime := a.Value.Time()
				a.Value = slog.StringValue(datetime.Format("2006-01-02 15:04:05.000"))
			}

			return a
		},
	}

	if logStyle == TextStyle {
		// TODO: 格式化时间，日志级别信息颜色输出
		handler = slog.NewTextHandler(output, opts)
	} else {
		handler = slog.NewJSONHandler(output, opts)
	}

	myHandler := Logger{Handler: handler}

	return slog.New(&myHandler)
}

// DefaultOutput 默认日志输出
func DefaultOutput(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    40, // megabytes
		MaxBackups: 30,
		MaxAge:     30,   //days
		Compress:   true, // disabled by default
	}
}

// toLevel 从字符串转换成slog.Leveler
func toLevel(level string) slog.Leveler {
	switch level {
	case "DEBUG":
		return slog.LevelDebug // -4
	case "INFO":
		return slog.LevelInfo // 0
	case "WARN":
		return slog.LevelWarn // 4
	case "ERROR":
		return slog.LevelError // 8
	default:
		return slog.LevelInfo // 0
	}
}
