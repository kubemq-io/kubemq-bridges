package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-bridges/pkg/metrics"
	"net/http"
	"sync"
	"time"
)

const (
	addRetryInterval = 1 * time.Second
)

type Service struct {
	bindings          sync.Map
	log               *logger.Logger
	exporter          *metrics.Exporter
	currentCtx        context.Context
	currentCancelFunc context.CancelFunc
	bindingStatus     sync.Map
	cfg               *config.Config
}

func New() (*Service, error) {
	s := &Service{
		bindings:      sync.Map{},
		log:           logger.NewLogger("binding-service"),
		bindingStatus: sync.Map{},
	}
	var err error
	s.exporter, err = metrics.NewExporter()
	if err != nil {
		return nil, fmt.Errorf("failed to to initialized metrics exporter, %w", err)
	}
	return s, nil
}

func NewForExternal() (*Service, error) {
	s := &Service{
		bindings:      sync.Map{},
		log:           logger.NewLogger("bridges-service"),
		bindingStatus: sync.Map{},
	}
	return s, nil
}
func (s *Service) Start(ctx context.Context, cfg *config.Config) error {
	s.cfg = cfg
	s.currentCtx, s.currentCancelFunc = context.WithCancel(ctx)
	if len(cfg.Bindings) == 0 {
		return nil
	}
	for _, bindingCfg := range cfg.Bindings {
		go func(ctx context.Context, cfg config.BindingConfig, logLevel string) {
			err := s.Add(ctx, cfg, logLevel)
			if err == nil {
				return
			} else {
				s.log.Errorf("failed to initialized binding, %s", err.Error())
			}
			count := 0
			for {
				select {
				case <-time.After(addRetryInterval):
					count++
					err := s.Add(ctx, cfg, logLevel)
					if err != nil {
						s.log.Errorf("failed to initialized binding: %s, attempt: %d, error: %s", cfg.Name, count, err.Error())
					} else {
						return
					}
				case <-ctx.Done():
					return
				}
			}

		}(s.currentCtx, bindingCfg, cfg.LogLevel)

	}
	return nil
}

func (s *Service) Stop() {
	if s.currentCancelFunc != nil {
		s.currentCancelFunc()
	}
	s.bindings.Range(func(key, value interface{}) bool {
		binder := value.(*Binder)
		err := s.Remove(binder.name)
		if err != nil {
			s.log.Error(err)
		}
		return true
	})

}
func (s *Service) Add(ctx context.Context, cfg config.BindingConfig, logLevel string) error {

	binder := NewBinder()
	status := newStatus(cfg)
	s.bindingStatus.Store(cfg.Name, status)
	err := binder.Init(ctx, cfg, s.exporter, logLevel)
	if err != nil {
		return err
	}
	err = binder.Start(ctx)
	if err != nil {
		return err
	}
	s.bindings.Store(cfg.Name, binder)
	status.Ready = true
	s.bindingStatus.Store(cfg.Name, status)
	return nil
}

func (s *Service) Remove(name string) error {
	val, ok := s.bindings.Load(name)
	if !ok {
		return fmt.Errorf("binding %s not found", name)
	}
	binder := val.(*Binder)
	err := binder.Stop()
	if err != nil {
		return err
	}
	s.bindings.Delete(name)
	s.bindingStatus.Delete(name)
	return nil
}

func (s *Service) PrometheusHandler() http.Handler {
	return s.exporter.PrometheusHandler()
}
func (s *Service) Stats() []*metrics.Report {
	return s.exporter.Store.List()
}
func (s *Service) GetStatus() []*Status {
	var list []*Status
	for _, binding := range s.cfg.Bindings {
		val, ok := s.bindingStatus.Load(binding.Name)
		if ok {
			list = append(list, val.(*Status))
		}
	}
	return list
}
