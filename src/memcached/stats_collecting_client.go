package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	stats "github.com/lyft/gostats"
)

type statsCollectingClient struct {
	c Client

	multiGetSuccess  stats.Counter
	multiGetError    stats.Counter
	incrementSuccess stats.Counter
	incrementMiss    stats.Counter
	incrementError   stats.Counter
	addSuccess       stats.Counter
	addError         stats.Counter
	addNotStored     stats.Counter
	keysRequested    stats.Counter
	keysFound        stats.Counter
}

func CollectStats(c Client, scope stats.Scope) Client {
	return statsCollectingClient{
		c:                c,
		multiGetSuccess:  scope.NewCounterWithTags("multiget", map[string]string{"code": "success"}),
		multiGetError:    scope.NewCounterWithTags("multiget", map[string]string{"code": "error"}),
		incrementSuccess: scope.NewCounterWithTags("increment", map[string]string{"code": "success"}),
		incrementMiss:    scope.NewCounterWithTags("increment", map[string]string{"code": "miss"}),
		incrementError:   scope.NewCounterWithTags("increment", map[string]string{"code": "error"}),
		addSuccess:       scope.NewCounterWithTags("add", map[string]string{"code": "success"}),
		addError:         scope.NewCounterWithTags("add", map[string]string{"code": "error"}),
		addNotStored:     scope.NewCounterWithTags("add", map[string]string{"code": "not_stored"}),
		keysRequested:    scope.NewCounter("keys_requested"),
		keysFound:        scope.NewCounter("keys_found"),
	}
}

func (s statsCollectingClient) GetMulti(keys []string) (map[string]*memcache.Item, error) {
	s.keysRequested.Add(uint64(len(keys)))

	results, err := s.c.GetMulti(keys)

	if err != nil {
		s.multiGetError.Inc()
	} else {
		s.keysFound.Add(uint64(len(results)))
		s.multiGetSuccess.Inc()
	}

	return results, err
}

func (s statsCollectingClient) Increment(key string, delta uint64) (newValue uint64, err error) {
	newValue, err = s.c.Increment(key, delta)
	switch err {
	case memcache.ErrCacheMiss:
		s.incrementMiss.Inc()
	case nil:
		s.incrementSuccess.Inc()
	default:
		s.incrementError.Inc()
	}
	return
}

func (s statsCollectingClient) Add(item *memcache.Item) error {
	err := s.c.Add(item)

	switch err {
	case memcache.ErrNotStored:
		s.addNotStored.Inc()
	case nil:
		s.addSuccess.Inc()
	default:
		s.addError.Inc()
	}

	return err
}
