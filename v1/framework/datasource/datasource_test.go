package datasource

import (
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	frameworkModels "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/models/commons"
)

const (
	uri     = "http://localhost:8080/api/v1/auth/sendOtp"
	method  = "SendOtp()"
	timeout = 60000 * time.Millisecond
	dsName  = "verifyOtpDS"
)

func TestCreateNewDataSource_Success(t *testing.T) {
	dsConf := &frameworkModels.DataSourceConfig{URI: uri, Timeout: timeout, Method: method, DsName: dsName}
	ds := CreateNewDataSource(dsConf, nil, DummyFilter(), DummyFilter())
	assert.NotNil(t, ds)
}

func TestURI_Success(t *testing.T) {
	dsConf := getDataSourceConfig()
	ds := CreateNewDataSource(&dsConf, nil)
	u := ds.URI()
	assert.NotNil(t, u)
}

func TestName_Success(t *testing.T) {
	dsConf := getDataSourceConfig()
	ds := CreateNewDataSource(&dsConf, nil)
	n := ds.Name()
	assert.NotNil(t, n)
}

func TestTimeout_Success(t *testing.T) {
	dsConf := getDataSourceConfig()
	ds := CreateNewDataSource(&dsConf, nil)
	tout := ds.Timeout()
	assert.NotNil(t, tout)
}

func TestMethod_Success(t *testing.T) {
	dsConf := getDataSourceConfig()
	ds := CreateNewDataSource(&dsConf, nil)
	m := ds.Method()
	assert.NotNil(t, m)
}

func TestDataSource_GetFilters_Success(t *testing.T) {
	dsConf := getDataSourceConfig()
	ds := CreateNewDataSource(&dsConf, nil)
	f := ds.GetFilters()
	assert.NotNil(t, f)
}

func DummyFilter() Filter {
	return func(c echo.Context) error {
		return nil
	}
}

func getDataSourceConfig() frameworkModels.DataSourceConfig {
	return frameworkModels.DataSourceConfig{URI: uri, Timeout: timeout, Method: method, DsName: dsName}
}
