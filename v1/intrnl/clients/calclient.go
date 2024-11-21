package clients

import (
	cal "github.com/Allen-Career-Institute/common-protos/cal/v1"
	"github.com/Allen-Career-Institute/common-protos/cal/v1/request"
	"github.com/Allen-Career-Institute/common-protos/cal/v1/response"
	calTypes "github.com/Allen-Career-Institute/common-protos/learning_material/v1/types/enums"
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"
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
