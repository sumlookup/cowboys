package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/sumlookup/cowboys/dao/db"
	"github.com/sumlookup/cowboys/mocks"
	pb "github.com/sumlookup/cowboys/pb"
	"go.uber.org/mock/gomock"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"testing"
	"time"
)

// setup all required data for tests
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func Test_Run_Success(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))

	engineMock.EXPECT().Run(gomock.Any()).Return(nil)

	svc := CowboysService{
		Engine: engineMock,
	}

	rsp, err := svc.Run(ctx, &pb.RunRequest{})

	assert.NoError(t, err)
	assert.Equal(t, &pb.RunResponse{}, rsp)
}

func Test_Run_Failed(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))

	engineMock.EXPECT().Run(gomock.Any()).Return(fmt.Errorf("more than one cowboy is required to simulate this shooter, currently found : 1"))

	svc := CowboysService{
		Engine: engineMock,
	}

	_, err := svc.Run(ctx, &pb.RunRequest{})
	assert.EqualError(t, err, "more than one cowboy is required to simulate this shooter, currently found : 1")
}

func Test_ReloadDefaultCowboys_No_Error(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))
	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	dbmocks.EXPECT().DeleteAllCowboys(gomock.Any()).Return(nil)

	dbmocks.EXPECT().CreateManyCowboys(gomock.Any(), getDefaultCowboys()).Return(int64(5), nil)

	svc := CowboysService{
		Engine: engineMock,
		Dao:    dbmocks,
	}
	os.Setenv("INPUTS_PATH", "../res/inputs.json")
	defer os.Clearenv()

	_, err := svc.ReloadDefaultCowboys(ctx, &pb.ReloadDefaultCowboysRequest{})

	assert.NoError(t, err)
}
func Test_ReloadDefaultCowboys_Delete_Failed(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))
	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	dbmocks.EXPECT().DeleteAllCowboys(gomock.Any()).Return(errors.New("connect: connection refused"))

	svc := CowboysService{
		Engine: engineMock,
		Dao:    dbmocks,
	}

	_, err := svc.ReloadDefaultCowboys(ctx, &pb.ReloadDefaultCowboysRequest{})
	assert.EqualError(t, err, "connect: connection refused")
}

func Test_ReloadDefaultCowboys_ReadFile_Failed(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))
	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	dbmocks.EXPECT().DeleteAllCowboys(gomock.Any()).Return(nil)

	svc := CowboysService{
		Engine: engineMock,
		Dao:    dbmocks,
	}

	_, err := svc.ReloadDefaultCowboys(ctx, &pb.ReloadDefaultCowboysRequest{})
	assert.EqualError(t, err, "open res/inputs.json: no such file or directory")
}

func Test_ReloadDefaultCowboys_Create_New_Cowboys_Failed(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))
	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	dbmocks.EXPECT().DeleteAllCowboys(gomock.Any()).Return(nil)

	dbmocks.EXPECT().CreateManyCowboys(gomock.Any(), getDefaultCowboys()).Return(int64(0), errors.New("connect: connection refused"))

	svc := CowboysService{
		Engine: engineMock,
		Dao:    dbmocks,
	}

	os.Setenv("INPUTS_PATH", "../res/inputs.json")
	defer os.Clearenv()

	_, err := svc.ReloadDefaultCowboys(ctx, &pb.ReloadDefaultCowboysRequest{})
	assert.EqualError(t, err, "connect: connection refused")
}

func Test_ReloadDefaultCowboys_Incorrect_Json_In_File(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))
	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	dbmocks.EXPECT().DeleteAllCowboys(gomock.Any()).Return(nil)

	svc := CowboysService{
		Engine: engineMock,
		Dao:    dbmocks,
	}
	writeDir := "../res/test_inputs.json"
	err := os.WriteFile(writeDir, []byte("{;DS{"), 0644)
	if err != nil {
		t.Errorf("failed to write file : %v", err)
	}

	os.Setenv("INPUTS_PATH", writeDir)
	defer os.Clearenv()

	_, err = svc.ReloadDefaultCowboys(ctx, &pb.ReloadDefaultCowboysRequest{})
	assert.EqualError(t, err, "invalid character ';' looking for beginning of object key string")
	os.Remove("../res/test_inputs.json")
}

func Test_ShootAtRandom_Success(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))
	shooterId := uuid.New()
	gameId := uuid.New()
	shooterName := "tester"
	shooterDmg := int32(1)

	engineMock.EXPECT().ShootRandomCowboy(gomock.Any(), shooterId, gameId, shooterName, shooterDmg).
		Return(int32(2), nil)
	svc := CowboysService{
		Engine: engineMock,
		//Dao:    dbmocks,
	}

	rsp, err := svc.ShootAtRandom(ctx, &pb.ShootAtRandomRequest{
		ShooterId:     shooterId.String(),
		GameId:        gameId.String(),
		ShooterName:   shooterName,
		ShooterDamage: shooterDmg,
	})
	assert.NoError(t, err)
	assert.Equal(t, int32(2), rsp.GetReceiverHealth())
}

func Test_ShootAtRandom_Failure(t *testing.T) {
	ctx := context.Background()

	engineMock := mocks.NewMockCowboysEngine(gomock.NewController(t))
	shooterId := uuid.New()
	gameId := uuid.New()
	shooterName := "tester"
	shooterDmg := int32(1)

	errMsg := errors.New("connect: connection refused")

	engineMock.EXPECT().ShootRandomCowboy(gomock.Any(), shooterId, gameId, shooterName, shooterDmg).
		Return(int32(0), errMsg)
	svc := CowboysService{
		Engine: engineMock,
	}

	_, err := svc.ShootAtRandom(ctx, &pb.ShootAtRandomRequest{
		ShooterId:     shooterId.String(),
		GameId:        gameId.String(),
		ShooterName:   shooterName,
		ShooterDamage: shooterDmg,
	})

	assert.Equal(t, errMsg, err)
}

func Test_ShootAtRandom_Incorrect_Shooter_Id(t *testing.T) {
	ctx := context.Background()

	svc := CowboysService{}

	_, err := svc.ShootAtRandom(ctx, &pb.ShootAtRandomRequest{
		ShooterId: "some_string",
	})
	assert.Equal(t, ErrorArgumentInvalidShooterId, err)
}

func Test_ShootAtRandom_Incorrect_Game_Id(t *testing.T) {
	ctx := context.Background()

	svc := CowboysService{}

	_, err := svc.ShootAtRandom(ctx, &pb.ShootAtRandomRequest{
		ShooterId: uuid.NewString(),
		GameId:    "some_string",
	})
	assert.Equal(t, ErrorArgumentInvalidGameId, err)
}

func Test_GetGameLogs_Success(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	gameId := uuid.New()
	gl := &db.GameLog{
		ID:             uuid.New(),
		CreatedAt:      pgtype.Timestamptz{Time: time.Now()},
		GameID:         uuid.New(),
		ShooterID:      uuid.New(),
		ReceiverID:     uuid.New(),
		Damage:         1,
		ReceiverHealth: 2,
		ShooterHealth:  2,
	}

	dbmocks.EXPECT().ListGameLogsByGameId(gomock.Any(), db.ListGameLogsByGameIdParams{
		GameID:      gameId,
		QuerySort:   "asc",
		QueryOffset: 0,
		QueryLimit:  1,
	}).Return([]*db.GameLog{
		gl,
	}, nil)

	dbmocks.EXPECT().CountAllGameLogs(gomock.Any(), gameId).Return(int64(1), nil)

	svc := CowboysService{
		Dao: dbmocks,
	}

	rsp, err := svc.GetGameLogs(ctx, &pb.GetGameLogsRequest{
		GameId: gameId.String(),
		Sort:   "asc",
		Limit:  1,
		Offset: 0,
	})

	assert.NoError(t, err)
	assert.Equal(t, gl.ID.String(), rsp.GameLogs[0].GetId())
	assert.Equal(t, timestamppb.New(gl.CreatedAt.Time), rsp.GameLogs[0].GetCreatedAt())
	assert.Equal(t, gl.ShooterID.String(), rsp.GameLogs[0].GetShooterId())
	assert.Equal(t, gl.ReceiverID.String(), rsp.GameLogs[0].GetReceiverId())
	assert.Equal(t, gl.Damage, rsp.GameLogs[0].GetDamage())
	assert.Equal(t, gl.ReceiverHealth, rsp.GameLogs[0].GetReceiverHealth())
	assert.Equal(t, gl.ShooterHealth, rsp.GameLogs[0].GetShooterHealth())

}

func Test_GetGameLogs_GameId_Incorrect(t *testing.T) {
	ctx := context.Background()

	svc := CowboysService{}

	_, err := svc.GetGameLogs(ctx, &pb.GetGameLogsRequest{
		GameId: "some_string",
	})

	assert.Equal(t, ErrorArgumentInvalidGameId, err)

}

func Test_GetGameLogs_ListLogs_Failure(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	gameId := uuid.New()

	errMsg := errors.New("connect: connection refused")

	dbmocks.EXPECT().ListGameLogsByGameId(gomock.Any(), db.ListGameLogsByGameIdParams{
		GameID:      gameId,
		QuerySort:   "asc",
		QueryOffset: 0,
		QueryLimit:  1,
	}).Return(nil, errMsg)

	svc := CowboysService{
		Dao: dbmocks,
	}

	_, err := svc.GetGameLogs(ctx, &pb.GetGameLogsRequest{
		GameId: gameId.String(),
		Sort:   "asc",
		Limit:  1,
		Offset: 0,
	})

	assert.Equal(t, errMsg, err)

}

func Test_GetGameLogs_CountLogs_Failure(t *testing.T) {
	ctx := context.Background()

	dbmocks := mocks.NewMockQuerier(gomock.NewController(t))

	gameId := uuid.New()
	gl := &db.GameLog{
		ID:             uuid.New(),
		CreatedAt:      pgtype.Timestamptz{Time: time.Now()},
		GameID:         uuid.New(),
		ShooterID:      uuid.New(),
		ReceiverID:     uuid.New(),
		Damage:         1,
		ReceiverHealth: 2,
		ShooterHealth:  2,
	}

	dbmocks.EXPECT().ListGameLogsByGameId(gomock.Any(), db.ListGameLogsByGameIdParams{
		GameID:      gameId,
		QuerySort:   "asc",
		QueryOffset: 0,
		QueryLimit:  1,
	}).Return([]*db.GameLog{
		gl,
	}, nil)

	errMsg := errors.New("connect: connection refused")

	dbmocks.EXPECT().CountAllGameLogs(gomock.Any(), gameId).Return(int64(0), errMsg)

	svc := CowboysService{
		Dao: dbmocks,
	}

	_, err := svc.GetGameLogs(ctx, &pb.GetGameLogsRequest{
		GameId: gameId.String(),
		Sort:   "asc",
		Limit:  1,
		Offset: 0,
	})

	assert.Equal(t, errMsg, err)

}

func Test_HealthCheck(t *testing.T) {
	ctx := context.Background()

	svc := CowboysService{}

	rsp, err := svc.Check(ctx, &healthpb.HealthCheckRequest{Service: "cowboys"})
	if err != nil {
		t.Logf("err")
	}

	assert.NoError(t, err)
	assert.Equal(t, "SERVING", rsp.Status.String())
}

func getDefaultCowboys() []db.CreateManyCowboysParams {
	return []db.CreateManyCowboysParams{
		{
			Name:   "John",
			Health: 10,
			Damage: 1,
		},
		{
			Name:   "Bill",
			Health: 8,
			Damage: 2,
		},
		{
			Name:   "Sam",
			Health: 10,
			Damage: 1,
		},
		{
			Name:   "Peter",
			Health: 5,
			Damage: 3,
		},
		{
			Name:   "Philip",
			Health: 15,
			Damage: 1,
		},
	}
}
