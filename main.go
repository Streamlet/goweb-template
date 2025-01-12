package main

import (
	"context"
	"flag"
	"github.com/Streamlet/gohttp"
	"github.com/Streamlet/gosql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"goweb-template/common/webframe"
	"log"
	"os"
)

type commandLineArgs struct {
	sock      string
	port      uint
	log       string
	redisAddr string
	redisDb   int
	mysqlAddr string
	mysqlUser string
	mysqlDb   string
}

func parseArgs() commandLineArgs {
	var args commandLineArgs
	flag.StringVar(&args.sock, "sock", "", "unix sock file")
	flag.UintVar(&args.port, "port", 80, "listen port")
	flag.StringVar(&args.log, "log", "", "log file")
	flag.StringVar(&args.redisAddr, "redis-addr", "localhost:6379", "redis address")
	flag.IntVar(&args.redisDb, "redis-db", 0, "redis database")
	flag.StringVar(&args.mysqlAddr, "mysql-addr", "localhost:3306", "mysql address")
	flag.StringVar(&args.mysqlUser, "mysql-user", "root", "mysql user")
	flag.StringVar(&args.mysqlDb, "mysql-db", "", "mysql database")
	flag.Parse()
	return args
}

func checkArgs(args commandLineArgs) bool {
	if args.sock == "" && args.port == 0 {
		log.Println("Either sock or port must be specified.")
		return false
	}
	if args.redisAddr == "" {
		log.Println("redis address must be specified.")
		return false
	}
	if args.mysqlAddr == "" || args.mysqlUser == "" || args.mysqlDb == "" {
		log.Println("mysql address, user and database must be specified.")
		return false
	}
	return true
}

func initLog(logFile string) {
	if logFile == "" {
		return
	}
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println("failed to open log file", logFile)
		return
	}
	log.SetOutput(f)
}

func connectRedis(addr string, db int) *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	_, err := rc.Ping(context.Background()).Result()
	if err != nil {
		_ = rc.Close()
		log.Print("failed to connect to redis: ", err.Error())
		return nil
	}
	return rc
}

func connectMysql(mysqlAddr, mysqlUser, mysqlDb string) *gosql.Connection {
	db, err := gosql.Connect("mysql", mysqlUser+"@tcp("+mysqlAddr+")/"+mysqlDb+"?charset=latin1&loc=Local&parseTime=True")
	if err != nil {
		log.Print("failed to connect to db: ", err.Error())
		return nil
	}
	return db
}

func main() {
	args := parseArgs()
	if !checkArgs(args) {
		return
	}

	initLog(args.log)

	rc := connectRedis(args.redisAddr, args.redisDb)
	if rc == nil {
		return
	}

	db := connectMysql(args.mysqlAddr, args.mysqlUser, args.mysqlDb)
	if db == nil {
		return
	}

	application := gohttp.NewApplication[webframe.HttpContext](webframe.NewContextFactory(rc, db))
	registerHandlers(application)
	if args.sock != "" {
		application.ServeSock(args.sock)
	} else {
		application.ServePort(args.port)
	}
}
