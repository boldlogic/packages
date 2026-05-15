package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {

	goodOpts := ServerOpts{
		ReadTimeout:       60,
		ReadHeaderTimeout: 30,
		WriteTimeout:      60,
		IdleTimeout:       60,
		ShutdownTimeout:   20,
	}

	tests := []struct {
		name            string
		conf            ServerConfig
		wantErrorsCount int
	}{
		{
			name: "корректный_конфиг",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   80,
				Opts:         goodOpts,
			},
			wantErrorsCount: 0,
		},
		{
			name: "минимальный_в_диапазоне_listen_port",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   1,
				Opts:         goodOpts,
			},
			wantErrorsCount: 0,
		},
		{
			name: "не_задан_shutdown_timeout",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   80,
				Opts:         ServerOpts{},
			},
			wantErrorsCount: 1,
		},
		{
			name: "shutdown_timeout_0",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   8080,
				Opts: ServerOpts{
					ShutdownTimeout: 0,
				},
			},
			wantErrorsCount: 1,
		},
		{
			name: "отрицательный_listen_port",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   -80,
				Opts:         goodOpts,
			},
			wantErrorsCount: 1,
		},
		{
			name: "порт_0",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   0,
				Opts:         goodOpts,
			},
			wantErrorsCount: 1,
		},
		{
			name:            "пустой",
			conf:            ServerConfig{},
			wantErrorsCount: 2,
		},
		{
			name: "максимальный_в_диапазоне_listen_port",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   65535,
				Opts:         goodOpts,
			},
			wantErrorsCount: 0,
		},
		{
			name: "порт_65536",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   65536,
				Opts:         goodOpts,
			},
			wantErrorsCount: 1,
		},
		{
			name: "пустой_listen_ip",
			conf: ServerConfig{
				ListenIp:     "",
				ExternalHost: "localhost",
				ListenPort:   8080,
				Opts:         goodOpts,
			},
			wantErrorsCount: 0,
		},
		{
			name: "пустой_external_host",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "",
				ListenPort:   8080,
				Opts:         goodOpts,
			},
			wantErrorsCount: 0,
		},
		{
			name: "отрицательный_read_timeout",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   8080,
				Opts: ServerOpts{
					ReadTimeout:     -1,
					ShutdownTimeout: 20,
				},
			},
			wantErrorsCount: 1,
		},
		{
			name: "большой_shutdown_timeout_61",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   8080,
				Opts: ServerOpts{
					ShutdownTimeout: 61,
				},
			},
			wantErrorsCount: 1,
		},
		{
			name: "корректный_shutdown_timeout_60",
			conf: ServerConfig{
				ListenIp:     "127.0.0.1",
				ExternalHost: "localhost",
				ListenPort:   8080,
				Opts: ServerOpts{
					ShutdownTimeout: 60,
				},
			},
			wantErrorsCount: 0,
		},
		{
			name: "абсолютно_некорректный",
			conf: ServerConfig{
				ListenPort: 0,
				Opts: ServerOpts{
					ReadTimeout:       -60,
					ReadHeaderTimeout: -30,
					WriteTimeout:      -60,
					IdleTimeout:       -60,
					ShutdownTimeout:   0,
				},
			},
			wantErrorsCount: 6,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got := tt.conf.Validate()
			assert.Len(t, got, tt.wantErrorsCount)
		})
	}
}

func TestApplyDefaults(t *testing.T) {

	tests := []struct {
		name             string
		conf             ServerConfig
		wantListenIp     string
		wantExternalHost string
		wantListenPort   int
	}{
		{
			name:             "корректный",
			conf:             ServerConfig{ListenIp: "0.0.0.0", ListenPort: 80, ExternalHost: "example.ru"},
			wantListenIp:     "0.0.0.0",
			wantExternalHost: "example.ru",
			wantListenPort:   80,
		},
		{
			name:             "дефолт",
			conf:             ServerConfig{},
			wantListenIp:     defaultListenIp,
			wantExternalHost: defaultExternalHost,
			wantListenPort:   defaultListenPort,
		},
		{
			name:             "дефолт_из-за_пробелов",
			conf:             ServerConfig{ListenIp: " ", ListenPort: 80, ExternalHost: " "},
			wantListenIp:     "127.0.0.1",
			wantExternalHost: "localhost",
			wantListenPort:   80,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.conf.ApplyDefaults()
			assert.Equal(t, tt.conf.ListenIp, tt.wantListenIp)
			assert.Equal(t, tt.conf.ExternalHost, tt.wantExternalHost)
			assert.Equal(t, tt.conf.ListenPort, tt.wantListenPort)
		})
	}
}

func TestApplyDefaultsAndValidate(t *testing.T) {

	tests := []struct {
		name            string
		inConf          ServerConfig
		wantConf        ServerConfig
		wantErrorsCount int
	}{
		{
			name:   "дефолт",
			inConf: ServerConfig{},
			wantConf: ServerConfig{
				ListenIp:     defaultListenIp,
				ExternalHost: defaultExternalHost,
				ListenPort:   defaultListenPort,
				Opts:         ServerOpts{ShutdownTimeout: defaultShutdownTimeout},
			},
			wantErrorsCount: 0,
		},
		{
			name: "отрицательный_shutdown_timeout",
			inConf: ServerConfig{
				Opts: ServerOpts{ShutdownTimeout: -1},
			},
			wantConf: ServerConfig{
				ListenIp:     defaultListenIp,
				ExternalHost: defaultExternalHost,
				ListenPort:   defaultListenPort,
				Opts:         ServerOpts{ShutdownTimeout: -1},
			},
			wantErrorsCount: 1,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			tt.inConf.ApplyDefaults()
			got := tt.inConf.Validate()
			assert.Equal(t, tt.inConf.ListenIp, tt.wantConf.ListenIp)
			assert.Equal(t, tt.inConf.ExternalHost, tt.wantConf.ExternalHost)
			assert.Equal(t, tt.inConf.ListenPort, tt.wantConf.ListenPort)
			assert.Equal(t, tt.inConf.Opts, tt.wantConf.Opts)

			assert.Len(t, got, tt.wantErrorsCount)
		})
	}
}
