package observer

import (
	"context"
	"fmt"
	"sync"

	"github.com/victorwong171/go-utils/utils"
)

type Action[P any] func(ctx context.Context, params P) error

type Cfg[P any] struct {
	IsAsync bool
	Name    string
	Action  Action[P]
}

type Event[P any] interface {
	Register(observers ...Cfg[P])
	Emit(ctx context.Context, params P) error
}

type event[P any] struct {
	logger       utils.Logger
	mu           sync.RWMutex
	observerList []Cfg[P]
}

func NewEvent[P any](logger utils.Logger, observers ...Cfg[P]) Event[P] {
	return &event[P]{
		logger:       logger,
		mu:           sync.RWMutex{},
		observerList: observers,
	}
}

func (e *event[P]) Register(observers ...Cfg[P]) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.observerList = append(e.observerList, observers...)
}

func (e *event[P]) Emit(ctx context.Context, params P) error {
	for _, o := range e.observerList {
		if !o.IsAsync {
			e.logger.Info(fmt.Sprintf("[emit] sync event: %s", o.Name))
			if err := o.Action(ctx, params); err != nil {
				return err
			}
		} else {
			e.logger.Info(fmt.Sprintf("[emit] async event: %s", o.Name))
			go func(do Action[P], name string) {
				if err := do(utils.ContextCopy(ctx), params); err != nil {
					e.logger.Errorf(fmt.Sprintf("exec %s failed, err = %v", name, err))
				}
			}(o.Action, o.Name)
		}
	}
	return nil
}
