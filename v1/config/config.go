package config

import (
	"encoding/json"
	"fmt"
	commonmodels "github.com/Allen-Career-Institute/go-bff-commons/v1/internal/commons"
	dc "github.com/Allen-Career-Institute/go-kratos-commons/dynamicconfig/v1"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config App config struct
type Config struct {
	Server                          ServerConfig
	AppConfig                       AppConfig
	DynamicConfig                   dc.DynamicConfig
	Logger                          Logger
	GoPool                          GoPool
	DataSource                      any
	PlaylistFilenameConfig          PlaylistConfig
	LmmRedisSecretLocation          string
	AkamaiConfig                    AkamaiConfig
	WhiteListSubTypesForContentAuth WhiteListSubTypesForContentAuth
}

// ClientConfig Internal call config struct
type ClientConfig struct {
	Endpoint  string
	Timeout   time.Duration
	Conn      time.Duration
	Namespace string
}

type CircuitBreakerClientConfig struct {
	FailurePercentageThresholdWithinTimePeriod   uint          // failure percentage threshold before opening the circuit
	FailureMinExecutionThresholdWithinTimePeriod uint          // The number of executions must also exceed the failureExecutionThreshold within the failureThresholdingPeriod
	FailurePeriodThreshold                       time.Duration // failureThresholdingPeriod is the time period in which the failure rate is calculated
	SuccessThreshold                             uint          // number of successive successes before closing the circuit
	Delay                                        time.Duration // delay before retrying after a failure
}

type RetryClientConfig struct {
	MaxRetries int           // Number of retries
	Delay      time.Duration // delay interval between each retry
}

type PlaylistConfig struct {
	HLSFilename  string
	DASHFilename string
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// DataSourceConfig Data Source config struct
type GoPool struct {
	MaxConcurrentRoutines uint32
}

type SchedulingServiceConfig struct {
	FetchSchedulesDS        DataSourceConfig
	FetchSchedulesSummaryDS DataSourceConfig
}

type DataSourceConfig struct {
	URI       string
	Method    string
	Timeout   time.Duration
	DsName    string
	Resource  string
	Action    string
	PreloadDS []string
}

// ServerConfig Server config struct
type ServerConfig struct {
	Port                        string
	App                         App
	ReadTimeout                 time.Duration
	WriteTimeout                time.Duration
	JwtSecret                   string
	JwtSecretLocation           string
	AesEncryptionKey            string
	AesSecretIV                 string
	AesEncryptionSecretLocation string
}

type WhiteListSubTypesForContentAuth struct {
	SubTypes []string
}

type AkamaiConfig struct {
	BaseURL string
	Key     string
}

type App struct {
	Name    string
	Version string
}

type AppConfig struct {
	AppID           string
	ConfigName      string
	PollingInterval int64
}

// LoadConfig Load config file from given path
func LoadConfig() (*Config, error) {
	var config Config

	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath(getConfigDirectory())

	if err := loadSubConfig(v, &config, "server"); err != nil {
		return nil, err
	}

	if err := loadSubConfig(v, &config.DataSource, "datasource"); err != nil {
		return nil, err
	}

	if os.Getenv("ENV") != "" {
		setJwtSecretConfig(&config)
		setEncryptionKeys(&config)
	}

	return &config, nil
}

func setEncryptionKeys(config *Config) {
	aesEncryptionKey, aesSecretIV, err := readAesEncryptionSecrets(config.Server.AesEncryptionSecretLocation)
	if err != nil {
		fmt.Printf("Error reading EncryptionKeys %v", err)
	}

	if aesEncryptionKey != "" {
		config.Server.AesEncryptionKey = aesEncryptionKey

		log.Info("successfully retrieved aesEncryptionKey")
	}

	if aesSecretIV != "" {
		config.Server.AesSecretIV = aesSecretIV

		log.Info("successfully retrieved aesSecretIV")
	}
}

func readAesEncryptionSecrets(location string) (secretKey, encryptionIV string, err error) {
	byteValue, err := os.ReadFile(location)
	if err != nil {
		fmt.Printf("Error reading aesEncryptionSecrets %v", err)
		return "", "", err
	}

	var aesEncryptionSecrets commonmodels.AesEncryptionSecrets

	err = json.Unmarshal(byteValue, &aesEncryptionSecrets)
	if err != nil {
		fmt.Printf("error while unmarshalling %v", err)
		return "", "", err
	}

	return aesEncryptionSecrets.EncryptionSecretKey, aesEncryptionSecrets.EncryptionIv, nil
}

func getConfigDirectory() string {
	var dir string

	env := os.Getenv("ENV")
	log.Infof("Found Env as %v", env)

	if env != "" {
		dir = "/data/conf/" + env + "/"
		log.Infof("Found dir as %v", dir)
	} else {
		log.Infof("failed to get build environment, using local configs")

		dir = "./config/local/"
	}

	return dir
}

func loadSubConfig(v *viper.Viper, config interface{}, name string) error {
	v.SetConfigName(name)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("error in reading %s config: %w", name, err)
	}

	return parseConfig(config, v)
}

func parseConfig(configStruct interface{}, v *viper.Viper) error {
	err := v.Unmarshal(configStruct)
	if err != nil {
		log.Errorf("unable to decode into struct, %v", err.Error())
		return err
	}

	return nil
}

const (
	TimeoutSuffix     = ".timeout"
	EndPointSuffix    = ".endpoint"
	ConnTimeoutSuffix = ".conn"
	NamespaceSuffix   = ".namespace"
	TimeInMs          = "ms"
	TimeInSeconds     = "s"
	FileType          = "properties"
	LocalConfigName   = "client"
	JwtSecretKey      = "jwt_secret"
)

const (
	CircuitBreakerFailurePercentageThresholdSuffix       = ".cb.failure_percentage_threshold"
	CircuitBreakerMinExecutionThresholdSuffix            = ".cb.min_execution_threshold"
	CircuitBreakerFailurePeriodThresholdSuffix           = ".cb.failure_period_threshold"
	CircuitBreakerDelaySuffix                            = ".cb.delay"
	CircuitBreakerSuccessThresholdSuffix                 = ".cb.success_threshold"
	DefaultCircuitBreakerPercentageThreshold             = 50
	DefaultCircuitMinExecutionThreshold                  = 100
	DefaultCircuitBreakerFailurePeriodThresholdInSeconds = 60
	DefaultCircuitBreakerDelayMs                         = 30000
	DefaultCircuitBreakerSuccessThreshold                = 5
)

const (
	RetryMaxRetriesSuffix  = ".retry.max_retries"
	RetryDelaySuffix       = ".retry.delay"
	DefaultRetryMaxRetries = 3
	DefaultRetryDelayMs    = 1000
)

const (
	Base10    = 10
	BitSize32 = 32
)

const (
	ReadConfigErrorLog      = "error in reading server properties file config client: %v ,Error: %v"
	StringToIntParsingError = "error parsing %s for client: %v ,Error: %v, using defaults"
	DynConfigParsingError   = "error fetching %s aws config for client: %v ,Error: %v"
)

// nolint:gocognit //cannot be reduced further
func GetClientConfigs(client string, cnf *Config) ClientConfig {
	// Reading from server properties files
	if cnf.DynamicConfig == nil {
		v := viper.New()
		v.SetConfigType(FileType)
		v.AddConfigPath(getConfigDirectory())
		v.SetConfigName(LocalConfigName)

		if err := v.ReadInConfig(); err != nil {
			log.Errorf(ReadConfigErrorLog, client, err.Error())
		}

		timeoutConfig := v.GetString(join(client, TimeoutSuffix))

		timeout, err := time.ParseDuration(join(timeoutConfig, TimeInMs))
		if timeout == 0 || err != nil {
			timeout = 6000 * time.Millisecond
		}

		connConfig := v.GetString(join(client, ConnTimeoutSuffix))

		conn, err := time.ParseDuration(join(connConfig, TimeInMs))
		if conn == 0 || err != nil {
			conn = 6000 * time.Millisecond
		}

		return ClientConfig{
			Endpoint: v.GetString(join(client, EndPointSuffix)),
			Timeout:  timeout,
			Conn:     conn,
		}
	}

	// Reading from aws App Config
	connConfig, err := cnf.DynamicConfig.Get(client + ConnTimeoutSuffix)
	if err != nil {
		log.Errorf("error fetching conn aws config for client: %v ,Error: %v", client, err.Error())
	}

	conn, err := time.ParseDuration(connConfig + TimeInMs)
	if conn == 0 || err != nil {
		conn = 6000 * time.Millisecond
	}

	timeoutConfig, err := cnf.DynamicConfig.Get(client + TimeoutSuffix)
	if err != nil {
		log.Errorf("error fetching timeout aws config for client: %v ,Error: %v", client, err.Error())
	}

	timeout, err := time.ParseDuration(timeoutConfig + TimeInMs)
	if timeout == 0 || err != nil {
		timeout = 6000 * time.Millisecond
	}

	endpoint, err := cnf.DynamicConfig.Get(client + EndPointSuffix)
	if err != nil {
		log.Errorf("error fetching endpoint aws config for client: %v ,Error: %v", client, err.Error())
	}

	return ClientConfig{Endpoint: endpoint, Timeout: timeout, Conn: conn}
}

func GetCircuitBreakerClientConfigs(client string, cnf *Config) CircuitBreakerClientConfig {
	if cnf.DynamicConfig == nil {
		return readCircuitBreakerConfigFromLocalConfig(client)
	}

	return readCircuitBreakerConfigFromDynConfig(client, cnf)
}

func readCircuitBreakerConfigFromLocalConfig(client string) CircuitBreakerClientConfig {
	v := viper.New()
	v.SetConfigType(FileType)
	v.AddConfigPath(getConfigDirectory())
	v.SetConfigName(LocalConfigName)

	if err := v.ReadInConfig(); err != nil {
		log.Errorf(ReadConfigErrorLog, client, err.Error())
	}

	failurePercentageThresholdConfig := v.GetString(join(client, CircuitBreakerFailurePercentageThresholdSuffix))
	failurePercentageThresholdConfigInt, err := strconv.ParseUint(failurePercentageThresholdConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failurePercentageThresholdConfig", client, err)

		failurePercentageThresholdConfigInt = DefaultCircuitBreakerPercentageThreshold
	}

	failureMinExecutionThresholdConfig := v.GetString(join(client, CircuitBreakerMinExecutionThresholdSuffix))
	failureMinExecutionThresholdConfigInt, err := strconv.ParseUint(failureMinExecutionThresholdConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failureMinExecutionThresholdConfig", client, err)

		failureMinExecutionThresholdConfigInt = DefaultCircuitMinExecutionThreshold
	}

	failurePeriodThresholdConfig := v.GetString(join(client, CircuitBreakerFailurePeriodThresholdSuffix))
	failurePeriodThresholdConfigInDuration, err := time.ParseDuration(join(failurePeriodThresholdConfig, TimeInSeconds))

	if failurePeriodThresholdConfigInDuration == 0 || err != nil {
		log.Warnf(StringToIntParsingError, "failurePeriodThresholdConfig", client, err)

		failurePeriodThresholdConfigInDuration = DefaultCircuitBreakerFailurePeriodThresholdInSeconds * time.Second
	}

	successThresholdConfig := v.GetString(join(client, CircuitBreakerSuccessThresholdSuffix))
	successThresholdConfigInt, err := strconv.ParseUint(successThresholdConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "successThresholdConfig", client, err)

		successThresholdConfigInt = DefaultCircuitBreakerSuccessThreshold
	}

	delayConfig := v.GetString(join(client, CircuitBreakerDelaySuffix))
	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))

	if delayConfigInDuration == 0 || err != nil {
		delayConfigInDuration = DefaultCircuitBreakerDelayMs * time.Millisecond
	}

	return CircuitBreakerClientConfig{
		FailurePercentageThresholdWithinTimePeriod:   uint(failurePercentageThresholdConfigInt),
		FailureMinExecutionThresholdWithinTimePeriod: uint(failureMinExecutionThresholdConfigInt),
		FailurePeriodThreshold:                       failurePeriodThresholdConfigInDuration,
		SuccessThreshold:                             uint(successThresholdConfigInt),
		Delay:                                        delayConfigInDuration,
	}
}

func readCircuitBreakerConfigFromDynConfig(client string, cnf *Config) CircuitBreakerClientConfig {
	// Reading from aws App Config
	failurePercentageThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerFailurePercentageThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "failurePercentageThresholdConfig", client, err)
	}

	failurePercentageThresholdConfigInt, err := strconv.ParseUint(failurePercentageThresholdConfig, 10, 32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failurePercentageThresholdConfig", client, err)

		failurePercentageThresholdConfigInt = DefaultCircuitBreakerPercentageThreshold
	}

	failureMinExecutionThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerMinExecutionThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "failureMinExecutionThresholdConfig", client, err)
	}

	failureMinExecutionThresholdConfigInt, err := strconv.ParseUint(failureMinExecutionThresholdConfig, 10, 32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "failureMinExecutionThresholdConfig", client, err)

		failureMinExecutionThresholdConfigInt = DefaultCircuitMinExecutionThreshold
	}

	failurePeriodThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerFailurePeriodThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "failurePeriodThresholdConfig", client, err)
	}

	failurePeriodThresholdConfigInDuration, err := time.ParseDuration(join(failurePeriodThresholdConfig, TimeInSeconds))
	if failurePeriodThresholdConfigInDuration == 0 || err != nil {
		log.Warnf(StringToIntParsingError, "failurePeriodThresholdConfig", client, err)

		failurePeriodThresholdConfigInDuration = DefaultCircuitBreakerFailurePeriodThresholdInSeconds * time.Second
	}

	successThresholdConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerSuccessThresholdSuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "successThresholdConfig", client, err)
	}

	successThresholdConfigInt, err := strconv.ParseUint(successThresholdConfig, 10, 32)
	if err != nil {
		log.Warnf(StringToIntParsingError, "successThresholdConfig", client, err)

		successThresholdConfigInt = DefaultCircuitBreakerSuccessThreshold
	}

	delayConfig, err := cnf.DynamicConfig.Get(client + CircuitBreakerDelaySuffix)
	if err != nil {
		log.Warnf(DynConfigParsingError, "successThresholdConfig", client, err)
	}

	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))
	if delayConfigInDuration == 0 || err != nil {
		log.Warnf(StringToIntParsingError, "delayConfig", client, err)

		delayConfigInDuration = DefaultCircuitBreakerDelayMs * time.Millisecond
	}

	return CircuitBreakerClientConfig{
		FailurePercentageThresholdWithinTimePeriod:   uint(failurePercentageThresholdConfigInt),
		FailureMinExecutionThresholdWithinTimePeriod: uint(failureMinExecutionThresholdConfigInt),
		FailurePeriodThreshold:                       failurePeriodThresholdConfigInDuration,
		SuccessThreshold:                             uint(successThresholdConfigInt),
		Delay:                                        delayConfigInDuration,
	}
}

func GetRetryClientConfigs(client string, cnf *Config) RetryClientConfig {
	if cnf.DynamicConfig == nil {
		return readRetryConfigFromLocalConfig(client)
	}

	return readRetryFromDynConfig(client, cnf)
}

func readRetryConfigFromLocalConfig(client string) RetryClientConfig {
	v := viper.New()
	v.SetConfigType(FileType)
	v.AddConfigPath(getConfigDirectory())
	v.SetConfigName(LocalConfigName)

	if err := v.ReadInConfig(); err != nil {
		log.Errorf(ReadConfigErrorLog, client, err.Error())
	}

	maxRetriesConfig := v.GetString(join(client, RetryMaxRetriesSuffix))
	maxRetriesConfigInt, err := strconv.ParseInt(maxRetriesConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf("error parsing maxRetriesConfig for client: %v ,Error: %v, using defaults", client, err.Error())

		maxRetriesConfigInt = DefaultRetryMaxRetries
	}

	delayConfig := v.GetString(join(client, RetryDelaySuffix))
	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))

	if delayConfigInDuration == 0 || err != nil {
		log.Warnf("error parsing retry delayConfig for client: %v ,Error: %v, using defaults", client, err)

		delayConfigInDuration = DefaultRetryDelayMs * time.Millisecond
	}

	return RetryClientConfig{
		MaxRetries: int(maxRetriesConfigInt),
		Delay:      delayConfigInDuration,
	}
}

func readRetryFromDynConfig(client string, cnf *Config) RetryClientConfig {
	// Reading from aws App Config
	maxRetriesConfig, err := cnf.DynamicConfig.Get(client + RetryMaxRetriesSuffix)
	if err != nil {
		log.Warnf("error fetching failureThresholdConfig aws config for client: %v ,Error: %v", client, err.Error())
	}

	maxRetriesConfigInt, err := strconv.ParseInt(maxRetriesConfig, Base10, BitSize32)
	if err != nil {
		log.Warnf("error parsing maxRetriesConfig for client: %v ,Error: %v, using defaults", client, err.Error())

		maxRetriesConfigInt = DefaultRetryMaxRetries
	}

	delayConfig, err := cnf.DynamicConfig.Get(client + RetryDelaySuffix)
	if err != nil {
		log.Warnf("error fetching retry delayConfig aws config for client: %v ,Error: %v", client, err.Error())
	}

	delayConfigInDuration, err := time.ParseDuration(join(delayConfig, TimeInMs))
	if delayConfigInDuration == 0 || err != nil {
		log.Warnf("error parsing retry delayConfig for client: %v ,Error: %v, using defaults", client, err)

		delayConfigInDuration = DefaultRetryDelayMs * time.Millisecond
	}

	return RetryClientConfig{
		MaxRetries: int(maxRetriesConfigInt),
		Delay:      delayConfigInDuration,
	}
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}

	return sb.String()
}

func setJwtSecretConfig(config *Config) {
	jwtSecret, err := readJwtSecret(config.Server.JwtSecretLocation)
	if err != nil {
		log.Errorf("error reading jwt secret %v", err)
	}

	if jwtSecret != "" {
		config.Server.JwtSecret = jwtSecret

		log.Info("successfully retrieved jwt secret")
	}
}

func readJwtSecret(location string) (string, error) {
	var result map[string]interface{}

	byteValue, err := os.ReadFile(location)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return "", err
	}

	jwtSecret, ok := result[JwtSecretKey].(string)
	if !ok {
		return "", err
	}

	return jwtSecret, nil
}
