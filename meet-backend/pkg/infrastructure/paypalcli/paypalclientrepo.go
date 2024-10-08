package paypalcli

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/erodriguezg/meet/pkg/core/repository"
)

// https://developer.paypal.com/docs/checkout/standard/integrate/

type paypalPayClient struct {
	httpClient       *http.Client
	paypalBaseApiUrl string
	paypalClientId   string
	paypalAppSecret  string
}

type CreateOrderPayload struct {
	Intent        string              `json:"intent"`
	PurchaseUnits []OrderPurchaseUnit `json:"purchase_units"`
}

type OrderPurchaseUnit struct {
	Amount PurchaseUnitAmount `json:"amount"`
}

type PurchaseUnitAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

func NewPayPalPaymentClientRepository(
	httpClient *http.Client,
	paypalBaseApiUrl string,
	paypalClientId string,
	paypalAppSecret string) repository.PaymentClientRepository {
	return &paypalPayClient{
		httpClient:       httpClient,
		paypalBaseApiUrl: paypalBaseApiUrl,
		paypalClientId:   paypalClientId,
		paypalAppSecret:  paypalAppSecret,
	}
}

func (port *paypalPayClient) GetClientData() (map[string]any, error) {
	return map[string]any{
		"clientId": port.paypalClientId,
	}, nil
}

func (port *paypalPayClient) CreateOrder(value float64, currencyCode string) (orderID string, outputErr error) {
	paypalCreateOrderUrl := fmt.Sprintf("%s/v2/checkout/orders", port.paypalBaseApiUrl)
	accessToken, err := port.generateAccessToken()
	if err != nil {
		return "", err
	}

	payloadStruct := CreateOrderPayload{
		Intent: "CAPTURE",
		PurchaseUnits: []OrderPurchaseUnit{
			{
				Amount: PurchaseUnitAmount{
					CurrencyCode: currencyCode,
					Value:        fmt.Sprintf("%.2f", value),
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(&payloadStruct)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		paypalCreateOrderUrl,
		bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Add("Content-Type", "application/json")

	response, err := port.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		outputErr = response.Body.Close()
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	var jsonResult map[string]any
	err = json.NewDecoder(response.Body).Decode(&jsonResult)
	if err != nil {
		return "", err
	}

	return jsonResult["id"].(string), nil
}

func (port *paypalPayClient) CapturePayment(orderID string) (paymentDetail map[string]any, outputErr error) {
	paypalCapturePaymentUrl := fmt.Sprintf("%s/v2/checkout/orders/%s/capture", port.paypalBaseApiUrl, orderID)
	accessToken, err := port.generateAccessToken()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		http.MethodPost,
		paypalCapturePaymentUrl,
		nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	request.Header.Set("Content-Type", "application/json")

	response, err := port.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		outputErr = response.Body.Close()
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	var jsonResult map[string]any
	err = json.NewDecoder(response.Body).Decode(&jsonResult)
	if err != nil {
		return nil, err
	}

	return jsonResult, nil
}

// private

func (port *paypalPayClient) generateAccessToken() (accessToken string, outputErr error) {

	paypalTokenUrl := fmt.Sprintf("%s/v1/oauth2/token", port.paypalBaseApiUrl)

	authPlain := fmt.Sprintf("%s:%s", port.paypalClientId, port.paypalAppSecret)
	authB64 := base64.StdEncoding.EncodeToString([]byte(authPlain))

	payload := url.Values{}
	payload.Set("grant_type", "client_credentials")

	request, err := http.NewRequest(
		http.MethodPost,
		paypalTokenUrl,
		strings.NewReader(payload.Encode()))
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", authB64))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := port.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer func() {
		outputErr = response.Body.Close()
	}()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	var jsonResult map[string]any
	err = json.NewDecoder(response.Body).Decode(&jsonResult)
	if err != nil {
		return "", err
	}

	return jsonResult["access_token"].(string), nil
}
