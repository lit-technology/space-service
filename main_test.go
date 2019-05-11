package main

import (
	"context"
	"testing"
	"time"

	pb "github.com/philip-bui/space-service/protos"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	target = "127.0.0.1:8000"
)

func TestMain(t *testing.M) {
	go func() {
		main()
	}()
	time.Sleep(1 * time.Second)
	t.Run()
}

func TestServer(t *testing.T) {
	suite.Run(t, new(MainSuite))
}

type MainSuite struct {
	suite.Suite
	pb.SparkClient
	grpc_health_v1.HealthClient
	context.Context
}

func (s *MainSuite) SetupSuite() {
	cc, err := grpc.Dial(target, grpc.WithInsecure())
	s.NoError(err)
	s.NotNil(cc)
	s.SparkClient = pb.NewSparkClient(cc)
	s.HealthClient = grpc_health_v1.NewHealthClient(cc)
	s.NoError(err)
	s.Context = context.Background()
}

func (s *MainSuite) TestMethods() {

}

func (s *MainSuite) TestHealthCheck() {
	s.NotNil(s.HealthClient)
	resp, err := s.HealthClient.Check(s.Context, &grpc_health_v1.HealthCheckRequest{})
	s.NoError(err)
	s.NotNil(resp)
}
