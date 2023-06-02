package cmd

import (
	"context"
	"errors"
	"log"

	"github.com/kostyaBroTutor/auth_interceptor/example/service"
	"github.com/kostyaBroTutor/auth_interceptor/interceptor"
	"github.com/kostyaBroTutor/auth_interceptor/pkg/process"
	"github.com/kostyaBroTutor/auth_interceptor/proto"
)

func main() {
	authClientInstance := authClient{
		users: map[string]*interceptor.TokenInfo{
			"123": {
				UserID: "123",
				Roles:  []proto.Role{proto.Role_ADMIN_ROLE},
				Permissions: []proto.Permission{
					proto.Permission_READ_SOMETHING_PERMISSION,
					proto.Permission_WRITE_SOMETHING_PERMISSION,
					proto.Permission_CHANGE_SOMETHING_PERMISSION,
				},
			},
			"564": {
				UserID: "564",
			},
		},
	}

	closeGRPCServer, err := service.NewExampleGrpcServer(
		service.Config{
			ListenAddr:           "localhost:3415",
			AuthClient:           authClientInstance,
			ExampleServiceServer: service.NewExampleService(),
		},
	)
	if err != nil {
		log.Println("can not to ru")
		return
	}

	process.WaitForTermination()

	closeGRPCServer()
}

type authClient struct {
	users map[string]*interceptor.TokenInfo
}

func (a authClient) Auth(
	ctx context.Context, token string,
) (*interceptor.TokenInfo, error) {
	if tokenInfo, exists := a.users[token]; exists {
		return tokenInfo, nil
	}

	return nil, errors.New("user not found")
}
