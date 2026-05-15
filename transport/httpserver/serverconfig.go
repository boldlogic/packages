package httpserver

import (
	"fmt"
	"strings"
)

// ServerConfig описывает сетевые настройки HTTP-сервера.
type ServerConfig struct {
	ListenIp     string     `yaml:"listen_ip" json:"listen_ip"`
	ExternalHost string     `yaml:"external_host" json:"external_host"`
	ListenPort   int        `yaml:"listen_port" json:"listen_port"`
	Opts         ServerOpts `yaml:"server_opts" json:"server_opts"`
}

// ServerOpts описывает таймауты HTTP-сервера в секундах.
type ServerOpts struct {
	ReadTimeout       int `yaml:"read_timeout" json:"read_timeout"`
	ReadHeaderTimeout int `yaml:"read_header_timeout" json:"read_header_timeout"`
	WriteTimeout      int `yaml:"write_timeout" json:"write_timeout"`
	IdleTimeout       int `yaml:"idle_timeout" json:"idle_timeout"`
	ShutdownTimeout   int `yaml:"shutdown_timeout" json:"shutdown_timeout"`
}

const (
	defaultListenIp        string = "127.0.0.1"
	defaultExternalHost    string = "localhost"
	defaultListenPort      int    = 8080
	defaultShutdownTimeout int    = 20
)

// ApplyDefaults заполняет пустые значения конфигурации значениями по умолчанию.
func (conf *ServerConfig) ApplyDefaults() {
	if strings.TrimSpace(conf.ListenIp) == "" {
		conf.ListenIp = defaultListenIp
	}
	if strings.TrimSpace(conf.ExternalHost) == "" {
		conf.ExternalHost = defaultExternalHost
	}
	if conf.ListenPort == 0 {
		conf.ListenPort = defaultListenPort
	}
	if conf.Opts.ShutdownTimeout == 0 {
		conf.Opts.ShutdownTimeout = defaultShutdownTimeout
	}

}

// Validate проверяет значения конфигурации и возвращает список найденных ошибок.
func (conf *ServerConfig) Validate() []error {
	var errs []error
	if conf.ListenPort < 1 || conf.ListenPort > 65535 {
		errs = append(errs, fmt.Errorf("listen_port должен быть в диапазоне 1-65535"))
	}
	if conf.Opts.ReadTimeout < 0 {
		errs = append(errs, fmt.Errorf("read_timeout не может быть отрицательным"))
	}
	if conf.Opts.ReadHeaderTimeout < 0 {
		errs = append(errs, fmt.Errorf("read_header_timeout не может быть отрицательным"))
	}
	if conf.Opts.WriteTimeout < 0 {
		errs = append(errs, fmt.Errorf("write_timeout не может быть отрицательным"))
	}
	if conf.Opts.IdleTimeout < 0 {
		errs = append(errs, fmt.Errorf("idle_timeout не может быть отрицательным"))
	}
	if conf.Opts.ShutdownTimeout < 1 || conf.Opts.ShutdownTimeout > 60 {
		errs = append(errs, fmt.Errorf("shutdown_timeout должен быть в диапазоне от 1 до 60"))
	}

	return errs
}
