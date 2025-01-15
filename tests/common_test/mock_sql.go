package testcommon

import (
	"database/sql"
	"fmt"
	"regexp"
	"sync"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MockDb struct {
	MockGorm *gorm.DB
	Db       *sql.DB
	Mock     sqlmock.Sqlmock
}

var once sync.Once
var mockDb *MockDb

// InitSqlMock
// 初始化 单元测试 中需要的 mock 对象
func InitSqlMock() (md *MockDb) {
	db, mock, _ := sqlmock.New()

	dialector := mysql.New(mysql.Config{
		DriverName:                "mysql",
		DSN:                       "sqlmock_db_0",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})

	gormDB, _ := gorm.Open(dialector, &gorm.Config{})
	mockDb = &MockDb{MockGorm: gormDB, Db: db, Mock: mock}

	return mockDb
}

func GetMysqlMock() *MockDb {
	if mockDb == nil {
		InitSqlMock()
	}

	return mockDb
}

func (md *MockDb) NewRows(rows []string) *sqlmock.Rows {
	return sqlmock.NewRows(rows)
}

func (md *MockDb) Close() {
	md.Mock.ExpectClose()
	err := md.Db.Close()
	mockDb = nil
	if err != nil {
		fmt.Println("关闭 mock 链接失败", err.Error())
	}
}

func (md *MockDb) ExpectationsWereMet() error {
	err := md.Mock.ExpectationsWereMet()
	return err
}

func (md *MockDb) ExpectQuery(expectedSQL string) *sqlmock.ExpectedQuery {
	return md.Mock.ExpectQuery(regexp.QuoteMeta(expectedSQL))
}

func (md *MockDb) ExpectExec(expectedSQL string) *sqlmock.ExpectedExec {
	return md.Mock.ExpectExec(regexp.QuoteMeta(expectedSQL))
}
func (md *MockDb) ExpectBegin() *sqlmock.ExpectedBegin {
	return md.Mock.ExpectBegin()
}

func (md *MockDb) ExpectCommit() *sqlmock.ExpectedCommit {
	return md.Mock.ExpectCommit()
}

func (md *MockDb) ExpectRollback() *sqlmock.ExpectedRollback {
	return md.Mock.ExpectRollback()
}
