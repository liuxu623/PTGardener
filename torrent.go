package main

import (
	"log"
	"strings"
	"time"

	"github.com/imroc/req"

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
func getMoeCat() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"cookie":     cfg.MoeCatCookies,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://www.moecat.best/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
		goon = 0
	}
	return torrentList
}
func getPTHome() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"cookie":     cfg.PTHomeCookies,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://www.pthome.net/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
								torrent.Title1 = strings.Replace(torrent.Title1, "[email protected]", "", -1)
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
		goon = 0
	}
	return torrentList
}
func getPTer() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"cookie":     cfg.PTerCookies,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://pter.club/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
			continue
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
					if i == 5 {
						s.Find("a").Each(func(i int, s *goquery.Selection) {
							alt, ok := s.Find("img").Attr("alt")
							if ok {
								if alt == "download" {
									href, _ := s.Attr("href")
									href = strings.Split(href, "&")[0]
									torrent.URL = "https://pter.club/" + href
								}
							}
						})
					}
					if i == 9 {
						torrent.Living = s.Text()
					}
					if i == 10 {
						torrent.Size = s.Text()
					}
					if i == 11 {
						torrent.Four = "↑" + s.Text()
					}
					if i == 12 {
						torrent.Four += " ↓" + s.Text()
					}
					if i == 13 {
						torrent.Four += " ✓" + s.Text()
						torrentList = append(torrentList, torrent)
					}
				})
			}
		})
		goon = 0
	}
	return torrentList
}
func getHDStreet() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		headers := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Referer":    "https://hdstreet.club/index.php",
			"cookie":     cfg.HDStreetCookies,
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://hdstreet.club/torrents.php", headers)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
						torrent.Site = "HDStreet"
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
						torrent.URL = "https://hdstreet.club/" + torrent.URL
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
		goon = 0
	}
	for goon := 3; goon > 0; goon-- {
		headers := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Referer":    "https://hdstreet.club/index.php",
			"cookie":     cfg.HDStreetCookies,
		}
		resp, err := r.Get("https://hdstreet.club/torrentsasia.php", headers)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
						torrent.Site = "HDStreet"
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
						torrent.URL = "https://hdstreet.club/" + torrent.URL
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
		goon = 0
	}
	// for _, v := range torrentList {
	// log.Println(v.URL)
	// }
	return torrentList
}
func getCHDBits() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"cookie":     cfg.CHDBitsCookies,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://chdbits.co/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
						torrent.Site = "CHDBits"
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
								// s.Find("a").Remove()
								// s.Find("b").Remove()
								// cc := s.Text()
								// cc = strings.Replace(cc, "剩余时间：", "", -1)
								// cc = strings.TrimSpace(cc)
								torrent.Title2 = s.Find("font").Text()
							}
							if i == 1 {
								dd, _ := s.Find("a").Attr("href")
								torrent.URL = "https://chdbits.co/" + dd
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
		goon = 0
	}
	return torrentList
}
func getOurBits() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"cookie":     cfg.OurBitsCookies,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://ourbits.club/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
							}
						})
						torrent.Site = "OurBits"
					}
					if i == 2 {
						torrent.Title1 = s.Find("a").Eq(0).Text()
						torrent.Title1 = strings.Replace(torrent.Title1, "[email protected]", "", -1)
						s.Find("img").Each(func(i int, s *goquery.Selection) {
							alt, _ := s.Attr("alt")
							if alt != "Sticky" {
								torrent.Sales = alt
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
						s.Find("a").Eq(0).Remove()
						s.Find("span").Remove()
						cc := s.Text()
						cc = strings.Replace(cc, "剩余时间：", "", -1)
						cc = strings.TrimSpace(cc)
						torrent.Title2 = cc
					}
					if i == 6 {
						s.Find("a").Each(func(i int, s *goquery.Selection) {
							alt, ok := s.Find("img").Attr("alt")
							if ok {
								if alt == "download" {
									href, _ := s.Attr("href")
									torrent.URL = "https://ourbits.club/" + href
								}
							}
						})
					}
					if i == 8 {
						torrent.Living = s.Text()
					}
					if i == 9 {
						torrent.Size = s.Text()
					}
					if i == 10 {
						torrent.Four = "↑" + s.Text()
					}
					if i == 11 {
						torrent.Four += " ↓" + s.Text()
					}
					if i == 12 {
						torrent.Four += " ✓" + s.Text()
						torrentList = append(torrentList, torrent)
					}
				})
			}
		})
		goon = 0
	}
	return torrentList
}
func getHDSkey() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"cookie":     cfg.HDSkyCookies,
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://hdsky.me/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
						torrent.Ttype = strings.Split(torrent.Ttype, "/")[1]
						torrent.Site = "HDSky"
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
								cc = strings.Replace(cc, "[优惠剩余时间：]", "", -1)
								cc = strings.TrimSpace(cc)
								torrent.Title2 = cc
							}
							if i == 1 {
								dd, _ := s.Find("form").Eq(0).Attr("action")
								torrent.URL = "https://hdsky.me/" + dd
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
		goon = 0
	}
	// if len(torrentList)==0{
	// 	torbot.sendMessage(chatID,"")
	// }
	return torrentList
}
func getSSD() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"cookie":     cfg.SSDCookies,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://springsunday.net/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
								// torrent.Ttype = strings.Split(torrent.Ttype, "(")[1]
								// torrent.Ttype = strings.Replace(torrent.Ttype, ")", "", -1)
							}
						})
						torrent.Site = "SSD"
					}
					if i == 2 {
						torrent.Title1 = s.Find("a").Eq(0).Text()
						torrent.Title1 = strings.Replace(torrent.Title1, "[email protected]", "", -1)
						s.Find("img").Each(func(i int, s *goquery.Selection) {
							alt, _ := s.Attr("alt")
							if alt != "Super" && alt != "Sticky" {
								torrent.Sales = alt
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
						s.Find("a").Eq(0).Remove()
						s.Find("span").Remove()
						s.Find("b").Remove()
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
									torrent.URL = "https://springsunday.net/" + href
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
		goon = 0
	}
	return torrentList
}
func getFrds() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"cookie":     cfg.FrdsCookies,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://pt.keepfrds.com/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
							}
						})
						torrent.Site = "Frds"
					}
					if i == 2 {
						torrent.Title1 = s.Find("a").Eq(0).Text()
						torrent.Title1 = strings.Replace(torrent.Title1, "[email protected]", "", -1)
						s.Find("a").Eq(0).Remove()
						s.Find("span").Remove()
						s.Find("b").Remove()
						cc := s.Text()
						cc = strings.Replace(cc, "剩余时间：", "", -1)
						cc = strings.TrimSpace(cc)
						torrent.Title2 = cc
					}
					if i == 3 {
						torrent.URL, _ = s.Find("div").Eq(0).Find("a").Attr("href")
						torrent.URL = "https://pt.keepfrds.com/" + torrent.URL
						s.Find("div").Eq(1).Find("img").Each(func(i int, s *goquery.Selection) {
							alt, ok := s.Attr("title")
							if ok {
								if alt != "收藏" {
									torrent.Sales = alt
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
		goon = 0
	}
	return torrentList
}
func getMT() []Torrent {
	var torrentList []Torrent
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"Host":                      "pt.m-team.cc",
			"Connection":                "keep-alive",
			"Cache-Control":             "max-age=0",
			"Upgrade-Insecure-Requests": "1",
			"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-User":            "?1",
			// "DNT":                       "1",
			// "Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
			"Sec-Fetch-Site":  "none",
			"Referer":         "https://pt.m-team.cc/index.php",
			"Accept-Encoding": "gzip, deflate, br",
			// "Accept-Language":           "zh-CN,zh;q=0.9,zh-TW;q=0.8",
			"cookie": cfg.MTCookies,
		}
		req.SetTimeout(60 * time.Second)
		resp, err := r.Get("https://pt.m-team.cc/torrents.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp.Response())
		if err != nil {
			log.Println(err)
			continue
		}
		for vvv := 3; vvv > 0; vvv-- {
			if strings.Contains(doc.Text(), "請輸入驗證器上顯示的6位元驗證碼") {
				doc = MT2Verify(doc.Text())
				if goon == 1 {
					return torrentList
				}
				continue
			}

		}
		var torrent Torrent
		doc.Find("table.torrents tbody tr").Each(func(i int, s *goquery.Selection) {
			if i != 0 {
				s.Find("td").Each(func(i int, s *goquery.Selection) {
					if i == 0 {
						s.Find("img").Each(func(i int, s *goquery.Selection) {
							if i == 0 {
								torrent.Ttype, _ = s.Attr("title")
								torrent.Ttype = strings.Split(torrent.Ttype, "(")[0]
							}
						})
						torrent.Site = "MTeam"
					}
					if i == 3 {
						torrent.Title1 = s.Find("a").Eq(0).Text()
						torrent.Title1 = strings.Replace(torrent.Title1, "[email protected]", "", -1)
						torrent.Sales, _ = s.Find("img[class^=pro_]").Attr("alt")
						torrent.Sales += "△"
						salestime := s.Find("span").Text()
						if torrent.Sales == "△" {
							torrent.Sales = "无优惠"
						}
						salestime = strings.Replace(salestime, "限時：", "", 1)
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
					if i == 4 {
						s.Find("a").Each(func(i int, s *goquery.Selection) {
							alt, ok := s.Find("img").Attr("alt")
							if ok {
								if alt == "download" {
									href, _ := s.Attr("href")
									href = strings.Split(href, "&")[0]
									torrent.URL = "https://pt.m-team.cc/" + href
								}
							}
						})
					}
					if i == 6 {
						torrent.Living = s.Text()
					}
					if i == 7 {
						torrent.Size = s.Text()
					}
					if i == 8 {
						torrent.Four = "↑" + s.Text()
					}
					if i == 9 {
						torrent.Four += " ↓" + s.Text()
					}
					if i == 10 {
						torrent.Four += " ✓" + s.Text()
						torrent.Four = strings.Replace(torrent.Four, ",", "%2C", -1)
						torrentList = append(torrentList, torrent)
					}
				})
			}
		})
		goon = 0
	}
	for goon := 3; goon > 0; goon-- {
		header := req.Header{
			"Host":                      "pt.m-team.cc",
			"Connection":                "keep-alive",
			"Cache-Control":             "max-age=0",
			"Upgrade-Insecure-Requests": "1",
			"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36",
			"Sec-Fetch-Mode":            "navigate",
			"Sec-Fetch-User":            "?1",
			// "DNT":                       "1",
			// "Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
			"Sec-Fetch-Site":  "none",
			"Referer":         "https://pt.m-team.cc/index.php",
			"Accept-Encoding": "gzip, deflate, br",
			// "Accept-Language":           "zh-CN,zh;q=0.9,zh-TW;q=0.8",
			"cookie": cfg.MTCookies,
		}
		resp, err := r.Get("https://pt.m-team.cc/adult.php", header)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(10) * time.Second)
			if goon == 1 {
				return torrentList
			}
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
								torrent.Ttype = strings.Split(torrent.Ttype, "(")[0]
							}
						})
						torrent.Site = "MTeam"
					}
					if i == 3 {
						torrent.Title1 = s.Find("a").Eq(0).Text()
						torrent.Title1 = strings.Replace(torrent.Title1, "[email protected]", "", -1)
						torrent.Sales, _ = s.Find("img[class^=pro_]").Attr("alt")
						torrent.Sales += "△"
						salestime := s.Find("span").Text()
						if torrent.Sales == "△" {
							torrent.Sales = "无优惠"
						}
						salestime = strings.Replace(salestime, "限時：", "", 1)
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
					if i == 4 {
						s.Find("a").Each(func(i int, s *goquery.Selection) {
							alt, ok := s.Find("img").Attr("alt")
							if ok {
								if alt == "download" {
									href, _ := s.Attr("href")
									href = strings.Split(href, "&")[0]
									torrent.URL = "https://pt.m-team.cc/" + href
								}
							}
						})
					}
					if i == 6 {
						torrent.Living = s.Text()
					}
					if i == 7 {
						torrent.Size = s.Text()
					}
					if i == 8 {
						torrent.Four = "↑" + s.Text()
					}
					if i == 9 {
						torrent.Four += " ↓" + s.Text()
					}
					if i == 10 {
						torrent.Four += " ✓" + s.Text()
						torrent.Four = strings.Replace(torrent.Four, ",", "%2C", -1)
						torrentList = append(torrentList, torrent)
					}
				})
			}
		})
		goon = 0
	}
	return torrentList
}
