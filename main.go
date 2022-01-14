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
    itemlists := FetchLatestInfectors()

    date := itemlists[0].Date
    header := fmt.Sprintf("最新（%s時点）の感染者数速報です🎉 \n感染者数が多い順で並べてます👍 \n\n", date)
    text := header

    for _ , item := range itemlists {
        text = text + fmt.Sprintf("%s：%d（前日比 + %d 人）\n", item.NameJp, item.Npatients, item.Diff)
    }

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

    // テキストメッセージを生成する
    message := linebot.NewTextMessage(text)
    // テキストメッセージを友達登録しているユーザー全員に配信する
    if _, err := bot.BroadcastMessage(message).Do(); err != nil {
        log.Fatal(err)
    }
}

func FetchLatestInfectors() []ItemList{
    api_url := "https://s0oe7xjx7i.execute-api.ap-northeast-1.amazonaws.com/Prod/infectors"

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