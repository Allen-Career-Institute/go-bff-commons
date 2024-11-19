// nolint:funlen,wsl // This package is long, but it is required to register all the functions
package routes

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
	internal "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	apiMiddlewares "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/middleware"
	models "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/models/commons"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/httperr"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func MapDataSourceURIRoutes(e *echo.Echo, dsm framework.DatasourceMappingsManager, cfg *config.Config, log logger.Logger, mw *apiMiddlewares.Manager, meter metric.Meter, m internal.Mapper) {
	uriMap := dsm.GetDataSourcesByNameMap()

	var uri = "/api/v1"
	v1 := e.Group(uri)

	for _, ds := range uriMap {
		if ds.URI() == "" || ds.Method() == "" {
			log.Infof("registering ds name " + ds.Name())
			continue
		}

		log.Infof("registering uri : " + uri + ds.URI() + ", method : " + ds.Method() + " with ds_name " + ds.Name())

		switch ds.Method() {
		case "POST":
			v1.POST(ds.URI(), NewDataSourceExecutor(ds, dsm, ds.Name(), cfg, log, meter, m).ExecuteDataSource, mw.AuthNMiddleware(cfg), mw.RbacMiddleware(cfg, ds.Resource(), ds.Action()))
		case "GET":
			v1.GET(ds.URI(), NewDataSourceExecutor(ds, dsm, ds.Name(), cfg, log, meter, m).ExecuteDataSource, mw.AuthNMiddleware(cfg), mw.RbacMiddleware(cfg, ds.Resource(), ds.Action()))
		case "PUT":
			v1.PUT(ds.URI(), NewDataSourceExecutor(ds, dsm, ds.Name(), cfg, log, meter, m).ExecuteDataSource, mw.AuthNMiddleware(cfg), mw.RbacMiddleware(cfg, ds.Resource(), ds.Action()))
		case "DELETE":
			v1.DELETE(ds.URI(), NewDataSourceExecutor(ds, dsm, ds.Name(), cfg, log, meter, m).ExecuteDataSource, mw.AuthNMiddleware(cfg), mw.RbacMiddleware(cfg, ds.Resource(), ds.Action()))
		case "PATCH":
			v1.PATCH(ds.URI(), NewDataSourceExecutor(ds, dsm, ds.Name(), cfg, log, meter, m).ExecuteDataSource, mw.AuthNMiddleware(cfg), mw.RbacMiddleware(cfg, ds.Resource(), ds.Action()))
		default:
			log.Infof("unknown http method " + ds.Method())
		}
	}
}

type DataSourceExecutor struct {
	ds     *datasource.DataSource
	dsm    framework.DatasourceMappingsManager
	dsName string
	cnf    *config.Config
	logger logger.Logger
	meter  metric.Meter
	m      internal.Mapper
}

func NewDataSourceExecutor(ds *datasource.DataSource, dsm framework.DatasourceMappingsManager, dsName string, cfg *config.Config,
	log logger.Logger, meter metric.Meter, m internal.Mapper) *DataSourceExecutor {
	return &DataSourceExecutor{ds: ds, dsm: dsm, dsName: dsName, cnf: cfg, logger: log, meter: meter, m: m}
}

func (e *DataSourceExecutor) ExecuteDataSource(c echo.Context) error {
	e.logger.WithContext(c).Infof("Executing Filter for DS %s", e.ds.Name())

	requestCount, requestErr := e.m.GetCount(utils.BffDsMetricPrefix + utils.Count)
	if requestErr != nil {
		e.logger.WithContext(c).Errorf("error in sending metric for request count %s", requestErr)
		return requestErr
	}

	startTime := time.Now()

	filters := e.ds.GetFilters()
	for _, filter := range filters {
		err := filter(c)

		if err != nil {
			sc := httperr.FromError(err)
			requestCount.Add(c.Request().Context(), 1, getAddMetricTags(sc.Status(), e.dsName)...)

			return c.JSON(
				sc.Status(),
				internal.PopulateResponse(sc.Status(), sc.Error(), nil),
			)
		}
	}

	connTimeout := time.Duration(e.ds.Timeout()) * time.Millisecond

	toCtx, conCancel := utils.GetRequestCtxWithTimeout(c, connTimeout)
	withMD := utils.AddAuthHeaderAsMetadata(toCtx, c)
	c.SetRequest(c.Request().WithContext(withMD))

	defer conCancel()

	response, err := e.ds.ExecuteHandler(c, e.cnf)

	reqDuration, rdErr := e.m.GetDuration(utils.BffDsMetricPrefix + utils.Duration)
	if rdErr != nil {
		requestCount.Add(c.Request().Context(), 1, getAddMetricTags(response.Status, e.dsName)...)

		e.logger.WithContext(c).Errorf("error in sending metric for request duration from ds %s", rdErr)

		return rdErr
	}

	if err != nil {
		sc := httperr.FromError(err)

		requestCount.Add(c.Request().Context(), 1, getAddMetricTags(response.Status, e.dsName)...)
		reqDuration.Record(c.Request().Context(), time.Since(startTime).Milliseconds(), getRecordMetricTags(response.Status, e.dsName)...)
		c.SetResponse(&echo.Response{
			Writer:    c.Response().Writer,
			Status:    sc.Status(),
			Size:      c.Response().Size,
			Committed: c.Response().Committed,
		})

		return err
	}

	requestCount.Add(c.Request().Context(), 1, getAddMetricTags(response.Status, e.dsName)...)
	reqDuration.Record(c.Request().Context(), time.Since(startTime).Milliseconds(), getRecordMetricTags(response.Status, e.dsName)...)

	return c.JSON(response.Status, response)
}

func (e *DataSourceExecutor) ExecuteDataSourceFromDS(c echo.Context) (*models.DSResponse, error) {
	e.logger.WithContext(c).Debugf("Executing Filter for datasource from DS %s", e.ds.Name())

	startTime := time.Now()
	requestCount, requestErr := e.m.GetCount(utils.BffDsMetricPrefix + utils.Count)

	if requestErr != nil {
		e.logger.WithContext(c).Errorf("error in sending metric for request count from DS %s", requestErr)
		return nil, requestErr
	}
	filters := e.ds.GetFilters()

	for _, filter := range filters {
		err := filter(c)
		if err != nil {
			sc := httperr.FromError(err)
			requestCount.Add(c.Request().Context(), 1, getAddMetricTags(sc.Status(), e.dsName)...)

			return nil, err
		}
	}

	response, err := e.ds.ExecuteHandler(c, e.cnf)
	reqDuration, rdErr := e.m.GetDuration(utils.BffDsMetricPrefix + utils.Duration)

	if rdErr != nil {
		e.logger.WithContext(c).Errorf("error in sending metric for request duration from ds %s", rdErr)
		requestCount.Add(c.Request().Context(), 1, getAddMetricTags(response.Status, e.dsName)...)

		return nil, rdErr
	}

	if err != nil {
		e.logger.WithContext(c).Errorf("error while sending metric for successCount from DS %s", err)
		requestCount.Add(c.Request().Context(), 1, getAddMetricTags(response.Status, e.dsName)...)
		reqDuration.Record(c.Request().Context(), time.Since(startTime).Milliseconds(), getRecordMetricTags(response.Status, e.dsName)...)

		return nil, err
	}

	requestCount.Add(c.Request().Context(), 1, getAddMetricTags(response.Status, e.dsName)...)
	reqDuration.Record(c.Request().Context(), time.Since(startTime).Milliseconds(), getRecordMetricTags(response.Status, e.dsName)...)

	return (*models.DSResponse)(&response), nil
}

func getAddMetricTags(status int, dsName string) []metric.AddOption {
	return []metric.AddOption{
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.StatusCode,
			Value: attribute.IntValue(status),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.ServiceEnv,
			Value: attribute.StringValue(os.Getenv("ENV")),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.DataSourceName,
			Value: attribute.StringValue(dsName),
		}),
	}
}

func getRecordMetricTags(status int, dsName string) []metric.RecordOption {
	return []metric.RecordOption{
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.StatusCode,
			Value: attribute.IntValue(status),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.ServiceEnv,
			Value: attribute.StringValue(os.Getenv("ENV")),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.DataSourceName,
			Value: attribute.StringValue(dsName),
		}),
	}
}
