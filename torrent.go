package main

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//Torrent is
type Torrent struct {
	Title1 string
	Title2 string
	Status string
	Living string
	Size   string
	URL    string
	Ttype  string
	Update bool
	HasUp  bool
	Sales  string
	Site   string
	Four   string
}

//
func getTorMoeCat() []Torrent {
	goon := true
	var torrentList []Torrent
	for goon {
		resp, err := r.Get("https://www.moecat.best/torrents.php")
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
		}
		var torrent Torrent
		doc.Find("table.torrents>tbody>tr").Each(func(i int, s *goquery.Selection) {
			if i > 1 && i%2 == 0 {
				s.Find("tr>td").Each(func(i int, s *goquery.Selection) {
					if i == 0 {
						torrent.Ttype, _ = s.Find("a img").Attr("title")
						torrent.Site = "MoeCat"
					}
					if i == 1 {
						torrent.Sales, _ = s.Find("a+img").Attr("title")
						torrent.Sales += s.Find("a+img+b").Text()
						torrent.Sales = strings.Replace(torrent.Sales, "[", "", -1)
						torrent.Sales = strings.Replace(torrent.Sales, "]", "", -1)
						torrent.Sales = strings.Replace(torrent.Sales, "剩余：", "△", -1)
						if torrent.Sales == "" {
							torrent.Sales = "无优惠"
						}
					}
					if i == 2 {
						torrent.Title1, _ = s.Find("a").Attr("title")
						torrent.Title1 = strings.TrimSpace(torrent.Title1)
						s.Find("a").Remove()
						s.Find("b").Remove()
						torrent.Title2 = strings.TrimSpace(s.Text())
					}
					if i == 3 {
						torrent.URL, _ = s.Find("a").Attr("href")
						torrent.URL = "https://www.moecat.best/" + torrent.URL
					}
					if i == 4 {
						torrent.Four = "↑" + s.Text()
					}
					if i == 5 {
						torrent.Four += " ↓" + s.Text()
					}
					if i == 6 {
						torrent.Four += " ✓" + s.Text()
					}
					if i == 7 {
						torrent.Four += " ✗" + s.Text()
					}
				})
			}
			if i > 1 && i%2 == 1 {
				s.Find("tr>td").Each(func(i int, s *goquery.Selection) {
					if i == 1 {
						torrent.Living = s.Text()
					}
					if i == 2 {
						torrent.Size = s.Text()
						torrentList = append(torrentList, torrent)
					}
				})
			}
		})
		goon = false
	}
	return torrentList
}
func getPTHome() []Torrent {
	goon := true
	var torrentList []Torrent
	for goon {
		resp, err := r.Get("https://www.pthome.net/torrents.php")
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
		}
		var torrent Torrent
		doc.Find("table.torrents tbody tr").Each(func(i int, s *goquery.Selection) {
			if i != 0 {
				s.Find(".rowfollow").Each(func(i int, s *goquery.Selection) {
					if i == 0 {
						torrent.Ttype, _ = s.Find("a img").Attr("title")
						torrent.Site = "PTHome"
					}
					if i == 1 {
						torrent.Sales = ""
						s.Find("img").Each(func(i int, s *goquery.Selection) {
							hh, _ := s.Attr("alt")
							if hh != "Sticky" && hh != "download" && hh != "Unbookmarked" && hh != "Bookmarked" {
								if hh == "H&R" {
									hh = "H%26R "
								}
								torrent.Sales += hh
							}
						})
						torrent.Sales += "△"
						salestime := s.Find("span").Text()
						if torrent.Sales == "△" {
							torrent.Sales = "无优惠"
						}
						if strings.Contains(salestime, "[email") {
							torrent.Sales += strings.Split(salestime, "]")[1]
						} else {
							torrent.Sales += salestime
						}
						// torrent.Sales = strings.Replace(torrent.Sales, "[email"+`&#\d+;`+"protected]", "", -1)
						// log.Println(torrent.Sales)
						s.Find("td").Each(func(i int, s *goquery.Selection) {
							if i == 0 {
								torrent.Title1 = strings.TrimSpace(s.Find("a").Text())
								s.Find("a").Remove()
								s.Find("b").Remove()
								cc := s.Text()
								cc = strings.Replace(cc, "剩余时间：", "", -1)
								cc = strings.TrimSpace(cc)
								torrent.Title2 = cc
							}
							if i == 1 {
								dd, _ := s.Find("a").Attr("href")
								torrent.URL = "https://www.pthome.net/" + dd
							}
						})
					}
					if i == 3 {
						torrent.Living = s.Text()
					}
					if i == 4 {
						torrent.Size = s.Text()
					}
					if i == 5 {
						torrent.Four = "↑" + s.Text()
					}
					if i == 6 {
						torrent.Four += " ↓" + s.Text()
					}
					if i == 7 {
						torrent.Four += " ✓" + s.Text()
						torrentList = append(torrentList, torrent)
					}
				})
			}
		})
		goon = false
	}
	return torrentList
}
func getPTer() []Torrent {
	goon := true
	var torrentList []Torrent
	for goon {
		resp, err := r.Get("https://pter.club/torrents.php")
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
		}
		var torrent Torrent
		doc.Find("table.torrents tbody tr").Each(func(i int, s *goquery.Selection) {
			if i != 0 {
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					if i == 0 {
						s.Find("img").Each(func(i int, s *goquery.Selection) {
							if i == 0 {
								torrent.Ttype, _ = s.Attr("title")
								if torrent.Ttype == "音乐短片 (MV)" {
									torrent.Ttype = "MV"
								}
								if torrent.Ttype == "电影 (Movie)" {
									torrent.Ttype = "电影"
								}
								if torrent.Ttype == "电视剧 (TV Play)" {
									torrent.Ttype = "电视剧"
								}
								if torrent.Ttype == "动漫 (Anime)" {
									torrent.Ttype = "动漫"
								}
								if torrent.Ttype == "音乐 (Music)" {
									torrent.Ttype = "音乐"
								}
								if torrent.Ttype == "学习 (Study)" {
									torrent.Ttype = "学习"
								}
								if torrent.Ttype == "电子书 (Ebook)" {
									torrent.Ttype = "电子书"
								}
								if torrent.Ttype == "软件 (Software)" {
									torrent.Ttype = "软件"
								}
								if torrent.Ttype == "综艺 (TV Show)" {
									torrent.Ttype = "综艺"
								}
								if torrent.Ttype == "纪录片 (Documentary)" {
									torrent.Ttype = "纪录片"
								}
								if torrent.Ttype == "游戏 (Game)" {
									torrent.Ttype = "游戏"
								}
								if torrent.Ttype == "体育 (Sport)" {
									torrent.Ttype = "体育"
								}
								if torrent.Ttype == "其它 (Other)" {
									torrent.Ttype = "其它"
								}
							}
						})
						torrent.Site = "PTer"
					}
					if i == 2 {
						torrent.Title1 = s.Find("a").Eq(0).Text()
						torrent.Title1 = strings.Replace(torrent.Title1, "[email protected]", "", -1)
						torrent.Sales, _ = s.Find("a").Eq(1).Find("img").Attr("alt")
						torrent.Sales += "△"
						salestime := s.Find("span").Text()
						if torrent.Sales == "△" {
							torrent.Sales = "无优惠"
						}
						if strings.Contains(salestime, "[email") {
							torrent.Sales += strings.Split(salestime, "]")[1]
						} else {
							torrent.Sales += salestime
						}
						s.Find("a").Eq(0).Remove()
						s.Find("span").Remove()
						cc := s.Text()
						cc = strings.Replace(cc, "剩余时间：", "", -1)
						cc = strings.TrimSpace(cc)
						torrent.Title2 = cc
					}
					if i == 3 {
						s.Find("a").Each(func(i int, s *goquery.Selection) {
							alt, ok := s.Find("img").Attr("alt")
							if ok {
								if alt == "download" {
									href, _ := s.Attr("href")
									torrent.URL = "https://pter.club/" + href
								}
							}
						})
					}
					if i == 5 {
						torrent.Living = s.Text()
					}
					if i == 6 {
						torrent.Size = s.Text()
					}
					if i == 7 {
						torrent.Four = "↑" + s.Text()
					}
					if i == 8 {
						torrent.Four += " ↓" + s.Text()
					}
					if i == 9 {
						torrent.Four += " ✓" + s.Text()
						torrentList = append(torrentList, torrent)
					}
				})
			}
		})
		goon = false
	}
	return torrentList
}
