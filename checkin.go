package main

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
)

func checkinPTHome() {
	if cfg.PTHomeUsername != "" {
		checkinHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Referer":    "https://www.pthome.net/torrents.php",
		}
		resp, err := r.Get("https://www.pthome.net/attendance.php", checkinHeader)
		if err != nil {
			log.Println(err)
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
		}
		message := doc.Find("td.text").Text()
		torbot.sendMessage(chatID, "PTHome 签到："+message)
	}
}
func checkinPTer() {
	if cfg.PTerUsername != "" {
		checkinHeader := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Referer":    "https://pter.club/torrents.php",
		}
		resp, err := r.Get("https://pter.club/attendance-ajax.php", checkinHeader)
		if err != nil {
			log.Println(err)
		}
		hh, _ := resp.ToString()
		hh = strings.Split(hh, "message\":\"")[1]
		hh = strings.Replace(hh, "<p>", "", -1)
		hh = strings.Replace(hh, "</p>", "", -1)
		hh = strings.Replace(hh, "<b>", "", -1)
		hh = strings.Replace(hh, "</b>", "", -1)
		torbot.sendMessage(chatID, "PTer 签到："+hh)
	}
}
func checkinMoeCat() {
	if cfg.MoeCatUsername != "" {
		questionid := ""
		choice := ""
		param := req.Param{
			"questionid":  questionid,
			"choice[]":    choice,
			"usercomment": "此刻心情:无",
			"submit":      "提交",
		}
		resp, err := r.Post("https://www.moecat.best/bakatest.php", param)
		if err != nil {
			log.Println(err)
		}
		log.Println(resp.ToString())
	}
}
func checkinAtNine() {
	for {
		now := time.Now()
		next := now
		if now.Hour() > 9 {
			next = next.Add(time.Hour * 24)
		}
		next = time.Date(next.Year(), next.Month(), next.Day(), 9, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		go checkinPTHome()
		go checkinPTer()
		go checkinMoeCat()
	}
}
