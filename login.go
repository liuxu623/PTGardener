package main

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/imroc/req"
)

func loginDone() {
	log.Println("等待登录")
	loginOK.Wait()
	time.Sleep(time.Duration(2) * time.Second)
	torbot.sendMessage(chatID, "全部账号登录成功，请输入/get获取种子")
}
func loginMoeCat(username string, password string) {
	loginOK.Add(1)
	if cfg.MoeCatCookies != "" {
		torbot.sendMessage(chatID, "MoeCat 已设置Cookies")
		loginOK.Done()
		return
	}
	for goon := 3; goon > 0; goon-- {
		loginHeader := req.Header{
			"Referer":        "https://www.moecat.best/login.php",
			"Sec-Fetch-Mode": "cors",
			// "DNT":              "1",
			"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"X-Requested-With": "XMLHttpRequest",
		}
		loginParam := req.Param{
			"username": username,
			"password": password,
		}
		_, err := r.Post("https://www.moecat.best/takelogin.php", loginParam, loginHeader)
		if err != nil {
			torbot.sendMessage(chatID, "登录 MoeCat 出现错误，请重试")
			log.Println(err)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		torbot.sendMessage(chatID, "MoeCat 登录成功！")
		goon = 0
		loginOK.Done()
		// resp.ToFile("222.html")
	}
}
func loginPTHome(username string, password string) {
	loginOK.Add(1)
	if cfg.PTHomeCookies != "" {
		torbot.sendMessage(chatID, "PTHome 已设置Cookies")
		chanGoPTer <- 1
		loginOK.Done()
		return
	}
	for goon := 3; goon > 0; goon-- {
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
			goon = 0
			chanGoPTer <- 1
			loginOK.Done()
		} else {
			torbot.sendMessage(chatID, "登录 PTHome 出现错误，请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
	}
}
func loginPTer(username string, password string) {
	loginOK.Add(1)
	if cfg.PTerCookies != "" {
		torbot.sendMessage(chatID, "PTer 已设置Cookies")
		loginOK.Done()
		chanGoOurBits <- 1
		return
	}
	<-chanGoPTer
	for goon := 3; goon > 0; goon-- {
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
			torbot.sendMessage(chatID, "Pter  请输入两步验证码，请预留10秒以上过期时间")
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
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "PTer 登录成功！")
			goon = 0
			loginOK.Done()
			chanGoOurBits <- 1
		} else {
			torbot.sendMessage(chatID, "登录 PTer 出现错误,请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
	}
}
func loginHDStreet(username string, password string) {
	loginOK.Add(1)
	if cfg.HDStreetCookies != "" {
		torbot.sendMessage(chatID, "HDStreet 已设置Cookies")
		loginOK.Done()
		return
	}
	for goon := 3; goon > 0; goon-- {
		r.Get("https://hdstreet.club/login.php")
		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Origin":     "https://hdstreet.club",
			"Referer":    "https://hdstreet.club/login.php",
		}
		loginParam := req.Param{
			"logintype": "username",
			"username":  username,
			"password":  password,
		}
		_, err := r.Post("https://hdstreet.club/takelogin.php", loginHeader, loginParam)
		if err != nil {
			log.Println(err)
			torbot.sendMessage(chatID, "登录 HDStreet 出现错误，请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		torbot.sendMessage(chatID, "HDStreet 登录成功!")
		goon = 0
		loginOK.Done()
	}
}
func loginCHDBits(username string, password string) {
	loginOK.Add(1)
	if cfg.CHDBitsCookies != "" {
		torbot.sendMessage(chatID, "CHDBits 已设置Cookies")
		loginOK.Done()
		return
	}
	for goon := 3; goon > 0; goon-- {
		r.Get("https://chdbits.co/login.php")
		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Referer":    "https://chdbits.co/login.php",
		}
		loginParam := req.Param{
			"username": username,
			"password": password,
		}
		resp, err := r.Post("https://chdbits.co/takelogin.php", loginHeader, loginParam)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "CHDBits 登录成功!")
			goon = 0
			loginOK.Done()
		} else {
			torbot.sendMessage(chatID, "登录 CHDBits 出现错误，请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
	}
}
func loginOurBits(username string, password string) {
	loginOK.Add(1)
	if cfg.OurBitsCookies != "" {
		torbot.sendMessage(chatID, "OurBits 已设置Cookies")
		loginOK.Done()
		chanGoHDSky <- 1
		return
	}
	<-chanGoOurBits
	for goon := 3; goon > 0; goon-- {
		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		r.Get("https://ourbits.club/login.php", loginHeader)
		verifycode := ""
		if cfg.OurBitsVerify {
			torbot.sendMessage(chatID, "OurBits  请输入两步验证码，请预留10秒以上过期时间")
			verifycode = <-chanimage
		}
		loginParam := req.Param{}
		if cfg.OurBitsSSL {
			loginParam = req.Param{
				"username":   username,
				"password":   password,
				"2fa_code":   verifycode,
				"trackerssl": "yes",
			}
		} else {
			loginParam = req.Param{
				"username": username,
				"password": password,
				"2fa_code": verifycode,
			}
		}
		resp, err := r.Post("https://ourbits.club/takelogin.php", loginHeader, loginParam)
		if err != nil {
			log.Println(err)
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "OurBits 登录成功！")
			goon = 0
			loginOK.Done()
			chanGoHDSky <- 1
		} else {
			torbot.sendMessage(chatID, "登录 OurBits 出现错误，请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
	}
}
func loginHDSky(username string, password string) {
	loginOK.Add(1)
	if cfg.HDSkyCookies != "" {
		torbot.sendMessage(chatID, "HDSky 已设置Cookies")
		loginOK.Done()
		chanGoSSD <- 1
		return
	}
	<-chanGoHDSky
	for goon := 3; goon > 0; goon-- {
		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		resp, err := r.Get("https://hdsky.me/login.php", loginHeader)
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
		torbot.sendPhoto(chatID, "https://hdsky.me/"+src, "#请输入验证码")
		imagehash := strings.Split(src, "imagehash=")[1]
		image := <-chanimage
		if cfg.PTerVerify {
			torbot.sendMessage(chatID, "HDSky  请输入两步验证码，请预留10秒以上过期时间")
			verifycode = <-chanimage
		}
		loginParam := req.Param{
			"username":    username,
			"password":    password,
			"imagestring": image,
			"imagehash":   imagehash,
			"oneCode":     verifycode,
		}
		resp, err = r.Post("https://hdsky.me/takelogin.php", loginHeader, loginParam)
		if err != nil {
			log.Println(err)
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "HDSky 登录成功！")
			goon = 0
			loginOK.Done()
			chanGoSSD <- 1
		} else {
			torbot.sendMessage(chatID, "登录 HDSky 出现错误,请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
	}
}
func loginSSD(username string, password string) {
	loginOK.Add(1)
	if cfg.SSDCookies != "" {
		torbot.sendMessage(chatID, "SSD 已设置Cookies")
		loginOK.Done()
		return
	}
	<-chanGoSSD
	for goon := 3; goon > 0; goon-- {
		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Referer":    "https://springsunday.net/login.php",
		}
		resp, err := r.Get("https://springsunday.net/login.php", loginHeader)
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
		src, ok := doc.Find(".verify-image>img").Attr("src")
		if !ok {
			log.Println("未找到验证图片")
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		torbot.sendPhoto(chatID, "https://springsunday.net/"+src, "#请输入验证码")
		imagehash := strings.Split(src, "imagehash=")[1]
		image := <-chanimage
		loginParam := req.Param{
			"username":    username,
			"password":    password,
			"imagestring": image,
			"imagehash":   imagehash,
			"keeplogin":   "yes",
			"returnto":    "index.php",
		}
		resp, err = r.Post("https://springsunday.net/takelogin.php", loginHeader, loginParam)
		if err != nil {
			log.Println(err)
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "SSD 登录成功！")
			goon = 0
			loginOK.Done()
		} else {
			torbot.sendMessage(chatID, "登录 SSD 出现错误，请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
	}
}
func loginFrds(username string, password string) {
	loginOK.Add(1)
	if cfg.FrdsCookies != "" {
		torbot.sendMessage(chatID, "Frds 已设置Cookies")
		loginOK.Done()
		return
	}
	// for goon := 1; goon > 0; goon-- {
	// 	loginHeader := req.Header{
	// 		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
	// 	}
	// 	r.Get("https://pt.keepfrds.com/login.php", loginHeader)
	// 	verifycode := ""
	// 	if cfg.FrdsVerify {
	// 		torbot.sendMessage(chatID, "OurBits  请输入两步验证码，请预留10秒以上过期时间")
	// 		verifycode = <-chanimage
	// 	}
	// 	loginParam := req.Param{
	// 		"username":   username,
	// 		"password":   password,
	// 		"2fa_code":   verifycode,
	// 		"trackerssl": "yes",
	// 	}
	// 	resp, err := r.Post("https://ourbits.club/takelogin.php", loginHeader, loginParam)
	// 	if err != nil {
	// 		log.Println(err)
	// 		torbot.sendMessage(chatID, "网络错误，请稍等5秒")
	// 		time.Sleep(time.Duration(5) * time.Second)
	// 		continue
	// 	}
	// 	html, _ := resp.ToString()
	// 	if strings.Contains(html, "欢迎回来") {
	// 		torbot.sendMessage(chatID, "OurBits 登录成功！")
	// 		goon = 0
	// 		loginOK.Done()
	// 		chanGoHDSky <- 1
	// 	} else {
	// 		torbot.sendMessage(chatID, "登录 OurBits 出现错误，请重试")
	// 		time.Sleep(time.Duration(5) * time.Second)
	// 		continue
	// 	}
	// }
}
func loginOpenCD(username string, password string) {
	loginOK.Add(1)
	if cfg.OpenCDCookies != "" {
		torbot.sendMessage(chatID, "OpenCD 已设置Cookies")
		loginOK.Done()
		return
	}
	<-chanGoOpenCD
	for goon := 1; goon > 0; goon-- {
		loginHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Origin":     "https://www.open.cd",
			"Referer":    "https://www.o/login.php",
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
			torbot.sendMessage(chatID, "Pter  请输入两步验证码，请预留10秒以上过期时间")
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
			torbot.sendMessage(chatID, "网络错误，请稍等5秒")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		html, _ := resp.ToString()
		if strings.Contains(html, "欢迎回来") {
			torbot.sendMessage(chatID, "PTer 登录成功！")
			goon = 0
			loginOK.Done()
			chanGoOurBits <- 1
		} else {
			torbot.sendMessage(chatID, "登录 PTer 出现错误,请重试")
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
	}
}
func loginMT(username string, password string) {
	loginOK.Add(1)
	if cfg.MTCookies != "" {
		torbot.sendMessage(chatID, "MT 已设置Cookies")
		loginOK.Done()
		return
	}
	<-chanGoMT
	for goon := 1; goon > 0; goon-- {
		// loginHeader := req.Header{
		// 	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
		// 	"Accept-Encoding":           "gzip, deflate, br",
		// 	"Accept-Language":           "zh-CN,zh;q=0.9,zh-TW;q=0.8",
		// 	"Cache-Control":             "max-age=0",
		// 	"Connection":                "keep-alive",
		// 	"Content-Length":            "44",
		// 	"Content-Type":              "application/x-www-form-urlencoded",
		// 	"Cookie":                    cfg.MTCookies,
		// 	"DNT":                       "1",
		// 	"Host":                      "pt.m-team.cc",
		// 	"Origin":                    "https://pt.m-team.cc",
		// 	"Referer":                   "https://pt.m-team.cc/login.php",
		// 	"Sec-Fetch-Mode":            "navigate",
		// 	"Sec-Fetch-Site":            "same-origin",
		// 	"Sec-Fetch-User":            "?1",
		// 	"Upgrade-Insecure-Requests": "1",
		// 	"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		// }
		// if strings.Contains(html,"請輸入驗證器上顯示的6位元驗證碼"){
		// https://pt.m-team.cc/verify.php?returnto=%2Ftorrents.php
		// }
	}
}

//MT2Verify is
func MT2Verify(text string) *goquery.Document {
	header := req.Header{}
	if cfg.MTCookies != "" {
		header = req.Header{
			// "Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
			"Accept-Encoding": "gzip, deflate, br",
			// "Accept-Language":           "zh-CN,zh;q=0.9,zh-TW;q=0.8",
			"Cache-Control":  "max-age=0",
			"Connection":     "keep-alive",
			"Content-Length": "10",
			"Content-Type":   "application/x-www-form-urlencoded",
			"Cookie":         cfg.MTCookies,
			// "DNT":                       "1",
			"Host":                      "pt.m-team.cc",
			"Origin":                    "https://pt.m-team.cc",
			"Referer":                   "https://pt.m-team.cc/verify.php?returnto=%2Ftorrents.php",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "same-origin",
			"Sec-Fetch-User":            "?1",
			"Upgrade-Insecure-Requests": "1",
			"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.75 Safari/537.36",
		}
	} else {
		header = req.Header{
			// "Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
			"Accept-Encoding": "gzip, deflate, br",
			// "Accept-Language":           "zh-CN,zh;q=0.9,zh-TW;q=0.8",
			"Cache-Control":  "max-age=0",
			"Connection":     "keep-alive",
			"Content-Length": "10",
			"Content-Type":   "application/x-www-form-urlencoded",
			// "DNT":                       "1",
			"Host":                      "pt.m-team.cc",
			"Origin":                    "https://pt.m-team.cc",
			"Referer":                   "https://pt.m-team.cc/verify.php?returnto=%2Ftorrents.php",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-Site":            "same-origin",
			"Sec-Fetch-User":            "?1",
			"Upgrade-Insecure-Requests": "1",
			"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.75 Safari/537.36",
		}
	}
	text = strings.TrimSpace(text)
	text = strings.Replace(text, "M-Team - TP :: 登錄", "", 1)
	torbot.sendMessage(chatID, text)
	otp := <-chanimage
	param := req.Param{
		"otp": otp,
	}
	resp, err := r.Post("https://pt.m-team.cc/verify.php?returnto=%2Ftorrents.php", header, param)
	if err != nil {
		log.Println(err)
		return nil
	}
	doc, err := goquery.NewDocumentFromResponse(resp.Response())
	return doc
}
