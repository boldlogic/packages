package dbzap

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
)

const testDriverName = "dbzap_test_driver"

var testDriverRegistered bool
var testDriverCloseCount int

type testDriver struct{}

type testConn struct{}

func (d testDriver) Open(_ string) (driver.Conn, error) {
	return &testConn{}, nil
}

func (c *testConn) Prepare(_ string) (driver.Stmt, error) {
	return nil, errors.New("не реализовано")
}

func (c *testConn) Close() error {
	testDriverCloseCount++
	return nil
}

func (c *testConn) Begin() (driver.Tx, error) {
	return nil, errors.New("не реализовано")
}

func (c *testConn) Ping(ctx context.Context) error {
	return errors.New("ошибка ping")
}

func registerTestDriver() {
	if !testDriverRegistered {
		sql.Register(testDriverName, testDriver{})
		testDriverRegistered = true
	}
}

func TestNew(t *testing.T) {
	t.Run("logger_nil_ErrDBPing", func(t *testing.T) {
		registerTestDriver()

		prevOpen := sqlOpen
		sqlOpen = func(_ string, dsn string) (*sql.DB, error) {
			return sql.Open(testDriverName, dsn)
		}
		defer func() { sqlOpen = prevOpen }()

		_, err := New(context.Background(), "ignored", nil)
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrDBPing) {
			t.Fatalf("ошибка %v, ожидали ErrDBPing", err)
		}
	})
}

func TestOpenDB(t *testing.T) {
	t.Run("ping_ErrDBPing", func(t *testing.T) {
		registerTestDriver()
		testDriverCloseCount = 0

		prevOpen := sqlOpen
		sqlOpen = func(_ string, dsn string) (*sql.DB, error) {
			return sql.Open(testDriverName, dsn)
		}
		defer func() { sqlOpen = prevOpen }()

		_, err := openDB(context.Background(), "ignored")
		if err == nil {
			t.Fatal("ожидали ошибку")
		}
		if !errors.Is(err, ErrDBPing) {
			t.Fatalf("ошибка %v, ожидали ErrDBPing", err)
		}
		if testDriverCloseCount != 1 {
			t.Fatalf("ожидали один вызов Close, получили %d", testDriverCloseCount)
		}
	})
}
