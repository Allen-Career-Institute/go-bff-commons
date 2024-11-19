package clients

import (
	"errors"
	cal "github.com/Allen-Career-Institute/common-protos/cal/v1"
	"github.com/Allen-Career-Institute/common-protos/cal/v1/request"
	"github.com/Allen-Career-Institute/common-protos/cal/v1/response"
	calTypes "github.com/Allen-Career-Institute/common-protos/learning_material/v1/types/enums"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/dynamicconfig"
	"github.com/labstack/echo/v4"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
)

const (
	SourceAllen   = "ALLEN"
	SourceTencent = "TENCENT"
)

func (cm *ClientManager) getCalClient(c echo.Context) (cal.CalClient, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.CalServiceClient, cm.cfg)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Got error while initializing cal client %v", err)
		return nil, err
	}

	client := cal.NewCalClient(conn)

	return client, nil
}

func (cm *ClientManager) getFlashcardsClient(c echo.Context) (cal.FlashcardsClient, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.CalServiceClient, cm.cfg)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Got error while initializing flashcards client %v", err)
		return nil, err
	}

	client := cal.NewFlashcardsClient(conn)
	return client, nil
}

func (cm *ClientManager) GetPlaylistFileRequest(c echo.Context, cnf *config.Config, contentID string) *request.GetVideoPlaylistFileRequest {
	return &request.GetVideoPlaylistFileRequest{
		Id:       contentID,
		Codec:    calTypes.Codec_x264,
		Filename: cm.getPlaylistFileName(c, cnf),
	}
}

func (cm *ClientManager) getPlaylistFileName(c echo.Context, cnf *config.Config) string {
	var filename string

	deviceType := c.Request().Header.Get(utils.DeviceType)

	switch deviceType {
	case utils.DeviceTypeWeb, utils.DeviceTypeiOS:
		filename = cnf.PlaylistFilenameConfig.HLSFilename
	default:
		filename = cnf.PlaylistFilenameConfig.DASHFilename
	}

	return filename
}

func (cm *ClientManager) GetNACPlaylistFile(c echo.Context, cnf *config.Config, contentID string) (*response.GetNACVideoPlaylistFileResponse, error) {
	calClient, err := cm.getCalClient(c)
	if err != nil {
		return nil, err
	}

	req := &request.GetNACVideoPlaylistFileRequest{
		Id:       contentID,
		Codec:    calTypes.Codec_x264,
		Filename: cm.getPlaylistFileName(c, cnf),
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	req.Source = SourceAllen

	resp, err := calClient.GetNACVideoPlaylistFile(apiCtx, req)
	if err == nil {
		return resp, nil
	}

	cm.logger.WithContext(c).Errorf("error while getting playlist file: %s from cal service for source %s codec %s, error is %v", req.Id, req.Source, req.Codec, err)
	cm.logger.WithContext(c).Infof("falling back to allen x264")

	req.Source = SourceTencent

	resp, err = calClient.GetNACVideoPlaylistFile(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("error while getting playlist file: %s from cal service for source %s codec %s, error is %v", req.Id, req.Source, req.Codec, err)
		return nil, err
	}

	return resp, nil
}

func (cm *ClientManager) GetPlayList(c echo.Context, cnf *config.Config, req *request.GetVideoPlaylistFileRequest) (*response.GetVideoPlaylistFileResponse, error) {
	calClient, err := cm.getCalClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)
	playlistSourceDC, ok := dynamicconfig.GetInterfaceOrDefaultFromAppConfig(cm.logger, cm.cfg.DynamicConfig, "playlistSourceSequence", []string{SourceTencent, SourceAllen}).([]interface{})
	var playlistSourceSequence []string
	if !ok {
		cm.logger.WithContext(c).Errorf("failed to unmarshall playlistSourceSequence from dynamic config")
		playlistSourceSequence = []string{SourceTencent, SourceAllen}
	}

	for _, playlistSource := range playlistSourceDC {
		playlistSourceSequence = append(playlistSourceSequence, playlistSource.(string))
	}

	if len(playlistSourceSequence) == 0 {
		playlistSourceSequence = []string{SourceTencent, SourceAllen}
	}

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	req.Source = playlistSourceSequence[0]

	resp, err := calClient.GetVideoPlaylistFile(apiCtx, req)
	if err == nil {
		return resp, nil
	}

	cm.logger.WithContext(c).Errorf("error while getting playlist file: %s from cal service for source %s codec %s, error is %v", req.Id, req.Source, req.Codec, err)
	cm.logger.WithContext(c).Infof("falling back to allen x264")

	req.Source = playlistSourceSequence[1]

	resp, err = calClient.GetVideoPlaylistFile(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("error while getting playlist file: %s from cal service for source %s codec %s, error is %v", req.Id, req.Source, req.Codec, err)
		return nil, err
	}

	return resp, nil
}

func (cm *ClientManager) GetLearningMaterial(c echo.Context, cnf *config.Config, req *request.GetLearningMaterialRequest) (*response.GetLearningMaterialResponse, error) {
	calClient, err := cm.getCalClient(c)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := calClient.GetLearningMaterial(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return resp, nil
}

func (cm *ClientManager) FetchLearningContent(c echo.Context, cnf *config.Config, contentIds []string) ([]*response.GetLearningMaterialResponse, error) {
	// once they have bulk get call, we will use that. This is temporary
	result := make([]*response.GetLearningMaterialResponse, 0, len(contentIds))

	if len(contentIds) == 0 {
		return result, nil
	}

	for val := range contentIds {
		in := request.GetLearningMaterialRequest{
			Id:                   contentIds[val],
			PresignedUrlTtlInMin: 120,
		}

		contentResult, err := cm.GetLearningMaterial(c, cnf, &in)
		if err != nil {
			return nil, err
		}

		result = append(result, contentResult)
	}

	return result, nil
}

func (cm *ClientManager) GetNAC(c echo.Context, cnf *config.Config, req *request.GetNACRequest) (*response.GetNACResponse, error) {
	calClient, err := cm.getCalClient(c)

	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := calClient.GetNAC(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return resp, nil
}

func (cm *ClientManager) DownloadLearningMaterial(c echo.Context, cnf *config.Config, req *request.DownloadLearningMaterialRequest) (*response.DownloadLearningMaterialResponse, error) {
	calClient, err := cm.getCalClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := calClient.DownloadLearningMaterial(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("error while downloading learning material: %s from cal service, error is %v", req.Id, err)
		return nil, err
	}

	return resp, nil
}

func (cm *ClientManager) GetDownloadPresignedURL(c echo.Context, cnf *config.Config, req *request.GetDownloadPresignedURLRequest) (*response.GetDownloadPresignedURLResponse, error) {
	calClient, err := cm.getCalClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := calClient.GetDownloadPresignedURL(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while getting presigned url for ID : %v from cal service, error : %v", req.Id, err)
		return nil, err
	}

	return resp, nil
}

func (cm *ClientManager) GetBulkDownloadMaterial(c echo.Context, cnf *config.Config, req *request.GetBulkDownloadMaterialRequest) (*response.GetBulkDownloadMaterialResponse, error) {
	calClient, err := cm.getCalClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := calClient.GetBulkDownloadMaterial(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while getting Bulk Download Material from cal service, error : %v", err)
		return nil, err
	}

	return resp, nil
}

func (cm *ClientManager) CreateFlashcardsSession(c echo.Context, cnf *config.Config, req *request.CreateFlashcardSessionRequest) (*response.CreateFlashcardSessionResponse, error) {
	cm.logger.WithContext(c).Infof("Creating flashcard session for request %s", req.String())
	flashcardsClient, err := cm.getFlashcardsClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := flashcardsClient.CreateFlashcardSession(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while creating flashcard session for req: %v, error : %v",
			req.String(), err)
		return nil, err
	}
	if resp == nil {
		cm.logger.WithContext(c).Errorf("CreateFlashcardsSession response is nil")
		return nil, errors.New("nil response")
	}
	cm.logger.WithContext(c).Infof("CreateFlashcardsSession response %s", resp.String())
	return resp, nil
}

func (cm *ClientManager) GetFlashcards(c echo.Context, cnf *config.Config, req *request.GetFlashcardsRequest) (*response.GetFlashcardsResponse, error) {
	cm.logger.WithContext(c).Infof("Getting flashcards for request %s", req.String())
	flashcardsClient, err := cm.getFlashcardsClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := flashcardsClient.GetFlashcards(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while getting flashcards for req: %v, error : %v",
			req.String(), err)
		return nil, err
	}
	if resp == nil {
		cm.logger.WithContext(c).Errorf("GetFlashcards response is nil")
	}
	cm.logger.WithContext(c).Infof("GetFlashcards response %s", resp.String())
	return resp, nil
}

func (cm *ClientManager) GetFlashcardsCount(c echo.Context, cnf *config.Config, req *request.CreateFlashcardSessionRequest) (*response.GetFlashcardsCountResponse, error) {
	cm.logger.WithContext(c).Infof("Getting flashcards count for request %s", req.String())
	flashcardsClient, err := cm.getFlashcardsClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := flashcardsClient.GetFlashcardsCount(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while getting flashcards count for req: %v, error : %v",
			req.String(), err)
		return nil, err
	}
	if resp == nil {
		cm.logger.WithContext(c).Errorf("GetFlashcardsCount response is nil")
		return nil, errors.New("nil response")
	}
	cm.logger.WithContext(c).Infof("GetFlashcardsCount response %s", resp.String())
	return resp, nil
}

func (cm *ClientManager) GetFlashcardsSessionStats(c echo.Context, cnf *config.Config, req *request.GetFlashcardSessionStatsRequest) (*response.GetFlashcardSessionStatsResponse, error) {
	cm.logger.WithContext(c).Infof("Getting flashcards session stats for request %s", req.String())
	flashcardsClient, err := cm.getFlashcardsClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.CalServiceClient, cnf)

	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)
	defer apiCancel()

	resp, err := flashcardsClient.GetFlashcardSessionStats(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Errorf("Error while getting flashcards session stats for req: %v, error : %v",
			req.String(), err)
		return nil, err
	}
	if resp == nil {
		cm.logger.WithContext(c).Errorf("GetFlashcardsSessionStats response is nil")
		return nil, errors.New("nil response")
	}
	cm.logger.WithContext(c).Infof("GetFlashcardsSessionStats response %s", resp.String())
	return resp, nil
}
