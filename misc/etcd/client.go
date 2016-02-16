// client.go
package etcd

import (
	"errors"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"github.com/gbember/util"
)

type Client struct {
	ec     client.Client
	kapi   client.KeysAPI
	path   string
	logFun func(interface{}, ...interface{}) error
	handle func(key, value string) //处理函数
}

func NewClient(addrs []string, path string, handle func(key, value string), logFun func(interface{}, ...interface{}) error) (*Client, error) {
	cfg := client.Config{
		Endpoints:               addrs,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	kapi := client.NewKeysAPI(c)
	cl := &Client{
		ec:     c,
		kapi:   kapi,
		logFun: logFun,
		handle: handle,
	}
	return cl, nil
}

func (c *Client) Set(key, val string) error {
	_, err := c.kapi.Set(context.Background(), c.path+key, val, nil)
	return err
}

func (c *Client) Get(key string) (val string, err error) {
	var resp *client.Response
	resp, err = c.kapi.Get(context.Background(), c.path+key, nil)
	if err != nil {
		return
	}
	val = resp.Node.Value
	return
}

func (c *Client) Update(key, val string) error {
	_, err := c.kapi.Update(context.Background(), c.path+key, val)
	return err
}

func (c *Client) Delete(key string) error {
	_, err := c.kapi.Delete(context.Background(), c.path+key, nil)
	return err
}

func (c *Client) Create(key, value string) error {
	_, err := c.kapi.Create(context.Background(), c.path+key, value)
	return err
}

func (c *Client) LoadAll() error {
	resp, err := c.kapi.Get(context.Background(), c.path, &client.GetOptions{Recursive: true})
	if err != nil {
		return err
	}
	if !resp.Node.Dir {
		return errors.New("etcd client path is not dir")
	}
	return c.load_node_all(resp.Node)
}

func (c *Client) load_node_all(node *client.Node) error {
	if node.Dir {
		for _, n := range node.Nodes {
			err := c.load_node_all(n)
			if err != nil {
				return err
			}
		}
	} else {
		c.handle(node.Key, node.Value)
	}
	return nil
}

func (c *Client) Watch() {
	defer util.PrintPanicStack(c.logFun)
	watcher := c.kapi.Watcher(c.path, &client.WatcherOptions{
		Recursive: true,
	})
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			break
		}
		if res.Action == "set" || res.Action == "update" {
			c.handle(res.Node.Key, res.Node.Value)
		}
	}
}
