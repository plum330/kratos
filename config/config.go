package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"

	// init encoding
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"
	_ "github.com/go-kratos/kratos/v2/encoding/xml"
	_ "github.com/go-kratos/kratos/v2/encoding/yaml"
)

var (
	// ErrNotFound is key not found.
	ErrNotFound = errors.New("key not found")
	// ErrTypeAssert is type assert error.
	ErrTypeAssert = errors.New("type assert error")

	_ Config = (*config)(nil)
)

// Observer is config observer.
type Observer func(string, Value)

// Config is a config interface.
type Config interface {
	Load() error
	Scan(v interface{}) error
	Value(key string) Value
	Watch(key string, o Observer) error
	Close() error
}

type config struct {
	opts      options
	reader    Reader
	cached    sync.Map
	observers sync.Map
	watchers  []Watcher
	log       *log.Helper
}

// New new a config with options.
func New(opts ...Option) Config {
	options := options{
		logger: log.DefaultLogger,
		decoder: func(src *KeyValue, target map[string]interface{}) error {
			if src.Format == "" {
				target[src.Key] = src.Value
				return nil
			}
			if codec := encoding.GetCodec(src.Format); codec != nil {
				return codec.Unmarshal(src.Value, &target)
			}
			return fmt.Errorf("unsupported key: %s format: %s", src.Key, src.Format)
		},
		resolver: defaultResolver,
	}
	for _, o := range opts {
		o(&options)
	}
	return &config{
		opts:   options,
		reader: newReader(options),
		log:    log.NewHelper(options.logger),
	}
}

// defaultResolver resolve placeholder (${key:value} or $key) in map value
func defaultResolver(input map[string]interface{}) error {
	mapper := func(name string) string {
		args := strings.Split(strings.TrimSpace(name), ":")
		if v, has := readValue(input, args[0]); has {
			s, _ := v.String()
			return s
		} else if len(args) > 1 { // default value
			return args[1]
		}
		return ""
	}

	var resolve func(map[string]interface{}) error
	resolve = func(sub map[string]interface{}) error {
		for k, v := range sub {
			switch vt := v.(type) {
			case string:
				sub[k] = os.Expand(vt, mapper)
			case map[string]interface{}:
				if err := resolve(vt); err != nil {
					return err
				}
			case []interface{}:
				for i, iface := range vt {
					if s, ok := iface.(string); ok {
						vt[i] = os.Expand(s, mapper)
					}
				}
				sub[k] = vt
			}
		}
		return nil
	}
	return resolve(input)
}

func (c *config) watch(w Watcher) {
	for {
		kvs, err := w.Next()
		if err != nil {
			time.Sleep(time.Second)
			c.log.Errorf("Failed to watch next config: %v", err)
			continue
		}
		if err := c.reader.Merge(kvs...); err != nil {
			c.log.Errorf("Failed to merge next config: %v", err)
			continue
		}
		c.cached.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(Value)
			if n, ok := c.reader.Value(k); ok && !reflect.DeepEqual(n.Load(), v.Load()) {
				v.Store(n.Load())
				if o, ok := c.observers.Load(k); ok {
					o.(Observer)(k, v)
				}
			}
			return true
		})
	}
}

func (c *config) Load() error {
	for _, src := range c.opts.sources {
		kvs, err := src.Load()
		if err != nil {
			return err
		}
		if err := c.reader.Merge(kvs...); err != nil {
			c.log.Errorf("Failed to merge config source: %v", err)
			return err
		}
		w, err := src.Watch()
		if err != nil {
			c.log.Errorf("Failed to watch config source: %v", err)
			return err
		}
		go c.watch(w)
	}
	return nil
}

func (c *config) Value(key string) Value {
	if v, ok := c.cached.Load(key); ok {
		return v.(Value)
	}
	if v, ok := c.reader.Value(key); ok {
		c.cached.Store(key, v)
		return v
	}
	return &errValue{err: ErrNotFound}
}

func (c *config) Scan(v interface{}) error {
	data, err := c.reader.Source()
	if err != nil {
		return err
	}
	return unmarshalJSON(data, v)
}

func (c *config) Watch(key string, o Observer) error {
	if v := c.Value(key); v.Load() == nil {
		return ErrNotFound
	}
	c.observers.Store(key, o)
	return nil
}

func (c *config) Close() error {
	for _, w := range c.watchers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	return nil
}
