package capsule

import (
	"sync"
	"time"

	"github.com/gordon-zhiyong/beehive-api/pkg/conf"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	defaultName        = "default"
	defaultMaxIdle     = 10
	defaultMaxConns    = 20
	defaultMaxLifetime = 10 * time.Minute
	maxIdle            int
	maxConns           int
	maxLifetime        time.Duration
	defaultConn        string
	dbMutex            sync.RWMutex
	dbConf             *viper.Viper
)

var dbConnections = map[string]*gorm.DB{}

// DB default database connection
var DB *gorm.DB
var dbOnce sync.Once

// initDB initialization default setting and default database connection
func initDB() {
	dbOnce.Do(func() {
		dbConf = conf.DB

		setDefaultSetting()
		DB = NewDBConnect(defaultConn)
	})
}

// NewDBConnect new database connection
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

func dbConnect(name string) (*gorm.DB, error) {
	if maxIdle == 0 {
		setDefaultSetting()
	}

	dsn := getDSN(name)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get sqlDB fail")
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxConns)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	return db, nil
}

func setDefaultSetting() {
	defaultConn = dbConf.GetString("default")
	maxIdle = dbConf.GetInt("maxIdle")
	maxConns = dbConf.GetInt("maxConns")
	maxLifetime = dbConf.GetDuration("maxLifetime")

	if maxIdle == 0 {
		maxIdle = defaultMaxIdle
	}

	if maxConns == 0 {
		maxConns = defaultMaxConns
	}

	if defaultConn == "" {
		defaultConn = defaultName
	}

	if maxLifetime == 0 {
		maxLifetime = defaultMaxLifetime
	}
}

func getDSN(conn string) string {
	username := dbConf.GetString(conn + ".user")
	password := dbConf.GetString(conn + ".password")
	host := dbConf.GetString(conn + ".host")
	port := dbConf.GetString(conn + ".port")
	charset := dbConf.GetString(conn + ".charset")
	dbName := dbConf.GetString(conn + ".name")

	return username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=" + charset + "&parseTime=True&loc=Local"
}
