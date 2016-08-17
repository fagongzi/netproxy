package conf

// Conf conf
type Conf struct {
	APIAddr string   `json:"apiAddr,omitempty"`
	Proxys  []*Proxy `json:"proxys,omitempty"`
}

// Proxy proxy
type Proxy struct {
	Src            string `json:"src,omitempty"`
	Target         string `json:"target,omitempty"`
	TimeoutConnect int    `json:"timeoutConnect,omitempty"`
	TimeoutWrite   int    `json:"timeoutWrite,omitempty"`
}
