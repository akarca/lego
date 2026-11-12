package internal

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-acme/lego/v5/internal/tester/servermock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockBuilder() *servermock.Builder[*Client] {
	return servermock.NewBuilder[*Client](
		func(server *httptest.Server) (*Client, error) {
			client, err := NewClient("secret")
			if err != nil {
				return nil, err
			}

			client.BaseURL, _ = url.Parse(server.URL)
			client.HTTPClient = server.Client()

			return client, nil
		},
		servermock.CheckHeader().
			WithAccept("application/json").
			WithAuthorization("Bearer secret"),
	)
}

func TestClient_AddRecord(t *testing.T) {
	client := mockBuilder().
		Route("POST /dns/example.com/records",
			servermock.ResponseFromFixture("add_record.json").
				WithStatusCode(http.StatusCreated),
			servermock.CheckRequestJSONBodyFromFixture("add_record-request.json"),
		).
		Build(t)

	record := Record{
		Type:     "TXT",
		Hostname: "_acme-challenge",
		Text:     "ADw2sEd82DUgXcQ9hNBZThJs7zVJkR5v9JeSbAb9mZY",
		TTL:      "5m",
	}

	result, err := client.AddRecord(t.Context(), "example.com", record)
	require.NoError(t, err)

	expected := &Record{
		ID:       1234,
		Type:     "TXT",
		Hostname: "_acme-challenge",
		Text:     "ADw2sEd82DUgXcQ9hNBZThJs7zVJkR5v9JeSbAb9mZY",
		TTL:      "5m",
	}

	assert.Equal(t, expected, result)
}

func TestClient_AddRecord_error(t *testing.T) {
	client := mockBuilder().
		Route("POST /dns/example.com/records",
			servermock.ResponseFromFixture("error.json").
				WithStatusCode(http.StatusNotFound),
		).
		Build(t)

	_, err := client.AddRecord(t.Context(), "example.com", Record{Type: "TXT", Hostname: "_acme-challenge", Text: "x"})
	require.EqualError(t, err, "not_found: No DNS zone example.com in your account.")
}

func TestClient_DeleteRecord(t *testing.T) {
	client := mockBuilder().
		Route("DELETE /dns/example.com/records/TXT/1234",
			servermock.Noop().
				WithStatusCode(http.StatusNoContent),
		).
		Build(t)

	err := client.DeleteRecord(t.Context(), "example.com", "TXT", 1234)
	require.NoError(t, err)
}

func TestClient_DeleteRecord_error(t *testing.T) {
	client := mockBuilder().
		Route("DELETE /dns/example.com/records/TXT/1234",
			servermock.ResponseFromFixture("error.json").
				WithStatusCode(http.StatusNotFound),
		).
		Build(t)

	err := client.DeleteRecord(t.Context(), "example.com", "TXT", 1234)
	require.EqualError(t, err, "not_found: No DNS zone example.com in your account.")
}
