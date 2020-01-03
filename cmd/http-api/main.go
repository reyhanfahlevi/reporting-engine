package main

import (
	"github.com/olivere/elastic/v7"
	"github.com/reyhanfahlevi/pkg/go/log"
	"github.com/tokopedia/td-report-engine/app/api/http"
	"github.com/tokopedia/td-report-engine/config"
	esreportsvc "github.com/tokopedia/td-report-engine/internal/service/es-report"
	ucreporting "github.com/tokopedia/td-report-engine/internal/usecase/reporting"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.Get()

	/* initialize resource like db, elastic, redis, httpclient etc here */

	if err != nil {
		log.Fatal("failed connect to db. ", err)
	}

	/* initialize services */

	es, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(cfg.Elastic.Host))
	if err != nil {
		log.Fatal(err)
	}
	esReportSvc := esreportsvc.New(es)

	/* initialize usecase */

	ucReporting := ucreporting.New(esReportSvc)

	/* initialize http handler */

	serve := http.New(ucReporting)

	// run server
	serve.Run(cfg.App.Port)
}
