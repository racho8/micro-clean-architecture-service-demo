package composer

import (
	"context"
	"demo-service/common"
	"demo-service/proto/pb"
	sctx "github.com/viettranx/service-context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"github.com/hashicorp/consul/api"
)

type authClient struct {
	grpcAuthClient pb.AuthServiceClient
}

func (ac *authClient) IntrospectToken(ctx context.Context, accessToken string) (sub string, tid string, err error) {
	resp, err := ac.grpcAuthClient.IntrospectToken(ctx, &pb.IntrospectReq{AccessToken: accessToken})

	if err != nil {
		return "", "", err
	}

	return resp.Sub, resp.Tid, nil
}

// ComposeAuthRPCClient use only for middleware: get token info
func ComposeAuthRPCClient(serviceCtx sctx.ServiceContext) *authClient {
	configComp := serviceCtx.MustGet(common.KeyCompConf).(common.Config)

	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}

	service, _, err := consulClient.Health().Service("grpc-service", "", true, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(service) == 0 {
		log.Fatal("No healthy service instances found")
	}

	serviceAddress := service[0].Service.Address
	servicePort := service[0].Service.Port

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.Dial(fmt.Sprintf("%s:%d", serviceAddress, servicePort), opts)

	if err != nil {
		log.Fatal(err)
	}

	return &authClient{pb.NewAuthServiceClient(clientConn)}
}

func composeUserRPCClient(serviceCtx sctx.ServiceContext) pb.UserServiceClient {
	configComp := serviceCtx.MustGet(common.KeyCompConf).(common.Config)

	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err)
	}

	service, _, err := consulClient.Health().Service("grpc-service", "", true, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(service) == 0 {
		log.Fatal("No healthy service instances found")
	}

	serviceAddress := service[0].Service.Address
	servicePort := service[0].Service.Port

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.Dial(fmt.Sprintf("%s:%d", serviceAddress, servicePort), opts)

	if err != nil {
		log.Fatal(err)
	}

	return pb.NewUserServiceClient(clientConn)
}
