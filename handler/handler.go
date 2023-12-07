package handler

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	daodb "github.com/sumlookup/cowboys/dao/db"
	"github.com/sumlookup/cowboys/engine"
	pb "github.com/sumlookup/cowboys/pb"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type CowboysService struct {
	pb.UnimplementedCowboysServiceServer
	Dao    daodb.Querier
	Engine engine.CowboysEngine
	Db     *pgxpool.Pool
}

// Watch is required by the Healtcheck service
func (s *CowboysService) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

// Check is required by the Healtcheck service
func (s *CowboysService) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
