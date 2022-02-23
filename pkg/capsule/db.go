package capsule

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var prefix = "database"
var defaultName = "default"
var dbConnections = map[string]*gorm.DB{}
var DB *gorm.DB
var dbMutex sync.RWMutex

func init() {
	var err error
	DB, err = dbConnect(defaultName)
	if err != nil {
		panic(err)
	}
}

func dbConnect(name string) (*gorm.DB, error) {
	dsn := getDSN(name)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get sqlDB fail.")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	return db, nil
}

func NewDBConnect(name string) *gorm.DB {
	dbMutex.RLock()
	connection, exists := dbConnections[name]
	dbMutex.RUnlock()

	if exists {
		return connection
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()
	conn, err := dbConnect(name)
	if err != nil {
		panic(err)
	}

	dbConnections[name] = conn
	return conn
}

func getDSN(name string) string {
	return "b2c_oversea_test:BPSBgEo9ZwKn!cj@tcp(172.16.10.25:3306)/beehive?charset=utf8mb4&parseTime=True&loc=Local"
}
