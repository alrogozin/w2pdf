package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	flags "github.com/jessevdk/go-flags"
	yaml "github.com/kylelemons/go-gypsy/yaml"
	_ "github.com/mattn/go-oci8"
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
var mSTimeout int64
var config *yaml.File
var bot *tgbotapi.BotAPI

var sqlQueryNewData, sqlInsert string
var orclCstring string
var db *sql.DB

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
	for {
		// Open DB connection
		var err error
		db, err = sql.Open("oci8", orclCstring)
		if err != nil {
			log.Fatal(err)
		}

		// Read again the config file
		config, err = yaml.ReadFile(opts.Cfg)
		if err != nil {
			fmt.Println(err)
		}
		mSTimeout, err = config.GetInt("s_timeout")
		if err != nil {
			mSTimeout = 120
		}

		print(time.Now().Format("02.01.2006 15:04:05") + " Starting...")

		// Get new jopex messages
		MsgUnits := orclGetNewData()

		// Form tlgrm meessage and send them
		for _, mu := range MsgUnits {
			// Отправить сообщение в telegram
			sendTlgrmMessage(mu.Display())
			saveJopexData(mu)
		}

		// defer close database
		/*
			defer func() {
				err = db.Close()
				if err != nil {
					fmt.Println("Close error is not nil:", err)
				} else {
					println("DB connection closed")
				}
			}()
		*/
		err = db.Close()
		if err != nil {
			fmt.Println("Close error is not nil:", err)
		} else {
			println(strconv.Itoa(len(MsgUnits)) + " messages sent. Next in " + strconv.Itoa((int(mSTimeout))) + " sec")
		}

		time.Sleep(time.Duration(mSTimeout) * time.Second)
	}
}

// ------------------------------------------------------
func sendTlgrmMessage(mMessage string) {
	msg := tgbotapi.NewMessage(-mChannelChatID, mMessage)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println(err, "#4")
	}
}

// ------------------------------------------------------
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

// ------------------------------------------------------
func saveJopexData(mu MsgUnit) {
	// var result sql.Result
	// insert into MAIL_TASK_SEND_TG(task_id, content_id, hash, date_mng) values (:1, :2, :3, :4)
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
	_, err = db.ExecContext(ctx, sqlInsert, mu.id, mu.contentId, mu.hashVal, time.Now())
	cancel()
	if err != nil {
		fmt.Println("ExecContext error is not nil:", err)
		return
	}
}

// ------------------------------------------------------
func orclGetNewData() []MsgUnit {
	var retValue []MsgUnit
	rows, err := db.Query(sqlQueryNewData)
	if err != nil {
		fmt.Println(err.Error(), sqlQueryNewData)
		log.Fatalln("err:", err)
	}

	var (
		num_pp, header, content, min_dtime, max_dtime, date_end, tss_name, prs_name, from_name, is_mvk, date_send, hash_val string
		id, content_id                                                                                                      int
		// content string
	)

	for rows.Next() {
		if err = rows.Scan(&id, &num_pp, &header, &content, &min_dtime, &max_dtime, &date_end, &tss_name, &prs_name, &from_name, &is_mvk, &date_send, &content_id, &hash_val); err != nil {
			log.Fatalln("error fetching", err)
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
		}
		retValue = append(retValue, *mu)
	}
	return retValue
}
