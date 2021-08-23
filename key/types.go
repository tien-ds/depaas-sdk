package key

import "fmt"

type Key struct {
	Pub  string `json:"pub"`
	Pri  string `json:"pri"`
	Addr string `json:"addr"`
}

func (k *Key) String() string {
	return fmt.Sprintf("pub %s \npri %s \naddr %s", k.Pub, k.Pri, k.Addr)
}
