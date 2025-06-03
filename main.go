package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	ora "alrogozin/pkg/orcl"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	flags "github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	yaml "github.com/kylelemons/go-gypsy/yaml"

	// _ "github.com/mattn/go-oci8"
	"go.uber.org/zap"
	// "database/sql"
	// _ "github.com/mattn/go-oci8"
)

var opts struct {
	Cfg string `short:"c" long:"config" default:".\\config\\config.yaml" desciption:"A realative path to config file."`
}
var args = []string{}

var files = []os.FileInfo{}
var mToken string = ""
var mChannelChatID int64
var mChannelChatIDBil int64
var mSTimeout int64
var config *yaml.File
var bot *tgbotapi.BotAPI

var sqlQueryNewData, sqlInsert string
var orclCstring string

// var db *sql.DB
var orcl *sqlx.DB = nil

// MsgUnit is "id, num_pp, header, content, min_dtime, max_dtime, date_end, tss_name, prs_name, from_name, is_mvk, date_send string"
type MsgUnit struct {
	id        int
	numPP     string
	header    string
	content   string
	minDtime  string
	maxDtime  string
	dateEnd   string
	tssName   string
	prsName   string
	fromName  string
	isMvk     string
	dateSend  string
	contentId int
	hashVal   string
	SbsCode   string
}

func (mu MsgUnit) Display() string {
	var pos = strings.LastIndex(mu.numPP, ".")
	var numPPShort = mu.numPP
	if pos != -1 {
		numPPShort = mu.numPP[:pos]
	}

	res := "#" + numPPShort + " " + mu.prsName + " от " + mu.maxDtime + "\n"
	res += "Подсистема: " + mu.tssName + "\n"
	res += "Сообщение от " + mu.fromName + "\n\n"
	res += "\"" + mu.content + "\"\n"
	// res += "----------------------" + "\n"
	// println(res)
	return res
}

// func displayMsgUnit(r *)

// ------------------------------------------------------
func init() {
	// Инициализация - параметры, config
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	// common.InitConfig()

	// обработка параметров запуска
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
	}
	if opts.Cfg == "" {
		panic("Нет пути до конфиг.файла")
	}
	config, err = yaml.ReadFile(opts.Cfg)
	if err != nil {
		fmt.Println(err)
	}

	// Telegram bot
	mToken, err = config.Get("bot_token")
	if err != nil {
		fmt.Println(err)
	}

	// Chat_id of the private channel: @jopex_channel
	mChannelChatID, err = config.GetInt("ChatID")
	if err != nil {
		fmt.Println(err, "#2")
	}
	mChannelChatIDBil, err = config.GetInt("ChatID$Bil")
	if err != nil {
		fmt.Println(err, "#2$bil")
	}

	mSTimeout, err = config.GetInt("s_timeout")
	if err != nil {
		fmt.Println(err, "#2")
	}

	bot, err = tgbotapi.NewBotAPI(mToken)
	if err != nil {
		fmt.Println(err, "#3")
	}

	// Oracle
	sqlQueryNewData, err = config.Get("Sql_NewData")
	if err != nil {
		fmt.Println(err, "#4")
	}

	sqlInsert, err = config.Get("Sql_Insert")
	if err != nil {
		fmt.Println(err, "#5")
	}

	orclCstring, err = config.Get("orcl_cstring")
	if err != nil {
		fmt.Println(err, "#5")
	}

}

// ------------------------------------------------------
func main() {
	// Соединение с БД
	// mvk/mvkprod$@172.16.27.8:1521/billapp?PROTOCOL=TCP
	orcl = ora.GetConnection("BILLAPP", "172.16.27.8", 1521, "MVK", "MVKPROD$")
	defer orcl.Close()

	for {
		// Read again the config file
		config, err := yaml.ReadFile(opts.Cfg)
		if err != nil {
			fmt.Println(err)
		}

		mSTimeout, err = config.GetInt("s_timeout")
		if err != nil {
			mSTimeout = 120
		}

		zap.L().Info(time.Now().Format("02.01.2006 15:04:05") + " Starting...")

		// Get new jopex messages
		MsgUnits := orclGetNewData()
		zap.L().Sugar().Infof("%d", len(MsgUnits))

		// Form tlgrm meessage and send them
		for _, mu := range MsgUnits {
			// Отправить сообщение в telegram
			sendTlgrmMessage(mu.Display(), mu.SbsCode)
			saveJopexData(mu)
		}

		time.Sleep(time.Duration(mSTimeout) * time.Second)
	}
}

// ------------------------------------------------------
func sendTlgrmMessage(mMessage string, p_sbs_code string) {
	var msg tgbotapi.MessageConfig
	if p_sbs_code == "MVK" {
		msg = tgbotapi.NewMessage(-mChannelChatID, mMessage)
	} else {
		msg = tgbotapi.NewMessage(-mChannelChatIDBil, mMessage)
	}
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println(err, "#4")
	}
}

// ------------------------------------------------------
/*
func fillFileArray(pDirPath string) {
	tFiles, err := ioutil.ReadDir(pDirPath)
	if err != nil {
		panic(err)
	}

	for _, file := range tFiles {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".rtf" {
			files = append(files, file)
		}
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

}
*/
// ------------------------------------------------------
func saveJopexData(mu MsgUnit) {
	// var result sql.Result
	// insert into MAIL_TASK_SEND_TG(task_id, content_id, hash, date_mng) values (:1, :2, :3, :4)
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	_, err = orcl.ExecContext(ctx, sqlInsert, mu.id, mu.contentId, mu.hashVal, time.Now())
	cancel()
	if err != nil {
		fmt.Println("ExecContext error is not nil:", err)
		return
	}
}

// ------------------------------------------------------
func orclGetNewData() []MsgUnit {
	var retValue []MsgUnit
	rows, err := orcl.Query(sqlQueryNewData)
	if err != nil {
		fmt.Println(err.Error(), sqlQueryNewData)
		panic(err)
	}

	var (
		num_pp, header, content, min_dtime, max_dtime, date_end, tss_name, prs_name, from_name, is_mvk, date_send, hash_val, sbs_code string
		id, content_id                                                                                                                int
		// content string
	)

	for rows.Next() {
		if err = rows.Scan(&id, &num_pp, &header, &content, &min_dtime, &max_dtime, &date_end, &tss_name, &prs_name, &from_name, &is_mvk, &date_send, &content_id, &hash_val, &sbs_code); err != nil {
			panic(err)
		}
		mu := &MsgUnit{
			id:        id,
			numPP:     num_pp,
			header:    header,
			content:   content,
			minDtime:  min_dtime,
			maxDtime:  max_dtime,
			dateEnd:   date_end,
			tssName:   tss_name,
			dateSend:  date_send,
			prsName:   prs_name,
			fromName:  from_name,
			isMvk:     is_mvk,
			contentId: content_id,
			hashVal:   hash_val,
			SbsCode:   sbs_code,
		}
		retValue = append(retValue, *mu)
	}
	return retValue
}
