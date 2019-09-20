package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/gin-gonic/gin"

	"github.com/imroc/req"
)

//Config 配置文件
type Config struct {
	BotAPI           string
	BotSecret        string
	ProxyURL         string
	TitleStyle       string
	OnlyFree         bool
	FirstSend        bool
	DownloadMode     int
	DownloadFilePath string
	TRUsername       string
	TRPassword       string
	TRRequestURL     string
	TRRequestPort    uint16
	QBUsername       string
	QBPassword       string
	QBRequestURL     string
	QBRequestPort    int
	MoeCatUsername   string
	MoeCatPassword   string
	MoeCatCookies    string
	PTHomeUsername   string
	PTHomePassword   string
	PTHomeCookies    string
	PTerUsername     string
	PTerPassword     string
	PTerCookies      string
	PTerSSL          bool
	PTerVerify       bool
	HDStreetUsername string
	HDStreetPassword string
	HDStreetCookies  string
	CHDBitsUsername  string
	CHDBitsPassword  string
	CHDBitsCookies   string
	OurBitsUsername  string
	OurBitsPassword  string
	OurBitsCookies   string
	OurBitsSSL       bool
	OurBitsVerify    bool
	HDSkyUsername    string
	HDSkyPassword    string
	HDSkyCookies     string
	HDSkyVerify      bool
	SSDUsername      string
	SSDPassword      string
	SSDCookies       string
	FrdsUsername     string
	FrdsPassword     string
	FrdsCookies      string
	FrdsVerify       bool
	MTUsername       string
	MTPassword       string
	MTCookies        string
	MTVerify         bool
	MTTrackerSSL     bool
	MTTrackerIPv6    bool
	OpenCDUsername   string
	OpenCDPassword   string
	OpenCDCookies    string
	OpenCDSSL        bool
	OpenCDTrackerSSL bool
	HDHomeUsername   string
	HDHomePassword   string
	HDHomeCookies    string
}

//不走代理的客户端,用来访问PT
var r = req.New()
var loginOK sync.WaitGroup

// 可能走代理的客户端，用来连接电报
var rq = req.New()
var torbot BotAPI
var chanTorUp = make(chan *Update, 100)

//收验证码的chan
var chanimage = make(chan string)
var chatID int64
var cfg Config

//按顺序处理登录的阻塞
var chanGoPTer = make(chan int, 1)
var chanGoOurBits = make(chan int, 1)
var chanGoHDSky = make(chan int, 1)
var chanGoSSD = make(chan int, 1)
var chanGoFrds = make(chan int, 1)
var chanGoOpenCD = make(chan int, 1)
var chanGoMT = make(chan int, 1)

//MoeCatTorList is
var MoeCatTorList []Torrent

//PTHomeTorList is
var PTHomeTorList []Torrent

//PTerTorList is
var PTerTorList []Torrent

//HDStreetList is
var HDStreetList []Torrent

//CHDBitsList is
var CHDBitsList []Torrent

// OurBitsList is
var OurBitsList []Torrent

// HDSkyList is
var HDSkyList []Torrent

// SSDList is
var SSDList []Torrent

// FrdsList is
var FrdsList []Torrent

// OpenCDList is
var OpenCDList []Torrent

// MTList is
var MTList []Torrent

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

//用gin监听消息
func botUpdateListen() {
	torbot.API = cfg.BotAPI
	torbot.Update = make(chan *Update, 500)
	app := gin.Default()
	app.Any("/"+torbot.API, handler)
	// app.Run("0.0.0.0:8443")
	server := &http.Server{Handler: app}
	l, err := net.Listen("tcp4", "0.0.0.0:9388")
	if err != nil {
		log.Println(err)
	}
	err = server.Serve(l)
}

//将消息转为结构体发送到消息chan
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

//循环处理消息
func torBotOperator() {
	// c := cron.New()
	// c.AddFunc("0 */15 * * * ?", checkMoeCat)
	// c.AddFunc("0 1-59/15 * * * ?", checkPTHome)
	// c.Start()
	firstMsgCount := 1
	for u := range chanTorUp {
		// log.Println("已收到bot传来消息")
		if u.CallbackQuery == nil {
			if u.Message.Text == cfg.BotSecret {
				// log.Println("已确认口令通过")
				chatID = u.Message.Chat.ID
				if cfg.PTHomeUsername != "" || cfg.PTHomeCookies != "" {
					// log.Println("已启动登录协程")
					go loginPTHome(cfg.PTHomeUsername, cfg.PTHomePassword)
				} else {
					chanGoPTer <- 1
				}
				if cfg.MoeCatUsername != "" || cfg.MoeCatCookies != "" {
					loginMoeCat(cfg.MoeCatUsername, cfg.MoeCatPassword)
				}
				if cfg.HDStreetUsername != "" || cfg.HDStreetCookies != "" {
					loginHDStreet(cfg.HDStreetUsername, cfg.HDStreetPassword)
				}
				if cfg.CHDBitsUsername != "" || cfg.CHDBitsCookies != "" {
					loginCHDBits(cfg.CHDBitsUsername, cfg.CHDBitsPassword)
				}
				if cfg.PTerUsername != "" || cfg.PTerCookies != "" {
					go loginPTer(cfg.PTerUsername, cfg.PTerPassword)
				} else {
					chanGoOurBits <- 1
				}
				if cfg.OurBitsUsername != "" || cfg.OurBitsCookies != "" {
					go loginOurBits(cfg.OurBitsUsername, cfg.OurBitsPassword)
				} else {
					chanGoHDSky <- 1
				}
				if cfg.HDSkyUsername != "" || cfg.HDSkyCookies != "" {
					go loginHDSky(cfg.HDSkyUsername, cfg.HDSkyPassword)
				} else {
					chanGoSSD <- 1
				}
				if cfg.SSDUsername != "" || cfg.SSDCookies != "" {
					go loginSSD(cfg.SSDUsername, cfg.SSDPassword)
				}
				if cfg.FrdsUsername != "" || cfg.FrdsCookies != "" {
					go loginFrds(cfg.FrdsUsername, cfg.FrdsPassword)
				}
				if cfg.MTUsername != "" || cfg.MTCookies != "" {
					go loginMT(cfg.MTUsername, cfg.MTPassword)
				}

				go loginDone()
			}
			if u.Message.Text == "/get" {
				if chatID != 0 {
					go getAll()
				} else {
					torbot.sendMessage(u.Message.Chat.ID, "请先输入设定的BotSecret,然后按照提示登录账号")
					firstMsgCount--
				}
			}
			if chatID != 0 {
				if len(u.Message.Text) == 6 {
					chanimage <- u.Message.Text
				}
			} else {
				if firstMsgCount > 0 {
					torbot.sendMessage(u.Message.Chat.ID, "请先输入设定的BotSecret,然后按照提示登录账号")
					firstMsgCount--
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

//将种子发送到电报
func sendTorrents(torrents *[]Torrent) {
	for i, v := range *torrents {
		if cfg.OnlyFree {
			if !strings.Contains(v.Sales, "免费") && !strings.Contains(v.Sales, "Free") {
				v.Update = true
			}
		}
		if !v.Update {
			km := fmt.Sprintf("%s▽%s@CBD1|%s@CBD3,下载@CBD%s|%s@CBD21|移除@CBD4", v.Living, v.Size, v.Sales, v.URL, v.Four)
			torbot.sendMessage(chatID, titleEscape(v), HTML, NoPreView, InlineKM(km))
			log.Println(v.HasUp, v.Update, v.URL, v.Title1)
			(*torrents)[i].Update = true
			(*torrents)[i].Status = ""
		}
	}
}

func titleEscape(t Torrent) string {
	stat := ""
	if t.Status != "" {
		t.Status = strings.Replace(t.Status, "<", "", 1)
		stat = "<b>" + t.Status + "</b>"
	}
	title := cfg.TitleStyle
	title = strings.Replace(title, "$Site", t.Site, -1)
	title = strings.Replace(title, "$Type", t.Ttype, -1)
	title = strings.Replace(title, "$Title1", t.Title1, -1)
	title = strings.Replace(title, "$Title2", t.Title2, -1)
	title = strings.Replace(title, "$Change", "\n"+stat, -1)
	title = strings.Replace(title, "$change", stat, -1)
	title = strings.Replace(title, "$Living", t.Living, -1)
	title = strings.Replace(title, "$Size", t.Size, -1)
	title = strings.Replace(title, "$Sales", t.Sales, -1)
	title = strings.Replace(title, "$Four", t.Four, -1)
	title = strings.Replace(title, "$URL", t.URL, -1)
	title = strings.Replace(title, "\\n", "\n", -1)
	return title
}

//对比两次获得的页面是否有变动,修改标记是否发送
func checkUpdate(tlist *[]Torrent, oldlist *[]Torrent) {
	for _, oldt := range *oldlist {
		for j, newt := range *tlist {
			if newt.URL == oldt.URL && len(newt.URL) > 30 {
				if strings.Split(newt.Sales, "△")[0] != strings.Split(oldt.Sales, "△")[0] {
					(*tlist)[j].Status = oldt.Sales + "==>" + newt.Sales
					(*tlist)[j].Update = false
				} else {
					(*tlist)[j].Update = true
				}
			}
			if len(newt.URL) < 30 {
				(*tlist)[j].Update = true
			}
		}
	}
}

func checkMoeCat() {
	if cfg.MoeCatUsername != "" || cfg.MoeCatCookies != "" {
		log.Println("check MoeCat")
		tlist := getMoeCat()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &MoeCatTorList)
			sendTorrents(&tlist)
			MoeCatTorList = tlist
		}
	}
}
func checkPTHome() {
	if cfg.PTHomeUsername != "" || cfg.PTHomeCookies != "" {
		log.Println("check PTHome")
		tlist := getPTHome()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &PTHomeTorList)
			sendTorrents(&tlist)
			PTHomeTorList = tlist
		}
	}
}
func checkPTer() {
	if cfg.PTerUsername != "" || cfg.PTerCookies != "" {
		log.Println("check PTer")
		tlist := getPTer()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &PTerTorList)
			sendTorrents(&tlist)
			PTerTorList = tlist
		}
	}
}
func checkHDStreet() {
	if cfg.HDStreetUsername != "" || cfg.HDStreetCookies != "" {
		log.Println("check HDStreet")
		tlist := getHDStreet()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &HDStreetList)
			sendTorrents(&tlist)
			HDStreetList = tlist
		}
	}
}
func checkCHDBits() {
	if cfg.CHDBitsUsername != "" || cfg.CHDBitsCookies != "" {
		log.Println("check CHDBits")
		tlist := getCHDBits()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &CHDBitsList)
			sendTorrents(&tlist)
			CHDBitsList = tlist
		}
	}
}
func checkOurBits() {
	if cfg.OurBitsUsername != "" || cfg.OurBitsCookies != "" {
		log.Println("check OurBits")
		tlist := getOurBits()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &OurBitsList)
			sendTorrents(&tlist)
			OurBitsList = tlist
		}
	}
}
func checkHDSky() {
	if cfg.HDSkyUsername != "" || cfg.HDSkyCookies != "" {
		log.Println("check HDSky")
		tlist := getHDSkey()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &HDSkyList)
			sendTorrents(&tlist)
			HDSkyList = tlist
		}
	}
}
func checkSSD() {
	if cfg.SSDUsername != "" || cfg.SSDCookies != "" {
		log.Println("check SSD")
		tlist := getSSD()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &SSDList)
			sendTorrents(&tlist)
			SSDList = tlist
		}
	}
}
func checkFrds() {
	if cfg.FrdsUsername != "" || cfg.FrdsCookies != "" {
		log.Println("check Frds")
		tlist := getFrds()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &FrdsList)
			sendTorrents(&tlist)
			FrdsList = tlist
		}
	}
}
func checkMT() {
	if cfg.MTUsername != "" || cfg.MTCookies != "" {
		log.Println("check MTeam")
		tlist := getMT()
		if len(tlist) != 0 {
			checkUpdate(&tlist, &MTList)
			sendTorrents(&tlist)
			MTList = tlist
		}
	}
}

//定时运行检查
func checkAll() {
	log.Println("start checkall")
	d := time.Duration(time.Minute * 10)
	t := time.NewTicker(d)
	defer t.Stop()
	defer log.Println("chackAll 函数退出")
	for {
		<-t.C
		checkMoeCat()
		checkPTHome()
		checkPTer()
		checkHDStreet()
		checkCHDBits()
		checkOurBits()
		checkHDSky()
		checkSSD()
		checkFrds()
		checkMT()
		log.Println("check   OK")
	}
}
func getAll() {
	if cfg.HDStreetUsername != "" || cfg.HDStreetCookies != "" {
		HDStreetList = getHDStreet()
	}
	if cfg.MoeCatUsername != "" || cfg.MoeCatCookies != "" {
		MoeCatTorList = getMoeCat()
	}
	if cfg.PTHomeUsername != "" || cfg.PTHomeCookies != "" {
		PTHomeTorList = getPTHome()
	}
	if cfg.PTerUsername != "" || cfg.PTerCookies != "" {
		PTerTorList = getPTer()
	}
	if cfg.CHDBitsUsername != "" || cfg.CHDBitsCookies != "" {
		CHDBitsList = getCHDBits()
	}
	if cfg.OurBitsUsername != "" || cfg.OurBitsCookies != "" {
		OurBitsList = getOurBits()
	}
	if cfg.HDSkyUsername != "" || cfg.HDSkyCookies != "" {
		HDSkyList = getHDSkey()
	}
	if cfg.SSDUsername != "" || cfg.SSDCookies != "" {
		SSDList = getSSD()
	}
	if cfg.FrdsUsername != "" || cfg.FrdsCookies != "" {
		FrdsList = getFrds()
	}
	if cfg.MTUsername != "" || cfg.MTCookies != "" {
		MTList = getMT()
	}
	if cfg.FirstSend {
		sendTorrents(&MoeCatTorList)
		sendTorrents(&PTHomeTorList)
		sendTorrents(&PTerTorList)
		sendTorrents(&HDStreetList)
		sendTorrents(&CHDBitsList)
		sendTorrents(&OurBitsList)
		sendTorrents(&HDSkyList)
		sendTorrents(&SSDList)
		sendTorrents(&FrdsList)
		sendTorrents(&MTList)
		cfg.FirstSend = false
	}
	go checkAll()
}
