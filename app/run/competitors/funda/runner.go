package funda

import (
	confFunda "vbtor/app/conf/funda"
	"vbtor/app/run/competitors"

	"vbtor/modules/html_parser"
	httpRequester "vbtor/modules/http_requester"
	"vbtor/modules/logger"
	"vbtor/modules/printer"

	"go.uber.org/zap"
)

const (
	connection = "keep-alive"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:106.0) Gecko/20100101 Firefox/106.0"

	requestMonth = "ALL"
	pathToMusic = "./music/vbt_notification.mp3"
)

type runner struct {
	conf		confFunda.Configuration

	flats		map[string]struct{}

	request		httpRequester.IHTTPRequester
	
	printer		printer.IPrinter
	logger		logger.ILogger
}

func NewRunner(
	conf confFunda.Configuration,
	logger logger.ILogger,
) competitors.ICompetitorRunner {
	request := httpRequester.NewHTTPRequester(conf.Site + ":" + "443")
	request.SetHeaders(
		map[string]string{
			"Host": conf.Host,
			"Connection": connection,
			"User-Agent": userAgent,
			"Cookie": conf.Cookie,
		},
	)

	return &runner{
		conf: conf,
		flats: make(map[string]struct{}),
		request: request,
		printer: printer.NewPrinter(),
		logger: logger,
	}
}

func (r *runner) handler() ([]string, int) {
	currentCursor := r.conf.Cursor
	countAllFlats := 0

	r.request.ConfigureHTTP2Client()

	var newFlats []string
	for currentCursor != "" {
		r.request.SetRequestURI(r.conf.Domain + currentCursor)

		if err := r.request.Do(); err != nil {
			r.logger.Error(
				"unexpected error while doing request",
				zap.Error(err),
				zap.String("url", r.conf.Domain+currentCursor),
			)
			break
		}

		body := r.request.GetBody()

		htmlParser := html_parser.NewFundaHTMLParser()
		flats, err := htmlParser.ParseFlats(body, &currentCursor)
		if err != nil {
			r.logger.Error("parce flats error", zap.Error(err))
		}

		countAllFlats += len(flats)

		for _, flat := range flats {
			_, ok := r.flats[flat]
			if ok == false {
				r.flats[flat] = struct{}{}

				newFlats = append(newFlats, r.conf.Domain+flat)
			}
		}
	}
	return newFlats, countAllFlats
}

func (r *runner) Run(f competitors.CompetitorFunc) {
	flats, countAllFlats := r.handler()

	f("FUNDA", flats, countAllFlats, r.logger)
}

func (r *runner) SetFlatsFromFile() error {
	return competitors.GetFromFile(r.conf.PropertiesPath, r.flats)
}

func (r *runner) SaveToFile() error {
	return competitors.SaveToFile(r.conf.PropertiesPath, r.flats)
}