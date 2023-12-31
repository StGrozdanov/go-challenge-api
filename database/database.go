package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"outdoorsy-api/utils"
	"sync"
	"time"
)

type db struct {
	host     string
	user     string
	password string
	port     string
	database string
	DB       *sqlx.DB
}

var instance *db

func Init(hosts string, user string, password string, port string, database string) {
	var syncOnce sync.Once
	if instance == nil {
		syncOnce.Do(
			func() {
				instance = &db{
					host:     utils.GetRandomHost(hosts),
					user:     user,
					password: password,
					port:     port,
					database: database,
				}
				connect()
				go connectionFailureListener()
			},
		)
	}
}

func connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", instance.host, instance.port, instance.user, instance.password, instance.database)
	sqlxConnection, err := sqlx.Open("postgres", psqlInfo)
	if err == nil {
		instance.DB = sqlxConnection
		instance.DB.SetMaxOpenConns(10)
		instance.DB.SetMaxIdleConns(10)
		instance.DB.SetConnMaxIdleTime(30 * time.Second)
		instance.DB.SetConnMaxLifetime(30 * time.Second)
	} else {
		utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Error on connection attempt to the Database")
	}
}

func connectionFailureListener() {
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		if err := Ping(); err != nil {
			utils.GetLogger().WithFields(log.Fields{"error": err.Error()}).Error("Database is not responding")
			CloseConnection()
			connect()
		}
	}
}
