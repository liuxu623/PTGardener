package main

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/imroc/req"
)

func loginMoeCat(username string, password string) {
	loginHeader := req.Header{
		"Referer":          "https://www.moecat.best/login.php",
		"Sec-Fetch-Mode":   "cors",
		"DNT":              "1",
		"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		"X-Requested-With": "XMLHttpRequest",
	}
	loginParam := req.Param{
		"username": username,
		"password": password,
	}
	_, err := r.Post("https://www.moecat.best/takelogin.php", loginParam, loginHeader)
	if err != nil {
		log.Println(err)
	}
	torbot.sendMessage(chatID, "MoeCat 登录成功！")
	// resp.ToFile("222.html")
}
func loginPTHome(username string, password string) {
	goon := true
	for goon {

		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Referer":    "https://www.pthome.net/login.php",
		}
		resp, err := r.Get("https://www.pthome.net/login.php", loginHeader)
		if err != nil {
			log.Println(err)
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		src, ok := doc.Find("td.rowhead").Next().Find("img").Attr("src")
		if !ok {
			log.Println("未找到验证图片")
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		torbot.sendPhoto(chatID, "https://www.pthome.net/"+src, "#请输入验证码")
		imagehash := strings.Split(src, "imagehash=")[1]
		image := <-chanimage
		loginParam := req.Param{
			"username":    username,
			"password":    password,
			"imagestring": image,
			"imagehash":   imagehash,
		}
		resp, err = r.Post("https://www.pthome.net/takelogin.php", loginHeader, loginParam)
		if err != nil {
			log.Println(err)
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "PTHome 登录成功！")
			goon = false
			chanGoPTer <- 1
		} else {
			torbot.sendMessage(chatID, "验证码错误，请重试")
		}
	}
}
func loginPTer(username string, password string) {
	<-chanGoPTer
	goon := true
	for goon {
		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Origin":     "https://pter.club",
			"Referer":    "https://pter.club/login.php",
		}
		resp, err := r.Get("https://pter.club/login.php", loginHeader)
		if err != nil {
			log.Println(err)
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
		}
		src := ""
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			alt, ok := s.Attr("alt")
			if ok {
				if alt == "CAPTCHA" {
					src, _ = s.Attr("src")
				}
			}
		})
		verifycode := ""
		torbot.sendPhoto(chatID, "https://pter.club/"+src, "#请输入验证码")
		imagehash := strings.Split(src, "imagehash=")[1]
		image := <-chanimage
		if cfg.PTerVerify {
			torbot.sendMessage(chatID, "请输入两步验证码，请预留10秒以上过期时间")
			verifycode = <-chanimage
		}
		loginParam := req.Param{}
		if cfg.PTerSSL {
			loginParam = req.Param{
				"username":    username,
				"password":    password,
				"imagestring": image,
				"imagehash":   imagehash,
				"verify_code": verifycode,
				"trackerssl":  "yes",
			}
		} else {
			loginParam = req.Param{
				"username":    username,
				"password":    password,
				"imagestring": image,
				"imagehash":   imagehash,
				"verify_code": verifycode,
			}
		}
		resp, err = r.Post("https://pter.club/takelogin.php", loginHeader, loginParam)
		if err != nil {
			log.Println(err)
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "PTer 登录成功！")
			goon = false
		} else {
			torbot.sendMessage(chatID, "验证码或二步验证错误，不判断了，反正都得重新输入")
		}
	}
}
