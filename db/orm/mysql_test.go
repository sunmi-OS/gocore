package orm

import (
	"gorm.io/gorm"
	"os"
	"strconv"
	"testing"
)

func TestSetDbConn(t *testing.T) {
	SetDbConn("name", &gorm.DB{})
	ormDb := GetORM("name")
	if ormDb == nil {
		t.Errorf("orm db is nil")
		return
	}
	t.Logf("set db connection success")

}

func TestInitDbConn(t *testing.T) {
	go func() {
		// the dsn invalid，will panic
		defer func() {
			if err := recover(); err != nil {
				t.Logf("panic: %v", err)
			}
		}()
		InitDbConn("test", WithDebug(true), WithDsn("test"))
	}()
	t.Logf("init db connection success")

	go func() {
		// the dsn cannot connect db, will panic
		defer func() {
			if err := recover(); err != nil {
				t.Logf("panic: %v", err)
			}
		}()
		InitDbConn("test", WithDebug(true), WithDsn("user:passwd@tcp(127.0.0.1:3306)/test"))
	}()
	t.Logf("init db connection success")

	t.Logf("check normal db connection with dsn begin")
	dsn := os.Getenv("DSN")
	InitDbConn("", WithDebug(true), WithDsn(dsn))
	conn := GetORM("iot_cloud_printer")
	if conn == nil {
		t.Errorf("conn is nil")
		return
	}

	t.Logf("check normal db connection. init and get db connection success")

	t.Logf("check normal db connection with host port begin")
	os.Setenv("DSN", "")

	host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")
	InitDbConn("", WithDebug(true),
		WithHost(host),
		WithPort(port),
		WithUsername(user),
		WithPassword(password),
		WithDatabase(dbName))

	t.Logf("check normal db connection with host port. init and get db connection success")

	partners := GetORM("partners")
	t.Logf("partners info: %+v", partners)

	printer := GetORM("iot_cloud_printer")
	t.Logf("printer info: %+v", printer)
	if partners == nil || printer == nil {
		t.Errorf("partners or printer is nil")
		return
	}

	Close()
}
