package rijnland

import (
	"net/http"
	confRijnland "vbtor/app/conf/rijnland"
	"vbtor/app/run/competitors"
	httpRequester "vbtor/modules/http_requester"
	hhhtr1 "vbtor/modules/http_requester/1"
	"vbtor/modules/logger"
	"vbtor/modules/printer"

	"go.uber.org/zap"
)

type runner struct {
	conf		confRijnland.Configuration

	flats		map[string]struct{}

	request		httpRequester.IHTTPRequester
	
	printer		printer.IPrinter
	logger		logger.ILogger
}

func NewRunner(
	conf confRijnland.Configuration,
	logger logger.ILogger,
) competitors.ICompetitorRunner {
	request := hhhtr1.NewHTTP1Requester(http.MethodPost, conf.Domain + conf.Cursor, nil)
	request.SetHeaders(
		map[string]string{
			"Host": conf.Host,
			"Connection": conf.Connection,
			"User-Agent": conf.UserAgent,
			// "Cookie": conf.Cookie,
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

		competitors.ToFile("rijnland", r.request.GetBody())
		

		// body := r.request.GetBody()
	

		// htmlParser := html_parser.NewRijnlandHTMLParser()
		// flats, err := htmlParser.ParseFlats(body, &currentCursor)
		// if err != nil {
		// 	r.logger.Error("parce flats error", zap.Error(err))
		// }

		// countAllFlats += len(flats)

		// for _, flat := range flats {
		// 	_, ok := r.flats[flat]
		// 	if ok == false {
		// 		r.flats[flat] = struct{}{}

		// 		newFlats = append(newFlats, r.conf.Domain+flat)
		// 	}
		// }
		currentCursor = ""
	}
	return newFlats, countAllFlats
}

func (r *runner) Run(f competitors.CompetitorFunc) {
	flats, countAllFlats := r.handler()

	f("RIJNLAND", flats, countAllFlats, r.logger)
}

func (r *runner) SetFlatsFromFile() error {
	return competitors.GetFromFile(r.conf.PropertiesPath, r.flats)
}

func (r *runner) SaveToFile() error {
	return competitors.SaveToFile(r.conf.PropertiesPath, r.flats)
}