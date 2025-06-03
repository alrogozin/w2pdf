package common

import (
	"os"
	"path"
	"runtime"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

// 	mPort, _ := strconv.Atoi(dbParams["PORT"])
//	dataSourceName := goora.BuildUrl(dbParams["IP_ADDRESS"], mPort, dbParams["DBNAME"], dbParams["USER"], dbParams["PASSWORD"], nil)

type Param_Oracle_Type struct {
	Port            string
	Ip_Address      string
	DBName          string
	User            string
	Password        string
	SQL_get_last_tl string
	SQL_set_appinfo string
	SQL_get_qas     string
}

type Param_Common_Type struct {
	TimeSleepSec int
}

type TgnMData_Type struct {
	Abbr   string
	ChatId int64
	Token  string
}

type HTTPServer_Type struct {
	Host string
	Port string
}

var Param_Oracle Param_Oracle_Type
var Param_Common Param_Common_Type
var Param_HTTPServer HTTPServer_Type
var TgnMDataArr []TgnMData_Type = make([]TgnMData_Type, 0)

func GetAppDir() (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func InitConfig() {
	var config_file, password string

	AppDir, err := GetAppDir()
	if err != nil {
		panic(err)
	}
	config_file = path.Join(path.Join(AppDir, ".."), "config/app.json")

	password, exists := os.LookupEnv("PASSWORD")
	if exists {
		zap.L().Info(password)
	} else {
		panic("OS param PASSWORD does not exist")
	}

	var js = koanf.New(".")
	if err = js.Load(file.Provider(config_file), json.Parser()); err != nil {
		zap.L().Info(err.Error())
		DieOnError("Reading config file ", err)
	}

	Param_Oracle.DBName = js.String("Oracle.dbname")
	Param_Oracle.Ip_Address = js.String("Oracle.ip_address")
	Param_Oracle.Port = js.String("Oracle.port")
	Param_Oracle.User = js.String("Oracle.user")
	Param_Oracle.SQL_get_last_tl = js.String("Oracle.SQL_get_last_tl")
	Param_Oracle.SQL_set_appinfo = js.String("Oracle.SQL_set_appinfo")
	Param_Oracle.SQL_get_qas = js.String("Oracle.SQL_get_qas")
	Param_Oracle.Password = password

	Param_HTTPServer.Host = js.String("HTTPServer.host")
	Param_HTTPServer.Port = js.String("HTTPServer.port")

	Param_Common.TimeSleepSec = int(js.Int64("Common.time_sleep_sec"))
}
