package models

import (
	"TemplateScanner/util"
	"fmt"
	"github.com/fatih/color"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

//var BotToken = "1719475124:AAEVm_231231231nnpJpu6VPzvk_NDIa8n0U"
//var ChatId = 5111111192

func LoadScannerYaml(yamlPath string) (*ScannerInfo, error) {
	var test ScannerInfo
	var body []byte
	body, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("open Yaml File Err:%s", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("read Scanner Yaml File Err:%s", err)
	}
	err = yaml.Unmarshal(body, &test)

	if err != nil {
		return nil, fmt.Errorf("unmarshal Data Error:%s", err)
	}
	return &test, nil
}

func InterpolateWithArgs(CommandString string, args map[string]string) (string, error) {

	re := regexp.MustCompile(`#{.*}`)
	CommandString = strings.TrimSpace(CommandString)
	for key, value := range args {
		CommandString = strings.ReplaceAll(CommandString, "#{"+key+"}", fmt.Sprintf(value))
	}
	matches := re.FindAllString(CommandString, -1)
	if len(matches) != 0 {
		fmt.Printf(color.RedString("[!]Matching %s is unsuccessful and will be deleted\n"), matches)
		for _, match := range matches {
			//fmt.Println(match)
			CommandString = strings.ReplaceAll(CommandString, match, "")
		}
	}
	return CommandString, nil
}

func ExecuteCmd(osInfo, cmdString, toolDirectoryName string) (string, error) {
	//fmt.Println(toolDirectoryName)
	var cmd *exec.Cmd
	if runtime.GOOS != osInfo {
		return "", fmt.Errorf(color.RedString("need Platfrom:%s,But Current Platform:%s"), osInfo, runtime.GOOS)
	}
	cdToolDirectory := "cd " + toolDirectoryName
	switch runtime.GOOS {
	case "linux":
		file, err := os.CreateTemp("", "*.sh")
		if err != nil {
			return "", fmt.Errorf(color.RedString("creating temporary file: %w"), err)
		}
		//fmt.Printf("The script of path is at [%s]\n", file.Name())
		defer os.Remove(file.Name())
		if _, err := file.Write([]byte(cdToolDirectory + "\n" + cmdString)); err != nil {
			file.Close()

			return "", fmt.Errorf("writing command to file: %w", err)
		}
		if err := file.Close(); err != nil {
			return "", fmt.Errorf("closing bash script: %w", err)
		}
		cmd = exec.Command("sh", file.Name())
	case "windows":
		file, err := os.CreateTemp("", "*.ps1")
		if err != nil {
			return "", fmt.Errorf("creating temporary file: %w", err)
		}
		//fmt.Printf("The script of path is at [%s]\n", file.Name())
		defer os.Remove(file.Name())
		if _, err := file.Write([]byte(cdToolDirectory + "\r\n" + cmdString)); err != nil {
			file.Close()
			return "", err
		}
		if err := file.Close(); err != nil {
			return "", err
		}
		cmd = exec.Command("powershell", "-noprofile", file.Name())
	default:
		fmt.Errorf("current Platform:%s Is Not Suppor", runtime.GOOS)
	}

	cmd.Env = os.Environ()
	execResult, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execute Err:%s", err)
	}
	//fmt.Errorf(string(execResult))
	//fmt.Println("==========================执行的命令===============================")
	//fmt.Println(cmdString)
	//fmt.Println("==========================执行的结果===============================")
	//fmt.Println(string(execResult))
	//fmt.Println("==================================================================")
	return string(execResult), nil
}

func LoadAllScannerYaml() []*ScannerInfo {

	var ScannerInfoList []*ScannerInfo

	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalln(color.RedString("Can't Get Current Path!"))
	}
	//mark
	configPath := util.GetPath(currentPath, "config")
	regexString := "*.yaml"
	scannerYamlPaths, err := filepath.Glob(filepath.Join(configPath, regexString))
	if err != nil {
		log.Fatalln(color.RedString("[!]Can't Get ScannerYaml List!"))
	}
	//fmt.Println(configPath)

	if len(scannerYamlPaths) >= 0 {
		for _, scannerYamlPath := range scannerYamlPaths {
			temp, err := LoadScannerYaml(scannerYamlPath)
			if err != nil {
				fmt.Printf(color.RedString("[!]Err:%s When Load %s"), err, filepath.Base(scannerYamlPath))
			}
			ScannerInfoList = append(ScannerInfoList, temp)
		}
	} else {
		log.Fatalln(color.RedString("[!]ScannerYamlList Is Empty!"))
	}
	//fmt.Println(ScannerInfoList)
	return ScannerInfoList
}

func InitTelegramBot(BotToken string) (*tgbotapi.BotAPI, error) {
	var bot *tgbotapi.BotAPI
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		//log.Panic(err)
		return nil, err
	}
	bot.Debug = false
	fmt.Printf(color.GreenString("[!]Authorized on account %s\n"), bot.Self.UserName)
	return bot, nil
}

func ListTemplateScannerInfo() {
	//yellow := color.New(color.FgHiWhite, color.BgYellow, color.Bold).SprintFunc()
	scanInfoList := LoadAllScannerYaml()
	currentPath, _ := os.Getwd()
	for key, scanInfo := range scanInfoList {
		fmt.Fprintf(color.Output, color.YellowString("=====================================Load %dth Scanner Yaml=====================================\n"), key+1)
		fmt.Println(color.GreenString("[!]ScannerName:"), scanInfo.ScannerName)
		fmt.Println(color.GreenString("[!]ScannerDescription:"), strings.TrimSpace(scanInfo.Description))
		fmt.Println(color.GreenString("[!]ScannerToolDirectoryName:"), util.GetPath(currentPath, "scanTool", scanInfo.ToolDirectoryName))
		fmt.Println(color.GreenString("[!]Type:"), scanInfo.Type)
		fmt.Println(color.GreenString("[!]ScannerTemplateAuthor:"), scanInfo.ScannerTemplateAuthor)
		fmt.Println(color.GreenString("[!]TimeOut:"), scanInfo.TimeOut, "s")
		var supportList []string
		if _, ok := scanInfo.SupportedPlatform["windows"]; ok {
			supportList = append(supportList, "windows")
		}
		if _, ok := scanInfo.SupportedPlatform["linux"]; ok {
			supportList = append(supportList, "linux")
		}
		for _, supportPlatform := range supportList {
			fmt.Println(color.GreenString("[!]SupportedPlatform:"), supportPlatform)
			fmt.Println(color.GreenString("[!]StartCmd:\n"), strings.TrimSpace(scanInfo.SupportedPlatform[supportPlatform].StartCmd))
		}
	}
}

func StartScannerTask(target, BotToken string, chatId int64, reportFlag bool, scannerList []string, inputargs []map[string]string) {

	var args = make(map[string]string)
	var bot *tgbotapi.BotAPI
	var err error

	if inputargs != nil {
		for _, inputarg := range inputargs {
			for k, v := range inputarg {
				args[k] = v
			}
		}
	}

	//yellow := color.New(color.FgHiWhite, color.BgYellow, color.Bold).SprintFunc()
	//red := color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
	//green := color.New(color.FgWhite, color.BgGreen, color.Bold).SprintFunc()

	if reportFlag == true {
		bot, err = InitTelegramBot(BotToken)
		if err != nil {
			fmt.Printf(color.YellowString("[!]Connect To TelegramBot Failed,Err:%s\n"), err)
			reportFlag = false
		}
	}

	scannerInfoList := LoadAllScannerYaml()

	if scannerList != nil {
		scannerInfoList = SelectTemplateScanner(scannerInfoList, scannerList)
	}

	//fmt.Println("========================使用的模板=======================")
	//for _, v := range scannerInfoList {
	//	fmt.Println(v.ScannerName)
	//}
	//fmt.Println("========================================================")

	osInfo := runtime.GOOS
	currentPath, err := os.Getwd()

	if err != nil {
		fmt.Errorf(color.RedString("[!]getwd Err:%s"), err)
	}
	toolPath := util.GetPath(currentPath, "scanTool")
	reportPath := util.GetPath(currentPath, "logs")
	args["target"] = target
	args["toolPath"] = toolPath

	var wg sync.WaitGroup

	for key, scannerInfo := range scannerInfoList {
		args["random"] = strings.TrimSpace(fmt.Sprintf("%d", util.GenerateRandInt()+key))
		args["ToolDirectoryName"] = scannerInfo.ToolDirectoryName
		args["ScannerName"] = scannerInfo.ScannerName
		//fmt.Println(args["random"])
		//fmt.Println("====================输入的参数============================")
		//for k, v := range args {
		//	fmt.Println(k, ":", v)
		//}
		//fmt.Println("========================================================")
		taskCmd, err := InterpolateWithArgs(scannerInfo.SupportedPlatform[osInfo].StartCmd, args)
		if err != nil {
			fmt.Printf(color.RedString("[!]Interpolate With Args Err:%s\n"), err)
		}
		fmt.Printf(color.GreenString("[!]The Task Of [%s] Scanning [%s] Is Staring!\n"), args["ScannerName"], target)
		//fmt.Println(taskCmd)
		wg.Add(1)
		go func(reportFlag bool, regexString, scannerName, toolPath, toolDirectoryName string) {
			defer wg.Done()
			defer fmt.Printf(color.GreenString("[!]The Task Of [%s] Scanning [%s] Is Completed!\n"), scannerName, target)
			if len(scannerInfo.SupportedPlatform[osInfo].StartCmd) != 0 {
				toolDirectoryName := util.GetPath(toolPath, toolDirectoryName)
				//fmt.Println(toolPath, toolDirectoryName)
				result, err := ExecuteCmd(osInfo, taskCmd, toolDirectoryName)
				if err != nil {
					fmt.Println(err)
				}
				//fmt.Println(color.GreenString(result))
				timeString := time.Now().Format("20060102150405")
				logName := scannerName + "_" + timeString + "_" + regexString + ".log"
				util.SavePrintAsLog(result, reportPath, logName)
			}
			//报送给tg
			toolPath = util.GetPath(toolPath, toolDirectoryName)
			//fmt.Println(bot, int64(ChatId), toolPath, args["random"])
			util.SendScannerReportToTelegramBot(reportFlag, bot, chatId, scannerName, toolPath, regexString)
		}(reportFlag, args["random"], args["ScannerName"], args["toolPath"], args["ToolDirectoryName"])

	}
	wg.Wait()
}

func StartScannerSingleTask(wg2 *util.WaitGroup, bot *tgbotapi.BotAPI, target string, taskNumber, chatId int64, reportFlag bool, scannerList []string, inputargs []map[string]string) {

	defer wg2.Done()
	var args = make(map[string]string)
	var err error

	if inputargs != nil {
		for _, inputarg := range inputargs {
			for k, v := range inputarg {
				args[k] = v
			}
		}
	}

	//yellow := color.New(color.FgHiWhite, color.BgYellow, color.Bold).SprintFunc()
	//red := color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
	green := color.New(color.FgWhite, color.BgGreen, color.Bold).SprintFunc()

	scannerInfoList := LoadAllScannerYaml()
	if scannerList != nil {
		scannerInfoList = SelectTemplateScanner(scannerInfoList, scannerList)
	}
	osInfo := runtime.GOOS
	currentPaht, err := os.Getwd()

	if err != nil {
		fmt.Errorf(color.RedString("[!]getwd Err:%s"), err)
	}
	toolPath := util.GetPath(currentPaht, "scanTool")
	reportPath := util.GetPath(currentPaht, "logs")
	args["target"] = target
	args["toolPath"] = toolPath

	var wg sync.WaitGroup

	fmt.Fprintf(color.Output, green("[!]Task Number [%d] Is Starting!\n"), taskNumber+1)
	for key, scannerInfo := range scannerInfoList {
		args["random"] = strings.TrimSpace(fmt.Sprintf("%d", util.GenerateRandInt()+key))
		args["ToolDirectoryName"] = scannerInfo.ToolDirectoryName
		args["ScannerName"] = scannerInfo.ScannerName
		//fmt.Println(args["random"])
		taskCmd, err := InterpolateWithArgs(scannerInfo.SupportedPlatform[osInfo].StartCmd, args)
		if err != nil {
			fmt.Printf(color.RedString("[!]Interpolate With Args Err:%s\n"), err)
		}
		fmt.Printf(color.GreenString("[!]The Task Of [%s] Scanning [%s] Is Staring!\n"), args["ScannerName"], target)
		//fmt.Println(taskCmd)
		wg.Add(1)
		go func(reportFlag bool, regexString, scannerName, toolPath, toolDirectoryName string, taskNumber int64) {
			defer wg.Done()
			defer fmt.Printf(color.GreenString("[!]The Task Of [%s] Scanning [%s] Is Completed!\n"), scannerName, target)
			if len(scannerInfo.SupportedPlatform[osInfo].StartCmd) != 0 {
				toolDirectoryName := util.GetPath(toolPath, toolDirectoryName)
				result, err := ExecuteCmd(osInfo, taskCmd, toolDirectoryName)
				if err != nil {
					fmt.Println(err)
				}
				timeString := time.Now().Format("20060102150405")
				logName := scannerName + "_" + timeString + "_" + regexString + ".log"
				util.SavePrintAsLog(result, reportPath, logName)
			}
			//报送给tg
			toolPath = util.GetPath(toolPath, toolDirectoryName)
			//fmt.Println(bot, int64(ChatId), toolPath, args["random"])
			util.SendScannerReportToTelegramBot(reportFlag, bot, chatId, scannerName, toolPath, regexString)
		}(reportFlag, args["random"], args["ScannerName"], args["toolPath"], args["ToolDirectoryName"], taskNumber)

	}
	wg.Wait()
}

func StartScannerTaskFromFile(filePath, BotToken string, thread, chatId int64, reportFlag bool, scannerList []string, inputargs []map[string]string) {
	currentPath, err := os.Getwd()
	if err != nil {
		fmt.Errorf(color.RedString("[!]getwd Err:%s"), err)
	}
	targetFile := util.GetPath(currentPath, filePath)

	targetList, err := util.ReadLinesFromFile(targetFile)
	if err != nil {
		fmt.Printf(color.RedString("%s"), err)
	}
	var wg = util.NewWaitGroup(int(thread))
	var bot *tgbotapi.BotAPI
	if reportFlag == true {
		bot, err = InitTelegramBot(BotToken)
		if err != nil {
			fmt.Printf(color.YellowString("[!]Connect To TelegramBot Failed,Err:%s\n"), err)
			reportFlag = false
		}
	}
	for key, target := range targetList {
		wg.BlockAdd()
		go func(target, BotToken string, key, chatId int64, reportFlag bool) {
			//StartScannerTask(target, BotToken, chatId, reportFlag)
			StartScannerSingleTask(wg, bot, target, key, chatId, reportFlag, scannerList, inputargs)
		}(target, BotToken, int64(key), chatId, reportFlag)

	}
	wg.Wait()
}

func SelectTemplateScanner(infoList []*ScannerInfo, selectInfoList []string) []*ScannerInfo {
	var selectedScannerInfo []*ScannerInfo
	for _, scanInfo := range infoList {
		for _, selectInfo := range selectInfoList {
			if strings.ToLower(scanInfo.ScannerName) == strings.ToLower(selectInfo) {
				selectedScannerInfo = append(selectedScannerInfo, scanInfo)
			}
		}
	}
	return selectedScannerInfo
}
