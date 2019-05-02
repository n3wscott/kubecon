package cache

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/n3wscott/kubecon/airport/pkg/events"
	"log"
	"strconv"
	"strings"
)

type Cache interface {
	Reset() error
	Get()
	Set()

	SetCarrierRoute(carrier string, route events.CarrierOfferData)
	GetCarrierRoute(carrier string) events.CarrierOfferData

	SetProductCount(retailer string, product events.Product, count int)
	GetProductCount(retailer string, product events.Product) int
	AdjustProductCount(retailer string, product events.Product, by int) int
}

var _ Cache = (*memcacheClient)(nil)

type memcacheClient struct {
	client *memcache.Client
}

func NewCache(client *memcache.Client) Cache {
	return &memcacheClient{
		client: client,
	}
}

func (c *memcacheClient) Reset() error {
	return c.client.DeleteAll()
}

func (c *memcacheClient) Get() {
	it, err := c.client.Get("foo")
	if err != nil {
		log.Printf("failed to get on memcached, %s", err.Error())
	}
	log.Println("Memcache get:", it)

}

func (c *memcacheClient) Set() {
	if err := c.client.Set(&memcache.Item{Key: "foo", Value: []byte("my value")}); err != nil {
		log.Fatalf("failed to set on memcached, %s", err.Error())
	}
}

// Carrier

func carrierKey(carrier string) string {
	return strings.ToLower(fmt.Sprintf("carrier-%s", carrier))
}

func (c *memcacheClient) SetCarrierRoute(carrier string, route events.CarrierOfferData) {

}

func (c *memcacheClient) GetCarrierRoute(carrier string) events.CarrierOfferData {
	return nil
}

// Products

func productKey(retailer string, product events.Product) string {
	return strings.ToLower(fmt.Sprintf("product-%s-%s", retailer, product))
}

func (c *memcacheClient) SetProductCount(retailer string, product events.Product, count int) {
	if err := c.client.Set(&memcache.Item{Key: productKey(retailer, product), Value: []byte(strconv.Itoa(count))}); err != nil {
		log.Fatalf("failed to set on memcached, %s", err.Error())
	}
}

func (c *memcacheClient) GetProductCount(retailer string, product events.Product) int {
	it, err := c.client.Get(productKey(retailer, product))
	if err != nil {
		log.Printf("failed to get on memcached, %s", err.Error())
		return 0
	}
	count, err := strconv.Atoi(string(it.Value))
	if err != nil {
		log.Printf("failed to convert count to int, %s", err.Error())
		return 0
	}
	return count
}

func (c *memcacheClient) AdjustProductCount(retailer string, product events.Product, by int) int {
	count := c.GetProductCount(retailer, product)
	count += by
	c.SetProductCount(retailer, product, count)
	return count
}
