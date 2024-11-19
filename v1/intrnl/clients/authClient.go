package clients

import (
	clients "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients/constants"

	auth "github.com/Allen-Career-Institute/common-protos/authentication/v1"
	authReq "github.com/Allen-Career-Institute/common-protos/authentication/v1/request"
	authRes "github.com/Allen-Career-Institute/common-protos/authentication/v1/response"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/labstack/echo/v4"
)

func (cm *ClientManager) GetAuthClient(c echo.Context) (auth.AuthenticationClient, error) {
	conn, err := cm.Grpc.GetConn(c, cm.logger, clients.AuthenticationServiceClient, cm.cfg)
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthenticationClient(conn)

	return client, nil
}

func (cm *ClientManager) RefreshToken(c echo.Context, cnf *config.Config, req *authReq.RefreshTokenRequest) (*authRes.RefreshTokenResponse, error) {
	authClient, err := cm.GetAuthClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.AuthenticationServiceClient, cnf)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	sores, err := authClient.RefreshTokenAuthentication(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return sores, nil
}

func (cm *ClientManager) Logout(c echo.Context, cnf *config.Config, req *authReq.LogoutRequest) (*authRes.LogoutResponse, error) {
	authClient, err := cm.GetAuthClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.AuthenticationServiceClient, cnf)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	sores, err := authClient.Logout(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return sores, nil
}

func (cm *ClientManager) VerifyOTP(c echo.Context, cnf *config.Config, req *authReq.VerifyOtpRequest) (*authRes.LoginResponse, error) {
	authClient, err := cm.GetAuthClient(c)
	if err != nil {
		return nil, err
	}

	conf := config.GetClientConfigs(clients.AuthenticationServiceClient, cnf)
	apiCtx, apiCancel := utils.GetRequestCtxWithTimeout(c, conf.Timeout)

	defer apiCancel()

	sores, err := authClient.VerifyOtp(apiCtx, req)
	if err != nil {
		cm.logger.WithContext(c).Error(err)
		return nil, err
	}

	return sores, nil
}
