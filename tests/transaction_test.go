package tests

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io/ioutil"
	"log"
	"net/http"
	"simple-banking-app/internal/dtos"
	"strings"
	"testing"
)

var (
	baseUrl   = "http://localhost:3000"
	authToken string
)

func init() {
	setup()
}

func setup() {
	payload := fmt.Sprintf(`{"email":"sam@mail.com", "password":"password1"}`)
	res, _, err := runRequest("/auth/login", "POST", payload, map[string]string{}, "")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var loginResponse struct {
		Data dtos.LoginResp `json:"data"`
	}
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		log.Fatalln(err)
	}
	authToken = loginResponse.Data.Token
}

func TestCreateTransaction(t *testing.T) {
	createTransactionPath := "/transaction/create-transaction"
	t.Run("successful credit transaction", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"CREDIT", "amount":1000.67}`)

		res, _, err := runRequest(createTransactionPath, "POST", body, map[string]string{}, authToken)

		utils.AssertEqual(t, nil, err, "err should be nil")
		utils.AssertEqual(t, fiber.StatusOK, res.StatusCode, "Status code")
	})

	t.Run("successful debit transaction", func(t *testing.T) {
		body := fmt.Sprintf(`{"type":"DEBIT", "amount":100.23}`)

		res, _, err := runRequest(createTransactionPath, "POST", body, map[string]string{}, authToken)

		utils.AssertEqual(t, nil, err, "err should be nil")
		utils.AssertEqual(t, fiber.StatusOK, res.StatusCode, "Status code")
	})
}

func runRequest(path, method, body string, query map[string]string, token string) (*http.Response, map[string]interface{}, error) {
	url := baseUrl + path
	reqBody := strings.NewReader(body)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	// defer res.Body.Close()

	// var response map[string]interface{}
	// err = json.NewDecoder(res.Body).Decode(&response)
	// if err != nil {
	// 	return res, nil, err
	// }

	return res, nil, nil
}
