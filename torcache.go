package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type TorCache struct {
	stopped         bool
	lastRetrieval   *time.Time
	badguys         *map[string]string
	wg              *sync.WaitGroup
	pattern         *regexp.Regexp
	lookupLocations []string
}

func CreateTorCache(lookup []string) *TorCache {
	pat := regexp.MustCompile("([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})")
	retVal := &TorCache{
		stopped:         false,
		lastRetrieval:   nil,
		badguys:         nil,
		pattern:         pat,
		wg:              nil,
		lookupLocations: lookup,
	}
	return retVal
}

func (cache *TorCache) performLookup() {
	defer cache.wg.Done()
	newMap := make(map[string]string)
	for _, url := range cache.lookupLocations {
		if strings.Index(url, "htt") != 0 {
			continue
		}
		fmt.Println("Looking up " + url)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		matches := cache.pattern.FindAll(body, -1)
		if matches != nil {
			for i := range matches {
				m := matches[i]
				newMap[string(m[:])] = url
			}
		}
	}
	cache.badguys = &newMap
}

func (cache *TorCache) getCache() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("performLookup() failed!")
		}
	}()
	if cache.wg == nil {
		var localWG sync.WaitGroup
		cache.wg = &localWG
		cache.wg.Add(1)
		n := time.Now()
		cache.lastRetrieval = &n
		go cache.performLookup()
	}
	cache.wg.Wait()
	return nil
}
func (cache *TorCache) GetNumberOfIPs() int {
	if cache.badguys == nil {
		return 0
	}
	return len(*cache.badguys)
}
func (cache *TorCache) Check(host string) (string, error) {
	if cache.lastRetrieval == nil || time.Since(*cache.lastRetrieval) > (time.Second*30) {
		cErr := cache.getCache()
		if cErr != nil {
			return "", cErr
		}
	}

	res, err := net.LookupHost(host)
	if err == nil {
		for _, v := range res {
			match, contains := (*cache.badguys)[v]
			if contains {
				return match, nil
			}
		}
		return "", nil
	} else {
		return "", err
	}
}
