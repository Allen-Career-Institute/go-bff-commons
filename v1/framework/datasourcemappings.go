package framework

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	framework "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
	commonmodels "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"

	"github.com/labstack/echo/v4"
)

type DataSourceMappings struct {
	nameToDS map[string]*framework.DataSource
	logger   logger.Logger
}

func NewDataSourceMappings(log logger.Logger) *DataSourceMappings {
	return &DataSourceMappings{
		logger:   log,
		nameToDS: make(map[string]*framework.DataSource),
	}
}

func (dm *DataSourceMappings) RegisterDataSource(ds *framework.DataSource) bool {
	name := ds.Name()
	_, exists := dm.nameToDS[name]

	if exists {
		dm.logger.Errorf("Data Source Already exist with name %s , uri %s , handler %s  ", name, ds.URI(), ds.Method())
		return false
	}

	dm.nameToDS[ds.Name()] = ds

	return true
}

func (dm *DataSourceMappings) GetDataSourceByName(name string) *framework.DataSource {
	ds := dm.nameToDS[name]
	return ds
}

func (dm *DataSourceMappings) GetDataSourcesByNameMap() map[string]*framework.DataSource {
	return dm.nameToDS
}

func (dm *DataSourceMappings) GetSharedDS(c echo.Context, dsName string, cnf *config.Config) *commonmodels.DSResponse {
	dm.logger.WithContext(c).Infof("fetching %v data from shared data", dsName)

	var data *commonmodels.DSResponse

	sharedData, ok := c.Get(utils.SharedDataSource).(map[string]*commonmodels.DSResponse)
	dm.logger.WithContext(c).Debugf("CONTEXT_in_datasourcemappings %+v", c)

	if c.Get(utils.SharedDataSource) != nil && !ok {
		dm.logger.WithContext(c).Errorf("error while type conversion for DS, %v", dsName)
	}

	data = sharedData[dsName]
	if data == nil {
		ds := dm.nameToDS[dsName]
		dm.logger.WithContext(c).Infof("sharedData not found, executing handler for ds: %s", dsName)

		processedData, err := ds.ExecuteHandler(c, cnf)
		if err != nil {
			dm.logger.WithContext(c).Errorf("error while executing ds: %v", dsName)
			return nil
		}

		data = &processedData
	}

	return data
}
