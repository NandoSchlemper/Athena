package api

import (
	"athena/domain"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	ListaVeiculos() (*domain.Response, error)
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
	return base64.StdEncoding.EncodeToString([]byte(w.Username + ":" + w.Password))
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
func (w *TrackerAPIClient) ListaVeiculos() (*domain.Response, error) {
	payload, err := w.config.getJsonPayload()

	if err != nil {
		return nil, fmt.Errorf("again brother?")
	}

	var response domain.Response

	req, err := http.NewRequest(
		"POST",
		w.config.getBaseUrl()+"/lista_veiculos",
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request %w", err)
	}

	req.Header.Add("Authorization", "Basic "+w.config.getAuth())
	req.Header.Add("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("no body friendo")
	}

	jsonBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("brother just failed")
	}

	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("no json bro")
	}

	return &response, nil

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
