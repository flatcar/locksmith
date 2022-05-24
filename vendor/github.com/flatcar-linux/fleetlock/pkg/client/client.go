package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPClient interface holds the required Post method
// to send FleetLock requests.
type HTTPClient interface {
	// Do send a `body` payload to the URL.
	Do(*http.Request) (*http.Response, error)
}

// Payload is the content to send
// to the FleetLock server.
type Payload struct {
	// ClientParams holds the parameters specific to the
	// FleetLock client.
	ClientParams *ClientParams `json:"client_params"`
}

// Client params is the object holding the
// ID and the group for each client.
type ClientParams struct {
	// ID is the client identifer. (e.g node name or UUID)
	ID string `json:"id"`
	// Group is the reboot-group of the client.
	Group string `json:"group"`
}

// Client holds the params related to the host
// in order to interact with the Fleet-Lock URL.
type Client struct {
	URL   string
	group string
	id    string
	http  HTTPClient
}

func (c *Client) generateRequest(endpoint string) (*http.Request, error) {
	payload := &Payload{
		ClientParams: &ClientParams{
			ID:    c.id,
			Group: c.group,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling the payload: %w", err)
	}

	j := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", c.URL, endpoint), j)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	headers := make(http.Header)
	headers.Add("fleet-lock-protocol", "true")
	req.Header = headers

	return req, nil
}

func handleResponse(resp *http.Response) error {
	statusType := resp.StatusCode / 100

	switch statusType {
	case 2:
		return nil
	case 3, 4, 5:
		// We try to extract an eventual error.
		r := bufio.NewReader(resp.Body)
		body, err := ioutil.ReadAll(r)
		if err != nil {
			return fmt.Errorf("reading body: %w", err)
		}

		resp.Body.Close()

		e := &Error{}
		if err := json.Unmarshal(body, &e); err != nil {
			return fmt.Errorf("unmarshalling error: %v", err)
		}

		return fmt.Errorf("fleetlock error: %s", e.String())
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// RecursiveLock tries to reserve (lock) a slot for rebooting
func (c *Client) RecursiveLock() error {
	req, err := c.generateRequest("v1/pre-reboot")
	if err != nil {
		return fmt.Errorf("generating request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("doing the request: %w", err)
	}

	return handleResponse(resp)
}

// UnlockIfHeld tries to release (unlock) a slot that it was previously holding
func (c *Client) UnlockIfHeld() error {
	req, err := c.generateRequest("v1/steady-state")
	if err != nil {
		return fmt.Errorf("generating request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("doing the request: %w", err)
	}

	return handleResponse(resp)
}

// New builds a Fleet-Lock client.
func New(URL, group, ID string, c HTTPClient) (*Client, error) {
	if _, err := url.ParseRequestURI(URL); err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}

	return &Client{
		URL:   URL,
		http:  c,
		group: group,
		id:    ID,
	}, nil
}
