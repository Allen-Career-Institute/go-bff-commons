package framework

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	framework "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models"

	"github.com/labstack/echo/v4"
)

type DatasourceMappingsManager interface {
	RegisterDataSource(ds *framework.DataSource) bool
	GetDataSourceByName(name string) *framework.DataSource
	GetDataSourcesByNameMap() map[string]*framework.DataSource
	GetSharedDS(c echo.Context, dsName string, cnf *config.CommonConfig) *models.DSResponse
}
