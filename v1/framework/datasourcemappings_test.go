package framework

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	framework "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

const (
	URI     = "http://localhost:8080/api/v1/auth/sendOtp"
	method  = "SendOtp()"
	timeout = 60000 * time.Millisecond
	dsName  = "sendOtpDS"
)

func TestNewDataSourceMappings(t *testing.T) {
	dsm := createDSMapping()
	assert.NotNil(t, dsm)
}

func TestRegisterDataSource_Success(t *testing.T) {
	dsm := createDSMapping()
	dsConf := getDataSourceConfig(timeout, method, dsName)
	ds := framework.CreateNewDataSource(&dsConf, nil)
	registered := dsm.RegisterDataSource(ds)
	assert.Equal(t, true, registered)
}

func TestGetDataSourceByName_Success(t *testing.T) {
	dsm := createDSMapping()
	dsConf := getDataSourceConfig(timeout, method, dsName)
	ds := framework.CreateNewDataSource(&dsConf, nil)
	_ = dsm.RegisterDataSource(ds)
	derivedDS := dsm.GetDataSourceByName(dsName)
	assert.NotNil(t, derivedDS)
}

func TestGetDataSourcesByNameMap_Success(t *testing.T) {
	dsm := createDSMapping()
	dsConf := getDataSourceConfig(timeout, method, dsName)
	ds := framework.CreateNewDataSource(&dsConf, nil)
	_ = dsm.RegisterDataSource(ds)
	m := dsm.GetDataSourcesByNameMap()
	assert.NotNil(t, m)
}

func createDSMapping() *DataSourceMappings {
	log := getLogger()
	dsm := NewDataSourceMappings(log)

	return dsm
}

func getDataSourceConfig(timeout time.Duration, method, dsName string) config.DataSourceConfig {
	return config.DataSourceConfig{URI: URI, Timeout: timeout, Method: method, DsName: dsName}
}

func getLogger() *logger.APILogger {
	cfg := getConfig()

	appLogger := logger.NewAPILogger(&cfg)
	appLogger.InitLogger()
	appLogger.Infof("Initialize Logger with LogLevel: %s", cfg.Logger.Level)

	return appLogger
}

func getConfig() config.Config {
	cfg := &config.Config{}
	return *cfg
}
