package Notification

import (
	"errors"
	"fmt"
	"net/http"
)

func SendWarning(num string) error {
	res, err := http.Get(fmt.Sprintf("http://localhost:3000/warn?number=%s", num))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}

type ReqBody struct {
	Numbers    []string `json:"numbers"`
	Hospital   string   `json:"hospital"`
	Addr       string   `json:"addr"`
	BloodGroup string   `json:"blood_group"`
}
