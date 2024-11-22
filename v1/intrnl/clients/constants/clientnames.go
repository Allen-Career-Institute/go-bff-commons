package clients

const (
	AuthenticationServiceClient   = "authenticationServiceConfig"
	PageServiceClient             = "pageServiceConfig"
	ResourceServiceClient         = "resourceServiceConfig"
	UserServiceClient             = "userServiceConfig"
	AuthorizationPdpServiceClient = "authorizationPdpConfig"
	CalServiceClient              = "calServiceConfig"
)

const UserClientConnError = "Error while creating grpc connection with user service: %v"
