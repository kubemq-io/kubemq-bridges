package events_store

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-go"

	"github.com/stretchr/testify/require"
	"testing"

	"time"
)

type mockEventStoreReceiver struct {
	host    string
	port    int
	channel string
	timeout time.Duration
}

func (m *mockEventStoreReceiver) run(ctx context.Context) (*kubemq.EventStoreReceive, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(m.host, m.port),
		kubemq.WithClientId("response-id"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		return nil, err
	}
	errCh := make(chan error, 1)
	eventStoreCh, err := client.SubscribeToEventsStore(ctx, m.channel, "", errCh, kubemq.StartFromNewEvents())
	if err != nil {
		return nil, err
	}
	select {
	case eventStore := <-eventStoreCh:
		return eventStore, nil
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
		cfg          config.Metadata
		mockReceiver *mockEventStoreReceiver
		req          interface{}
		wantResp     *kubemq.EventStoreReceive
		wantErr      bool
	}{
		{
			name: "event-request",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
				},
			},
			mockReceiver: &mockEventStoreReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events_store1",
				timeout: 10 * time.Second,
			},
			req: kubemq.NewEvent().
				SetBody([]byte("data")).
				SetMetadata("metadata").
				SetChannel("events_store1").
				SetId("id"),
			wantResp: &kubemq.EventStoreReceive{
				Id:       "id",
				Channel:  "events_store1",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "event-store request",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
				},
			},
			mockReceiver: &mockEventStoreReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events_store2",
				timeout: 10 * time.Second,
			},
			req: &kubemq.EventStoreReceive{
				Id:        "id",
				Sequence:  1,
				Timestamp: time.Time{},
				Channel:   "events_store2",
				Metadata:  "metadata",
				Body:      []byte("data"),
				ClientId:  "",
				Tags:      nil,
			},
			wantResp: &kubemq.EventStoreReceive{
				Id:       "id",
				Channel:  "events_store2",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "command request",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
				},
			},
			mockReceiver: &mockEventStoreReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events_store3",
				timeout: 10 * time.Second,
			},
			req: &kubemq.CommandReceive{
				Id:       "id",
				Channel:  "events_store3",
				Metadata: "metadata",
				Body:     []byte("data"),
				Tags:     nil,
			},
			wantResp: &kubemq.EventStoreReceive{
				Id:       "id",
				Channel:  "events_store3",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "query request",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
				},
			},
			mockReceiver: &mockEventStoreReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events_store4",
				timeout: 10 * time.Second,
			},
			req: &kubemq.QueryReceive{
				Id:       "id",
				Channel:  "events_store4",
				Metadata: "metadata",
				Body:     []byte("data"),
				Tags:     nil,
			},
			wantResp: &kubemq.EventStoreReceive{
				Id:       "id",
				Channel:  "events_store4",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "queue request",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
				},
			},
			mockReceiver: &mockEventStoreReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events_store5",
				timeout: 10 * time.Second,
			},
			req: kubemq.NewQueueMessage().
				SetId("id").
				SetChannel("events_store5").
				SetMetadata("metadata").
				SetBody([]byte("data")),
			wantResp: &kubemq.EventStoreReceive{
				Id:       "id",
				Channel:  "events_store5",
				Metadata: "metadata",
				Body:     []byte("data"),
				ClientId: "response-id",
				Tags:     nil,
			},
			wantErr: false,
		},
		{
			name: "bad request - invalid type",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:50000",
				},
			},
			mockReceiver: &mockEventStoreReceiver{
				host:    "localhost",
				port:    50000,
				channel: "events_store",
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
			recRequestCh := make(chan *kubemq.EventStoreReceive, 1)
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
			err := target.Init(ctx, tt.cfg)
			require.NoError(t, err)
			_, err = target.Do(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			select {
			case gotRequest := <-recRequestCh:
				require.EqualValues(t, tt.wantResp.Id, gotRequest.Id)
				require.EqualValues(t, tt.wantResp.Metadata, gotRequest.Metadata)
				require.EqualValues(t, tt.wantResp.Body, gotRequest.Body)
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
		name    string
		cfg     config.Metadata
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address":    "localhost:50000",
					"client_id":  "client_id",
					"auth_token": "some-auth token",
					"channels":   "some-channel",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad connection",
			cfg: config.Metadata{
				Name: "kubemq-target",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost:40000",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()

			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}
