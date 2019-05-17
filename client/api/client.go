package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/olekukonko/tablewriter"

	"github.com/christarazi/gravitational-challenge/client/util"
	"github.com/christarazi/gravitational-challenge/config"
	"github.com/christarazi/gravitational-challenge/models"
)

// Client is a struct for interacting with the API.
type Client struct {
	args []string
}

// NewClient creates an instance of the Client struct.
func NewClient(args []string) *Client {
	return &Client{
		args: args,
	}
}

// Start sends a POST request to the /start endpoint.
func (c *Client) Start() (uint64, error) {
	data, err := marshal(models.StartRequest{
		Command: c.args[0],
		Args:    c.args[1:],
	})
	if err != nil {
		return 0, err
	}

	uri := fmt.Sprintf("http://0.0.0.0:%d/start", config.Port)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("Error getting response: %v", err)
	}

	defer resp.Body.Close()

	err = checkHTTPStatusCode(resp)
	if err != nil {
		return 0, err
	}

	sr := &models.StartResponse{}
	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		return 0, fmt.Errorf("Error decoding response: %v", err)
	}

	return sr.ID, nil
}

// Stop sends a POST request to the /stop endpoint.
func (c *Client) Stop(id uint64) error {
	data, err := marshal(models.StopRequest{ID: id})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("http://0.0.0.0:%d/stop", config.Port)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("Error getting response: %v", err)
	}

	defer resp.Body.Close()

	err = checkHTTPStatusCode(resp)
	if err != nil {
		return err
	}

	return nil
}

// Status sends a GET request to the /status endpoint.
func (c *Client) Status() (string, error) {
	if len(c.args) == 0 {
		return allStatus()
	}

	id, err := util.ConvertAndValidateID(c.args[0])
	if err != nil {
		return "", err
	}

	return status(id)
}

func status(id uint64) (string, error) {
	uri := fmt.Sprintf("http://0.0.0.0:8080/status/%d", id)

	resp, err := do(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	sr := &models.StatusResponse{}
	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		return "", fmt.Errorf("Error decoding response: %v", err)
	}

	return sr.Status, nil
}

func allStatus() (string, error) {
	uri := fmt.Sprintf("http://0.0.0.0:%d/status", config.Port)

	resp, err := do(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	asr := &models.AllStatusResponse{}
	err = json.NewDecoder(resp.Body).Decode(asr)
	if err != nil {
		return "", fmt.Errorf("Error decoding response: %v", err)
	}

	buffer := &bytes.Buffer{}
	table := tablewriter.NewWriter(buffer)
	table.SetHeader([]string{"ID", "Command", "Args", "Status"})

	for _, v := range asr.Jobs {
		table.Append([]string{
			strconv.FormatUint(v.ID, 10),
			v.Command,
			strings.Join(v.Args, ","),
			v.Status})
	}

	table.Render()

	return buffer.String(), nil
}

func do(uri string) (*http.Response, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("Error getting response: %v", err)
	}

	err = checkHTTPStatusCode(resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func marshal(payload interface{}) ([]byte, error) {
	var data []byte

	data, err := json.Marshal(payload)
	if err != nil {
		errwrap.Wrapf("Error marshalling request: {{err}}", err)
	}

	return data, err
}

func checkHTTPStatusCode(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading body of response: %v", err)
		}

		return fmt.Errorf("Server returned %d: %v", resp.StatusCode, string(body))
	}

	return nil
}
