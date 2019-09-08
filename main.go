package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hekmon/transmissionrpc"

	"github.com/BurntSushi/toml"

	"github.com/gin-gonic/gin"

	"github.com/imroc/req"
)

//Config 配置文件
type Config struct {
	BotAPI         string
	BotSecret      string
	ProxyURL       string
	OnlyFree       bool
	FirstSend      bool
	TRMode         int
	TRFilePath     string
	TRUsername     string
	TRPassword     string
	TRRequestURL   string
	TRRequestPort  uint16
	MoeCatUsername string
	MoeCatPassword string
	PTHomeUsername string
	PTHomePassword string
	PTerUsername   string
	PTerPassword   string
	PTerSSL        bool
	PTerVerify     bool
}

var r = req.New()
var rq = req.New()
var torbot BotAPI
var chanTorUp = make(chan *Update, 100)
var chanimage = make(chan string)
var chatID int64
var cfg Config
var chanGoPTer = make(chan int)

//MoeCatTorList is
var MoeCatTorList []Torrent

//PTHomeTorList is
var PTHomeTorList []Torrent

//PTerTorList is
var PTerTorList []Torrent

func main() {
	if _, err := toml.DecodeFile("./config.toml", &cfg); err != nil {
		log.Println(err)
	}
	if cfg.ProxyURL != "" {
		rq.SetProxyUrl(cfg.ProxyURL)
	}
	go torBotOperator()
	botUpdateListen()
}
func botUpdateListen() {
	torbot.API = cfg.BotAPI
	torbot.Update = make(chan *Update, 500)
	app := gin.Default()
	app.Any("/"+torbot.API, handler)
	app.Run(":9388")
}
func handler(c *gin.Context) {
	var u Update
	reader := c.Request.Body
	body, err := ioutil.ReadAll(reader)
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	// log.Println(string(prettyJSON.Bytes()))
	err = json.Unmarshal(body, &u)
	if err != nil {
		log.Println(err)
	}
	chanTorUp <- &u
}
func torBotOperator() {
	// c := cron.New()
	// c.AddFunc("0 */15 * * * ?", checkMoeCat)
	// c.AddFunc("0 1-59/15 * * * ?", checkPTHome)
	// c.Start()
	for u := range chanTorUp {
		if u.CallbackQuery == nil {
			if u.Message.Text == cfg.BotSecret {
				chatID = u.Message.Chat.ID
				if cfg.PTHomeUsername != "" {
					go loginPTHome(cfg.PTHomeUsername, cfg.PTHomePassword)
				} else {
					chanGoPTer <- 1
				}
				if cfg.MoeCatUsername != "" {
					loginMoeCat(cfg.MoeCatUsername, cfg.MoeCatPassword)
				}
				if cfg.PTerUsername != "" {
					go loginPTer(cfg.PTerUsername, cfg.PTerPassword)
				}
			}
			if u.Message.Text == "/get" {
				if chatID != 0 {
					MoeCatTorList = getTorMoeCat()
					PTHomeTorList = getPTHome()
					PTerTorList = getPTer()
					if cfg.FirstSend {
						sendTorrents(&MoeCatTorList)
						sendTorrents(&PTHomeTorList)
						sendTorrents(&PTerTorList)
					}
					go checkAll()
				}
			}
			if chatID != 0 {
				if len(u.Message.Text) == 6 {
					chanimage <- u.Message.Text
				}
			}

		} else if u.CallbackQuery.Data == "4" {
			torbot.deleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
		} else if strings.HasPrefix(u.CallbackQuery.Data, "http") {
			if len(u.CallbackQuery.Data) > 30 {
				downloadTorrent(u.CallbackQuery.Data)
				torbot.deleteMessage(u.CallbackQuery.Message.Chat.ID, u.CallbackQuery.Message.MessageID)
			}
		}
	}
}
func downloadTorrent(URL string) {
	resp, err := r.Get(URL)
	if err != nil {
		log.Println(err)
	}
	hh := resp.Response().Header
	bb, _ := url.QueryUnescape(hh["Content-Disposition"][0])
	cc := strings.Split(bb, "=")
	if cfg.TRMode == 0 {
		err = resp.ToFile(cfg.TRFilePath + cc[1])
		if err != nil {
			log.Println(err)
		}
	}
	if cfg.TRMode == 1 {
		input, _ := resp.ToBytes()
		transmissionbt, err := transmissionrpc.New(cfg.TRRequestURL, cfg.TRUsername, cfg.TRPassword, &transmissionrpc.AdvancedConfig{
			HTTPS: false,
			Port:  cfg.TRRequestPort,
		})
		if err != nil {
			log.Println(err)
		}
		base64Torrent := base64.StdEncoding.EncodeToString(input)
		tr, err := transmissionbt.TorrentAdd(&transmissionrpc.TorrentAddPayload{
			MetaInfo: &base64Torrent,
			// DownloadDir: &cfg.TRFilePath,
		})
		if err != nil {
			log.Println(err)
		}
		log.Println(tr)

	}
}
func sendTorrents(torrents *[]Torrent) {
	for i, v := range *torrents {
		if cfg.OnlyFree {
			if !strings.Contains(v.Sales, "免费") || !strings.Contains(v.Sales, "Free") {
				v.Update = true
			}
		}
		if !v.Update || v.HasUp {
			km := fmt.Sprintf("%s▽%s@CBD1|%s@CBD3,下载@CBD%s|%s@CBD21|移除@CBD4", v.Living, v.Size, v.Sales, v.URL, v.Four)
			stat := ""
			if v.Status != "" {
				stat = "\n<b>" + v.Status + "</b>"
				(*torrents)[i].Status = ""
			}
			torbot.sendMessage(chatID, v.Site+"〄"+v.Ttype+"▷"+v.Title1+"\n"+v.Title2+stat, HTML, InlineKM(km))
			log.Println(v.HasUp, v.Update, v.URL, v.Title1)
			(*torrents)[i].Update = true
			if len(v.URL) > 30 {
				(*torrents)[i].HasUp = false
			}
		}
	}
}
func checkUpdate(tlist *[]Torrent, oldlist *[]Torrent) {
	for _, oldt := range *oldlist {
		for j, newt := range *tlist {
			if newt.URL == oldt.URL && len(oldt.URL) > 30 {
				log.Println(oldt.URL)
				if strings.Split(newt.Sales, "△")[0] != strings.Split(oldt.Sales, "△")[0] {
					(*tlist)[j].Status = oldt.Sales + "==>" + newt.Sales
					(*tlist)[j].Update = false
				} else if len(oldt.URL) < 30 {
					(*tlist)[j].HasUp = true
				} else {
					(*tlist)[j].Update = true
				}
			} else if len(newt.URL) < 30 {
				(*tlist)[j].HasUp = true
			}
		}
	}
}
func checkMoeCat() {
	log.Println("check MoeCat")
	if cfg.MoeCatUsername != "" {
		tlist := getTorMoeCat()
		checkUpdate(&tlist, &MoeCatTorList)
		sendTorrents(&tlist)
		MoeCatTorList = tlist
	}
}
func checkPTHome() {
	log.Println("check PTHome")
	if cfg.PTHomeUsername != "" {
		tlist := getPTHome()
		checkUpdate(&tlist, &PTHomeTorList)
		sendTorrents(&tlist)
		PTHomeTorList = tlist
	}
}
func checkPTer() {
	log.Println("check PTer")
	if cfg.PTerUsername != "" {
		tlist := getPTer()
		checkUpdate(&tlist, &PTerTorList)
		sendTorrents(&tlist)
		PTerTorList = tlist
	}
}
func checkAll() {
	log.Println("start checkall")
	d := time.Duration(time.Minute * 15)
	t := time.NewTicker(d)
	defer t.Stop()
	defer log.Println("chackAll 函数退出")
	for {
		<-t.C
		checkMoeCat()
		checkPTHome()
		checkPTer()
		log.Println("check   OK")
	}
}
