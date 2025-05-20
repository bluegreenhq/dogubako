package config

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type MySQLConfig struct {
	Host            string
	Port            int
	DBName          string
	User            string
	Password        string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}

func NewMySQLConfig(host string, port int, dbName, user, password string) *MySQLConfig {
	return &MySQLConfig{
		Host:            host,
		Port:            port,
		DBName:          dbName,
		User:            user,
		Password:        password,
		ConnMaxLifetime: 30 * time.Minute,
		MaxOpenConns:    35, // max_connections=45想定で、10余裕を持たせて設定
		MaxIdleConns:    35,
	}
}

var reDatabasePath = regexp.MustCompile(`mysql://(.+):(.+)@tcp\((.+):(.+)\)/(.+)`)

func NewMySQLConfigWithDatabasePath(dbPath string) (*MySQLConfig, error) {
	result := reDatabasePath.FindStringSubmatch(dbPath)
	user := result[1]
	password := result[2]
	host := result[3]

	port, err := strconv.Atoi(result[4])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	database := result[5]

	return NewMySQLConfig(host, port, database, user, password), nil
}

func (config MySQLConfig) DataSourceName() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", config.User, config.Password, config.Host, config.Port, config.DBName)
}
