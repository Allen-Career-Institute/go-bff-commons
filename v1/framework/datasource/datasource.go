package datasource

import (
	"github.com/Allen-Career-Institute/common-protos/authorization/v1/types"
	"github.com/labstack/echo/v4"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
)

type (
	DataSource struct {
		filters   []Filter
		handler   HandlerFunc
		name      string
		uri       string
		timeout   int64
		method    string
		resource  types.ResourceTypes
		action    types.Action
		preloadDS []string
	}

	// Filter defines a function to process middleware.
	Filter func(ctx echo.Context) error

	// HandlerFunc defines a function to serve HTTP requests.
	HandlerFunc func(ctx echo.Context, cnf *config.Config) (models.DSResponse, error)
)

func CreateNewDataSource(dsConf models.DataSourceConfig, handler HandlerFunc, m ...Filter) (ds *DataSource) {
	resource := types.ResourceTypes_RESOURCE_UNSPECIFIED

	if dsConf.Resource != "" {
		enumValue, ok := types.ResourceTypes_value[dsConf.Resource]
		if !ok {
			panic("Invalid resource: " + dsConf.Resource + utils.ConfigureString + dsConf.DsName)
		}

		resource = types.ResourceTypes(enumValue)
	}

	action := types.Action_ACTION_UNSPECIFIED

	if dsConf.Action != "" {
		enumValue, ok := types.Action_value[dsConf.Action]
		if !ok {
			panic("Invalid action: " + dsConf.Action + utils.ConfigureString + dsConf.DsName)
		}

		action = types.Action(enumValue)
	}

	ds = &DataSource{
		handler:   handler,
		name:      dsConf.DsName,
		uri:       dsConf.URI,
		timeout:   int64(dsConf.Timeout),
		method:    dsConf.Method,
		filters:   addFilters(m...),
		resource:  resource,
		action:    action,
		preloadDS: dsConf.PreloadDS,
	}

	return ds
}

func addFilters(filter ...Filter) []Filter {
	var filters []Filter
	for i := len(filter) - 1; i >= 0; i-- {
		filters = append(filters, filter[i])
	}

	return filters
}

func (dsm *DataSource) URI() string {
	return dsm.uri
}

func (dsm *DataSource) Name() string {
	return dsm.name
}

func (dsm *DataSource) Timeout() int64 {
	return dsm.timeout
}

func (dsm *DataSource) Method() string {
	return dsm.method
}

func (dsm *DataSource) Resource() types.ResourceTypes {
	return dsm.resource
}

func (dsm *DataSource) Action() types.Action {
	return dsm.action
}

func (dsm *DataSource) ExecuteHandler(c echo.Context, cnf *config.Config) (models.DSResponse, error) {
	return dsm.handler(c, cnf)
}

func (dsm *DataSource) GetFilters() []Filter {
	return dsm.filters
}

func (dsm *DataSource) GetPreloadDS() []string {
	return dsm.preloadDS
}
