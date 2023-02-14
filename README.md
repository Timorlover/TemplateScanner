# TemplateScanner
## 描述
通过编写模板来集成各种扫描器进行批量扫描。
![image](https://user-images.githubusercontent.com/116296194/218671217-ee0f5462-df0a-4121-bf6e-0e8e596a7a67.png)
![image](https://user-images.githubusercontent.com/116296194/218671275-277bde67-ba20-4419-9705-7f0c9276dae5.png)
## 帮助
```NAME:
   TemplateScanner - A new cli application

USAGE:
   TemplateScanner [global options] command [command options] [arguments...]

VERSION:
   1.0

AUTHOR:
   爱吃火锅的提莫

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --target value                                 Scan Target,eg:-target https://www.exmaple.com
   --file value                                   Scan Targets,eg:-file urls.txt
   --token value                                  TelegramBot Token,eg:-token 172323232124:AAEVm_Sv32CCnnpJpu32332zvk_NDIa83230U
   --id value                                     TelegramBot ChatId,eg:-id 5932323032 (default: 0)
   --thread value                                 Number of concurrent executions,eg:-thread 3 (default: 3)
   --list                                         List scanner template info,-list (default: false)
   --templatelist value [ --templatelist value ]  Select scanner template to scan,eg: xray,nuclei
   --args value                                   Replace the input parameters into the template,eg:-args proxy=http://127.0.0.1:8080,$HOME=/tmp
   --help, -h                                     show help (default: false)
   --version, -v                                  print the version (default: false)
   ```
![image](https://user-images.githubusercontent.com/116296194/218672012-2d03919a-32a2-44ee-a18f-e2fbda81ef60.png)
## 效果
![image](https://user-images.githubusercontent.com/116296194/218671544-ab81625b-22f5-470f-b27c-b85215571ac5.png)
