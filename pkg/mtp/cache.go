package mtp

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type cacheHandler struct {
	c *cache.Cache
}

type CError struct {
	Code uint32
	Msg  string
}

type CInstance struct {
	Path        string
	UniqueKeys  map[string]string
	OpIsSuccess bool
	OpErrStr    string
}
type CParamSetResult struct {
	Path        string
	OpIsSuccess bool
	OpErrStr    string
}

func (m *Mtp) cacheInit() error {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": m.Cfg.Cache.ServerAddr,
		},
		HeartbeatFrequency: time.Hour,
	})

	c := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	m.cacheH.c = c
	return nil
}

func (m *Mtp) cacheSetError(epId string, msgId string, data *CError) error {
	if m.cacheH.c == nil {
		log.Println("Cache is not initialized")
		return errors.New("Cache not initalized")
	}

	item := &cache.Item{
		Ctx:   context.TODO(),
		Key:   epId + msgId + "e",
		TTL:   time.Hour,
		Value: data,
	}

	log.Printf("Cache set key:%v, value:%+v\n", item.Key, item.Value)
	if err := m.cacheH.c.Set(item); err != nil {
		log.Println("Error in writing to Cache, err:", err)
		return err
	}
	return nil
}

func (m *Mtp) cacheGetError(epId string, msgId string) (*CError, error) {
	if m.cacheH.c == nil {
		log.Println("Cache is not initialized")
		return nil, errors.New("Cache not initalized")
	}

	key := epId + msgId + "e"
	data := &CError{}

	log.Println("Getting from cache with key:", key)
	if err := m.cacheH.c.Get(context.TODO(), key, data); err != nil {
		log.Println("Cache No data found, err:", err)
		return nil, err
	}
	log.Printf("Cache hit errorMsg: %v\n", data.Msg)
	return data, nil
}

func (m *Mtp) cacheSetInstance(epId string, msgId string, data *CInstance) error {

	if m.cacheH.c == nil {
		log.Println("Cache is not initialized")
		return errors.New("Cache not initalized")
	}

	item := &cache.Item{
		Ctx:   context.TODO(),
		Key:   epId + msgId,
		TTL:   time.Hour,
		Value: data,
	}

	if err := m.cacheH.c.Set(item); err != nil {
		log.Printf("Cache set key:%v, value:%+v\n", item.Key, item.Value)
		if err := m.cacheH.c.Set(item); err != nil {
			log.Println("Error in writing to Cache, err:", err)
			return err
		}
	}
	return nil
}

func (m *Mtp) cacheGetInstance(epId string, msgId string) (*CInstance, error) {
	if m.cacheH.c == nil {
		log.Println("Cache is not initialized")
		return nil, errors.New("Cache not initalized")
	}

	key := epId + msgId
	data := &CInstance{}

	log.Println("Getting from cache with key:", key)
	if err := m.cacheH.c.Get(context.TODO(), key, data); err != nil {
		log.Println("Error in reading from Cache, err:", err)
		return nil, err
	}
	log.Printf("Cache hit instance: %+v\n", data)
	return data, nil
}

func (m *Mtp) cacheSetParamSetResult(epId string, msgId string, data *CParamSetResult) error {

	if m.cacheH.c == nil {
		log.Println("Cache is not initialized")
		return errors.New("Cache not initalized")
	}

	item := &cache.Item{
		Ctx:   context.TODO(),
		Key:   epId + msgId,
		TTL:   time.Hour,
		Value: data,
	}

	log.Println("Setting set param result into cache")
	if err := m.cacheH.c.Set(item); err != nil {
		log.Printf("Cache set key:%v, value:%+v\n", item.Key, item.Value)
		if err := m.cacheH.c.Set(item); err != nil {
			log.Println("Error in writing to Cache, err:", err)
			return err
		}
	}
	return nil
}

func (m *Mtp) cacheGetParamSetResult(epId string, msgId string) (*CParamSetResult, error) {
	if m.cacheH.c == nil {
		log.Println("Cache is not initialized")
		return nil, errors.New("Cache not initalized")
	}

	key := epId + msgId
	data := &CParamSetResult{}

	log.Println("Getting from cache with key:", key)
	if err := m.cacheH.c.Get(context.TODO(), key, data); err != nil {
		log.Println("Error in reading from Cache, err:", err)
		return nil, err
	}
	log.Printf("Cache hit setParamResult: %+v\n", data)
	return data, nil
}
