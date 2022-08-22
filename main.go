package main

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	officialConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"log"
	"net/http"
	"os"
	"time"
	"wechat_tpl_msg/conf"
)

var official *officialaccount.OfficialAccount

func init() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	conf.Parse()

	official = wechat.NewWechat().GetOfficialAccount(&officialConfig.Config{
		AppID:     conf.DefaultConfig.WechatOfficial.AppID,
		AppSecret: conf.DefaultConfig.WechatOfficial.AppSecret,
		Cache:     cache.NewMemory(),
	})

	if conf.DefaultConfig.Mod == "test" {
		log.Println("当前是测试模式，将立即发送一条消息并退出，如需定时发送请将 mod 值改为其他任意值，只要不是 test 就行")
		sendTplMsg()
		os.Exit(0)
	}
}

var cronEntryId cron.EntryID

func main() {
	c := cron.New()
	var err error
	cronEntryId, err = c.AddFunc(conf.DefaultConfig.Cron, func() {
		sendTplMsg()
		log.Println("执行成功, 下次执行时间是", c.Entry(cronEntryId).Next.String())
	})
	if err != nil {
		log.Panicf("添加定时任务失败 err:%s", err.Error())
	}
	c.Start()
	log.Println("启动成功, 下次执行时间是", c.Entry(cronEntryId).Next.String())
	select {}
}

func sendTplMsg() {
	//今日天气 ：{{weath.DATA}}
	//今日温度 ：{{tem.DATA}}
	//在一起：{{totalLoveDay.DATA}}天
	//生日还有：{{birthDay.DATA}}天
	//{{qinghua.DATA}}
	weath := getWeather()
	qinghua := qingHuaMsg()
	for _, openId := range conf.DefaultConfig.WechatOfficial.OpenIds {
		_, err := official.GetTemplate().Send(&message.TemplateMessage{
			ToUser:     openId,
			TemplateID: conf.DefaultConfig.WechatOfficial.TemplateID,
			Data: map[string]*message.TemplateDataItem{
				"weath":        {Value: weath.Wea},
				"tem":          {Value: weath.Tem},
				"totalLoveDay": {Value: fmt.Sprintf("%d", conf.DefaultConfig.GetLoveDay())},
				"birthDay":     {Value: fmt.Sprintf("%d", conf.DefaultConfig.GetBirthDay())},
				"qinghua":      {Value: qinghua, Color: conf.DefaultConfig.Colors.Qinghua},
			},
		})

		if err != nil {
			log.Printf("发送模版消息失败 openId=[%s] err:%s\n", openId, err.Error())
		}
	}
}

// 获取情话消息
func qingHuaMsg() string {
	resp, err := http.Get("https://api.uomg.com/api/rand.qinghua")
	if err != nil {
		log.Println("获取情话消息失败 err:", err.Error())
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("获取情话消息失败 status:", resp.Status)
		return ""
	}

	res := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("获取情话消息失败 decode err:", err.Error())
		return ""
	}
	return res["content"].(string)
}

type weather struct {
	Nums       int    `json:"nums"`
	Cityid     string `json:"cityid"`
	City       string `json:"city"`
	Date       string `json:"date"`
	Week       string `json:"week"`
	UpdateTime string `json:"update_time"`
	Wea        string `json:"wea"`
	WeaImg     string `json:"wea_img"`
	Tem        string `json:"tem"`
	TemDay     string `json:"tem_day"`
	TemNight   string `json:"tem_night"`
	Win        string `json:"win"`
	WinSpeed   string `json:"win_speed"`
	WinMeter   string `json:"win_meter"`
	Air        string `json:"air"`
	Pressure   string `json:"pressure"`
	Humidity   string `json:"humidity"`
}

func getWeather() weather {
	uri := fmt.Sprintf("https://www.yiketianqi.com/free/day?appid=%s&appsecret=%s&unescape=1&cityid=%s",
		conf.DefaultConfig.Yiketianqi.Appid, conf.DefaultConfig.Yiketianqi.Appsecret, conf.DefaultConfig.Yiketianqi.Cityid)
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("获取天气消息失败 err:", err.Error())
		return weather{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("获取天气消息失败 status:", resp.Status)
		return weather{}
	}

	res := weather{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Println("获取天气消息失败 decode err:", err.Error())
		return weather{}
	}
	return res
}
