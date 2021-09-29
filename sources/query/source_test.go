package query

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/middleware"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
	"github.com/kubemq-io/kubemq-go"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type mockTarget struct {
	setResponse interface{}
	setError    error
	delay       time.Duration
}

func (m *mockTarget) Do(ctx context.Context, request interface{}) (interface{}, error) {
	time.Sleep(m.delay)
	return m.setResponse, m.setError
}

func setupSource(ctx context.Context, targets []middleware.Middleware) (*Source, error) {
	source := New()
	err := source.Init(ctx, config.Metadata{

		"address":                    "localhost:50000",
		"client_id":                  "responseid",
		"auth_token":                 "some-auth token",
		"channel":                    "queries",
		"group":                      "group",
		"auto_reconnect":             "true",
		"reconnect_interval_seconds": "1",
		"max_reconnects":             "0",
		"sources":                    "1",
	}, config.Metadata{}, nil)
	if err != nil {
		return nil, err
	}
	err = source.Start(ctx, targets)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	return source, nil
}
func sendQuery(t *testing.T, ctx context.Context, req *kubemq.Query, sendChannel string, timeout time.Duration) (*kubemq.QueryResponse, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	require.NoError(t, err)
	return client.SetQuery(req).SetChannel(sendChannel).SetTimeout(timeout).Send(ctx)

}
func TestClient_processQuery(t *testing.T) {
	tests := []struct {
		name     string
		target   middleware.Middleware
		req      *kubemq.Query
		wantResp *kubemq.QueryResponse
		timeout  time.Duration
		sendCh   string
		wantErr  bool
	}{
		{
			name: "request - command target - executed",
			target: &mockTarget{
				setResponse: &kubemq.CommandResponse{
					CommandId:        "id",
					ResponseClientId: "responseid",
					Executed:         true,
					ExecutedAt:       time.Unix(1000, 0),
					Error:            "",
					Tags:             nil,
				},
				setError: nil,
				delay:    0,
			},
			req: kubemq.NewQuery().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				ResponseClientId: "responseid",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Error:            "",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "queries",
			wantErr: false,
		},
		{
			name: "request - command target - not executed",
			target: &mockTarget{
				setResponse: &kubemq.CommandResponse{
					CommandId:        "id",
					ResponseClientId: "responseid",
					Executed:         false,
					Error:            "some-error",
					Tags:             nil,
				},
				setError: nil,
				delay:    0,
			},
			req: kubemq.NewQuery().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				ResponseClientId: "responseid",
				Executed:         false,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "some-error",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "queries",
			wantErr: false,
		},
		{
			name: "request - query target - executed",
			target: &mockTarget{
				setResponse: &kubemq.QueryResponse{
					QueryId:          "id",
					Executed:         true,
					ExecutedAt:       time.Unix(1000, 0),
					Metadata:         "some-metadata",
					ResponseClientId: "responseid",
					Body:             nil,
					CacheHit:         false,
					Error:            "",
					Tags:             nil,
				},
				setError: nil,
				delay:    0,
			},
			req: kubemq.NewQuery().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Metadata:         "some-metadata",
				ResponseClientId: "responseid",
				Body:             nil,
				CacheHit:         false,
				Error:            "",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "queries",
			wantErr: false,
		},
		{
			name: "request - query target - not executed",
			target: &mockTarget{
				setResponse: &kubemq.QueryResponse{
					QueryId:          "id",
					Executed:         false,
					ExecutedAt:       time.Time{},
					Metadata:         "",
					ResponseClientId: "responseid",
					Body:             nil,
					CacheHit:         false,
					Error:            "some-error",
					Tags:             nil,
				},
			},
			req: kubemq.NewQuery().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				ResponseClientId: "responseid",
				Executed:         false,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "some-error",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "queries",
			wantErr: false,
		},
		{
			name: "request - command target - error",
			target: &mockTarget{
				setResponse: nil,
				setError:    fmt.Errorf("some-error"),
				delay:       0,
			},
			req: kubemq.NewQuery().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				ResponseClientId: "responseid",
				Executed:         false,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "some-error",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "queries",
			wantErr: false,
		},
		{
			name: "request - other target type - executed",
			target: &mockTarget{
				setResponse: nil,
				setError:    nil,
				delay:       0,
			},
			req: kubemq.NewQuery().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				ResponseClientId: "responseid",
				Executed:         true,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "queries",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c, err := setupSource(ctx, []middleware.Middleware{tt.target})
			require.NoError(t, err)
			defer func() {
				_ = c.Stop()
			}()
			gotResp, err := sendQuery(t, ctx, tt.req, tt.sendCh, tt.timeout)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
		})
	}
}

//
func TestClient_Init(t *testing.T) {
	tests := []struct {
		name       string
		connection config.Metadata
		wantErr    bool
	}{
		{
			name: "init",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"concurrency":                "1",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "2",
			},
			wantErr: false,
		},
		{
			name: "init - error",
			connection: config.Metadata{
				"address": "localhost",
			},
			wantErr: true,
		},
		{
			name: "init - bad channel",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "",
				"group":                      "",
				"concurrency":                "1",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "2",
			},
			wantErr: true,
		},
		{
			name: "init - bad Source",
			connection: config.Metadata{
				"address":                    "localhost",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "2",
			},
			wantErr: true,
		},
		{
			name: "init - bad reconnect interval",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "-1",
				"max_reconnects":             "0",
				"sources":                    "2",
			},
			wantErr: true,
		},
		{
			name: "init - bad sources",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "-1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.connection, config.Metadata{}, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
