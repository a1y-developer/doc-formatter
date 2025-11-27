package options

import (
	"net"
	"strconv"

	authpb "github.com/a1y/doc-formatter/api/grpc/auth/v1"
	"github.com/a1y/doc-formatter/internal/auth"
	"github.com/a1y/doc-formatter/internal/auth/handler"
	"github.com/a1y/doc-formatter/internal/auth/infra/persistence"
	"github.com/a1y/doc-formatter/internal/auth/manager/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"k8s.io/kubectl/pkg/util/i18n"
)

type AuthOptions struct {
	Port     int
	Database DatabaseOptions
}

func NewAuthOptions() *AuthOptions {
	return &AuthOptions{
		Port:     DefaultPort,
		Database: DatabaseOptions{},
	}
}

func (o *AuthOptions) Complete(args []string) {}

func (o *AuthOptions) Validate() error {
	return nil
}

func (o *AuthOptions) Config() (*auth.Config, error) {
	cfg := auth.NewConfig()
	if err := o.Database.ApplyTo(&cfg.DB); err != nil {
		return nil, err
	}
	cfg.Port = o.Port
	return cfg, nil
}

func (o *AuthOptions) AddFlags(cmd *cobra.Command) {
	port, err := strconv.Atoi(PortEnv)
	if err != nil {
		port = DefaultPort
	}
	cmd.Flags().IntVarP(&o.Port, "port", "p", port,
		i18n.T("specify the port for the auth service to listen on"))
	o.Database.AddFlags(cmd.Flags())
}

func (o *AuthOptions) Run() error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	userRepository := persistence.NewUserRepository(config.DB)
	userManager := user.NewUserManager(userRepository)
	authHandler, err := handler.NewHandler(userManager)
	if err != nil {
		return err
	}

	lis, _ := net.Listen("tcp", ":"+strconv.Itoa(config.Port))
	server := grpc.NewServer()
	authpb.RegisterAuthServiceServer(server, authHandler)

	logrus.Infof("Auth service running at :%d", config.Port)

	err = server.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
