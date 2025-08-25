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

func (w *TrackerAPIConfig) getJsonPayload() ([]byte, error) {
	jsonPayload, err := json.Marshal(w.DefaultPayload)
	if err != nil {
		return nil, fmt.Errorf("bro just failed %w", err)
	}
	return jsonPayload, nil
}

func (w *TrackerAPIConfig) getBaseUrl() string {
	return w.BaseURL
}

func (w *TrackerAPIConfig) getAuth() string {
	return base64.StdEncoding.EncodeToString([]byte(w.Username + ":" + w.Password))
}

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

func (w *TrackerAPIClient) ListaVeiculos() (*domain.Response, error) {
	payload, err := w.config.getJsonPayload()
	if err != nil {
		return nil, fmt.Errorf("failed. sus. owo %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		w.config.getBaseUrl()+"/lista_veiculos",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("no request 4 u %w", err)
	}

	req.Header.Add("Authorization", "Basic "+w.config.getAuth())
	req.Header.Add("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("u're a failure %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("no fucking body: %w", err)
	}

	// Debug: print the raw response
	// fmt.Printf("Raw response: %s\n", string(jsonBytes))
	// obs: nunca mais fazer isso

	var response domain.Response
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("FAILED TO JSON BRUH, FAILED!: %w. Response: %s", err, string(jsonBytes))
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
	cfg := &TrackerAPIConfig{
		Username: os.Getenv("WRSAT_USER"),
		Password: os.Getenv("WRSAT_PASSWORD"),
		BaseURL:  os.Getenv("WRSAT_BASE_URL"),
		Timeout:  timeout,
	}
	cfg.SetDefaultTracker()
	return cfg
}
