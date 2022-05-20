package notifyunit

import (
	"log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	Bot 	*tgbotapi.BotAPI
	ChatID	int64
}

func (tn *TelegramNotifier) Init(api_key string, chat_id int64) {
	bot, err := tgbotapi.NewBotAPI(api_key)
	if err != nil {
		log.Panic(err)
	}
	tn.Bot = bot
//	tn.Bot.Debug = true
	tn.ChatID = chat_id
}
	

func (tn *TelegramNotifier) Notify(msg string) {
	tn.Bot.Send(tgbotapi.NewMessage(tn.ChatID, msg))
}


/*
func main() {
	tn := new(TelegramNotifier)
	tn.Init("5132582314:AAEVW_xw4JvRPWb-HQ9_VU8i62L3-QFJQ8U", 5055419231)
	tn.Notify("hello")
}
*/
