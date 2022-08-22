package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var DefaultConfig = &Conf{
	Colors: Colors{
		Qinghua: "#550038",
	},
}

type Conf struct {
	Mod            string         `json:"mod"`
	Cron           string         `json:"cron"`
	LoveStartDate  string         `json:"love_start_date"`
	LoveStartTime  time.Time      `json:"-"`
	BirthDate      string         `json:"birth_date"`
	BirthTime      time.Time      `json:"-"`
	WechatOfficial WechatOfficial `json:"wechat_official"`
	Yiketianqi     Yiketianqi     `json:"yiketianqi"`

	Colors Colors `json:"colors"`
}

type Colors struct {
	Qinghua string `json:"qinghua"`
}

func (c *Conf) GetLoveDay() int {
	return int(time.Now().Sub(c.LoveStartTime).Hours() / 24.0)
}

func (c *Conf) GetBirthDay() int {
	birth, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%s", getCurrentDate().Year(), c.BirthDate))
	if err != nil {
		log.Println("GetBirthDay 错误", birth)
		return -1111
	}
	if getCurrentDate().Sub(birth) > 0 {
		birth, _ = time.Parse("2006-01-02", fmt.Sprintf("%d-%s", getCurrentDate().Year()+1, c.BirthDate))
	}
	return int(birth.Sub(getCurrentDate()).Hours() / 24.0)
}

func getCurrentDate() time.Time {
	nowStr := time.Now().Format("2006-01-02")
	now, _ := time.Parse("2006-01-02", nowStr)
	return now
}

// WechatOfficial 微信公众号配置
type WechatOfficial struct {
	AppID      string   `json:"app_id"`      // appid
	AppSecret  string   `json:"app_secret"`  // appsecret
	OpenIds    []string `json:"open_ids"`    // 要接受消息的公众号
	TemplateID string   `json:"template_id"` // 必须, 模版ID
}

// Yiketianqi www.yiketianqi.com/free/day 配置
type Yiketianqi struct {
	Appid     string `json:"appid" `
	Appsecret string `json:"appsecret"`
	Cityid    string `json:"cityid"`
}

func Parse() {
	confBytes, err := ioutil.ReadFile("./config_local.json")
	if err != nil {
		confBytes, err = ioutil.ReadFile("./config.json")
		if err != nil {
			log.Panicf("解析配置文件出错 err: %s", err.Error())
		}
	}
	if err := json.Unmarshal(confBytes, DefaultConfig); err != nil {
		log.Panicf("解析配置文件 Unmarshal 出错  err: %s", err.Error())
	}

	if DefaultConfig.LoveStartTime, err = time.Parse("2006-01-02", DefaultConfig.LoveStartDate); err != nil {
		log.Panicf("解析配置文件 love_start_date 出错  err: %s", err.Error())
	}
	if DefaultConfig.BirthTime, err = time.Parse("01-02", DefaultConfig.BirthDate); err != nil {
		log.Panicf("解析配置文件 birth_date 出错  err: %s", err.Error())
	}
}
