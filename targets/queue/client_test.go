package queue

import (
	"context"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
	"github.com/kubemq-io/kubemq-go"

	"github.com/stretchr/testify/require"
	"testing"

	"time"
)

type mockQueueReceiver struct {
	host    string
	port    int
	channel string
	timeout int32
}

func (m *mockQueueReceiver) run(ctx context.Context) (*kubemq.QueueMessage, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(m.host, m.port),
		kubemq.WithClientId("response-id"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		return nil, err
	}

	queueMessages, err := client.ReceiveQueueMessages(ctx, &kubemq.ReceiveQueueMessagesRequest{
		RequestID:           "id",
		ClientID:            uuid.New().String(),
		Channel:             m.channel,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     m.timeout,
		IsPeak:              false,
	})
	if err != nil {
		return nil, err
	}
	if len(queueMessages.Messages) == 0 {
		return nil, nil
	}
	return queueMessages.Messages[0], nil
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name         string
		connection   config.Metadata
		mockReceiver *mockQueueReceiver
		req          interface{}
		wantResp     *kubemq.QueueMessage
		wantErr      bool
	}{
		{
			name: "event-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues1",
				timeout: 5,
			},
			req: kubemq.NewEvent().
				SetBody([]byte("data")).
				SetMetadata("metadata").
				SetChannel("queues1").
				SetId("id"),
			wantResp: kubemq.NewQueueMessage().
				SetMetadata("metadata").
				SetId("id").
				SetBody([]byte("data")),
			wantErr: false,
		},
		{
			name: "event-store-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues2",
				timeout: 5,
			},
			req: &kubemq.EventStoreReceive{
				Id:        "id",
				Sequence:  1,
				Timestamp: time.Time{},
				Channel:   "queues2",
				Metadata:  "metadata",
				Body:      []byte("data"),
				ClientId:  "",
				Tags:      nil,
			},
			wantResp: kubemq.NewQueueMessage().
				SetMetadata("metadata").
				SetId("id").
				SetBody([]byte("data")),
			wantErr: false,
		},
		{
			name: "command-store-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues3",
				timeout: 5,
			},
			req: &kubemq.CommandReceive{
				Id:         "id",
				ResponseTo: "some-response",
				Channel:    "queues3",
				Metadata:   "metadata",
				Body:       []byte("data"),
				Tags:       nil,
			},
			wantResp: kubemq.NewQueueMessage().
				SetMetadata("metadata").
				SetId("id").
				SetBody([]byte("data")),
			wantErr: false,
		},
		{
			name: "query-store-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues4",
				timeout: 5,
			},
			req: &kubemq.QueryReceive{
				Id:         "id",
				ResponseTo: "some-response",
				Channel:    "queues4",
				Metadata:   "metadata",
				Body:       []byte("data"),
				Tags:       nil,
			},
			wantResp: kubemq.NewQueueMessage().
				SetMetadata("metadata").
				SetId("id").
				SetBody([]byte("data")),
			wantErr: false,
		},
		{
			name: "query-store-request",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues5",
				timeout: 5,
			},
			req: kubemq.NewQueueMessage().
				SetId("id").
				SetChannel("queues5").
				SetMetadata("metadata").
				SetBody([]byte("data")),
			wantResp: kubemq.NewQueueMessage().
				SetMetadata("metadata").
				SetId("id").
				SetBody([]byte("data")),
			wantErr: false,
		},
		{
			name: "bad request - invalid type",
			connection: map[string]string{
				"address": "localhost:50000",
			},
			mockReceiver: &mockQueueReceiver{
				host:    "localhost",
				port:    50000,
				channel: "queues",
				timeout: 5,
			},
			req:      "bad-format",
			wantResp: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			recRequestCh := make(chan *kubemq.QueueMessage, 1)
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
			err := target.Init(ctx, tt.connection, "", nil)
			require.NoError(t, err)
			_, err = target.Do(ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			select {
			case gotRequest := <-recRequestCh:
				require.EqualValues(t, tt.wantResp.MessageID, gotRequest.MessageID)
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
		name       string
		connection config.Metadata
		wantErr    bool
	}{
		{
			name: "init",
			connection: map[string]string{
				"address":            "localhost:50000",
				"client_id":          "client_id",
				"auth_token":         "some-auth token",
				"channels":           "some-channel",
				"expiration_seconds": "0",
				"delay_seconds":      "0",
				"max_receive_count":  "1",
				"dead_letter_queue":  "",
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
				"address":            "localhost:40000",
				"client_id":          "client_id",
				"auth_token":         "some-auth token",
				"channels":           "some-channel",
				"expiration_seconds": "-1",
				"delay_seconds":      "0",
				"max_receive_count":  "0",
				"dead_letter_queue":  "",
			},
			wantErr: true,
		},
		{
			name: "init - bad expiration",
			connection: map[string]string{
				"address":            "localhost:50000",
				"client_id":          "client_id",
				"auth_token":         "some-auth token",
				"channels":           "some-channel",
				"expiration_seconds": "-1",
				"delay_seconds":      "0",
				"max_receive_count":  "0",
				"dead_letter_queue":  "",
			},
			wantErr: true,
		},
		{
			name: "init - bad delay",
			connection: map[string]string{
				"address":            "localhost:50000",
				"client_id":          "client_id",
				"auth_token":         "some-auth token",
				"channels":           "some-channel",
				"expiration_seconds": "0",
				"delay_seconds":      "-1",
				"max_receive_count":  "0",
				"dead_letter_queue":  "",
			},
			wantErr: true,
		},
		{
			name: "init - bad max receive count",
			connection: map[string]string{
				"address":            "localhost:50000",
				"client_id":          "client_id",
				"auth_token":         "some-auth token",
				"channels":           "some-channel",
				"expiration_seconds": "0",
				"delay_seconds":      "0",
				"max_receive_count":  "-1",
				"dead_letter_queue":  "",
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
