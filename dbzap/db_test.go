package dbzap

import "testing"

func TestDBConfig_ApplyDefaults(t *testing.T) {
	t.Run("пустой_конфиг", func(t *testing.T) {
		db := DBConfig{}
		db.ApplyDefaults()

		if db.Host != "localhost" {
			t.Fatalf("Host=%q, ожидали localhost", db.Host)
		}
		if db.Driver != "sqlserver" {
			t.Fatalf("Driver=%q, ожидали sqlserver", db.Driver)
		}
	})

	t.Run("host_и_driver_уже_заданы", func(t *testing.T) {
		db := DBConfig{Host: "db.example", Driver: "sqlserver"}
		db.ApplyDefaults()

		if db.Host != "db.example" {
			t.Fatalf("Host=%q, ожидали db.example", db.Host)
		}
		if db.Driver != "sqlserver" {
			t.Fatalf("Driver=%q, ожидали sqlserver", db.Driver)
		}
	})
}

func TestDBConfig_ApplySecretsFromEnv(t *testing.T) {
	t.Run("задан_DB_PASSWORD", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "from-env")
		t.Setenv("MSSQL_SA_PASSWORD", "")
		t.Setenv("DB_USER", "")

		db := DBConfig{Password: "from-yaml"}
		db.ApplySecretsFromEnv()

		if db.Password != "from-env" {
			t.Fatalf("Password=%q, ожидали from-env", db.Password)
		}
	})

	t.Run("задан_MSSQL_SA_PASSWORD", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "")
		t.Setenv("MSSQL_SA_PASSWORD", "from-mssql-env")
		t.Setenv("DB_USER", "")

		db := DBConfig{Password: "from-yaml"}
		db.ApplySecretsFromEnv()

		if db.Password != "from-mssql-env" {
			t.Fatalf("Password=%q, ожидали from-mssql-env", db.Password)
		}
	})

	t.Run("задан_DB_USER", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "")
		t.Setenv("MSSQL_SA_PASSWORD", "")
		t.Setenv("DB_USER", "env-user")

		db := DBConfig{User: "yaml-user"}
		db.ApplySecretsFromEnv()

		if db.User != "env-user" {
			t.Fatalf("User=%q, ожидали env-user", db.User)
		}
	})

	t.Run("переменные_окружения_пустые", func(t *testing.T) {
		t.Setenv("DB_PASSWORD", "")
		t.Setenv("MSSQL_SA_PASSWORD", "")
		t.Setenv("DB_USER", "")

		db := DBConfig{User: "u", Password: "p"}
		db.ApplySecretsFromEnv()

		if db.User != "u" || db.Password != "p" {
			t.Fatalf("User=%q Password=%q, ожидали u и p", db.User, db.Password)
		}
	})
}

func TestDBConfig_Validate(t *testing.T) {
	t.Run("все_обязательные_поля_заданы", func(t *testing.T) {
		db := DBConfig{
			Server:   "host",
			Name:     "db",
			User:     "u",
			Password: "p",
			Driver:   "sqlserver",
		}

		errs := db.Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})

	t.Run("пустой_конфиг", func(t *testing.T) {
		var db DBConfig

		errs := db.Validate()
		if len(errs) != 5 {
			t.Fatalf("ожидали 5 ошибок, получили %d: %v", len(errs), errs)
		}
	})

	t.Run("driver_не_sqlserver", func(t *testing.T) {
		db := DBConfig{
			Server:   "host",
			Name:     "db",
			User:     "u",
			Password: "p",
			Driver:   "postgres",
		}

		errs := db.Validate()
		if len(errs) != 1 {
			t.Fatalf("ожидали 1 ошибку, получили %v", errs)
		}
	})
}

func TestDBConfig_ApplyDefaultsThenValidate(t *testing.T) {
	t.Run("driver_не_задан", func(t *testing.T) {
		db := DBConfig{
			Server:   "srv",
			Name:     "n",
			User:     "u",
			Password: "p",
		}

		db.ApplyDefaults()
		errs := db.Validate()
		if len(errs) != 0 {
			t.Fatalf("ожидали 0 ошибок, получили %v", errs)
		}
	})
}
