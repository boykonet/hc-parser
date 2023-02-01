package vbt

import (
	"sync"

	confVBT "vbtor/app/conf/vbt"
	"vbtor/modules/html_parser"
	httpRequester "vbtor/modules/http_requester"
	"vbtor/modules/logger"
	"vbtor/modules/printer"

	"vbtor/app/run/competitors"

	"go.uber.org/zap"
)

type runner struct {
	wg    *sync.WaitGroup
	mutex *sync.Mutex

	conf confVBT.Configuration

	flats map[string]struct{}

	requests map[string]httpRequester.IHTTPRequester

	printer printer.IPrinter

	logger logger.ILogger
}

func NewRunner(
	conf confVBT.Configuration,
	logger logger.ILogger,
) competitors.ICompetitorRunner {
	requests := make(map[string]httpRequester.IHTTPRequester, 2)

	rDecember := httpRequester.NewHTTPRequester(conf.Site + ":" + "443")
	rDecember.SetCookie(conf.Cookie.December)

	rJanuary := httpRequester.NewHTTPRequester(conf.Site + ":" + "443")
	rJanuary.SetCookie(conf.Cookie.January)

	requests["DECEMBER"] = rDecember
	requests["JANUARY"] = rJanuary

	return &runner{
		wg:    &sync.WaitGroup{},
		mutex: &sync.Mutex{},

		conf: conf,

		flats: make(map[string]struct{}, 0),

		requests: requests,

		printer: printer.NewPrinter(),

		logger: logger,
	}
}

func (r *runner) handler(
	httpRequester httpRequester.IHTTPRequester,
	flatsChan chan []string,
	countAllFlatsChan chan int,
	month string,
) {
	currentCursor := r.conf.Cursor
	countAllFlats := 0

	httpRequester.ConfigureHTTP2Client()

	var newFlats []string
	for currentCursor != "" {
		httpRequester.SetRequestURI(r.conf.Domain + currentCursor)

		if err := httpRequester.Do(); err != nil {
			r.logger.Error(
				"unexpected error while doing request",
				zap.Error(err),
				zap.String("url", r.conf.Domain + currentCursor),
			)
			break
		}

		body := httpRequester.GetBody()

		htmlParser := html_parser.NewVbtHTMLParser()
		flats, _ := htmlParser.ParseFlats(body, &currentCursor)

		countAllFlats += len(flats)

		for _, flat := range flats {
			_, ok := r.flats[flat]
			if ok == false {
				r.mutex.Lock()
				r.flats[flat] = struct{}{}
				r.mutex.Unlock()

				newFlats = append(newFlats, r.conf.Domain+flat)
			}
		}
	}

	flatsChan <- newFlats
	close(flatsChan)

	countAllFlatsChan <- countAllFlats
	close(countAllFlatsChan)

	r.wg.Done()
}

func (r *runner) Run(f competitors.CompetitorFunc) {
	var flatsChans []chan []string
	var countAllFlatsChans []chan int

	lenRequests := len(r.requests)

	for i := 0; i < lenRequests; i++ {
		flatsChans = append(flatsChans, make(chan []string))
		countAllFlatsChans = append(countAllFlatsChans, make(chan int))
	}

	i := 0
	for httpRequesterName, httpRequester := range r.requests {
		r.wg.Add(1)

		go r.handler(httpRequester, flatsChans[i], countAllFlatsChans[i], httpRequesterName)
		i++
	}

	var flats []string
	var countAllFlats int

	for i := 0; i < len(r.requests); i++ {
		flats = append(flats, <-flatsChans[i]...)
		countAllFlats += <- countAllFlatsChans[i]
	}

	r.wg.Wait()

	f("VB&T", flats, countAllFlats, r.logger)
}

func (r *runner) SetFlatsFromFile() error {
	return competitors.GetFromFile(r.conf.PropertiesPath, r.flats)
}

func (r *runner) SaveToFile() error {
	return competitors.SaveToFile(r.conf.PropertiesPath, r.flats)
}
