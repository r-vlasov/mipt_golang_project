package main
import (
	notify "monitoring/notify"
	notifyunit "monitoring/notify/units"
	config	"monitoring/config"
	target "monitoring/target"
	"os"
	"fmt"
	"time"
	"strconv"
	"path/filepath"
)



func main() {
	fmt.Println("... Start monitoring ...\n")

	fmt.Println("... Parse main config ...")
	configuration := new(config.MainConfig)
	fp, _ := filepath.Abs("./config.yml")
	if configuration.Load(fp) != nil {
		panic("Can't parse main config")
	}

	fmt.Println("... Initialize notificator ...")
	notifier := new(notify.Notifier)

	fmt.Println("... Add telegram ...")
	tgnotifier := new(notifyunit.TelegramNotifier)

	tg_apikey, exists := os.LookupEnv("TG_APIKEY")
	if !exists {
		panic("Set up TG_APIKEY environment variable")
	}

	tg_chatid_raw, exists := os.LookupEnv("TG_CHAT_ID")
	if !exists {
		panic("Set up TG_CHAT_ID environment variable")
	}
	tg_chatid, err := strconv.Atoi(tg_chatid_raw)
	if err != nil {
		panic("Set up correct TG_CHAT_ID environment variable")
	}

	tgnotifier.Init(tg_apikey, int64(tg_chatid))

	notifier.Init(configuration.NotificationDelay, "[Custom monitoring] #> ", []notify.NotifierUnitInterface{
		&notifyunit.StdoutNotifier{},  // debug - mirror telegram -> stdout
		tgnotifier,
	})
	
	fmt.Println("... Run Notification ...")
	go notifier.RunNotify()




	fmt.Println("... Initialize HTTP-target ...")
	var p target.HttpTarget
	p.Init(configuration.Targets.Http)
	go p.Monitor(10 * time.Second, notifier.GetAlertChannel()) // 10 - per target delay


	fmt.Println("... Initialize PING-target ...")
	var t target.PingTarget
        t.Init(configuration.Targets.Ping)
        go t.Monitor(10 * time.Second, notifier.GetAlertChannel()) // 10 - per target delay


	fmt.Println("... Initialize DNS-target ...")
	var d target.DnsTarget
        d.Init(configuration.Targets.Dns)
        go d.Monitor(10 * time.Second, notifier.GetAlertChannel()) // 10 - per target delay

	/* not yet implemented
	RIP protocol
	*/

	// sleep forewer
	time.Sleep(time.Duration(1 << 63 - 1))



	/* future
	targets := []string{"HTTP", "PING", "DNS"}
	confs := []interface{}{
		configuration.Targets.Http,
		configuration.Targets.Ping,
		configuration.Targets.Dns,
	}
	items := []interface{}{
		target.HttpTarget{},
		target.PingTarget{},
		target.DnsTarget{},
	}

	for i := 0; i < len(targets); i++ {
		fmt.Println("... Initialize " + targets[i] + "-target ...")
		items[i].Init(confs[i])
		go items[i].Monitor(10 * time.Second, notifier.GetAlertChannel())
	}
	*/


}
