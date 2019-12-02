package es_report

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	es_report "github.com/tokopedia/reporting-engine/internal/model/es-report"
)

// Report struct for report package
type Report struct {
	es *elastic.Client
}

// New instantiate the report
func New(es *elastic.Client) *Report {
	return &Report{
		es: es,
	}
}

// StoreReport store the report data
func (r *Report) StoreReport(ctx context.Context, param es_report.ParamReporting) error {
	// check and create index if not available
	err := r.createIndex(ctx, param.Index, param.Mapping)
	if err != nil {
		return err
	}

	_, err = r.es.Index().Index(param.Index).Type(param.Type).BodyJson(param.Data).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

// GetReport get report data
func (r *Report) GetReports(ctx context.Context, param es_report.ParamGetReports) ([]map[string]interface{}, int64, error) {
	var (
		totalData   int64
		reports     = []map[string]interface{}{}
		filter      []elastic.Query
		rangeFilter []elastic.Query
	)

	for k, v := range param.Filter {
		filter = append(filter, elastic.NewTermQuery(k, v))
	}

	for k, v := range param.RangeFilter {
		rangeFilter = append(rangeFilter, elastic.NewRangeQuery(k).Gte(v.From).Lte(v.To))
	}

	search := r.es.Search(param.Index).Type(param.Type).
		Query(
			elastic.NewBoolQuery().
				Filter(filter...).
				Filter(rangeFilter...)).
		From(param.From).
		Size(param.Size).
		Sort("created_time", false)

	resp, err := search.Do(ctx)
	if err != nil && !elastic.IsNotFound(err) {
		return reports, totalData, err
	}

	if elastic.IsNotFound(err) {
		return reports, totalData, nil
	}

	for _, h := range resp.Hits.Hits {
		tmp := make(map[string]interface{})

		err = json.Unmarshal(*h.Source, &tmp)
		if err != nil {
			continue
		}

		reports = append(reports, tmp)
	}

	totalData = resp.TotalHits()
	return reports, totalData, nil
}

func (r *Report) checkIndex(ctx context.Context, index string) (bool, error) {
	exist, err := r.es.IndexExists(index).Do(ctx)
	if err != nil {
		return exist, err
	}

	return exist, nil
}

func (r *Report) createIndex(ctx context.Context, index, mappings string) error {
	check, err := r.checkIndex(ctx, index)
	if err != nil {
		return err
	}

	if check {
		return nil
	}

	svc := r.es.CreateIndex(index)

	if mappings != "" {
		svc.Body(mappings)
	}

	result, err := svc.Do(ctx)
	if err != nil {
		return err
	}

	if !result.Acknowledged {
		return errors.New("index not acknowledged")
	}

	return nil
}
