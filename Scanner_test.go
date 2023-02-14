package main

import (
	"TemplateScanner/models"
	"TemplateScanner/util"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strings"
	"testing"
)

var BotToken = "1719205124:AAEVm_SvvixsCCnnpJpu6VPzvk_NDIa8n0U"
var ChatId = 598498092

var scannerYaml *models.ScannerInfo

func TestScannerStruct(t *testing.T) {
	data, err := models.LoadScannerYaml("C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\config\\xray.yaml")
	if err != nil {
		fmt.Println("ok")
	}
	fmt.Println(data.Description)
	fmt.Println(data.ScannerTemplateAuthor)
	fmt.Println(data.ScannerName)
	fmt.Println(data.SupportedPlatform["windows"].StartCmd)
	fmt.Println(data.TimeOut)
	fmt.Println("==========================")
	fmt.Println(data.Type)
}

func TestGenerateRandInt(t *testing.T) {
	fmt.Println(util.GenerateRandInt())
}

func TestExecuteScannerTask(t *testing.T) {

	BotToken := "1719205124:AAEVm_SvvixsCCnnpJpu6VPzvk_NDIa8n0U"
	ChatId := 598498092

	if err := util.ConnectTeleBotTest(BotToken, int64(ChatId)); err == false {
		log.Fatal("Connect to the bot failed!Check your BotToken or ChatId!")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	var args = make(map[string]string)

	currentPaht, err := os.Getwd()
	if err != nil {
		fmt.Errorf("getwd Err:%s", err)
	}

	toolPath := util.GetPath(currentPaht, "scanTool")
	args["toolPath"] = toolPath
	args["random"] = strings.TrimSpace(fmt.Sprintf("%d", util.GenerateRandInt()))
	args["target"] = "http://180.63.89.98/cgi-bin/login.php"

	defer func() {
		toolPath = util.GetPath(toolPath, "BBScan")
		fmt.Println(toolPath)
		regexFileString, err := util.GetScannerReport(toolPath, args["random"])
		if err != nil {
			fmt.Printf("Get Scanner Report Err:%s", err)
		}
		listLines, err := util.ReadLinesFromFile(regexFileString)
		if err != nil {
			fmt.Printf("ReadLines From Reports Err:%s", err)
		}

		for _, line := range listLines {
			util.SendMessageBot(bot, int64(ChatId), line)
		}

		util.UploadFileToBot(bot, int64(ChatId), regexFileString)
	}()

	data, err := models.LoadScannerYaml("C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\config\\bbscan.yaml")
	if err != nil {
		fmt.Println("ok")
	}

	args["ToolDirectoryName"] = data.ToolDirectoryName
	taskCmd, err := models.InterpolateWithArgs(data.SupportedPlatform["windows"].StartCmd, args)
	if err != nil {
		fmt.Printf("Err:%s", err)
	}

	fmt.Println(taskCmd)
	resutl, err := models.ExecuteCmd("windows", taskCmd, "bbscan")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resutl)
}

func TestSendMsgOrFileToTelegram(t *testing.T) {
	BotToken := "1719205124:AAEVm_SvvixsCCnnpJpu6VPzvk_NDIa8n0U"
	ChatId := 598498092

	if err := util.ConnectTeleBotTest(BotToken, int64(ChatId)); err == false {
		log.Fatal("Connect to the bot failed!Check your BotToken or ChatId!")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	util.SendMessageBot(bot, int64(ChatId), "test")
	util.UploadFileToBot(bot, int64(ChatId), "C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\scanTool\\nuclei_73673810134.json")
}

func TestGetScannerReport(t *testing.T) {
	result, err := util.GetScannerReport("C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\scanTool\\", "58543488311763")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestReadLinesFromReports(t *testing.T) {
	results, err := util.ReadLinesFromFile("C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\scanTool\\nuclei_29435209385.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(results)
}

func TestLoadAllScannerYaml(t *testing.T) {
	data := models.LoadAllScannerYaml()
	for _, v := range data {
		fmt.Println(v.ScannerName)
	}
}

func TestSendScannerReportToTelegramBot(t *testing.T) {
	reportFlag := true
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	fmt.Println("===============================================================")
	//fmt.Println(bot)
	//fmt.Println(int64(ChatId))
	//util.SendMessageBot(bot, int64(ChatId), "test")
	util.SendScannerReportToTelegramBot(reportFlag, bot, int64(ChatId), "xray", "C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\scanTool\\BBScan", "120656406525")
}

func TestStartScannerTask(t *testing.T) {
	BotToken := "1719205124:AAEVm_SvvixsCCnnpJpu6VPzvk_NDIa8n0U"
	chatId := 598498092
	models.StartScannerTask("http://180.235.39.192", BotToken, int64(chatId), false, []string{"xray"}, nil)
}

func TestStartScannerTaskFromFile(t *testing.T) {
	BotToken := "1719205124:AAEVm_SvvixsCCnnpJpu6VPzvk_NDIa8n0U"
	chatId := 598498092
	models.StartScannerTaskFromFile("urls.txt", BotToken, 3, int64(chatId), true, nil, nil)
}

func TestInitial(t *testing.T) {
	scanInfoList := models.LoadAllScannerYaml()
	for _, scanInfo := range scanInfoList {
		fmt.Println(scanInfo.ScannerName)
	}
}

func TestListTemplateScannerInfo(t *testing.T) {
	models.ListTemplateScannerInfo()
}

func TestSavePrintAsLog(t *testing.T) {
	util.SavePrintAsLog("test", "C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\reports", "bbscan_2342342.log")
}

func TestCheckAndCreateFile(t *testing.T) {
	util.CheckAndCreateFile("1")
}

func TestIsEmptyFile(t *testing.T) {
	fmt.Println(util.IsEmptyFile("C:\\Users\\zhtty\\GolandProjects\\BugBoutryScanner\\TemplateScanner\\scanTool\\BBScan\\bbscan_388211791783.txt"))
}

func TestParse(t *testing.T) {
	var parse models.Config

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
	}
	parse1, err := parse.Parse(data)
	fmt.Println(parse1)
	fmt.Println(parse1.InputArguments)
	for k, v := range parse1.InputArguments {
		fmt.Println(k, v)
		for k1, v1 := range v {
			fmt.Println(k1, v1)
		}
	}
	fmt.Println(parse1.LoadTemplate)

}

func TestLoadConfigYaml(t *testing.T) {
	data, err := models.LoadConfigYaml()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
	fmt.Println(data.LoadTemplate)
	fmt.Println(data.InputArguments)
	fmt.Println(data.Timeout)
	fmt.Println(data.Telegram.ChatID)
	fmt.Println(data.Telegram.Token)
	fmt.Println(data.Thread)

}

func TestParseMapFlag(t *testing.T) {
	fmt.Println(util.ParseMapFlag("A=‘dsd’,C=D,E=G"))
}
