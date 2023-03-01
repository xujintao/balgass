package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/xujintao/balgass/cmd/client_spider/conf"
)

var client *http.Client

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (urlProxy *url.URL, err error) {
				for _, task := range conf.Tasks {
					if task.Host == req.URL.Host && task.Proxy != "" {
						urlProxy, err = url.Parse(task.Proxy)
						if err != nil {
							return nil, fmt.Errorf("proxy string parse failed, %v", err)
						}
						switch urlProxy.Scheme {
						case "http", "https", "socks5":
						default:
							return nil, fmt.Errorf("proxy scheme invalid")
						}
						return urlProxy, nil
					}
				}
				return
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func main() {
	var wg sync.WaitGroup
	for _, task := range conf.Tasks {
		wg.Add(1)
		go func(task *conf.Task) {
			defer wg.Done()
			u := url.URL{
				Scheme: task.Scheme,
				Host:   task.Host,
			}
			do := func(host string) {
				for _, path := range task.Paths {
					u.Path = path
					req, err := http.NewRequest("GET", u.String(), nil)
					if err != nil {
						log.Println(err)
						return
					}
					req.Host = host
					req.Header.Set("User-Agent", "another agent")
					res, err := client.Do(req)
					if err != nil {
						log.Println(err)
						return
					}
					res.Body.Close()
					log.Printf("request[%s %s %s] host[%s] status[%d] body_bytes_recv[%d]", req.Method, path, req.Proto, host, res.StatusCode, res.ContentLength)
				}
			}
			if len(task.Hosts) != 0 {
				for _, host := range task.Hosts {
					do(host)
					time.Sleep(1 * time.Second)
				}
			} else {
				do(task.Host)
			}
		}(task)
	}
	wg.Wait()
}
