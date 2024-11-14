package dynamicconfig

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"

	dc "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1"
	"github.com/labstack/echo/v4"
)

func GetOrDefaultFromAppConfig(dynamicConfig dc.DynamicConfig, key, defaultValue string) string {
	if dynamicConfig == nil {
		return defaultValue
	}

	url, err := dynamicConfig.Get(key)
	if err != nil {
		return defaultValue
	}

	if url != utils.EmptyString {
		return url
	}

	return defaultValue
}

func GetInterfaceOrDefaultFromAppConfig(log logger.Logger, dynamicConfig dc.DynamicConfig, key string, defaultValue interface{}) interface{} {
	if dynamicConfig == nil {
		return defaultValue
	}

	valInConf, err := dynamicConfig.GetAsInterface(key)
	if err != nil {
		log.Errorf("Error occurred while fetching dynamic-config for key %s, Err:- %v", key, err)
		return defaultValue
	}

	if valInConf != nil {
		log.Infof("Successfully fetched dynamic-config for key %s", key)
		return valInConf
	}

	return defaultValue
}

func GetFeatureFlag(ctx echo.Context, log logger.Logger, dynamicConfig dc.DynamicConfig, key string) bool {
	if dynamicConfig == nil {
		return false
	}

	valInConf, err := dynamicConfig.GetAsInterface(key)

	if err != nil {
		log.WithContext(ctx).Debugf("GetFeatureFlag: Error occurred while fetching dynamic-config for key %s, Err:- %v", key, err)
		return false
	}

	if valInConf != nil {
		if val, ok := valInConf.(bool); ok {
			return val
		}

		log.WithContext(ctx).Errorf("GetFeatureFlag: Error occurred while type casting dynamic-config for key %s, value %v", key, valInConf)
	}

	return false
}
