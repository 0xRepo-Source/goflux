package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Transport is an abstraction for underlying transport (ssh, quic, http).
type Transport interface {
	Dial(addr string) error
	Listen(addr string) error
}

// ChunkData represents chunk data being transferred.
type ChunkData struct {
	Path     string `json:"path"`
	ChunkID  int    `json:"chunk_id"`
	Data     []byte `json:"data"`
	Checksum string `json:"checksum"`
	Total    int    `json:"total"` // total number of chunks
}

// HTTPClient is an HTTP-based transport client.
type HTTPClient struct {
	BaseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		client:  &http.Client{},
	}
}

func (h *HTTPClient) Dial(addr string) error {
	h.BaseURL = addr
	return nil
}

func (h *HTTPClient) Listen(addr string) error {
	return fmt.Errorf("HTTPClient cannot listen")
}

// UploadChunk uploads a single chunk.
func (h *HTTPClient) UploadChunk(chunk ChunkData) error {
	data, err := json.Marshal(chunk)
	if err != nil {
		return err
	}

	resp, err := h.client.Post(h.BaseURL+"/upload", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s", string(body))
	}
	return nil
}

// Download downloads a file.
func (h *HTTPClient) Download(path string) ([]byte, error) {
	resp, err := h.client.Get(h.BaseURL + "/download?path=" + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed: %s", string(body))
	}

	return io.ReadAll(resp.Body)
}

// List lists files at a path.
func (h *HTTPClient) List(path string) ([]string, error) {
	resp, err := h.client.Get(h.BaseURL + "/list?path=" + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list failed: status %d", resp.StatusCode)
	}

	var files []string
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}
	return files, nil
}
