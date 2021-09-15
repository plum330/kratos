package wrr

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/selector"
	"github.com/go-kratos/kratos/v2/selector/node/direct"
)

const (
	// Name is wrr balancer name
	Name = "wrr"
)

var _ selector.Balancer = &Balancer{} // Name is balancer name

// WithFilters with select filters
func WithFilters(filters ...selector.Filter) Option {
	return func(o *options) {
		o.filters = filters
	}
}

// Option is random builder option.
type Option func(o *options)

// options is random builder options
type options struct {
	filters []selector.Filter
}

// Balancer is a random balancer.
type Balancer struct {
	mu            sync.Mutex
	currentWeight map[string]float64
}

// New random a selector.
func New(opts ...Option) selector.Selector {
	var options options
	for _, o := range opts {
		o(&options)
	}

	return &selector.Default{
		Balancer: &Balancer{
			currentWeight: make(map[string]float64),
		},
		NodeBuilder: &direct.Builder{},
		Filters:     options.filters,
	}
}

// Pick pick a weighted node.
func (p *Balancer) Pick(_ context.Context, nodes []selector.WeightedNode) (selector.WeightedNode, selector.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, selector.ErrNoAvailable
	}
	var totalWeight float64
	var selected selector.WeightedNode
	var selectWeight float64

	// nginx wrr load balancing algorithm: http://blog.csdn.net/zhangskd/article/details/50194069
	p.mu.Lock()
	for _, node := range nodes {
		totalWeight += node.Weight()
		cwt := p.currentWeight[node.Address()]
		// current += effectiveWeight
		cwt += node.Weight()
		p.currentWeight[node.Address()] = cwt
		if selected == nil || selectWeight < cwt {
			selectWeight = cwt
			selected = node
		}
	}
	p.currentWeight[selected.Address()] = selectWeight - totalWeight
	p.mu.Unlock()

	d := selected.Pick()
	return selected, d, nil
}
