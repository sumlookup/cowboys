package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	daodb "github.com/sumlookup/cowboys/dao/db"
	pb "github.com/sumlookup/cowboys/pb"
	"os"
)

// Run runs the game simulation at current game mode
func (s *CowboysService) Run(ctx context.Context, req *pb.RunRequest) (*pb.RunResponse, error) {
	log := logrus.WithContext(ctx)

	err := s.Engine.Run(ctx)
	if err != nil {
		log.Errorf("failed while running game engine : %v", err)
		return nil, err
	}

	return &pb.RunResponse{}, nil
}

func (s *CowboysService) ReloadDefaultCowboys(ctx context.Context, req *pb.ReloadDefaultCowboysRequest) (*pb.ReloadDefaultCowboysResponse, error) {
	log := logrus.WithContext(ctx)
	_, err := s.Dao.DeleteAllCowboys(ctx)
	if err != nil {
		log.Errorf("ReloadDefaultCowboys : Failed to delete all cowboys : %v", err)
		return nil, err
	}

	b, err := os.ReadFile("./res/inputs.json")
	if err != nil {
		log.Errorf("ReloadDefaultCowboys : Failed while reading file - err : %v", err)
		return nil, err
	}

	cw := []daodb.CreateManyCowboysParams{}
	err = json.Unmarshal(b, &cw)
	if err != nil {
		log.Errorf("ReloadDefaultCowboys : Failed while unmarshaling defaults")
		return nil, err
	}
	cc := []*pb.Cowboy{}
	inserted := []*daodb.Cowboy{}
	rsp := s.Dao.CreateManyCowboys(ctx, cw)
	queryError := []string{}
	rsp.Query(func(i int, cowboys []*daodb.Cowboy, err error) {
		if err != nil {
			log.Errorf("ReloadDefaultCowboys : encountered an error while creating many cowboys : %v", err)
			queryError = append(queryError, err.Error())
		} else {
			inserted = append(inserted, cowboys...)
		}

	})

	for _, v := range inserted {
		cc = append(cc, &pb.Cowboy{
			Name:   v.Name,
			Health: v.Health,
			Damage: v.Damage,
		})
	}

	return &pb.ReloadDefaultCowboysResponse{Cowboys: cc}, nil

}

func (s *CowboysService) GetCowboyByName(ctx context.Context, req *pb.GetCowboyByNameRequest) (*pb.GetCowboyByNameResponse, error) {

	cw, err := s.Dao.GetSingleCowboyByName(ctx, req.GetName())
	if err != nil {
		return nil, err
	}
	return &pb.GetCowboyByNameResponse{
		Cowboy: &pb.Cowboy{
			Name:   cw.Name,
			Health: cw.Health,
			Damage: cw.Damage,
		},
	}, nil
}

// Test
func (s *CowboysService) Test(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	return &pb.TestResponse{}, nil
}

func (s *CowboysService) ShootAtRandom(ctx context.Context, req *pb.ShootAtRandomRequest) (*pb.ShootAtRandomResponse, error) {
	log := logrus.WithContext(ctx)
	shooterId, err := uuid.Parse(req.GetShooterId())
	if err != nil {
		log.Errorf("failed to parse shooter_id : %v", err)
		return nil, err
	}

	health, err := s.Engine.ShootRandomCowboy(ctx, shooterId, req.GetShooterName(), req.GetShooterDamage())
	if err != nil {
		log.Errorf("failed shooting at target : %v", err)
		return nil, err
	}

	return &pb.ShootAtRandomResponse{
		ReceiverHealth: health,
	}, nil
}
