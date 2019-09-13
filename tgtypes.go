package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/technoweenie/multipartstreamer"
)

//BotAPI is method of bot
type BotAPI struct {
	API    string
	Update chan *Update
}

//参数的常量标签
const (
	Mute          = "disable_notification"     //回复不提醒
	NoPreView     = "disable_web_page_preview" //发送的链接不预览
	Markdown      = "parse_Markdown"           //以Markdown解析内容
	HTML          = "parse_HTML"               //以HTML解析内容
	RszKB         = "resize_keyboard"          //键盘自适应缩小
	OneKB         = "one_time_keyboard"        //键盘按下消失
	SelKB         = "selective_keyboard"       //键盘特定显示
	RemoveKM      = "remove_keyboard"          //移除特定用户键盘
	RemoveAllKM   = "remove_all_keyboard"      //移除所有用户键盘
	ForceReply    = "force_reply"              //发送消息并让对方直接回复
	AllForceReply = "all_force_reply"          //发送消息并让所有人直接回复
)

//REPLY 是回复的ID
func REPLY(id int) string {
	return "REPLY" + strconv.Itoa(id)
}

//创建一个按钮
func makeBtn(text ...string) string {
	kb := `{"text":"` + text[0] + `"`
	for _, v := range text {
		if v == "ReqC" {
			kb += `,"request_contact":true`
		}
		if v == "ReqL" {
			kb += `,"request_location":true`
		}
		if strings.HasPrefix(v, "URL") {
			kb += `,"url":"` + strings.TrimPrefix(v, "URL") + `"`
		}
		if strings.HasPrefix(v, "CBD") {
			kb += `,"callback_data":"` + strings.TrimPrefix(v, "CBD") + `"`
		}
		if strings.HasPrefix(v, "SIQ") {
			kb += `,"switch_inline_query":"` + strings.TrimPrefix(v, "SIQ") + `"`
		}
		if strings.HasPrefix(v, "SIC") {
			kb += `,"switch_inline_query_current_chat":"` + strings.TrimPrefix(v, "SIC") + `"`
		}
	}
	kb += `}`
	return kb
}

//ReplyKM 创造回复键盘
func ReplyKM(str string, option ...string) string {
	strArry := strings.Split(str, ",")
	var kb string
	kb += `[`
	for i, v := range strArry {
		lineArry := strings.Split(v, "|")
		var line string
		line += `[`
		for i, v := range lineArry {
			kv := strings.Split(v, ":")
			line += makeBtn(kv...)
			if i+1 < len(lineArry) {
				line += `,`
			}
		}
		line += `]`
		kb += line
		if i+1 < len(strArry) {
			kb += `,`
		}
	}
	kb += `]`
	for _, v := range option {
		if v == "resize_keyboard" {
			kb += `,"resize_keyboard":true`
		}
		if v == "one_time_keyboard" {
			kb += `,"one_time_keyboard":true`
		}
		if v == "selective_keyboard" {
			kb += `,"selective":true`
		}
	}
	return "is_a_reply_keyboard_markup" + kb
}

//InlineKM 创造inline键盘
func InlineKM(str string) string {
	strArry := strings.Split(str, ",")
	var kb string
	kb += `[`
	for i, v := range strArry {
		lineArry := strings.Split(v, "|")
		var line string
		line += `[`
		for i, v := range lineArry {
			kv := strings.Split(v, "@")
			line += makeBtn(kv...)
			if i+1 < len(lineArry) {
				line += `,`
			}
		}
		line += `]`
		kb += line
		if i+1 < len(strArry) {
			kb += `,`
		}
	}
	kb += `]`
	return "is_a_inline_keyboard_markup" + kb
}

func parseArg(arg []string) string {
	var queryStr string
	for _, v := range arg {
		if v == "disable_notification" {
			queryStr += `&disable_notification=true`
		}
		if strings.Contains(v, "REPLY") {
			queryStr += `&reply_to_message_id=` + strings.TrimPrefix(v, "REPLY")
		}
		if v == "disable_web_page_preview" {
			queryStr += `&disable_web_page_preview=true`
		}
		if v == "parse_Markdown" {
			queryStr += `&parse_mode=Markdown`
		}
		if v == "parse_HTML" {
			queryStr += `&parse_mode=HTML`
		}
		if strings.HasPrefix(v, "is_a_reply_keyboard_markup") {
			queryStr += `&reply_markup={"keyboard":` + strings.TrimPrefix(v, "is_a_reply_keyboard_markup") + `}`
		}
		if v == "remove_keyboard" {
			queryStr += `&reply_markup={"remove_keyboard":true,"selective":true}`
		}
		if v == "remove_all_keyboard" {
			queryStr += `&reply_markup={"remove_keyboard":true,"selective":false}`
		}
		if v == "force_reply" {
			queryStr += `&reply_markup={"force_reply":true,"selective":true}`
		}
		if v == "all_force_reply" {
			queryStr += `&reply_markup={"force_reply":true,"selective":false}`
		}
		if strings.HasPrefix(v, "is_a_inline_keyboard_markup") {
			queryStr += `&reply_markup={"inline_keyboard":` + strings.TrimPrefix(v, "is_a_inline_keyboard_markup") + `}`
		}
		if strings.HasPrefix(v, "#") {
			queryStr += `&caption=` + url.QueryEscape(strings.TrimPrefix(v, "#"))
		}
	}
	return queryStr
}

//sendMessage
//HTML|Markdown   parse_mode
//Mute    NoPreview    REPLY(id) ReplyKN()
func (bot *BotAPI) sendMessage(chatID int64, text string, arg ...string) {
	for goon:=3; goon > 0; goon-- {

		msg := fmt.Sprintf(`?chat_id=%d&text=%s`, chatID, url.QueryEscape(text)) + parseArg(arg)
		resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/sendMessage" + msg)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(30) * time.Second)
			continue
		}
		hh, err := resp.ToString()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Duration(1) * time.Minute)
			continue
		}
		if strings.Contains(hh, "false") {
			log.Println(resp.ToString())
			log.Println(msg)
		}
		goon = 0
	}
}
func (bot *BotAPI) deleteMessage(chatID int64, msgID int) {
	msg := fmt.Sprintf(`?chat_id=%d&message_id=%d`, chatID, msgID)
	resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/deleteMessage" + msg)
	html, err := resp.ToString()
	if err != nil {
		log.Println(err)
	}
	log.Println(html)
}
func (bot *BotAPI) editMessageText(chatID int64, msgID int, text string, arg ...string) {
	msg := fmt.Sprintf(`?text=%s&chat_id=%d&message_id=%d`, url.QueryEscape(text), chatID, msgID) + parseArg(arg)
	rq.Post("https://api.telegram.org/bot" + bot.API + "/editMessageText" + msg)
}
func (bot *BotAPI) editMessageCaption(chatID int64, msgID int, arg ...string) {
	msg := fmt.Sprintf(`?chat_id=%d&message_id=%d`, chatID, msgID) + parseArg(arg)
	rq.Post("https://api.telegram.org/bot" + bot.API + "/editMessageText" + msg)
}
func (bot *BotAPI) editMessageReplyMarkup(chatID int64, msgID int, arg ...string) {
	msg := fmt.Sprintf(`?chat_id=%d&message_id=%d`, chatID, msgID) + parseArg(arg)
	rq.Post("https://api.telegram.org/bot" + bot.API + "/editMessageReplyMarkup" + msg)
}
func (bot *BotAPI) editMessagePhoto(chatID int64, messageID int, media string, caption string, arg ...string) {
	if !strings.Contains(media, ".") || strings.HasPrefix(media, "http://") || strings.HasPrefix(media, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&message_id=%d`, chatID, messageID) + parseArg(arg)
		jsonstr := fmt.Sprintf(`{"type":"photo","media":"%s"`, media)
		if caption != "" {
			jsonstr += `,"caption":"` + caption + `"`
		}
		jsonstr += `}`
		msg += `&media=` + url.QueryEscape(jsonstr)
		resp, _ := rq.Post("https://api.telegram.org/bot" + bot.API + "/editMessageMedia" + msg)
		str, _ := resp.ToString()
		fmt.Println(str)
		return
	}
	// id := fmt.Sprintf("%d", chatID)
	msg := fmt.Sprintf(`?chat_id=%d&message_id=%d`, chatID, messageID) + parseArg(arg)
	jsonstr := fmt.Sprintf(`{"type":"photo","media":"attach://photo"`)
	if caption != "" {
		jsonstr += `,"caption":"` + caption + `"`
	}
	jsonstr += `}`
	msg += `&media=` + url.QueryEscape(jsonstr)
	ms := multipartstreamer.New()
	fileHandle, err := os.Open(media)
	errPrint(err, e.Tag(1))
	defer fileHandle.Close()
	fi, err := os.Stat(media)
	errPrint(err, e.Tag(2))
	ms.WriteReader("photo", fileHandle.Name(), fi.Size(), fileHandle)
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/editMessageMedia"+msg, nil)
	req.Close = true
	errPrint(err, e.Tag(3))
	ms.SetupRequest(req)
	http.DefaultClient.Do(req)
}
func (bot *BotAPI) forwardMessage(chatID int64, fromChatID int64, messageID int, arg ...string) {
	msg := fmt.Sprintf(`?chat_id=%d&from_chat_id=%d&message_id=%d`, chatID, fromChatID, messageID) + parseArg(arg)
	rq.Post("https://api.telegram.org/bot" + bot.API + "/forwardMessage" + msg)
}

func (bot *BotAPI) sendPhoto(chatID int64, input string, arg ...string) {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&photo=%s`, chatID, url.QueryEscape(input)) + parseArg(arg)
		resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/sendPhoto" + msg)
		if err != nil {
			fmt.Println(err)
		}
		html, _ := resp.ToString()
		fmt.Println(html)
		return
	}
	msg := fmt.Sprintf(`?chat_id=%d`, chatID) + parseArg(arg)
	ms := multipartstreamer.New()
	fileHandle, err := os.Open(input)
	errPrint(err, e.Tag(1))
	defer fileHandle.Close()
	fi, err := os.Stat(input)
	errPrint(err, e.Tag(2))
	ms.WriteReader("photo", fileHandle.Name(), fi.Size(), fileHandle)
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/sendPhoto"+msg, nil)
	req.Close = true
	errPrint(err, e.Tag(3))
	ms.SetupRequest(req)
	http.DefaultClient.Do(req)
}
func (bot *BotAPI) sendAudio(chatID int64, input string, arg ...string) {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&audio=%s`, chatID, input) + parseArg(arg)
		rq.Post("https://api.telegram.org/bot" + bot.API + "/sendAudio" + msg)
	} else {
		// id := fmt.Sprintf("%d", chatID)
		msg := fmt.Sprintf(`?chat_id=%d`, chatID) + parseArg(arg)
		ms := multipartstreamer.New()
		// params := make(map[string]string)
		// params["chat_id"] = id
		// ms.WriteFields(params)
		fileHandle, err := os.Open(input)
		errPrint(err, e.Tag(1))
		defer fileHandle.Close()
		fi, err := os.Stat(input)
		errPrint(err, e.Tag(2))
		ms.WriteReader("audio", fileHandle.Name(), fi.Size(), fileHandle)
		req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/sendAudio"+msg, nil)
		req.Close = true
		errPrint(err, e.Tag(3))
		ms.SetupRequest(req)
		resp, err := http.DefaultClient.Do(req)
		errPrint(err, e.Tag(5))
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		errPrint(err, e.Tag(0))
		fmt.Println(string(bytes))
	}
}
func (bot *BotAPI) sendDocument(chatID int64, input string, arg ...string) {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&document=%s`, chatID, input) + parseArg(arg)
		resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/sendDocument" + msg)
		errPrint(err, e.Tag(0))
		status, err := resp.ToString()
		log.Println(status)
	} else {
		// id := fmt.Sprintf("%d", chatID)
		msg := fmt.Sprintf(`?chat_id=%d`, chatID) + parseArg(arg)
		ms := multipartstreamer.New()
		// params := make(map[string]string)
		// params["chat_id"] = id
		// ms.WriteFields(params)
		fileHandle, err := os.Open(input)
		errPrint(err, e.Tag(1))
		defer fileHandle.Close()
		fi, err := os.Stat(input)
		errPrint(err, e.Tag(2))
		ms.WriteReader("document", fileHandle.Name(), fi.Size(), fileHandle)
		req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/sendDocument"+msg, nil)
		req.Close = true
		errPrint(err, e.Tag(3))
		ms.SetupRequest(req)
		resp, err := http.DefaultClient.Do(req)
		errPrint(err, e.Tag(5))
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		errPrint(err, e.Tag(0))
		fmt.Println(string(bytes))
	}
}
func (bot *BotAPI) sendVideo(chatID int64, input string, arg ...string) {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&video=%s`, chatID, input) + parseArg(arg)
		resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/sendVideo" + msg)
		errPrint(err, e.Tag(0))
		status, err := resp.ToString()
		log.Println(status)
	} else {
		// id := fmt.Sprintf("%d", chatID)
		msg := fmt.Sprintf(`?chat_id=%d`, chatID) + parseArg(arg)
		ms := multipartstreamer.New()
		// params := make(map[string]string)
		// params["chat_id"] = id
		// ms.WriteFields(params)
		fileHandle, err := os.Open(input)
		errPrint(err, e.Tag(1))
		defer fileHandle.Close()
		fi, err := os.Stat(input)
		errPrint(err, e.Tag(2))
		ms.WriteReader("video", fileHandle.Name(), fi.Size(), fileHandle)
		req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/sendVideo"+msg, nil)
		req.Close = true
		errPrint(err, e.Tag(3))
		ms.SetupRequest(req)
		resp, err := http.DefaultClient.Do(req)
		errPrint(err, e.Tag(5))
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		errPrint(err, e.Tag(0))
		fmt.Println(string(bytes))
	}
}
func (bot *BotAPI) sendAnimation(chatID int64, input string, arg ...string) {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&animation=%s`, chatID, input) + parseArg(arg)
		resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/sendAnimation" + msg)
		errPrint(err, e.Tag(0))
		status, err := resp.ToString()
		log.Println(status)
	} else {
		// id := fmt.Sprintf("%d", chatID)
		msg := fmt.Sprintf(`?chat_id=%d`, chatID) + parseArg(arg)
		ms := multipartstreamer.New()
		// params := make(map[string]string)
		// params["chat_id"] = id
		// ms.WriteFields(params)
		fileHandle, err := os.Open(input)
		errPrint(err, e.Tag(1))
		defer fileHandle.Close()
		fi, err := os.Stat(input)
		errPrint(err, e.Tag(2))
		ms.WriteReader("animation", fileHandle.Name(), fi.Size(), fileHandle)
		req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/sendAnimation"+msg, nil)
		req.Close = true
		errPrint(err, e.Tag(3))
		ms.SetupRequest(req)
		resp, err := http.DefaultClient.Do(req)
		errPrint(err, e.Tag(5))
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		errPrint(err, e.Tag(0))
		fmt.Println(string(bytes))
	}
}
func (bot *BotAPI) sendVoice(chatID int64, input string, arg ...string) {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&voice=%s`, chatID, input) + parseArg(arg)
		resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/sendVoice" + msg)
		errPrint(err, e.Tag(0))
		status, err := resp.ToString()
		log.Println(status)
	} else {
		// id := fmt.Sprintf("%d", chatID)
		msg := fmt.Sprintf(`?chat_id=%d`, chatID) + parseArg(arg)
		ms := multipartstreamer.New()
		// params := make(map[string]string)
		// params["chat_id"] = id
		// ms.WriteFields(params)
		fileHandle, err := os.Open(input)
		errPrint(err, e.Tag(1))
		defer fileHandle.Close()
		fi, err := os.Stat(input)
		errPrint(err, e.Tag(2))
		ms.WriteReader("voice", fileHandle.Name(), fi.Size(), fileHandle)
		req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/sendVoice"+msg, nil)
		req.Close = true
		errPrint(err, e.Tag(3))
		ms.SetupRequest(req)
		resp, err := http.DefaultClient.Do(req)
		errPrint(err, e.Tag(5))
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		errPrint(err, e.Tag(0))
		fmt.Println(string(bytes))
	}
}
func (bot *BotAPI) sendVideoNote(chatID int64, input string, arg ...string) {
	if !strings.Contains(input, ".") || strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		msg := fmt.Sprintf(`?chat_id=%d&video_note=%s`, chatID, input) + parseArg(arg)
		resp, err := rq.Post("https://api.telegram.org/bot" + bot.API + "/sendVideoNote" + msg)
		errPrint(err, e.Tag(0))
		status, err := resp.ToString()
		log.Println(status)
	} else {
		// id := fmt.Sprintf("%d", chatID)
		msg := fmt.Sprintf(`?chat_id=%d`, chatID) + parseArg(arg)
		ms := multipartstreamer.New()
		// params := make(map[string]string)
		// params["chat_id"] = id
		// ms.WriteFields(params)
		fileHandle, err := os.Open(input)
		errPrint(err, e.Tag(1))
		defer fileHandle.Close()
		fi, err := os.Stat(input)
		errPrint(err, e.Tag(2))
		ms.WriteReader("video_note", fileHandle.Name(), fi.Size(), fileHandle)
		req, err := http.NewRequest("POST", "https://api.telegram.org/bot"+bot.API+"/sendVideoNote"+msg, nil)
		req.Close = true
		errPrint(err, e.Tag(3))
		ms.SetupRequest(req)
		resp, err := http.DefaultClient.Do(req)
		errPrint(err, e.Tag(5))
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		errPrint(err, e.Tag(0))
		fmt.Println(string(bytes))
	}
}

// Update is an update response, from GetUpdates.
type Update struct {
	UpdateID           int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query"`
	Poll               *Poll               `json:"poll" form:"poll"`
}

// Message is returned by almost every request, and contains data about
// almost anything.
type Message struct {
	MessageID             int                `json:"message_id"`
	From                  *User              `json:"from"` // optional
	Date                  int                `json:"date"`
	Chat                  *Chat              `json:"chat"`
	ForwardFrom           *User              `json:"forward_from"`            // optional
	ForwardFromChat       *Chat              `json:"forward_from_chat"`       // optional
	ForwardFromMessageID  int                `json:"forward_from_message_id"` // optional
	ForwardSignature      string             `json:"forward_signature" form:"forward_signature"`
	ForwardSenderName     string             `json:"forward_sender_name" form:"forward_sender_name"`
	ForwardDate           int                `json:"forward_date"`     // optional
	ReplyToMessage        *Message           `json:"reply_to_message"` // optional
	EditDate              int                `json:"edit_date"`        // optional
	MediaGroupID          string             `json:"media_gruop_id" form:"media_gruop_id"`
	AuthorDignature       string             `json:"author_signature" form:"author_signature"`
	Text                  string             `json:"text"`             // optional
	Entities              *[]MessageEntity   `json:"entities"`         // optional
	CaptionEntities       *[]MessageEntity   `json:"caption_entities"` // optional
	Audio                 *Audio             `json:"audio"`            // optional
	Document              *Document          `json:"document"`         // optional
	Animation             *Animation         `json:"animation"`        // optional
	Game                  *Game              `json:"game"`             // optional
	Photo                 *[]PhotoSize       `json:"photo"`            // optional
	Sticker               *Sticker           `json:"sticker"`          // optional
	Video                 *Video             `json:"video"`            // optional
	VideoNote             *VideoNote         `json:"video_note"`       // optional
	Voice                 *Voice             `json:"voice"`            // optional
	Caption               string             `json:"caption"`          // optional
	Contact               *Contact           `json:"contact"`          // optional
	Location              *Location          `json:"location"`         // optional
	Venue                 *Venue             `json:"venue"`            // optional
	Poll                  *Poll              `json:"poll" form:"poll"`
	NewChatMembers        *[]User            `json:"new_chat_members"`        // optional
	LeftChatMember        *User              `json:"left_chat_member"`        // optional
	NewChatTitle          string             `json:"new_chat_title"`          // optional
	NewChatPhoto          *[]PhotoSize       `json:"new_chat_photo"`          // optional
	DeleteChatPhoto       bool               `json:"delete_chat_photo"`       // optional
	GroupChatCreated      bool               `json:"group_chat_created"`      // optional
	SuperGroupChatCreated bool               `json:"supergroup_chat_created"` // optional
	ChannelChatCreated    bool               `json:"channel_chat_created"`    // optional
	MigrateToChatID       int64              `json:"migrate_to_chat_id"`      // optional
	MigrateFromChatID     int64              `json:"migrate_from_chat_id"`    // optional
	PinnedMessage         *Message           `json:"pinned_message"`          // optional
	Invoice               *Invoice           `json:"invoice"`                 // optional
	SuccessfulPayment     *SuccessfulPayment `json:"successful_payment"`      // optional
	//	PassportData          *PassportData      `json:"passport_data,omitempty"` // optional
}

// User is a user on Telegram.
type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`     // optional
	UserName     string `json:"username"`      // optional
	LanguageCode string `json:"language_code"` // optional
	IsBot        bool   `json:"is_bot"`        // optional
}

// Chat contains information about the place a message was sent.
type Chat struct {
	ID                  int64      `json:"id"`
	Type                string     `json:"type"`
	Title               string     `json:"title"`                          // optional
	UserName            string     `json:"username"`                       // optional
	FirstName           string     `json:"first_name"`                     // optional
	LastName            string     `json:"last_name"`                      // optional
	AllMembersAreAdmins bool       `json:"all_members_are_administrators"` // optional
	Photo               *ChatPhoto `json:"photo"`
	Description         string     `json:"description,omitempty"` // optional
	InviteLink          string     `json:"invite_link,omitempty"` // optional
	PinnedMessage       *Message   `json:"pinned_message"`        // optional
}

// Video contains information about a video.
type Video struct {
	FileID    string     `json:"file_id"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Duration  int        `json:"duration"`
	Thumbnail *PhotoSize `json:"thumb"`     // optional
	MimeType  string     `json:"mime_type"` // optional
	FileSize  int        `json:"file_size"` // optional
}

// Voice contains information about a voice.
type Voice struct {
	FileID   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"` // optional
	FileSize int    `json:"file_size"` // optional
}

// VideoNote contains information about a video.
type VideoNote struct {
	FileID    string     `json:"file_id"`
	Length    int        `json:"length"`
	Duration  int        `json:"duration"`
	Thumbnail *PhotoSize `json:"thumb"`     // optional
	FileSize  int        `json:"file_size"` // optional
}

// Contact contains information about a contact.
//
// Note that LastName and UserID may be empty.
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"` // optional
	UserID      int    `json:"user_id"`   // optional
}

// Location contains information about a place.
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Venue contains information about a venue, including its Location.
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID string   `json:"foursquare_id"` // optional
}

//Poll catains informaton about a poll
type Poll struct {
	ID       string        `json:"id"`
	Question string        `json:"question" form:"question"`
	Options  *[]PollOption `json:"options" form:"options"`
	IsClosed bool          `json:"is_closed" form:"is_closed"`
}

//PollOption contains information about one answer option in a poll
type PollOption struct {
	Text       string `json:"text" form:"text"`
	VoterCount int    `json:"voter_count" form:"voter_count"`
}

// Sticker contains information about a sticker.
type Sticker struct {
	FileID    string     `json:"file_id"`
	Width     int        `json:"width"`
	Height    int        `json:"height"`
	Thumbnail *PhotoSize `json:"thumb"`     // optional
	Emoji     string     `json:"emoji"`     // optional
	FileSize  int        `json:"file_size"` // optional
	SetName   string     `json:"set_name"`  // optional
}

// PhotoSize contains information about photos.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"` // optional
}

// Invoice contains basic information about an invoice.
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

// SuccessfulPayment contains basic information about a successful payment.
type SuccessfulPayment struct {
	Currency                string     `json:"currency"`
	TotalAmount             int        `json:"total_amount"`
	InvoicePayload          string     `json:"invoice_payload"`
	ShippingOptionID        string     `json:"shipping_option_id,omitempty"`
	OrderInfo               *OrderInfo `json:"order_info,omitempty"`
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"`
}

// OrderInfo represents information about an order.
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Email           string           `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// ShippingAddress represents a shipping address.
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// MessageEntity contains information about data in a Message.
type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	URL    string `json:"url"`  // optional
	User   *User  `json:"user"` // optional
}

// Game is a game within Telegram.
type Game struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Photo        []PhotoSize     `json:"photo"`
	Text         string          `json:"text"`
	TextEntities []MessageEntity `json:"text_entities"`
	Animation    Animation       `json:"animation"`
}

// Animation is a GIF animation demonstrating the game.
type Animation struct {
	FileID   string    `json:"file_id"`
	Thumb    PhotoSize `json:"thumb"`
	FileName string    `json:"file_name"`
	MimeType string    `json:"mime_type"`
	FileSize int       `json:"file_size"`
}

// Audio contains information about audio.
type Audio struct {
	FileID    string `json:"file_id"`
	Duration  int    `json:"duration"`
	Performer string `json:"performer"` // optional
	Title     string `json:"title"`     // optional
	MimeType  string `json:"mime_type"` // optional
	FileSize  int    `json:"file_size"` // optional
}

// Document contains information about a document.
type Document struct {
	FileID    string     `json:"file_id"`
	Thumbnail *PhotoSize `json:"thumb"`     // optional
	FileName  string     `json:"file_name"` // optional
	MimeType  string     `json:"mime_type"` // optional
	FileSize  int        `json:"file_size"` // optional
}

// PreCheckoutQuery contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             *User      `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// ShippingQuery contains information about an incoming shipping query.
type ShippingQuery struct {
	ID              string           `json:"id"`
	From            *User            `json:"from"`
	InvoicePayload  string           `json:"invoice_payload"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// CallbackQuery is data sent when a keyboard button with callback data
// is clicked.
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`           // optional
	InlineMessageID string   `json:"inline_message_id"` // optional
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`            // optional
	GameShortName   string   `json:"game_short_name"` // optional
}

// ChosenInlineResult is an inline query result chosen by a User
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            *User     `json:"from"`
	Location        *Location `json:"location"`
	InlineMessageID string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

// InlineQuery is a Query from Telegram for an inline request.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Location *Location `json:"location"` // optional
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
}

// ChatPhoto represents a chat photo.
type ChatPhoto struct {
	SmallFileID string `json:"small_file_id"`
	BigFileID   string `json:"big_file_id"`
}

// InlineKeyboardMarkup is a custom keyboard presented for an inline bot.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// ReplyKeyboardMarkup allows the Bot to set a custom keyboard.
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard"`   // optional
	OneTimeKeyboard bool               `json:"one_time_keyboard"` // optional
	Selective       bool               `json:"selective"`         // optional
}

// ReplyKeyboardRemove allows the Bot to hide a custom keyboard.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective"`
}

// InlineKeyboardButton is a button within a custom keyboard for
// inline query responses.
//
// Note that some values are references as even an empty string
// will change behavior.
//
// CallbackGame, if set, MUST be first button in first row.
type InlineKeyboardButton struct {
	Text                         string        `json:"text"`
	URL                          *string       `json:"url,omitempty"`                              // optional
	CallbackData                 *string       `json:"callback_data,omitempty"`                    // optional
	SwitchInlineQuery            *string       `json:"switch_inline_query,omitempty"`              // optional
	SwitchInlineQueryCurrentChat *string       `json:"switch_inline_query_current_chat,omitempty"` // optional
	CallbackGame                 *CallbackGame `json:"callback_game,omitempty"`                    // optional
	Pay                          bool          `json:"pay,omitempty"`                              // optional
}

// CallbackGame is for starting a game in an inline keyboard button.
type CallbackGame struct{}

// KeyboardButton is a button within a custom keyboard.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact"`
	RequestLocation bool   `json:"request_location"`
}

//Result struct
type Result struct {
	MessageID int          `json:"message_id" form:"message_id"`
	From      *User        `json:"from" form:"from"`
	Chat      *Chat        `json:"chat" form:"chat"`
	Date      int          `json:"date" form:"date"`
	Photo     []*PhotoSize `json:"photo" form:"photo"`
	Caption   string       `json:"caption" form:"caption"`
}

//Status is send status
type Status struct {
	OK     bool    `json:"ok" form:"ok"`
	Result *Result `json:"result" form:"result"`
}
