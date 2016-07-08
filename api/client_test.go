package api

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockHttpClient struct {
	GetFunc func(string) (*http.Response, error)
}

func (c mockHttpClient) Get(url string) (*http.Response, error) {
	return c.GetFunc(url)
}

func TestClientURL(t *testing.T) {
	client := NewClient("endpoint")
	assert.Equal(t, "http://endpoint", client.url)
}

func TestClientHttpError(t *testing.T) {
	client := NewClient("endpoint")
	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, errors.New("boom!")
		},
	}

	_, err := client.get("dummy", map[string]interface{}{})
	assert.EqualError(t, err, "boom!")
}

func TestClientReaderError(t *testing.T) {
	client := NewClient("endpoint")
	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(new(brokenReader)),
			}, nil
		},
	}

	_, err := client.get("dummy", map[string]interface{}{})
	assert.EqualError(t, err, "read error!")
}

func TestClientValidErrorResponse(t *testing.T) {
	client := NewClient("endpoint")
	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"error":"error happened"}`))),
			}, nil
		},
	}

	_, err := client.get("dummy", map[string]interface{}{})
	assert.EqualError(t, err, "error happened")
}

func TestClientInvalidErrorResponse(t *testing.T) {
	client := NewClient("endpoint")
	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`ugh`))),
			}, nil
		},
	}

	_, err := client.get("dummy", map[string]interface{}{})
	assert.Equal(t, err, ErrNotJSON)
}

func TestClientAddGroup(t *testing.T) {
	client := NewClient("endpoint")
	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			assert.Contains(t, url, "endpoint/api/group/add?")
			assert.Contains(t, url, "bootstrap-brokers=localhost%3A9092")
			assert.Contains(t, url, "group-id=foo")
			assert.Contains(t, url, "subscription=bar")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte{})),
			}, nil
		},
	}

	err := client.AddGroup("foo", "bar", "localhost:9092")
	assert.Nil(t, err)
}

func TestClientListGroups(t *testing.T) {
	client := NewClient("endpoint")
	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			assert.Contains(t, url, "endpoint/api/group/list")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`[]`))),
			}, nil
		},
	}

	groups, err := client.ListGroups()
	assert.Nil(t, err)
	assert.Empty(t, groups)

	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`ugh`))),
			}, nil
		},
	}

	groups, err = client.ListGroups()
	assert.Equal(t, ErrNotJSON, err)
	assert.Nil(t, groups)

	client.httpClient = mockHttpClient{
		GetFunc: func(url string) (*http.Response, error) {
			assert.Contains(t, url, "endpoint/api/group/list")

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{}`))),
			}, nil
		},
	}

	groups, err = client.ListGroups()
	assert.NotNil(t, err)
	assert.Nil(t, groups)
}

type brokenReader struct{}

func (r *brokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error!")
}
