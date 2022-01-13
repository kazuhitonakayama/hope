package main

import (
	"log"
	// "strings"
	// "os"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "github.com/line/line-bot-sdk-go/linebot"Cqq
)

type Infector struct {
	ItemList []ItemList `json:"itemList"`
}

type ItemList struct {
    Date      string `json:"date"`
    NameJp    string `json:"name_jp"`
    Npatients string `json:"npatients"`
}

func main() {
	// resp, err := http.Get("https://xcodaxq124.execute-api.ap-northeast-1.amazonaws.com/Prod/books")
    resp, err := http.Get("https://opendata.corona.go.jp/api/Covid19JapanAll")
    // エラーに値があればログに出力し終了する
    if err != nil {
        log.Fatal(err)
    }
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var infectors Infector
	if err := json.Unmarshal(body, &infectors); err != nil {
        log.Fatal(err)
    }
    
    for _, infector := range infectors.ItemList {
        fmt.Printf("%s \n", infector)
    }

    // // LINE Botクライアント生成する
    // // BOT にはチャネルシークレットとチャネルトークンを環境変数から読み込み引数に渡す
    // bot, err := linebot.New(
    //     os.Getenv("LINE_BOT_CHANNEL_SECRET"),
    //     os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
    // )
	// // エラーに値があればログに出力し終了する
    // if err != nil {
    //     log.Fatal(err)
    // }

    // // テキストメッセージを生成する
    // message := linebot.NewTextMessage("hello, world")
    // // テキストメッセージを友達登録しているユーザー全員に配信する
    // if _, err := bot.BroadcastMessage(message).Do(); err != nil {
    //     log.Fatal(err)
    // }
}