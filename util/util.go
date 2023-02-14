package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetPath(pathArg ...string) string {
	var newPath string
	newPath = filepath.Join(pathArg...)
	newPath = filepath.FromSlash(newPath)
	return newPath
}

func GenerateRandInt() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(1000000000000)
	_ = fmt.Sprintf("%d", n)
	return n
}

func ConnectTeleBotTest(BotToken string, ChatId int64) (result bool) {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Println(color.RedString("[!]BotToken authentication error!"))
		return false
	}
	msg := tgbotapi.NewMessage(ChatId, "[!]Connection Test!")
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(color.RedString("[!]ChatId is incorrect!"))
		return false
	}
	return true
}

func SendMessageBot(bot *tgbotapi.BotAPI, ChatId int64, message string) {

	msg := tgbotapi.NewMessage(ChatId, message)
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Println(color.RedString("[!]SendMessage To TG Failed!"))
	} else {
		fmt.Println(color.GreenString("[!]SendMessage To TG Success!"))
	}

}

func UploadFileToBot(bot *tgbotapi.BotAPI, ChatId int64, filePath string) {

	fileName := filepath.Base(filePath)
	data, _ := os.ReadFile(filePath)
	b := tgbotapi.FileBytes{Name: fileName, Bytes: data}
	msg := tgbotapi.NewDocumentUpload(ChatId, b)
	msg.Caption = "Vulnerability Report: " + fileName
	_, err := bot.Send(msg)
	if err != nil {
		fmt.Printf(color.RedString("[!]Upload %s Failed!\n"), fileName)
	} else {
		fmt.Printf(color.GreenString("[!]Upload %s Success!\n"), fileName)
	}

}

func GetScannerReport(toolPath, randString string) (string, error) {
	regexString := "*" + string(randString) + "*"
	scannerReport, err := filepath.Glob(filepath.Join(toolPath, regexString))
	if err != nil {
		return "", fmt.Errorf("get Scanner Report Err:%s", err)
	}
	//fmt.Println(scannerReport)
	if len(scannerReport) == 0 {
		return "", fmt.Errorf("no Report Generate")
	} else if len(scannerReport) == 1 {
		if IsEmptyFile(scannerReport[0]) {
			return scannerReport[0], nil
		} else {
			return "", fmt.Errorf("is Empty File")
		}

	} else {
		return "", fmt.Errorf("regular Match Error Multiple Files")
	}
}

func ReadLinesFromFile(reportFile string) ([]string, error) {
	var linesList []string
	f, err := os.Open(reportFile)
	if err != nil {
		return []string{}, fmt.Errorf("open File Err:%s", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		linesList = append(linesList, scanner.Text())
	}
	return linesList, nil
}

func SendScannerReportToTelegramBot(reportFlag bool, bot *tgbotapi.BotAPI, ChatId int64, scannerName, toolPath, randString string) {
	regexFileString, err := GetScannerReport(toolPath, randString)
	if err != nil {
		fmt.Printf(color.RedString("[!]Get Scanner Report [%s] With Err:%s\n"), scannerName, err)
	} else {
		fmt.Printf(color.GreenString("[!]The Results Of [%s] Was Saved At [%s]\n"), scannerName, regexFileString)
	}

	if reportFlag == true && err == nil {
		listLines, err := ReadLinesFromFile(regexFileString)
		if err != nil {
			fmt.Printf(color.RedString("[!]ReadLines From Reports Err:%s\n"), err)
		}

		for _, line := range listLines {
			//_ = line
			SendMessageBot(bot, ChatId, line)
		}
		//搁置
		//fmt.Println(regexFileString)
		UploadFileToBot(bot, ChatId, regexFileString)
	}
}

func SavePrintAsLog(result, fileDirectoryPath, reportName string) {

	filePath := filepath.Join(fileDirectoryPath, reportName)
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		fmt.Println(color.RedString("Mkdir ScannerReportLog Err:"), err)
		return
	}
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf(color.GreenString("[!]%s Is Logged!\n"), reportName)
	}
}

func CheckAndCreateFile(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		// 文件不存在，创建该文件
		fmt.Printf(color.RedString("%s Not Exist,Create!\n"), dirname)
		err = os.Mkdir(dirname, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsEmptyFile(filePath string) bool {
	// Open the file using the provided file path.
	file, err := os.Open(filePath)
	if err != nil {
		// Handle the error.
		fmt.Println(color.RedString("Open File Err:"), err)
		return false
	} else {
		fileInfo, err := file.Stat()
		file.Close()
		if err != nil {
			// Handle the error.
			fmt.Println(color.RedString("Handle The Report Err:"), err)
		}

		if fileInfo.Size() == 0 {
			// The file is empty.
			//os.ReadFile(filePath)
			return false
		}
		// The file is not empty.
		return true
	}

	// Check the file size. If it is zero, the file is empty.

}

func ParseMapFlag(flag string) []map[string]string {
	var m []map[string]string
	tmp := make(map[string]string)
	if flag == "" {
		return nil
	}

	// 将 flag 的值按照 "," 分割
	pairs := strings.Split(flag, ",")

	// 遍历每个键值对
	for _, pair := range pairs {
		// 将键值对按照 "=" 分割
		tmp = map[string]string{}
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			// 忽略无效的键值对
			continue
		} else {
			//fmt.Println(parts[0], parts[1])
			tmp[parts[0]] = parts[1]
		}
		m = append(m, tmp)
	}
	return m
}
