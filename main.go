package main

import (
	"context"
	"flag"
	"goweb/common/webframe"
	"goweb/handler"
	"goweb/setup"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Streamlet/gohttp"
	"github.com/Streamlet/gosql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type config struct {
	Log   string      `toml:"log"`
	Mysql mysqlConfig `toml:"mysql"`
	Redis redisConfig `toml:"redis"`
	Http  httpConfig  `toml:"http"`
}

type httpConfig struct {
	Unix string `toml:"unix"`
	Tcp  string `toml:"tcp"`
}

type mysqlConfig struct {
	Address  string `toml:"address"`
	Db       string `toml:"db"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type redisConfig struct {
	Address  string `toml:"address"`
	Db       int    `toml:"db"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

func initLog(logFile string) {
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Println("failed to open log file", logFile)
		return
	}
	log.SetOutput(f)
}

func connectMysql(address, db, user, password string) *gosql.Connection {
	userPassword := user
	if password != "" {
		userPassword += ":" + password
	}
	c, err := gosql.Connect("mysql", userPassword+"@tcp("+address+")/"+db+"?charset=latin1&loc=Local&parseTime=True&clientFoundRows=true")
	if err != nil {
		log.Print("failed to connect to db: ", err.Error())
		return nil
	}
	return c
}

func connectRedis(address string, db int, user, password string) *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     address,
		DB:       db,
		Username: user,
		Password: password,
	})
	_, err := rc.Ping(context.Background()).Result()
	if err != nil {
		_ = rc.Close()
		log.Print("failed to connect to redis: ", err.Error())
		return nil
	}
	return rc
}

type commandLineArgs struct {
	setup   bool
	config  string
	webroot string
	debug   bool
}

func parseArgs() commandLineArgs {
	var args commandLineArgs
	flag.StringVar(&args.config, "config", "config.toml", "config file")
	flag.BoolVar(&args.setup, "setup", false, "setup the system")
	flag.StringVar(&args.webroot, "webroot", "", "web root for debug")
	flag.BoolVar(&args.debug, "debug", false, "debug mode, showing detail error message")
	flag.Parse()
	return args
}

func main() {
	args := parseArgs()
	if args.config == "" {
		log.Println("Usage: metaccount --config <config file> [--setup]")
		return
	}

	conf := config{
		Http: httpConfig{
			Tcp: ":8080",
		},
		Mysql: mysqlConfig{
			Address: "localhost:3306",
			Db:      "test",
			User:    "root",
		},
		Redis: redisConfig{
			Address: "localhost:6379",
			Db:      0,
		},
	}
	if args.config != "" {
		if _, err := toml.DecodeFile(args.config, &conf); err != nil {
			log.Println("Failed to parse config file:", err.Error())
			return
		}
	}

	if conf.Log != "" {
		initLog(conf.Log)
	}

	if args.setup {
		db := connectMysql(conf.Mysql.Address, "", conf.Mysql.User, conf.Mysql.Password)
		if db == nil {
			return
		}

		setup.InteractiveSetup(db, conf.Mysql.Db)
		return
	}

	db := connectMysql(conf.Mysql.Address, conf.Mysql.Db, conf.Mysql.User, conf.Mysql.Password)
	if db == nil {
		return
	}

	rc := connectRedis(conf.Redis.Address, conf.Redis.Db, conf.Redis.User, conf.Redis.Password)
	if rc == nil {
		return
	}

	application := gohttp.NewApplication[webframe.HttpContext](webframe.NewContextFactory(rc, db, args.debug))
	handler.Registers(application, args.webroot)

	if conf.Http.Unix != "" {
		application.ServeUnix(conf.Http.Unix)
	} else {
		application.ServeTcp(conf.Http.Tcp)
	}
}
