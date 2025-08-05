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

type ITrackerAPIConfig interface {
	SetDefaultTracker()
	getTimeout() time.Duration
	getBaseUrl() string
	getAuth() string
	getJsonPayload() ([]byte, error)
}

type ITrackerAPIClient interface {
	ListaVeiculos() (*http.Response, error)
}

type TrackerAPIConfig struct {
	Username       string
	Password       string
	BaseURL        string
	Timeout        time.Duration
	DefaultPayload map[string]string
}

// getPayload implements IWrsatAPIConfig.
func (w *TrackerAPIConfig) getJsonPayload() ([]byte, error) {
	jsonPayload, err := json.Marshal(w.DefaultPayload)
	if err != nil {
		return nil, fmt.Errorf("payload to json failed: %w", err)
	}
	return jsonPayload, nil
}

// getBaseUrl implements IWrsatAPIConfig.
func (w *TrackerAPIConfig) getBaseUrl() string {
	return w.BaseURL
}

func (w *TrackerAPIConfig) getAuth() string {
	return w.Username + ":" + w.Password
}

// SetHeaders implements IWrsatAPIConfig.
func (w *TrackerAPIConfig) SetDefaultTracker() {
	w.DefaultPayload = map[string]string{
		"usuario":   w.Username,
		"senha":     w.Password,
		"ordem":     "ASC",
		"limit":     "100",
		"pagina":    "1",
		"descricao": "",
	}
}

func (w *TrackerAPIConfig) getTimeout() time.Duration {
	return w.Timeout
}

type TrackerAPIClient struct {
	client *http.Client
	config ITrackerAPIConfig
}

// listaVeiculos implements IWrsatAPIClient.
func (w *TrackerAPIClient) ListaVeiculos() (*http.Response, error) {
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

func NewTrackerAPIClient(cfg ITrackerAPIConfig) ITrackerAPIClient {
	return &TrackerAPIClient{
		client: &http.Client{
			Timeout: cfg.getTimeout(),
		},
		config: cfg,
	}
}

func NewTrackerAPIConfig(timeout time.Duration) ITrackerAPIConfig {
	return &TrackerAPIConfig{
		Username: os.Getenv("WRSAT_USER"),
		Password: os.Getenv("WRSAT_PASSWORD"),
		BaseURL:  os.Getenv("WRSAT_BASE_URL"),
		Timeout:  timeout,
	}
}
