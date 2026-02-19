package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"wild_bluebell/setting"
)

// 全局日志实例，供其他包使用
var lg *zap.Logger

// Init 初始化 zap 日志实例
// 参数:
//   - cfg: 日志配置结构体，包含日志级别、文件名、最大大小等
//   - mode: 运行时模式，"dev" 为开发模式，其他为生产模式
func Init(cfg *setting.LogConfig, mode string) (err error) {
	// 1. 创建日志写入器，支持日志文件轮转（按大小和时间）
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	// 2. 创建日志编码器，决定日志的输出格式
	encoder := getEncoder()

	// 3. 解析日志级别字符串（如 "debug", "info", "warn", "error"）为 zapcore.Level
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}

	// 4. 根据模式创建不同的 core（日志处理器链）
	var core zapcore.Core
	if mode == "dev" {
		// 开发模式：日志同时输出到文件和控制台
		// NewConsoleEncoder 使用更易读的格式输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			// 核心1：JSON格式写入文件
			zapcore.NewCore(encoder, writeSyncer, l),
			// 核心2：控制台输出（开发环境便于调试）
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// 生产模式：仅输出 JSON 格式到文件
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	// 5. 创建 logger 实例，AddCaller() 会记录调用者的文件名和行号
	lg = zap.New(core, zap.AddCaller())

	// 6. 将 logger 设置为全局默认日志实例
	// 这样可以直接调用 zap.L()、zap.S() 等全局方法记录日志
	zap.ReplaceGlobals(lg)
	zap.L().Info("init logger success")
	return
}

// getEncoder 创建日志编码器
// 返回一个 JSON 格式的编码器，配置了时间、级别、耗时等字段的格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder         // 时间格式：ISO8601
	encoderConfig.TimeKey = "time"                                // 时间字段名
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       // 级别用大写字母输出（如 INFO, DEBUG）
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder // 耗时单位：秒
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       // 调用位置：短格式（文件:行号）
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 创建日志写入器，支持日志文件轮转
// 使用 lumberjack 库实现日志切割和清理
// 参数:
//   - filename: 日志文件名
//   - maxSize: 单个日志文件最大大小（MB）
//   - maxBackup: 保留的日志备份数量
//   - maxAge: 日志文件保留天数
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件路径
		MaxSize:    maxSize,   // 达到此大小后自动切割（MB）
		MaxBackups: maxBackup, // 保留的历史文件数量
		MaxAge:     maxAge,    // 文件保留天数，超期的会被删除
	}
	// AddSync 将 io.Writer 包装为 WriteSyncer（支持并发安全）
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger Gin 框架的日志中间件
// 记录每个请求的详细信息：状态码、方法、路径、耗时、客户端IP等
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()             // 记录请求开始时间
		path := c.Request.URL.Path      // 请求路径
		query := c.Request.URL.RawQuery // 查询参数
		c.Next()                        // 执行后续处理链

		cost := time.Since(start) // 计算耗时
		lg.Info(path,             // 记录日志
			zap.Int("status", c.Writer.Status()),                                 // HTTP 状态码
			zap.String("method", c.Request.Method),                               // HTTP 方法
			zap.String("path", path),                                             // 请求路径
			zap.String("query", query),                                           // 查询参数
			zap.String("ip", c.ClientIP()),                                       // 客户端 IP
			zap.String("user-agent", c.Request.UserAgent()),                      // User-Agent
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()), // 错误信息
			zap.Duration("cost", cost),                                           // 请求耗时
		)
	}
}

// GinRecovery Gin 框架的恢复中间件
// 捕获处理链中可能出现的 panic，防止服务崩溃
// 参数:
//   - stack: 是否记录完整的堆栈信息
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer + recover 捕获 panic
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool

				// 检测是否是 "broken pipe" 错误（客户端断开连接）
				// 这种错误不需要记录堆栈，只需要记录错误信息即可
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "brokenpipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求的原始信息（用于日志记录）
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				if brokenPipe {
					// 客户端断开连接：只记录错误，不记录堆栈
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error)) // 让 Gin 处理错误
					c.Abort()            // 终止后续处理
					return
				}

				// 其他错误：根据 stack 参数决定是否记录堆栈
				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())), // 完整堆栈信息
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 返回 500 错误给客户端
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
