package main

import (
	"TemplateScanner/models"
	"TemplateScanner/util"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func init() {
	dirs := []string{"config", "logs", "scanTool"}
	for _, dir := range dirs {
		err := util.CheckAndCreateFile(dir)
		if err != nil {
			fmt.Println(color.RedString("Check Dirs Err:"), err)
		}
	}
}

func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "target",
			Value: "",
			Usage: "Scan Target,eg:-target https://www.exmaple.com",
		},
		&cli.StringFlag{
			Name:  "file",
			Value: "",
			Usage: "Scan Targets,eg:-file urls.txt",
		},
		&cli.StringFlag{
			Name:  "token",
			Value: "",
			Usage: "TelegramBot Token,eg:-token 172323232124:AAEVm_Sv32CCnnpJpu32332zvk_NDIa83230U",
		},
		&cli.Int64Flag{
			Name:  "id",
			Value: 0,
			Usage: "TelegramBot ChatId,eg:-id 5932323032",
		},
		&cli.Int64Flag{
			Name:  "thread",
			Value: 3,
			Usage: "Number of concurrent executions,eg:-thread 3",
		},
		&cli.BoolFlag{
			Name:  "list",
			Usage: "List scanner template info,-list",
		},
		&cli.StringSliceFlag{
			Name:  "templatelist",
			Usage: "Select scanner template to scan,eg: xray,nuclei",
		},
		&cli.StringFlag{
			Name:  "args",
			Usage: "Replace the input parameters into the template,eg:-args proxy=http://127.0.0.1:8080,$HOME=/tmp",
		},
	}

	app.Action = func(c *cli.Context) error {

		var token string
		var id int64
		var thread int64
		var templateList []string
		var inputArgs []map[string]string

		list := c.Bool("list")
		target := c.String("target")
		file := c.String("file")
		inputArgs = util.ParseMapFlag(c.String("args"))
		configData, err := models.LoadConfigYaml()
		if err != nil {
			fmt.Printf(color.RedString("Load ConfigYaml Err:%s\n", err))
			configData, _ = models.GenerateConfigYaml()
		}

		if (target != "" || file != "") && c.NumFlags() == 1 && err == nil {
			//fmt.Println(configData)
			id = int64(configData.Telegram.ChatID)
			token = configData.Telegram.Token
			thread = int64(configData.Thread)
			templateList = configData.LoadTemplate
			inputArgs = configData.InputArguments
		} else {
			id = c.Int64("id")
			token = c.String("token")
			thread = c.Int64("thread")
			templateList = c.StringSlice("templatelist")
		}

		if list {
			models.ListTemplateScannerInfo()
			os.Exit(1)
		}
		if target != "" {
			if token != "" && id != 0 && util.ConnectTeleBotTest(token, id) {
				models.StartScannerTask(target, token, id, true, templateList, inputArgs)
			} else {
				models.StartScannerTask(target, token, id, false, templateList, inputArgs)
			}
		} else if file != "" {
			if token != "" && id != 0 && util.ConnectTeleBotTest(token, id) {
				models.StartScannerTaskFromFile(file, token, thread, id, true, templateList, inputArgs)
			} else {
				models.StartScannerTaskFromFile(file, token, thread, id, false, templateList, inputArgs)
			}
		} else {
			fmt.Println(color.RedString("Please Input Right Target!"))
			fmt.Println(color.YellowString("Example:TemplateScanner.exe -target https://www.baidu.com -token 1718205114:AAEVm_SviixsyCnnpJpu6VPzvk_NDIa8n0U -id 594478092"))
			fmt.Println(color.YellowString("Example:TemplateScanner.exe -file urls.txt -token 1718205114:AAEVm_SviixsCCnn56pu6VPzvk_NDIa8n0U -id 592222092 -thread 3"))
		}
		return nil
	}

	app.Name = "TemplateScanner"
	app.Version = "1.0"
	app.Authors = []*cli.Author{{"爱吃火锅的提莫", ""}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
