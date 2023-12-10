package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	daodb "github.com/sumlookup/cowboys/dao/db"
	pb "github.com/sumlookup/cowboys/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
)

// Run runs the game simulation
// Be default it checks for available cowboys in DB and spawns a process(goroutine) for each one
func (s *CowboysService) Run(ctx context.Context, req *pb.RunRequest) (*pb.RunResponse, error) {
	log := logrus.WithContext(ctx)

	err := s.Engine.Run(ctx)
	if err != nil {
		log.Errorf("Run : failed while running game engine : %v", err)
		return nil, err
	}

	return &pb.RunResponse{}, nil
}

// ReloadDefaultCowboys Deletes any remaining cowboys
// Takes data from ./res/inputs.json and loads in DB
func (s *CowboysService) ReloadDefaultCowboys(ctx context.Context, req *pb.ReloadDefaultCowboysRequest) (*pb.ReloadDefaultCowboysResponse, error) {
	log := logrus.WithContext(ctx)
	log.Infof("Reloading to default cowboys")

	err := s.Dao.DeleteAllCowboys(ctx)
	if err != nil {
		log.Errorf("ReloadDefaultCowboys : Failed to delete all cowboys : %v", err)
		return nil, err
	}

	b, err := os.ReadFile(getInputsResPath())
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
	_, err = s.Dao.CreateManyCowboys(ctx, cw)
	if err != nil {
		log.Errorf("ReloadDefaultCowboys : encountered an error while creating many cowboys : %v", err)
		return nil, err
	}

	for _, v := range cw {
		cc = append(cc, &pb.Cowboy{
			Name:   v.Name,
			Health: v.Health,
			Damage: v.Damage,
		})
	}
	log.Infof("Reloading Successfull")

	return &pb.ReloadDefaultCowboysResponse{Cowboys: cc}, nil

}

// ShootAtRandom accepts shooter_id (cowboy id) and performs DB transaction to get a random alive cowboy and shoot at it
func (s *CowboysService) ShootAtRandom(ctx context.Context, req *pb.ShootAtRandomRequest) (*pb.ShootAtRandomResponse, error) {
	log := logrus.WithContext(ctx)
	shooterId, err := uuid.Parse(req.GetShooterId())
	if err != nil {
		log.Errorf("ShootAtRandom : failed to parse shooter_id : %v", err)
		return nil, ErrorArgumentInvalidShooterId
	}
	gameId, err := uuid.Parse(req.GetGameId())
	if err != nil {
		log.Errorf("ShootAtRandom : failed to parse game_id : %v", err)
		return nil, ErrorArgumentInvalidGameId
	}

	health, err := s.Engine.ShootRandomCowboy(ctx, shooterId, gameId, req.GetShooterName(), req.GetShooterDamage())
	if err != nil {
		log.Errorf("ShootAtRandom : failed shooting at target : %v", err)
		return nil, err
	}

	return &pb.ShootAtRandomResponse{
		ReceiverHealth: health,
	}, nil
}

// GetGameLogs if running intermediate game at the end you get a game_id
// By providing the game id you can check each log for each cowboy shot
func (s *CowboysService) GetGameLogs(ctx context.Context, req *pb.GetGameLogsRequest) (*pb.GetGameLogsResponse, error) {

	log := logrus.WithContext(ctx)

	gameId, err := uuid.Parse(req.GetGameId())
	if err != nil {
		log.Errorf("GetGameLogs : failed to parse game_id : %v", err)
		return nil, ErrorArgumentInvalidGameId
	}

	logs, err := s.Dao.ListGameLogsByGameId(ctx, daodb.ListGameLogsByGameIdParams{
		GameID:      gameId,
		QuerySort:   getSort(req.GetSort()),
		QueryLimit:  getPageSize(req.GetLimit()),
		QueryOffset: getOffset(req.GetOffset()) * getPageSize(req.GetLimit()),
	})
	if err != nil {
		log.Errorf("GetGameLogs : failed to list game logs : %v", err)
		return nil, err
	}

	count, err := s.Dao.CountAllGameLogs(ctx, gameId)
	if err != nil {
		log.Errorf("GetGameLogs : failed counting all game logs : %v", err)
		return nil, err
	}

	gl := []*pb.GameLog{}
	for _, l := range logs {
		gl = append(gl, &pb.GameLog{
			Id:             l.ID.String(),
			CreatedAt:      timestamppb.New(l.CreatedAt.Time),
			ShooterId:      l.ShooterID.String(),
			ReceiverId:     l.ReceiverID.String(),
			Damage:         l.Damage,
			ShooterHealth:  l.ShooterHealth,
			ReceiverHealth: l.ReceiverHealth,
		})
	}

	return &pb.GetGameLogsResponse{
		GameLogs:   gl,
		Page:       req.GetOffset(),
		TotalCount: int32(count),
	}, nil
}
