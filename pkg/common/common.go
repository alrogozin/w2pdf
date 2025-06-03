package common

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/saintfish/chardet"
	"go.uber.org/zap"
)

func DieOnError(msg string, err error) {
	if err != nil {
		zap.L().Error(msg + " " + err.Error())
		os.Exit(1)
	}
}

func GetFileExt(p_file_name string) string {
	mFileExt := path.Ext(p_file_name)
	return strings.Trim(strings.ToLower(mFileExt), " ")
}

func UTC_string_to_char(p_string string) string {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	layout := "2006-01-02 15:04:05 UTC"
	t, err := time.Parse(layout, p_string)
	return t.In(location).Format("02.01.2006 15:04:05")
}

func MakeBase64(p_data []byte, p_file_name string) string {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(p_data)))
	base64.StdEncoding.Encode(dst, p_data)
	// zap.L().Info("----------------------------")
	// zap.L().Info(string(dst))
	// zap.L().Info("----------------------------")

	// err := os.WriteFile(p_file_name, []byte(dst), 0644)
	// if err != nil {
	// 	zap.L().Panic(err.Error())
	// }

	return string(dst)
}

func Decode_From_Base64(p_base64_string string) ([]byte, error) {
	dst, err := base64.StdEncoding.DecodeString(p_base64_string)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// if Charset != "UTF-8" {
func DetectCodePage(str string) string {
	// ------------------------------------------
	// Перекодировка в UTF-8 (NB!!! препролагаю, что файл в Win-1251)
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest([]byte(str))
	if err == nil {
		fmt.Printf("Detected charset is %s\n", result.Charset)
	}
	return result.Charset
}
