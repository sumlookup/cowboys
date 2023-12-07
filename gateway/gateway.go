package gateway

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	clientGw "github.com/sumlookup/cowboys/pb"
	"google.golang.org/grpc/metadata"
)

func NewGwServer(ctx context.Context, grpc grpc.ClientConnInterface, port string) *http.Server {
	log.Infof("port - %v", port)
	log.Infof("Starting HTTP API on %s. Env: %s", port, os.Getenv("ENV"))

	// create http server
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithMetadata(func(c context.Context, req *http.Request) metadata.MD {
			return metadata.Pairs("x-forwarded-method", req.Method)
		}),
	)

	err := clientGw.RegisterCowboysServiceHandlerClient(ctx, mux, clientGw.NewCowboysServiceClient(grpc))
	if err != nil {
		log.Fatalf("Could not register external service handler for client: %s", err.Error())
	}

	addr := fmt.Sprintf(":%s", port)
	hmux := http.NewServeMux()
	hmux.HandleFunc("/healthz", healthzServer()) // this is accessible only internally
	hmux.Handle("/", mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	return &http.Server{
		Addr:    addr,
		Handler: c.Handler(hmux),
	}
}

// healthzServer returns a simple health handler which returns ok.
func healthzServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "ok")
	}
}
