package Notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func SendBulkMessage(body ReqBody) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("http://localhost:3000/"), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}

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
	Users      []User `json:"users"`
	Hospital   string `json:"hospital"`
	Addr       string `json:"addr"`
	Type       int    `json:"type"`
	BloodGroup string `json:"blood_group"`
	Unit       int    `json:"unit"`
}

type User struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}
