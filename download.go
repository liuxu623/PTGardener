package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/hekmon/transmissionrpc"
	"github.com/imroc/req"
)

//MoeCatPasskey is
var MoeCatPasskey string

//PTHomePasskey is
var PTHomePasskey string

//PTerPasskey is
var PTerPasskey string

//HDStreetPasskey is
var HDStreetPasskey string

//CHDBitsPasskey is
var CHDBitsPasskey string

//OurBitsPasskey is
var OurBitsPasskey string

//HDSkyPasskey is
var HDSkyPasskey string

// SSDPasskey is
var SSDPasskey string

// FrdsPasskey is
var FrdsPasskey string

// MTPasskey is
var MTPasskey string

func downloadTorrent(URL string) {
	log.Println(URL)
	var passkey string
	// resp, err := r.Get(URL)
	// if err != nil {
	// 	log.Println(err)
	// }
	// hh := resp.Response().Header
	// bb, _ := url.QueryUnescape(hh["Content-Disposition"][0])
	// cc := strings.Split(bb, "=")
	// if cfg.DownloadMode == 0 {
	// 	err = resp.ToFile(cfg.DownloadFilePath + cc[1])
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }
	if strings.Contains(URL, "www.moecat.best") {
		if MoeCatPasskey == "" {
			header := req.Header{
				"cookie": cfg.MoeCatCookies,
			}
			respp, err := r.Get("https://www.moecat.best/usercp.php", header)
			if err != nil {
				log.Println(err)
			}
			re := regexp.MustCompile(`(?m)passkey=(.+?)<`)
			html, _ := respp.ToString()
			MoeCatPasskey = re.FindStringSubmatch(html)[1]
		}
		passkey = MoeCatPasskey
	}
	if strings.Contains(URL, "hdstreet.club") {
		if HDStreetPasskey == "" {
			header := req.Header{
				"cookie": cfg.HDStreetCookies,
			}
			respp, err := r.Get("https://hdstreet.club/usercp.php", header)
			if err != nil {
				log.Println(err)
			}
			re := regexp.MustCompile(`(?m)passkey=(.+?)<`)
			html, _ := respp.ToString()
			HDStreetPasskey = re.FindStringSubmatch(html)[1]
		}
		passkey = HDStreetPasskey
	}
	if strings.Contains(URL, "pthome.net") {
		if PTHomePasskey == "" {
			header := req.Header{
				"cookie": cfg.PTHomeCookies,
			}
			respp, err := r.Get("https://www.pthome.net/usercp.php", header)
			if err != nil {
				log.Println(err)
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密钥" {
					PTHomePasskey = s.Find("td").Eq(1).Text()
				}
			})
		}
		passkey = PTHomePasskey
	}
	if strings.Contains(URL, "pter.club") {
		if PTerPasskey == "" {
			header := req.Header{
				"cookie": cfg.PTerCookies,
			}
			respp, err := r.Get("https://pter.club/usercp.php", header)
			if err != nil {
				log.Println(err)
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密钥" {
					PTerPasskey = s.Find("td").Eq(1).Text()
				}
			})
		}
		passkey = PTerPasskey
	}
	if strings.Contains(URL, "chdbits.co") {
		if CHDBitsPasskey == "" {
			header := req.Header{
				"cookie": cfg.CHDBitsCookies,
			}
			respp, err := r.Get("https://chdbits.co/usercp.php", header)
			if err != nil {
				log.Println(err)
				log.Println("Get CHDBits Usercp Failed")
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密钥" {
					CHDBitsPasskey = s.Find("td").Eq(1).Text()
				}
			})
		}
		passkey = CHDBitsPasskey
	}
	if strings.Contains(URL, "ourbits.club") {
		if OurBitsPasskey == "" {
			header := req.Header{
				"cookie": cfg.OurBitsCookies,
			}
			respp, err := r.Get("https://ourbits.club/usercp.php", header)
			if err != nil {
				log.Println(err)
				log.Println("Get OurBits Usercp Failed")
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密钥" {
					OurBitsPasskey = s.Find("td").Eq(1).Text()
				}
			})
		}
		passkey = OurBitsPasskey
	}
	if strings.Contains(URL, "hdsky.me") {
		if HDSkyPasskey == "" {
			header := req.Header{
				"cookie": cfg.HDSkyCookies,
			}
			respp, err := r.Get("https://hdsky.me/usercp.php", header)
			if err != nil {
				log.Println(err)
				log.Println("Get HDSky Usercp Failed")
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密钥" {
					HDSkyPasskey = s.Find("td").Eq(1).Text()
				}
			})
		}
		passkey = HDSkyPasskey
	}
	if strings.Contains(URL, "springsunday.net") {
		if SSDPasskey == "" {
			header := req.Header{
				"cookie": cfg.SSDCookies,
			}
			respp, err := r.Get("https://springsunday.net/usercp.php", header)
			if err != nil {
				log.Println(err)
				log.Println("Get SSD Usercp Failed")
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密钥" {
					SSDPasskey = s.Find("td").Eq(1).Text()
				}
			})
		}
		passkey = SSDPasskey
	}
	if strings.Contains(URL, "keepfrds.com") {
		if FrdsPasskey == "" {
			header := req.Header{
				"cookie": cfg.FrdsCookies,
			}
			respp, err := r.Get("https://pt.keepfrds.com/usercp.php", header)
			if err != nil {
				log.Println(err)
				log.Println("Get Frds Usercp Failed")
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密钥" {
					FrdsPasskey = s.Find("td").Eq(1).Text()
				}
			})
		}
		passkey = FrdsPasskey
	}
	if strings.Contains(URL, "m-team.cc") {
		if MTPasskey == "" {
			header := req.Header{
				"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
				"Accept-Encoding":           "gzip, deflate, br",
				"Accept-Language":           "zh-CN,zh;q=0.9,zh-TW;q=0.8",
				"Cache-Control":             "max-age=0",
				"Connection":                "keep-alive",
				"Cookie":                    cfg.MTCookies,
				"DNT":                       "1",
				"Host":                      "pt.m-team.cc",
				"Referer":                   "https://pt.m-team.cc/torrents.php",
				"Sec-Fetch-Mode":            "navigate",
				"Sec-Fetch-Site":            "none",
				"Sec-Fetch-User":            "?1",
				"Upgrade-Insecure-Requests": "1",
				"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			}
			respp, err := r.Get("https://pt.m-team.cc/usercp.php", header)
			if err != nil {
				log.Println(err)
				log.Println("Get MTeam Usercp Failed")
			}
			doc, err := goquery.NewDocumentFromResponse(respp.Response())
			doc.Find("tr").Each(func(i int, s *goquery.Selection) {
				if s.Find("td").Eq(0).Text() == "密匙" {
					MTPasskey = s.Find("td").Eq(1).Text()
					if cfg.MTTrackerIPv6 {
						MTPasskey += "&ipv6=1"
					}
					if cfg.MTTrackerSSL {
						MTPasskey += "&https=1"
					}
				}
			})
		}
		passkey = MTPasskey
	}
	if cfg.DownloadMode == 1 {
		//input, _ := resp.ToBytes()
		transmissionbt, err := transmissionrpc.New(cfg.TRRequestURL, cfg.TRUsername, cfg.TRPassword, &transmissionrpc.AdvancedConfig{
			HTTPS: false,
			Port:  cfg.TRRequestPort,
		})
		if err != nil {
			log.Println(err)
		}
		// base64Torrent := base64.StdEncoding.EncodeToString(input)
		downlink := URL + "&passkey=" + passkey
		_, err = transmissionbt.TorrentAdd(&transmissionrpc.TorrentAddPayload{
			Filename: &downlink,
			// MetaInfo: &base64Torrent,
			// DownloadDir: &cfg.TRFilePath,
		})
		if err != nil {
			log.Println(err)
		}
		// log.Println(*tr.ID)
		// log.Println(*tr.Name)
		// log.Println(*tr.HashString)
	}
	if cfg.DownloadMode == 2 {
		loginparam := req.Param{
			"username": cfg.QBUsername,
			"password": cfg.QBPassword,
		}
		resp, err := r.Post("http://"+cfg.QBRequestURL+":"+strconv.Itoa(cfg.QBRequestPort)+"/api/v2/auth/login", loginparam)
		if err != nil {
			log.Println(err)
		}

		torparam := req.Param{
			"urls": URL + "&passkey=" + passkey,
		}
		resp, err = r.Post("http://"+cfg.QBRequestURL+":"+strconv.Itoa(cfg.QBRequestPort)+"/api/v2/torrents/add", torparam)
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.ToString())
	}

}

// TrReannounce is
func TrReannounce(i int64) {
	transmissionbt, err := transmissionrpc.New(cfg.TRRequestURL, cfg.TRUsername, cfg.TRPassword, &transmissionrpc.AdvancedConfig{
		HTTPS: false,
		Port:  cfg.TRRequestPort,
	})
	if err != nil {
		log.Println(err)
	}
	for {
		time.Sleep(time.Duration(i) * time.Second)
		log.Println("======获取更多Peer=======")
		err = transmissionbt.TorrentReannounceRecentlyActive()
		if err != nil {
			log.Println(err)
		}
	}
}
