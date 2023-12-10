package handler

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

var (
	DEFAULT_PAGE_SIZE   = int32(100)
	DEFAULT_INPUTS_PATH = "res/inputs.json"

	// Handler
	ErrorArgumentInvalidShooterId  = status.Errorf(codes.InvalidArgument, "invalid shooter_id provided, should be valid uuid format.")
	ErrorArgumentInvalidReceiverId = status.Errorf(codes.InvalidArgument, "invalid receiver_id provided, should be valid uuid format.")
	ErrorArgumentInvalidGameId     = status.Errorf(codes.InvalidArgument, "invalid game_id provided, should be valid uuid format.")
)

func getPageSize(reqSize int32) int32 {
	if reqSize < 0 {
		return 0
	}

	if reqSize > 0 && reqSize < DEFAULT_PAGE_SIZE {
		return reqSize
	}

	return DEFAULT_PAGE_SIZE
}

func getOffset(reqSize int32) int32 {
	if reqSize < 0 {
		return 0
	}
	return reqSize
}

func getInputsResPath() string {
	op := os.Getenv("INPUTS_PATH")
	if op == "" {
		return DEFAULT_INPUTS_PATH
	}
	return op
}

func getSort(sort string) string {
	if sort != "desc" {
		return "asc"
	}
	return "desc"
}
