package dbzap

import (
	"fmt"
	"os"
)

// DBConfig описывает параметры подключения к базе данных и настройки,
// которые могут быть загружены из конфигурационного файла или окружения.
type DBConfig struct {
	Driver   string `yaml:"driver,omitempty" json:"driver,omitempty"`
	Server   string `yaml:"server,omitempty" json:"server,omitempty"`
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Name     string `yaml:"db_name" json:"db_name"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
	SSLMode  string `yaml:"ssl_mode,omitempty" json:"ssl_mode,omitempty"`
}

// ApplyDefaults заполняет значения по умолчанию для полей конфигурации,
// если они не были заданы явно.
func (db *DBConfig) ApplyDefaults() {
	if db.Host == "" {
		db.Host = "localhost"
	}
	if db.Driver == "" {
		db.Driver = "sqlserver"
	}
}

// ApplySecretsFromEnv подставляет секреты и чувствительные параметры из
// переменных окружения, если они заданы.
//
// Приоритет для пароля: `DB_PASSWORD`, затем `MSSQL_SA_PASSWORD`.
// Пользователь может быть переопределён через `DB_USER`.
func (db *DBConfig) ApplySecretsFromEnv() {
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		db.Password = v
	} else if v := os.Getenv("MSSQL_SA_PASSWORD"); v != "" {
		db.Password = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		db.User = v
	}
}

// Validate проверяет обязательные поля конфигурации и возвращает список
// ошибок валидации. Пустой результат означает, что конфигурация корректна.
func (db *DBConfig) Validate() []error {
	var errs []error
	if db.Name == "" {
		errs = append(errs, fmt.Errorf("не заполнено 'db_name'"))
	}
	if db.User == "" {
		errs = append(errs, fmt.Errorf("не заполнено 'user'"))
	}
	if db.Password == "" {
		errs = append(errs, fmt.Errorf("не заполнено 'password'"))
	}
	if db.Server == "" {
		errs = append(errs, fmt.Errorf("не заполнено 'server'"))
	}
	if db.Driver != "sqlserver" {
		errs = append(errs, fmt.Errorf("неподдерживаемый 'driver'"))
	}

	return errs
}

// GetDSN формирует строку подключения к SQL Server на основе текущей
// конфигурации.
func (db *DBConfig) GetDSN() (dsn string) {
	return fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&trustServerCertificate=true",
		db.User, db.Password, db.Server, db.Name)
}
