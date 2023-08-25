package main

// ISUCON的な参考: https://github.com/isucon/isucon12-qualify/blob/main/webapp/go/isuports.go#L336
// sqlx的な参考: https://jmoiron.github.io/sqlx/

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
)

const (
	listenPort = 12345
)

var (
	dbConn *sqlx.DB
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func loadDBDialConfigFromOSEnv() (*mysql.Config, error) {
	conf := mysql.NewConfig()
	return conf, nil
}

func connectDB() (*sqlx.DB, error) {
	const (
		networkTypeEnvKey = "ISUCON13_MYSQL_DIALCONFIG_NET"
		addrEnvKey        = "ISUCON13_MYSQL_DIALCONFIG_ADDRESS"
		userEnvKey        = "ISUCON13_MYSQL_DIALCONFIG_USER"
		passwordEnvKey    = "ISUCON13_MYSQL_DIALCONFIG_PASSWORD"
		dbNameEnvKey      = "ISUCON13_MYSQL_DIALCONFIG_DATABASE"
		parseTimeEnvKey   = "ISUCON13_MYSQL_DIALCONFIG_PARSETIME"
	)

	conf := mysql.NewConfig()

	// 環境変数がセットされていなかった場合でも一旦動かせるように、デフォルト値を入れておく
	// この挙動を変更して、エラーを出すようにしてもいいかもしれない
	conf.Net = "tcp"
	conf.Addr = "127.0.0.1:3306"
	conf.User = "isucon"
	conf.Passwd = "isucon"
	conf.DBName = "isupipe"
	conf.ParseTime = true

	if v, ok := os.LookupEnv(networkTypeEnvKey); ok {
		conf.Net = v
	}
	if v, ok := os.LookupEnv(addrEnvKey); ok {
		conf.Addr = v
	}
	if v, ok := os.LookupEnv(userEnvKey); ok {
		conf.User = v
	}
	if v, ok := os.LookupEnv(passwordEnvKey); ok {
		conf.Passwd = v
	}
	if v, ok := os.LookupEnv(dbNameEnvKey); ok {
		conf.DBName = v
	}
	if v, ok := os.LookupEnv(parseTimeEnvKey); ok {
		parseTime, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse environment variable '%s' as bool: %+v", parseTimeEnvKey, err)
		}
		conf.ParseTime = parseTime
	}

	return sqlx.Open("mysql", conf.FormatDSN())
}

func main() {
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(echolog.DEBUG)
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// top
	e.GET("/tag", nil)

	// livestream
	e.GET("/livestream", getLivestreamsHandler)
	e.GET("/livestream/:livestream_id", getLivestreamHandler)
	e.GET("/livestream/:livestream_id/comment", nil)
	e.GET("/livestream/user/:userid", nil)
	e.GET("/reaction", nil)
	// movie
	e.GET("/movie", nil)

	// user
	e.POST("/user", nil)
	e.POST("/login", nil)
	e.GET("/user/:user_id", nil)
	e.GET("/user/:user_id/channel", nil)
	e.POST("/user/:user_id/channel/:channel_id/subscribe", nil)
	e.POST("/user/:user_id/channel/:channel_id/unsubscribe", nil)
	e.GET("/channel", nil)
	e.GET("/channel/:channel_id/subscribers", nil)
	e.GET("/channel/:channel_id/movie", nil)
	e.POST("/channel", nil)
	e.PUT("/channel/:channel_id", nil)
	e.DELETE("/channel/:channel_id", nil)

	// DB接続
	conn, err := connectDB()
	if err != nil {
		e.Logger.Fatalf("failed to connect db: %v", err)
		return
	}
	conn.SetMaxOpenConns(10)
	defer conn.Close()
	dbConn = conn

	// HTTPサーバ起動
	listenAddr := net.JoinHostPort("", strconv.Itoa(listenPort))
	if err := e.Start(listenAddr); err != nil {
		e.Logger.Fatal(err)
	}
}
