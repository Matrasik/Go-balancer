package balancer

import (
	"net/url"
)

type BackendServerInfo struct {
	Id        int      `json:"id"`
	Address   *url.URL `json:"-"`
	UrlString string   `json:"address"`
	isAlive   bool
}
