package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type IWrsatAPIConfig interface {
	SetAPIVariables()
	getTimeout() time.Duration
	getBaseUrl() string
	getAuth() string
	getJsonPayload() ([]byte, error)
}

type IWrsatAPIClient interface {
	ListaVeiculos() (*http.Response, error)
}

type WrsatAPIConfig struct {
	Username       string
	Password       string
	BaseURL        string
	Timeout        time.Duration
	DefaultPayload map[string]string
}

// getPayload implements IWrsatAPIConfig.
func (w *WrsatAPIConfig) getJsonPayload() ([]byte, error) {
	jsonPayload, err := json.Marshal(w.DefaultPayload)
	if err != nil {
		return nil, fmt.Errorf("payload to json failed: %w", err)
	}
	return jsonPayload, nil
}

// getBaseUrl implements IWrsatAPIConfig.
func (w *WrsatAPIConfig) getBaseUrl() string {
	return w.BaseURL
}

func (w *WrsatAPIConfig) getAuth() string {
	return w.Username + ":" + w.Password
}

// SetHeaders implements IWrsatAPIConfig.
func (w *WrsatAPIConfig) SetAPIVariables() {
	w.DefaultPayload = map[string]string{
		"usuario":   w.Username,
		"senha":     w.Password,
		"ordem":     "ASC",
		"limit":     "100",
		"pagina":    "1",
		"descricao": "",
	}
}

func (w *WrsatAPIConfig) getTimeout() time.Duration {
	return w.Timeout
}

type WrsatAPIClient struct {
	client *http.Client
	config IWrsatAPIConfig
}

// listaVeiculos implements IWrsatAPIClient.
func (w *WrsatAPIClient) ListaVeiculos() (*http.Response, error) {
	payload, _ := w.config.getJsonPayload()
	req, err := http.NewRequest(
		"POST",
		w.config.getBaseUrl()+"/lista_veiculos",
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request %w", err)
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(w.config.getAuth())))
	req.Header.Add("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func NewWrsatAPIClient(cfg IWrsatAPIConfig) IWrsatAPIClient {
	return &WrsatAPIClient{
		client: &http.Client{
			Timeout: cfg.getTimeout(),
		},
		config: cfg,
	}
}

func NewWrsatAPIConfig(timeout time.Duration) IWrsatAPIConfig {
	return &WrsatAPIConfig{
		Username: os.Getenv("WRSAT_USER"),
		Password: os.Getenv("WRSAT_PASSWORD"),
		BaseURL:  os.Getenv("WRSAT_BASE_URL"),
		Timeout:  timeout,
	}
}
