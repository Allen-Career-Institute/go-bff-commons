package dynamicconfig

import (
	"context"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1/factory"

	dc "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1"
	dconf "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1/configs"
)

func NewDynamicConfig(cfg *config.Config, l logger.Logger) (dc.DynamicConfig, error) {
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
