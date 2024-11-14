package dynamicconfig

import (
	"context"

	dc "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1"
	dconf "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1/configs"
	factory "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1/factory"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

func NewDynamicConfig(cfg *config.CommonConfig, l logger.Logger) (dc.DynamicConfig, error) {
	ctx := context.Background()
	dcFactory := factory.DynamicConfigHandler{}

	dynamicConfig, err := dcFactory.GetDynamicConfig(ctx, dc.APP_CONFIG)
	if err != nil {
		return nil, err
	}

	conf := dconf.Configuration{
		AppID:           cfg.AppConfig.AppID,
		ConfigName:      cfg.AppConfig.ConfigName,
		PollingInterval: cfg.AppConfig.PollingInterval,
	}

	err = dynamicConfig.Init(&conf)
	if err != nil {
		l.Errorf("error in initializing aws configs, %v", err)
		return nil, err
	}

	return dynamicConfig, nil
}
