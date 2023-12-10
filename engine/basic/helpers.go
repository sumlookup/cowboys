package basic

import (
	"bytes"
	"context"
	"fmt"

	//"github.com/google/uuid"
	//"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func postHttpRequest(ctx context.Context, url, body string) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.WithContext(ctx).Errorf("can't prepare request to %s, %s", url, err.Error())
		return fmt.Errorf("can't prepare request to '%s': %v", url, err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("faile while making http post : %v", err)
		return fmt.Errorf("cannot make http request: %v", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Errorf("http request received statuscode = %v", response.StatusCode)
		return fmt.Errorf("http request received statuscode = %v", response.StatusCode)
	}

	return nil
}
