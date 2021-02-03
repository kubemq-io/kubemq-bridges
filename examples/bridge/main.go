package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/pkg/uuid"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(1000) // n will be between 0 and 10
	time.Sleep(time.Duration(n) * time.Millisecond)
}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clientA, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30501),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientB, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30502),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientC, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30503),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}
	clientD, err := kubemq.NewClient(context.Background(),
		kubemq.WithAddress("localhost", 30504),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		errCh := make(chan error)
		queriesCh, err := clientB.SubscribeToQueries(ctx, "queries", "", errCh)
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case err := <-errCh:
				log.Fatal(err)
				return
			case query, more := <-queriesCh:
				if !more {
					log.Println("client b, done")
					return
				}
				log.Printf("client B: Query Received:\nId %s\nChannel: %s\nMetadata: %s\nBody: %s\n", query.Id, query.Channel, query.Metadata, query.Body)
				randomSleep()
				err := clientB.NewResponse().
					SetRequestId(query.Id).
					SetResponseTo(query.ResponseTo).
					SetExecutedAt(time.Now()).
					SetMetadata("response from client b").
					SetBody([]byte("got your query, you are good to go")).
					Send(ctx)
				if err != nil {
					log.Fatal(err)
				}
			case <-ctx.Done():
				return
			}
		}

	}()
	go func() {
		errCh := make(chan error)
		queriesCh, err := clientC.SubscribeToQueries(ctx, "queries", "", errCh)
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case err := <-errCh:
				log.Fatal(err)
				return
			case query, more := <-queriesCh:
				if !more {
					log.Println("client c, done")
					return
				}
				log.Printf("client C: Query Received:\nId %s\nChannel: %s\nMetadata: %s\nBody: %s\n", query.Id, query.Channel, query.Metadata, query.Body)
				randomSleep()
				err := clientC.NewResponse().
					SetRequestId(query.Id).
					SetResponseTo(query.ResponseTo).
					SetExecutedAt(time.Now()).
					SetMetadata("response from client c").
					SetBody([]byte("got your query, you are good to go")).
					Send(ctx)
				if err != nil {
					log.Fatal(err)
				}
			case <-ctx.Done():
				return
			}
		}

	}()
	go func() {
		errCh := make(chan error)
		queriesCh, err := clientD.SubscribeToQueries(ctx, "queries", "", errCh)
		if err != nil {
			log.Fatal(err)
		}
		for {
			select {
			case err := <-errCh:
				log.Fatal(err)
				return
			case query, more := <-queriesCh:
				if !more {
					log.Println("client d, done")
					return
				}
				randomSleep()
				log.Printf("client D: Query Received:\nId %s\nChannel: %s\nMetadata: %s\nBody: %s\n", query.Id, query.Channel, query.Metadata, query.Body)
				err := clientD.NewResponse().
					SetRequestId(query.Id).
					SetResponseTo(query.ResponseTo).
					SetExecutedAt(time.Now()).
					SetMetadata("response from client d").
					SetBody([]byte("got your query, you are good to go")).
					Send(ctx)
				if err != nil {
					log.Fatal(err)
				}
			case <-ctx.Done():
				return
			}
		}

	}()
	// give some time to connect a receiver
	time.Sleep(1 * time.Second)
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	counter := 0
	for {
		counter++
		response, err := clientA.NewQuery().
			SetId("some-query-id").
			SetChannel("queries").
			SetMetadata("").
			SetBody([]byte("hello kubemq - sending from client a query, please reply")).
			SetTimeout(10 * time.Second).
			Send(ctx)
		if err != nil {
			log.Println(fmt.Sprintf("error sending query %d, error: %s", counter, err))
		} else {
			log.Printf("Response Received:\nQueryID: %s\nExecutedAt:%s\nMetadata: %s\nBody: %s\n", response.QueryId, response.ExecutedAt, response.Metadata, response.Body)
		}

		select {
		case <-gracefulShutdown:
			break
		default:
			time.Sleep(time.Second)
		}
	}
}
