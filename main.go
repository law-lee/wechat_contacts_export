package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/eatmoreapple/openwechat"
)

func main() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式

	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	defer reloadStorage.Close()
	if err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	if err != nil {
		fmt.Println("get friends err: %v", err)
	}
	data := make([][]string, 0)
	header := []string{"昵称", "备注名称", "性别", "省", "市", "签名"}
	data = append(data, header)
	var sex string
	for _, f := range friends {
		if f.Sex == 1 {
			sex = "男"
		} else {
			sex = "女"
		}
		profile := []string{f.NickName, f.RemarkName, sex, f.Province, f.City, f.Signature}
		data = append(data, profile)
	}
	fileName := fmt.Sprintf("wechat_friends_%s", time.Now().Format("2006_01_02_15_04_05"))
	err = writeCSV(fileName, data)
	if err != nil {
		fmt.Printf("write csv err: %v", err)
	}
	// 获取所有的群组
	//groups, err := self.Groups()
	//fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	//bot.Block()
}

func writeCSV(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range data {
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
