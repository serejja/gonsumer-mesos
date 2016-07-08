package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/serejja/gonsumer-mesos/framework"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	groupAddEndpointURL  = "/api/group/add"
	groupListEndpointURL = "/api/group/list"
)

var (
	ErrNotJSON = errors.New("Server returned non-JSON response.")
)

type Client struct {
	url        string
	httpClient httpClient
}

func NewClient(url string) *Client {
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	return &Client{
		url:        url,
		httpClient: netHttpClient{},
	}
}

func (c *Client) AddGroup(groupID string, subscription string, bootstrapBrokers string) error {
	_, err := c.get(groupAddEndpointURL, map[string]interface{}{
		framework.ParamGroupID:          groupID,
		framework.ParamSubscription:     subscription,
		framework.ParamBootstrapBrokers: bootstrapBrokers,
	})

	return err
}

func (c *Client) ListGroups() ([]*framework.Group, error) {
	rawGroups, err := c.get(groupListEndpointURL, nil)
	if err != nil {
		return nil, err
	}

	var groups []*framework.Group
	err = json.Unmarshal(rawGroups, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *Client) get(endpoint string, params map[string]interface{}) ([]byte, error) {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, fmt.Sprint(value))
	}
	queryString := values.Encode()

	url := fmt.Sprintf("%s%s?%s", c.url, endpoint, queryString)
	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		errorResponse := new(struct {
			Error string
		})

		err = json.Unmarshal(responseBody, &errorResponse)
		if err != nil {
			return nil, ErrNotJSON
		}

		return nil, errors.New(errorResponse.Error)
	}

	return responseBody, nil
}

type httpClient interface {
	Get(url string) (*http.Response, error)
}

type netHttpClient struct{}

func (netHttpClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}
