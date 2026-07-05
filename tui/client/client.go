package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const defaultBaseURL = "https://gowordle.alwaysdata.net"

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(backendUrl string) *Client {
	baseURL := backendUrl
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

type wordResponse struct {
	ID    int    `json:"id"`
	Word  string `json:"word"`
	Error string `json:"error"`
}

func (c *Client) RandomWord() (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/words/random", nil)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	var result wordResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, "", err
	}
	if resp.StatusCode != http.StatusOK {
		if result.Error != "" {
			return 0, "", fmt.Errorf("%s", result.Error)
		}
		return 0, "", fmt.Errorf("fetch word failed (%d)", resp.StatusCode)
	}
	return result.ID, result.Word, nil
}

func (c *Client) SaveGame(wordID int, attempts json.RawMessage, won bool, duration int) error {
	body, err := json.Marshal(map[string]any{
		"word_id":  wordID,
		"attempts": attempts,
		"won":      won,
		"duration": duration,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/games", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("save game failed (%d)", resp.StatusCode)
	}
	return nil
}

func (c *Client) Login(email, password string) (string, error) {
	body, err := json.Marshal(loginRequest{Email: email, Password: password})
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Post(c.baseURL+"/api/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result authResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		if result.Error != "" {
			return "", fmt.Errorf("%s", result.Error)
		}
		return "", fmt.Errorf("login failed (%d)", resp.StatusCode)
	}
	return result.Token, nil
}

func (c *Client) Register(username, email, password string) (string, error) {
	body, err := json.Marshal(registerRequest{Username: username, Email: email, Password: password})
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Post(c.baseURL+"/api/auth/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result authResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusCreated {
		if result.Error != "" {
			return "", fmt.Errorf("%s", result.Error)
		}
		return "", fmt.Errorf("registration failed (%d)", resp.StatusCode)
	}
	return result.Token, nil
}
