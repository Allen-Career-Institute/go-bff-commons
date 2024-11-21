// nolint: funlen,gocritic,gofmt // this package is required
package utils

import (
	"fmt"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type FullLottieHeaderConfig struct {
	FullLottieHeaderLightURL string
	FullLottieHeaderDarkURL  string
}

type PodcastConfig struct {
	ContentID                        string
	Title                            string
	Subtitle                         string
	Thumbnail                        string
	CarouselTitle                    string
	AutoPlay                         bool
	Muted                            bool
	ShowMinimalControlsUntilPrepared bool
	FullScreenEnabled                bool
	ShowSeekBar                      bool
	ShowMuteButton                   bool
	ShowPlaybackSpeed                bool
	ShowSeekBarTimer                 bool
	ShowBrightnessIndicator          bool
	PodcastCollector                 []PodcastCollector
}

type PodcastCollector struct {
	State            string
	Title            string
	UploadedDate     string
	PublishedDate    string
	WatchedDuration  string
	DisplayViewCount string
	ViewCount        int
	DisplayLikeCount string
	LikeCount        int
	ThumbnailDark    string
	ThumbnailLight   string
	ContentID        string
	VideoImageURL    string
	VideoURL         string
}

type LottieHeaderConfig struct {
	LottieHeaderLightURL string
	LottieHeaderDarkURL  string
}

const (
	LottieHeaderDarkURLAssetKey      = "assets.carecorner.lottie_header_dark_url"
	LottieHeaderLightURLAssetKey     = "assets.carecorner.lottie_header_light_url"
	FullLottieHeaderDarkURLAssetKey  = "assets.carecorner.full_lottie_header_dark_url"
	FullLottieHeaderLightURLAssetKey = "assets.carecorner.full_lottie_header_light_url"
)

const (
	PodcastAssetKey                          = "assets.lmm.podcast.contentID"
	TitleAssetKey                            = "assets.lmm.podcast.title"
	SubtitleAssetKey                         = "assets.lmm.podcast.subtitle"
	ThumbnailAssetKey                        = "assets.lmm.podcast.thumbnail"
	AutoPlayAssetKey                         = "assets.lmm.podcast.autoPlay"
	MutedAssetKey                            = "assets.lmm.podcast.muted"
	ShowMinimalControlsUntilPreparedAssetKey = "assets.lmm.podcast.showMinimalControlsUntilPrepared"
	FullScreenEnabledAssetKey                = "assets.lmm.podcast.fullScreenEnabled"
	ShowSeekBarAssetKey                      = "assets.lmm.podcast.showSeekBar"
	ShowMuteButtonAssetKey                   = "assets.lmm.podcast.showMuteButton"
	ShowPlaybackSpeedAssetKey                = "assets.lmm.podcast.showPlaybackSpeed"
	ShowSeekBarTimerAssetKey                 = "assets.lmm.podcast.showSeekBarTimer"
	ShowBrightnessIndicatorAssetKey          = "assets.lmm.podcast.showBrightnessIndicator"
	PodcastDataAssetKey                      = "assets.lmm.podcast.collection"
	PodcastPrefixAssetKey                    = "assets.lmm.podcast"
	DateLayout                               = "2006-01-02"
	UpcomingState                            = "UPCOMING"
	StateAssetKey                            = "assets.lmm.podcast.state"
	PodcastTitleAssetKey                     = "assets.lmm.podcast.title"
	UploadedDateAssetKey                     = "assets.lmm.podcast.uploaded_date"
	PublishedDateAssetKey                    = "assets.lmm.podcast.published_date"
	WatchedDurationAssetKey                  = "assets.lmm.podcast.watched_duration"
	DisplayViewCountAssetKey                 = "assets.lmm.podcast.display_view_count"
	ViewCountAssetKey                        = "assets.lmm.podcast.view_count"
	DisplayLikeCountAssetKey                 = "assets.lmm.podcast.display_like_count"
	LikeCountAssetKey                        = "assets.lmm.podcast.like_count"
	ThumbnailDarkAssetKey                    = "assets.lmm.podcast.thumbnail_dark"
	ThumbnailLightAssetKey                   = "assets.lmm.podcast.thumbnail_light"
	ContentIDAssetKey                        = "assets.lmm.podcast.content_id"
	VideoImageURLAssetKey                    = "assets.lmm.podcast.video_image_url"
	CarouselTitleAssetKey                    = "assets.lmm.podcast.carousel_title"
)

func GetConfigValues(cnf *config.Config) (PodcastConfig, error) {
	var podcastConfig PodcastConfig

	contentID, err := cnf.DynamicConfig.Get(PodcastAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.ContentID = contentID

	title, err := cnf.DynamicConfig.Get(TitleAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.Title = title

	subtitle, err := cnf.DynamicConfig.Get(SubtitleAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.Subtitle = subtitle

	thumbnail, err := cnf.DynamicConfig.Get(ThumbnailAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.Thumbnail = thumbnail

	autoPlayStr, err := cnf.DynamicConfig.Get(AutoPlayAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	autoPlay, err := strconv.ParseBool(autoPlayStr)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.AutoPlay = autoPlay

	return podcastConfig, nil
}

// nolint:funlen // TODO: check by resp. team if it can be reduced further
func GetCarousalConfigValues(cnf *config.Config) (PodcastConfig, error) {
	var podcastConfig PodcastConfig

	// AutoPlay
	autoPlay, err := GetBoolValue(cnf, AutoPlayAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.AutoPlay = autoPlay

	// Muted
	muted, err := GetBoolValue(cnf, MutedAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.Muted = muted

	// ShowMinimalThumbnailControls
	showMinimalControlsUntilPrepared, err := GetBoolValue(cnf, ShowMinimalControlsUntilPreparedAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.ShowMinimalControlsUntilPrepared = showMinimalControlsUntilPrepared

	// FullScreenEnabled
	fullScreenEnabled, err := GetBoolValue(cnf, FullScreenEnabledAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.FullScreenEnabled = fullScreenEnabled

	// ShowSeekBar
	showSeekBar, err := GetBoolValue(cnf, ShowSeekBarAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.ShowSeekBar = showSeekBar

	// ShowMuteButton
	showMuteButton, err := GetBoolValue(cnf, ShowMuteButtonAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.ShowMuteButton = showMuteButton

	// ShowPlaybackSpeed
	showPlaybackSpeed, err := GetBoolValue(cnf, ShowPlaybackSpeedAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.ShowPlaybackSpeed = showPlaybackSpeed

	// ShowSeekBarTimer
	showSeekBarTimer, err := GetBoolValue(cnf, ShowSeekBarTimerAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.ShowSeekBarTimer = showSeekBarTimer

	// ShowBrightnessIndicator
	showBrightnessIndicator, err := GetBoolValue(cnf, ShowBrightnessIndicatorAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.ShowBrightnessIndicator = showBrightnessIndicator

	// Carousel Title
	carouselTitle, err := cnf.DynamicConfig.Get(CarouselTitleAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.CarouselTitle = carouselTitle

	// PodcastCollector
	podcastData, err := cnf.DynamicConfig.GetAsInterface(PodcastDataAssetKey)
	if err != nil {
		return podcastConfig, err
	}

	collectors, err := mapPodcastCollectors(podcastData)
	if err != nil {
		return podcastConfig, err
	}

	podcastConfig.PodcastCollector = collectors

	return podcastConfig, nil
}

// nolint: funlen,gocritic // Owner to reduce number of lines
func mapPodcastCollectors(data interface{}) ([]PodcastCollector, error) {
	collectors := make([]PodcastCollector, 0)

	// Convert the interface{} to []interface{}
	dataSlice, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data type: %T", data)
	}

	// Convert []interface{} to []map[string]interface{}
	mapSlice := make([]map[string]interface{}, 0, len(dataSlice))

	for _, item := range dataSlice {
		mapItem, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid data type in slice: %T", item)
		}

		mapSlice = append(mapSlice, mapItem)
	}

	// Iterate over the data and create PodcastCollector structs
	for _, item := range mapSlice {
		collector := PodcastCollector{
			State:            getStringValue(item, StateAssetKey),
			Title:            getStringValue(item, PodcastTitleAssetKey),
			UploadedDate:     getStringValue(item, UploadedDateAssetKey),
			PublishedDate:    getStringValue(item, PublishedDateAssetKey),
			WatchedDuration:  getStringValue(item, WatchedDurationAssetKey),
			DisplayViewCount: getStringValue(item, DisplayViewCountAssetKey),
			ViewCount:        getIntValue(item, ViewCountAssetKey),
			DisplayLikeCount: getStringValue(item, DisplayLikeCountAssetKey),
			LikeCount:        getIntValue(item, LikeCountAssetKey),
			ThumbnailDark:    getStringValue(item, ThumbnailDarkAssetKey),
			ThumbnailLight:   getStringValue(item, ThumbnailLightAssetKey),
			ContentID:        getStringValue(item, ContentIDAssetKey),
			VideoImageURL:    getStringValue(item, VideoImageURLAssetKey),
		}
		collectors = append(collectors, collector)
	}

	// Sort the collectors
	sort.Slice(collectors, func(i, j int) bool {
		// Move the "UPCOMING" state to the top
		state1 := collectors[i].State
		if state1 == UpcomingState {
			return true
		}

		state2 := collectors[j].State
		if state2 == UpcomingState {
			return false
		}

		// Parse the uploaded dates
		iDate, err := time.Parse(DateLayout, collectors[i].UploadedDate)
		if err != nil {
			return false
		}

		jDate, err := time.Parse(DateLayout, collectors[j].UploadedDate)
		if err != nil {
			return false
		}

		// Sort by uploaded date in descending order
		return iDate.After(jDate)
	})

	return collectors, nil
}

func getIntValue(data map[string]interface{}, key string) int {
	value, ok := data[key]
	if !ok {
		return 0
	}

	intValue, err := strconv.Atoi(strings.ReplaceAll(value.(string), ",", ""))
	if err != nil {
		return 0
	}

	return intValue
}
func getStringValue(m map[string]interface{}, key string) string {
	value, ok := m[key]
	if !ok {
		return ""
	}

	stringValue, ok := value.(string)
	if !ok {
		return ""
	}

	return stringValue
}

func GetBoolValue(cnf *config.Config, key string) (bool, error) {
	strValue, err := cnf.DynamicConfig.Get(key)
	if err != nil {
		return false, err
	}

	boolValue, err := strconv.ParseBool(strValue)

	if err != nil {
		return false, err
	}

	return boolValue, nil
}

func GetLottieHeaderConfigValues(cnf *config.Config) (*LottieHeaderConfig, error) {
	lottieHeaderDarkURL, err := cnf.DynamicConfig.Get(LottieHeaderDarkURLAssetKey)
	if err != nil {
		return nil, err
	}

	lottieHeaderLightURL, err := cnf.DynamicConfig.Get(LottieHeaderLightURLAssetKey)
	if err != nil {
		return nil, err
	}

	lottieHeaderConfig := &LottieHeaderConfig{
		LottieHeaderDarkURL:  lottieHeaderDarkURL,
		LottieHeaderLightURL: lottieHeaderLightURL,
	}

	return lottieHeaderConfig, nil
}

func FetchConfigStringValues(cnf *config.Config, key string) ([]string, error) {
	configValues, err := cnf.DynamicConfig.GetAsInterface(key)
	if err != nil {
		log.Errorf("error occurred while fetching config values for key: %s, err: %+v", key, err)
		return nil, err
	}
	var allowedValues []string
	if configValues != nil {
		for _, excludedID := range configValues.([]interface{}) {
			if excludedID != nil {
				excludedIDStr, ok := excludedID.(string)
				if ok {
					allowedValues = append(allowedValues, excludedIDStr)
				}
			}
		}
	}
	return allowedValues, nil
}

func GetFullLottieHeaderConfigValues(cnf *config.Config) (*FullLottieHeaderConfig, error) {
	lottieHeaderDarkURL, err := cnf.DynamicConfig.Get(FullLottieHeaderDarkURLAssetKey)
	if err != nil {
		return nil, err
	}

	lottieHeaderLightURL, err := cnf.DynamicConfig.Get(FullLottieHeaderLightURLAssetKey)
	if err != nil {
		return nil, err
	}

	lottieHeaderConfig := &FullLottieHeaderConfig{
		FullLottieHeaderDarkURL:  lottieHeaderDarkURL,
		FullLottieHeaderLightURL: lottieHeaderLightURL,
	}

	return lottieHeaderConfig, nil
}
