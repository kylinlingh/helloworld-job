package logger

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

/*
使用方法：

普通用法：
- log.Ctx(ctx).Trace().Msg("using preconfigured auth function")
- log.Ctx(ctx).Info().Msgf("using %s datastore engine", opts.Engine)

额外携带参数：使用 Int()，Strs()等函数
- log.Ctx(ctx).Trace().Int("preshared-keys-count", len(c.PresharedSecureKey)).Msg("using gRPC auth with preshared key(s)")
- log.Ctx(ctx).Info().Strs("files", opts.BootstrapFiles).Msg("initializing datastore from bootstrap files")
- log.Ctx(ctx).Info().Stringer("timeout", gracePeriod).Msg("starting shutdown grace period")

携带自定义的对象：
- log.Ctx(ctx).Info().EmbedObject(nscc).Msg("configured namespace cache")

将 map 展开：
- log.Ctx(ctx).Info().Fields(helpers.Flatten(c.DebugMap())).Msg("configuration") -> func Flatten(debugMap map[string]any) map[string]any

携带错误：
- log.Ctx(cmd.Context()).Fatal().Err(err).Msg("failed to create gRPC job")
- log.Ctx(ctx).Error().Err(fmt.Errorf("unable to marshall termination log: %w", err)).Msg("failed to downloadfrom termination log")
*/

var Logger zerolog.Logger

// 必须要提供 init 函数，以便在 main 函数执行前先初始化一个 log 对象，才能在解析出配置文件前打印信息
func init() {
	//SetGlobalLogger(zerolog.Nop())
	SetGlobalLogger(zerolog.New(os.Stdout))
}

func New(level string, runMode string, dir string, id string) {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	case "trace":
		l = zerolog.TraceLevel
	default:
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)

	_, err := os.Stat(dir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		os.Mkdir(dir, os.ModePerm)
	}

	logFileName := fmt.Sprintf("%s.log", id)
	logFile, _ := os.OpenFile(path.Join(dir, logFileName), os.O_CREATE|os.O_WRONLY, 0644)

	if runMode == "dev" {
		// 格式化日志，更容易阅读
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano}
		consoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		consoleWriter.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		consoleWriter.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}

		multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile)

		skipFrameCount := 0
		SetGlobalLogger(
			zerolog.New(multiWriter).
				Level(l).
				With().
				Timestamp().
				CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
				Logger())
	} else {
		skipFrameCount := 1
		consoleWriter := os.Stdout
		multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile)
		SetGlobalLogger(
			zerolog.New(multiWriter).
				Level(l).
				With().
				Timestamp().
				CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).
				Logger())
	}
}

func SetGlobalLogger(logger zerolog.Logger) {
	Logger = logger
	zerolog.DefaultContextLogger = &Logger
}

func GetLogger() *zerolog.Logger {
	return &Logger
}

func With() zerolog.Context { return Logger.With() }

func Err(err error) *zerolog.Event { return Logger.Err(err) }

func Trace() *zerolog.Event { return Logger.Trace() }

func Debug() *zerolog.Event { return Logger.Debug() }

func Info() *zerolog.Event { return Logger.Info() }

func Warn() *zerolog.Event { return Logger.Warn() }

func Error() *zerolog.Event { return Logger.Error() }

func Fatal() *zerolog.Event { return Logger.Fatal() }

func WithLevel(level zerolog.Level) *zerolog.Event { return Logger.WithLevel(level) }

func Log() *zerolog.Event { return Logger.Log() }

func Ctx(ctx context.Context) *zerolog.Logger { return zerolog.Ctx(ctx) }

func Print(v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}
