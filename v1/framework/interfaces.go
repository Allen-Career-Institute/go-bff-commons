package framework

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	framework "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
	commonmodels "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"

	"github.com/labstack/echo/v4"
)

type DatasourceMappingsManager interface {
	RegisterDataSource(ds *framework.DataSource) bool
	GetDataSourceByName(name string) *framework.DataSource
	GetDataSourcesByNameMap() map[string]*framework.DataSource
	GetSharedDS(c echo.Context, dsName string, cnf *config.Config) *commonmodels.DSResponse
}
