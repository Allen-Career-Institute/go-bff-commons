package utils

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"testing"

	v1 "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetFullLottieHeaderConfigValues(t *testing.T) {
	mockConfig := v1.MockDynamicConfig{}
	mockConfig.On("Get", FullLottieHeaderDarkURLAssetKey).Return("https://example.com/dark-lottie.json", nil)
	mockConfig.On("Get", FullLottieHeaderLightURLAssetKey).Return("https://example.com/light-lottie.json", nil)

	c := &config.Config{DynamicConfig: &mockConfig}

	result, err := GetFullLottieHeaderConfigValues(c)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://example.com/dark-lottie.json", result.FullLottieHeaderDarkURL)
	assert.Equal(t, "https://example.com/light-lottie.json", result.FullLottieHeaderLightURL)

	mockConfig.AssertExpectations(t)
}
