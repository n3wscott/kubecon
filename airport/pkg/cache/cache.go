package cache

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/n3wscott/kubecon/airport/pkg/events"
	"log"
	"strconv"
	"strings"
)

type Cache interface {
	Reset() error

	SetCarrierRoute(carrier string, route events.CarrierOfferData)
	GetCarrierRoute(carrier string) events.CarrierOfferData

	SetWarehouseOffers(warehouse string, offers events.CustomerOfferData)
	GetWarehouseOffers(warehouse string) events.CustomerOfferData

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

// Carrier

func carrierKey(carrier string) string {
	return strings.ToLower(fmt.Sprintf("carrier-%s", carrier))
}

func (c *memcacheClient) SetCarrierRoute(carrier string, route events.CarrierOfferData) {
	value, err := json.Marshal(route)
	if err != nil {
		log.Println("failed to marshal route,", err)
		return
	}

	if err := c.client.Set(&memcache.Item{Key: carrierKey(carrier), Value: value}); err != nil {
		log.Println("failed to set on memcached,", err.Error())
	}
}

func (c *memcacheClient) GetCarrierRoute(carrier string) events.CarrierOfferData {
	route := make(events.CarrierOfferData, 0)

	it, err := c.client.Get(carrierKey(carrier))
	if err != nil {
		log.Printf("failed to get on memcached, %s", err.Error())
		return route
	}

	err = json.Unmarshal(it.Value, &route)
	if err != nil {
		log.Println("failed to unmarshal route,", err)
	}
	return route
}

// Warehouse

func warehouseKey(warehouse string) string {
	return strings.ToLower(fmt.Sprintf("warehouse-%s", warehouse))
}

func (c *memcacheClient) SetWarehouseOffers(warehouse string, offers events.CustomerOfferData) {
	value, err := json.Marshal(offers)
	if err != nil {
		log.Println("failed to marshal offers,", err)
		return
	}

	if err := c.client.Set(&memcache.Item{Key: warehouseKey(warehouse), Value: value}); err != nil {
		log.Println("failed to set on memcached,", err.Error())
	}
}
func (c *memcacheClient) GetWarehouseOffers(warehouse string) events.CustomerOfferData {
	offers := make(events.CustomerOfferData, 0)

	it, err := c.client.Get(warehouseKey(warehouse))
	if err != nil {
		log.Printf("failed to get on memcached, %s", err.Error())
		return offers
	}

	err = json.Unmarshal(it.Value, &offers)
	if err != nil {
		log.Println("failed to unmarshal offers,", err)
	}
	return offers
}

// Products

func productKey(retailer string, product events.Product) string {
	return strings.ToLower(fmt.Sprintf("product-%s-%s", retailer, product))
}

func (c *memcacheClient) SetProductCount(retailer string, product events.Product, count int) {
	if err := c.client.Set(&memcache.Item{Key: productKey(retailer, product), Value: []byte(strconv.Itoa(count))}); err != nil {
		log.Println("failed to set on memcached,", err.Error())
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
