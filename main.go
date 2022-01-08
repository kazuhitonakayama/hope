package main

import (
	"log"
	// "strings"
	// "os"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "github.com/line/line-bot-sdk-go/linebot"
)

type Book struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Contents string `json:"contents"`
}

func main() {
	resp, err := http.Get("https://gcal8havkk.execute-api.ap-northeast-1.amazonaws.com/Prod/books")
    // エラーに値があればログに出力し終了する
    if err != nil {
        log.Fatal(err)
    }
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var books []Book
	if err := json.Unmarshal(body, &books); err != nil {
        log.Fatal(err)
    }
	fmt.Println(string(body))

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