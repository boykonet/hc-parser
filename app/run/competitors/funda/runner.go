package funda

import (
	"vbtor/app/conf/funda"
	"vbtor/app/run/competitors"

	"vbtor/modules/html_parser"
	httpRequester "vbtor/modules/http_requester"
	"vbtor/modules/logger"
	"vbtor/modules/printer"

	"go.uber.org/zap"
)

const (
	fundaDomain = "https://www.funda.nl"
	fundaSSLURI = "funda.nl:443"
	fundaCursor = "/en/huur/amsterdam,diemen,amstelveen,haarlem,zaandam,alkmaar/0-2000/3-dagen/"
	// fundaCursor = "/en/huur/heel-nederland/0-2000/3-dagen/"
	connection = "keep-alive"
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:106.0) Gecko/20100101 Firefox/106.0"
	host = "www.funda.nl"

	requestMonth = "ALL"
	pathToMusic = "./music/vbt_notification.mp3"
)

type runner struct {
	flats		map[string]struct{}

	request		httpRequester.IHTTPRequester
	
	printer		printer.IPrinter
	logger		logger.ILogger
}

func NewRunner(
	conf funda.Configuration,
	logger logger.ILogger,
) competitors.ICompetitorRunner {
	request := httpRequester.NewHTTPRequester(fundaSSLURI)
	request.SetHeaders(
		map[string]string{
			"Host": host,
			"Connection": connection,
			"User-Agent": userAgent,
			"Cookie": conf.Cookie,
		},
	)

	return &runner{
		flats: make(map[string]struct{}),
		request: request,
		printer: printer.NewPrinter(),
		logger: logger,
	}
}

// func toFile(path string, body []byte) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return fmt.Errorf("error open file: %v", err)
// 	}
// 	defer file.Close()

// 	reader := bytes.NewReader(body)

// 	_, err = io.Copy(file, reader)
// 	if err != nil {
// 		return fmt.Errorf("error copy file: %v", err)
// 	}
// 	return nil
// }

func (r *runner) handler() ([]string, int) {
	currentCursor := fundaCursor
	countAllFlats := 0

	r.request.ConfigureHTTP2Client()

	var newFlats []string
	for currentCursor != "" {
		r.request.SetRequestURI(fundaDomain + currentCursor)

		if err := r.request.Do(); err != nil {
			r.logger.Error(
				"unexpected error while doing request",
				zap.Error(err),
				zap.String("url", fundaDomain+currentCursor),
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

				newFlats = append(newFlats, fundaDomain+flat)
			}
		}
	}
	return newFlats, countAllFlats
}

func (r *runner) Run(f competitors.CompetitorFunc) {
	flats, countAllFlats := r.handler()

	f("FUNDA", flats, countAllFlats, r.logger)
}