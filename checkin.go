package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
)

func checkinPTHome() {
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
func checkinPTer() {
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
func checkinMoeCat() {

}
