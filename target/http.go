package targethandlers

import (
	"time"
	"sync"
	"fmt"
	"net"
	"context"
	"net/http"
	neturl "net/url"
	config "monitoring/config"
)

type HttpTarget struct {
	EndPoints 	[]config.HttpTargetConfig
	Client		*http.Client
}


func (p *HttpTarget) Init(raw_conf []config.HttpTargetConfig) {
	p.EndPoints = raw_conf
	p.Client = &http.Client{ // interesting
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
        		return http.ErrUseLastResponse
		},
    	}
}

func (p *HttpTarget) sendAlert(name string, url string, problem string, alertchannel chan<- map[string]string) {
	msg := fmt.Sprintf("[HTTP] %s:%s is not available - %s", name, url, problem)
	alertchannel <- map[string]string{
		name + ":" + url : msg,
	}
}

func (p *HttpTarget) runnerNotifier(
		wg *sync.WaitGroup, 
		conf config.HttpTargetConfig, 
		alertchannel chan<- map[string]string) {

	defer wg.Done()
	timeout := conf.Timeout
	repeat := conf.Repeat
	status_code := conf.StatusCode
	name := conf.Name

	url := conf.Url
	_, err := neturl.Parse(url)
	if err != nil {
		panic("Url parse error")
	}


	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		panic("Can't form request")
	}

	for i := 0; i < repeat; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
		defer cancel()
		req = req.WithContext(ctx)

		response, err := p.Client.Do(req)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				if i == repeat - 1 {
					p.sendAlert(name, url, "http timeout exceed", alertchannel)
					return
				}
			}
			if i == repeat - 1 {
				p.sendAlert(name, url, "http response error", alertchannel)
				return
			}
			continue
		}
		if response.StatusCode != status_code && i == repeat - 1 {
			p.sendAlert(name, url, "http status code is not correct", alertchannel)
			return
		}
		break
	}
}

func (p *HttpTarget) Monitor(delay time.Duration, alertchannel chan<- map[string]string) {
	var wg sync.WaitGroup

	for {
		for _, httpr := range p.EndPoints {
			wg.Add(1)
			go p.runnerNotifier(&wg, httpr, alertchannel)
		}
		wg.Wait()
		time.Sleep(delay)
	}
}

/*
func main() {
	var p HttpTarget;
	a := map[string]interface{} {
		"Url" : "http://www.google.com/",
		"Repeat" : 2,
		"Timeout" : 2,
		"StatusCode" : 200,
	}

	b := map[string]interface{} {
		"Url" : "https://www.mipt.ru/",
		"Repeat" : 2,
		"Timeout" : 2,
		"StatusCode" : 200,
	}

	c := map[string]interface{} {
		"Url" : "http://lyalichikimchen.com/",
		"Repeat" : 2,
		"Timeout" : 2,
		"StatusCode" : 200,
	}


	p.Init([]map[string]interface{}{a,b,c})
	p.Monitor(10 * time.Second)
}
*/
