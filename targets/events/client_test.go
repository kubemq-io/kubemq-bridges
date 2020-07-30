package events

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-go"

	"github.com/stretchr/testify/require"
	"testing"

	"time"
)

type mockEventReceiver struct {
	host    string
	port    int
	channel string
	timeout time.Duration
}

func (m *mockEventReceiver) run(ctx context.Context) (*kubemq.Event, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(m.host, m.port),
		kubemq.WithClientId("response-id"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		return nil, err
	}
	errCh := make(chan error, 1)
	eventCh, err := client.SubscribeToEvents(ctx, m.channel, "", errCh)
	if err != nil {
		return nil, err
	}
	select {
	case event := <-eventCh:
		return event, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, nil

	case <-time.After(m.timeout):
		return nil, fmt.Errorf("timeout")
	}

}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		connection   config.Metadata
		mockReceiver *mockEventReceiver
		req          interface{}
		wantResp     interface{}
		wantErr      bool
	}{
		{
			name: "event-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockEventReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events1",
				timeout: 10 * time.Second,
			},
			req: kubemq.NewEvent().
				SetBody([]byte("data")).
				SetMetadata("metadata").
				SetChannel("events1").
				SetId("id"),
			wantResp: &kubemq.Event{
				Id:       "id",
				Channel:  "events1",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "event-store request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockEventReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events2",
				timeout: 10 * time.Second,
			},
			req: &kubemq.EventStoreReceive{
				Id:        "id",
				Sequence:  1,
				Timestamp: time.Time{},
				Channel:   "events2",
				Metadata:  "metadata",
				Body:      []byte("data"),
				ClientId:  "",
				Tags:      nil,
			},
			wantResp: &kubemq.Event{
				Id:       "id",
				Channel:  "events2",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "command request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockEventReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events3",
				timeout: 10 * time.Second,
			},
			req: &kubemq.CommandReceive{
				Id:       "id",
				Channel:  "events3",
				Metadata: "metadata",
				Body:     []byte("data"),
				Tags:     nil,
			},
			wantResp: &kubemq.Event{
				Id:       "id",
				Channel:  "events3",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "query request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockEventReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events4",
				timeout: 10 * time.Second,
			},
			req: &kubemq.QueryReceive{
				Id:       "id",
				Channel:  "events4",
				Metadata: "metadata",
				Body:     []byte("data"),
				Tags:     nil,
			},
			wantResp: &kubemq.Event{
				Id:       "id",
				Channel:  "events4",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "queue request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockEventReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events5",
				timeout: 10 * time.Second,
			},
			req: kubemq.NewQueueMessage().
				SetId("id").
				SetChannel("events5").
				SetMetadata("metadata").
				SetBody([]byte("data")),
			wantResp: &kubemq.Event{
				Id:       "id",
				Channel:  "events5",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "bad request - invalid type",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockEventReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events",
				timeout: 10 * time.Second,
			},
			req:      "bad-format",
			wantResp: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			recRequestCh := make(chan *kubemq.Event, 1)
			recErrCh := make(chan error, 1)
			go func() {
				gotRequest, err := tt.mockReceiver.run(ctx)
				select {
				case recErrCh <- err:
				case recRequestCh <- gotRequest:
				}
			}()
			time.Sleep(time.Second)
			target := New()
			err := target.Init(ctx, tt.connection)
			require.NoError(t, err)
			_, err = target.Do(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			select {
			case gotRequest := <-recRequestCh:
				require.EqualValues(t, tt.wantResp, gotRequest)
			case err := <-recErrCh:
				require.NoError(t, err)
			case <-ctx.Done():
				require.NoError(t, ctx.Err())
			}
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
				"address":    "localhost:50000",
				"client_id":  "client_id",
				"auth_token": "some-auth token",
				"channels":   "some-channel",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.connection); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
