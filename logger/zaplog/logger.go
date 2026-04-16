package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var fileWriteSyncers sync.Map

type cachedFileWriteSyncer struct {
	mu sync.Mutex
	f  *os.File
}

func (s *cachedFileWriteSyncer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.f.Write(p)
}

func (s *cachedFileWriteSyncer) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.f.Sync()
}

// writeSyncers возвращает список синхронизаторов записи для вывода логов.
// Если OutputFile пустой, возвращает синхронизатор для stdout.
// Если файл недоступен для открытия, также возвращается stdout.
func writeSyncers(cfg Config) []zapcore.WriteSyncer {
	if cfg.OutputFile == "" {
		return []zapcore.WriteSyncer{zapcore.Lock(os.Stdout)}
	}

	if syncer, ok := fileWriteSyncers.Load(cfg.OutputFile); ok {
		return []zapcore.WriteSyncer{syncer.(zapcore.WriteSyncer)}
	}

	f, err := os.OpenFile(cfg.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return []zapcore.WriteSyncer{zapcore.Lock(os.Stdout)}
	}

	syncer := &cachedFileWriteSyncer{f: f}
	actual, loaded := fileWriteSyncers.LoadOrStore(cfg.OutputFile, syncer)
	if loaded {
		_ = f.Close()
		return []zapcore.WriteSyncer{actual.(zapcore.WriteSyncer)}
	}

	return []zapcore.WriteSyncer{syncer}
}

// New создаёт новый экземпляр zap.Logger на основе переданной конфигурации.
// Логгер создаётся с включённым caller skip для автоматического добавления
// информации о месте вызова в каждое сообщение лога.
func New(cfg Config) *zap.Logger {
	lvl := zapcore.InfoLevel
	if parsedLevel, err := zapcore.ParseLevel(cfg.Level); err == nil {
		lvl = parsedLevel
	}

	atomicLevel := zap.NewAtomicLevelAt(lvl)
	writers := writeSyncers(cfg)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	var encoder zapcore.Encoder
	if strings.ToUpper(cfg.Format) == "JSON" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), atomicLevel)
	return zap.New(core, zap.AddCaller())
}
