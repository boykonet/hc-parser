package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"vbtor/modules/html_parser"
	"vbtor/modules/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.com/dgrr/http2"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/valyala/fasthttp"
)

var (
	pathToNotificationMusic = "./music/vbt_notification.mp3"

	url                    = "https://vbtverhuurmakelaars.nl"
	userAgent              = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:106.0) Gecko/20100101 Firefox/106.0"
	accept                 = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
	acceptLang             = "en-US,en;q=0.5"
	acceptEncoding         = "gzip, deflate, br"
	connection             = "keep-alive"
	upgradeUncecureRequest = "1"
	secFetchDest           = "document"
	secFetchMode           = "navigate"
	secFetchSite           = "none"
	te                     = "trailers"
)

type configuration struct {
	Cookie string `yaml: cookie`
	Min int `yaml: min`
	Max int `yaml: max`
	// req.Header.Add("User-Agent", userAgent)
	// req.Header.Add("Accept", accept)
	// req.Header.Add("Accept-Language", acceptLang)
	// req.Header.Add("Accept-Encoding", acceptEncoding)
	// req.Header.Add("Connection", connection)
	// req.Header.Add("Cookie", conf.Cookie)
	// req.Header.Add("Upgrade-Insecure-Requests", upgradeUncecureRequest)
	// req.Header.Add("Sec-Fetch-Dest", secFetchDest)
	// req.Header.Add("Sec-Fetch-Mode", secFetchMode)
	// req.Header.Add("Sec-Fetch-Site", secFetchSite)
	// req.Header.Add("TE", te)
}

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

type textColors struct {
	Reset, Red, Green, Yellow, Blue string
}

func initColors() textColors {
	return textColors{
		Reset: reset,
		Red: red,
		Green: green,
		Yellow: yellow,
		Blue: blue,
	}
}

func parceConfiguration(pathToFile string) (*configuration, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	conf := &configuration{}

	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, file); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(buff.Bytes(), conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func run(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}

func printMessage(currentCounter int, newFlats []string, firstCursor string, log logger.ILogger) {
	colors := initColors()

	fmt.Println(colors.Yellow + "VB&T Housing Corporation parsing" + colors.Reset)

	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	fmt.Println("	Current counter:", currentCounter)
	fmt.Println("	View all items:", colors.Blue + url + firstCursor + colors.Reset)
	if len(newFlats) > 0 {
		fmt.Println("	Status:",
		colors.Green + strconv.Itoa(len(newFlats)) + " new items found" + colors.Reset)
		for i, flat := range newFlats {
			var color string
			if i % 2 == 1 {
				color = colors.Yellow
			} else {
				color = colors.Blue
			}
			fmt.Println()
			fmt.Println("		Link:", color + flat + colors.Reset)
		}
		if err := run(pathToNotificationMusic); err != nil {
			log.Error("falat error play music", zap.Error(err))
		}
	} else {
		fmt.Println("	Status: no new items found")
	}
	fmt.Println()
}

func main() {
	log := logger.NewLogger()
	
	// Parse configuration file with cookie and min/max intervals
	conf, err := parceConfiguration("configurations/conf.yaml")
	if err != nil {
		log.Fatal("fatal error parse configuration", zap.Error(err))
	}

	req := fasthttp.AcquireRequest()
	req.Header.SetCookie("Cookie", conf.Cookie)

	hc := &fasthttp.HostClient{
		Addr: "vbtverhuurmakelaars.nl:443",
	}
	err = http2.ConfigureClient(hc, http2.ClientOpts{})
	if err != nil {
		log.Fatal(fmt.Sprintf("%s doesn't support http/2", hc.Addr), zap.Error(err))
	}

	firstCursor := "/en/woningen"
	addrs := make(map[string]struct{}, 0)

	for {
		var newFlats []string
		cursor := firstCursor
	
		for cursor != "" {
			req.SetRequestURI(url + cursor)
			resp := fasthttp.AcquireResponse()
			err = hc.Do(req, resp)
			if err != nil {
				log.Error("error request", zap.Error(err))
			}
		
			htmlParser := html_parser.NewHTMLParser(resp.Body())
			flats := htmlParser.ParseFlats(&cursor)
	 
			for _, flat := range flats {
				_, ok := addrs[flat]
				if ok == false {
					addrs[flat] = struct{}{}
					newFlats = append(newFlats, url + flat)
				}
			}
		}

		// Print message
		printMessage(len(addrs), newFlats, firstCursor, log)

		// Getting time for sleep
		timeSleep := rand.Intn(conf.Max - conf.Min) + conf.Min
		rand.Seed(time.Now().UnixNano())
		
		// Sleep on this time
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
}
