package datasource

import (
	"errors"
	"fmt"
	configPkg "github.com/Adesubomi/magic-ayo-api/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	POSTGRES string = "postgres"
	MYSQL    string = "mysql"
)

func ConnectDatabase(conf *configPkg.Config) (*gorm.DB, error) {

	connection := conf.Database.Connection
	if connection == POSTGRES {
		postgresConn, err := postgresConnection(conf)
		if err != nil {
			return nil, err
		}
		return postgresConn, nil
	} else if connection == MYSQL {
		mysqlConn, err := mysqlConnection(conf)
		if err != nil {
			return nil, err
		}
		return mysqlConn, nil
	}

	return nil, errors.New("could not connect to any database")
}

func mysqlConnection(conf *configPkg.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.DbName,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf(" ?? Connection to database failed: %v\n", err)
		return database, err
	}

	return database, nil
}

func postgresConnection(conf *configPkg.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Africa/Lagos",
		conf.Database.Host,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.DbName,
		conf.Database.Port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return database, err
	}

	return database, nil
}
