package main

import (
	"log"
	// // "strings"
	"os"
	"encoding/json"
	"fmt"
	// "database/sql"
	// "fmt"
	"io"
	// "log"
	// "encoding/json"
	// _ "github.com/go-sql-driver/mysql"
	"net/http"
	// "net/url"
	"github.com/line/line-bot-sdk-go/linebot"
)

type ItemList struct {
    Date      string `json:"date"`
    NameJp    string `json:"name_jp"`
    Npatients int `json:"npatients"`
    Diff int `json:"diff"`
	CreatedAt string `json:"created_at"`
}

func main() {
    // LINE Botクライアント生成する
    // BOT にはチャネルシークレットとチャネルトークンを環境変数から読み込み引数に渡す
    bot, err := linebot.New(
        os.Getenv("LINE_BOT_CHANNEL_SECRET"),
        os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
    )
	// エラーに値があればログに出力し終了する
    if err != nil {
        log.Fatal(err)
    }

    first_text := ExecuteFirstMessage()
    second_text := ExecuteSecondMessage()

    // テキストメッセージを生成する
    message := linebot.NewTextMessage(first_text)
    // テキストメッセージを友達登録しているユーザー全員に配信するa
    if _, err := bot.BroadcastMessage(message).Do(); err != nil {
        log.Fatal(err)
    }
    // テキストメッセージを生成する
    second_message := linebot.NewTextMessage(second_text)
    // テキストメッセージを友達登録しているユーザー全員に配信するa
    if _, err := bot.BroadcastMessage(second_message).Do(); err != nil {
        log.Fatal(err)
    }

    // Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func ExecuteFirstMessage() string{
    itemlists := FetchLatestInfectors()

    date := itemlists[0].Date
    header := fmt.Sprintf("最新（%s時点）の感染者数速報です🎉 \n感染者数が多い順で並べてます👍 \n\n", date)
    text := header

    for _ , item := range itemlists {
        text = text + fmt.Sprintf("%s：%d（前日比 + %d 人）\n", item.NameJp, item.Npatients, item.Diff)
    }
    return text
}

func ExecuteSecondMessage() string{
    dangerslists := FetchDangers()

    header_second := fmt.Sprintf("前日からの感染者数の増加数が高いトップ10です、、、 \n気をつけてね、、 \n\n")
    text_second := header_second

    for _ , item := range dangerslists {
        text_second = text_second + fmt.Sprintf("%s：%d（前日比 + %d 人）\n", item.NameJp, item.Npatients, item.Diff)
    }
    return text_second
}

func FetchLatestInfectors() []ItemList{
    api_url := "https://jbft55gtp3.execute-api.ap-northeast-1.amazonaws.com/Prod/infectors"

    resp, err := http.Get(api_url)
    // エラーに値があればログに出力し終了する
    if err != nil {
        log.Fatal(err)
    }
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var itemlist []ItemList
	if err := json.Unmarshal(body, &itemlist); err != nil {
        log.Fatal(err)
    }

    return itemlist
}

func FetchDangers() []ItemList{
    api_url := "https://jbft55gtp3.execute-api.ap-northeast-1.amazonaws.com/Prod/dangers"

    resp, err := http.Get(api_url)
    // エラーに値があればログに出力し終了する
    if err != nil {
        log.Fatal(err)
    }
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var itemlist []ItemList
	if err := json.Unmarshal(body, &itemlist); err != nil {
        log.Fatal(err)
    }

    return itemlist
}