package query

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"

	"github.com/kubemq-io/kubemq-go"

	"github.com/stretchr/testify/require"
	"testing"

	"time"
)

type mockQueryReceiver struct {
	host           string
	port           int
	channel        string
	executionDelay time.Duration
	executionError error
	executionTime  int64
}

func (m *mockQueryReceiver) run(ctx context.Context, t *testing.T) error {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(m.host, m.port),
		kubemq.WithClientId("response-id"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		return err
	}
	errCh := make(chan error, 1)
	queryCh, err := client.SubscribeToQueries(ctx, m.channel, "", errCh)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case query := <-queryCh:
				time.Sleep(m.executionDelay)
				queryResponse := client.R().SetRequestId(query.Id).SetResponseTo(query.ResponseTo).SetExecutedAt(time.Unix(m.executionTime, 0))
				if m.executionError != nil {
					queryResponse.SetError(m.executionError)
				} else {
					queryResponse.SetBody(query.Body)
					queryResponse.SetMetadata(query.Metadata)
				}
				err := queryResponse.Send(ctx)
				require.NoError(t, err)
			case err := <-errCh:
				require.NoError(t, err)
			case <-ctx.Done():
				return
			}

		}
	}()
	time.Sleep(time.Second)
	return nil
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		connection   config.Metadata
		mockReceiver *mockQueryReceiver
		req          interface{}
		wantResp     interface{}
		wantErr      bool
	}{
		{
			name: "event-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries1",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req: kubemq.NewEvent().
				SetBody([]byte("data")).
				SetMetadata("metadata").
				SetChannel("queries1").
				SetId("id"),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Metadata:         "metadata",
				ResponseClientId: "response-id",
				Body:             []byte("data"),
				CacheHit:         false,
				Error:            "",
				Tags:             nil,
			},
			wantErr: false,
		},
		{
			name: "event-store-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries2",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req: &kubemq.EventStoreReceive{
				Id:        "id",
				Sequence:  1,
				Timestamp: time.Time{},
				Channel:   "queries2",
				Metadata:  "metadata",
				Body:      []byte("data"),
				ClientId:  "",
				Tags:      nil,
			},
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Metadata:         "metadata",
				ResponseClientId: "response-id",
				Body:             []byte("data"),
				CacheHit:         false,
				Error:            "",
				Tags:             nil,
			},
			wantErr: false,
		},
		{
			name: "command-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries3",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req: &kubemq.QueryReceive{
				Id:         "id",
				ResponseTo: "some-response",
				Channel:    "queries3",
				Metadata:   "metadata",
				Body:       []byte("data"),
				Tags:       nil,
			},
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Metadata:         "metadata",
				ResponseClientId: "response-id",
				Body:             []byte("data"),
				CacheHit:         false,
				Error:            "",
				Tags:             nil,
			},
			wantErr: false,
		},
		{
			name: "query-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries4",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req: &kubemq.QueryReceive{
				Id:         "id",
				ResponseTo: "some-response",
				Channel:    "queries4",
				Metadata:   "metadata",
				Body:       []byte("data"),
				Tags:       nil,
			},
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Metadata:         "metadata",
				ResponseClientId: "response-id",
				Body:             []byte("data"),
				CacheHit:         false,
				Error:            "",
				Tags:             nil,
			},
			wantErr: false,
		},
		{
			name: "queue-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries5",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req: kubemq.NewQueueMessage().
				SetId("id").
				SetChannel("queries5").
				SetMetadata("metadata").
				SetBody([]byte("data")),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Metadata:         "metadata",
				ResponseClientId: "response-id",
				Body:             []byte("data"),
				CacheHit:         false,
				Error:            "",
				Tags:             nil,
			},
			wantErr: false,
		},
		{
			name: "event-request-overwrite channel ",
			connection: map[string]string{
				"address":         "localhost:50000",
				"default_channel": "queries6",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries6",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req: kubemq.NewEvent().
				SetBody([]byte("data")).
				SetMetadata("metadata").
				SetChannel("events-channel").
				SetId("id"),
			wantResp: &kubemq.QueryResponse{
				QueryId:          "id",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Metadata:         "metadata",
				ResponseClientId: "response-id",
				Body:             []byte("data"),
				CacheHit:         false,
				Error:            "",
				Tags:             nil,
			},
			wantErr: false,
		},
		{
			name: "bad request - invalid type",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries",
				executionDelay: 0,
				executionError: nil,
				executionTime:  1000,
			},
			req:      "bad-format",
			wantResp: nil,
			wantErr:  true,
		},
		{
			name: "event-request- query error ",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "queries7",
				executionDelay: 0,
				executionError: fmt.Errorf("some-error"),
				executionTime:  1000,
			},
			req: kubemq.NewEvent().
				SetBody([]byte("data")).
				SetMetadata("metadata").
				SetChannel("events-channel").
				SetId("id"),
			wantResp: nil,
			wantErr:  true,
		},
		{
			name: "event-request - query timeout",
			connection: map[string]string{
				"address":         "localhost:50000",
				"timeout_seconds": "1",
			},
			mockReceiver: &mockQueryReceiver{
				host:           "localhost",
				port:           50000,
				channel:        "none-query-channel",
				executionDelay: 2 * time.Second,
				executionError: nil,
				executionTime:  1000,
			},
			req: kubemq.NewEvent().
				SetBody([]byte("data")).
				SetMetadata("metadata").
				SetChannel("events-channel").
				SetId("id"),
			wantResp: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := tt.mockReceiver.run(ctx, t)
			require.NoError(t, err)
			target := New()
			err = target.Init(ctx, tt.connection, "", nil)
			require.NoError(t, err)
			gotResp, err := target.Do(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
		})
	}
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name       string
		connection config.Metadata
		wantErr    bool
	}{
		{
			name: "init",
			connection: map[string]string{
				"address":         "localhost:50000",
				"client_id":       "client_id",
				"auth_token":      "some-auth token",
				"default_channel": "some-channel",
				"timeout_seconds": "100",
			},
			wantErr: false,
		},
		{
			name: "init - error",
			connection: map[string]string{
				"address": "localhost",
			},
			wantErr: true,
		},
		{
			name: "init - bad connection",
			connection: map[string]string{
				"address": "localhost:40000",
			},
			wantErr: true,
		},
		{
			name: "init - bad timeout seconds",
			connection: map[string]string{
				"address":         "localhost:50000",
				"timeout_seconds": "-1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.connection, "", nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
