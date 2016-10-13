package conf

import (
	"encoding/json"
	"io"
)

// Conf conf
type Conf struct {
	APIAddr string       `json:"apiAddr,omitempty"`
	Units   []*ProxyUnit `json:"units,omitempty"`
}

// Ctl ctl lossRate and so on
type Ctl struct {
	Address string   `json:"address"`
	In      *CtlUnit `json:"in"`
	Out     *CtlUnit `json:"out"`
}

// CtlUnit CtlUnit
type CtlUnit struct {
	LossRate int `json:"lossRate"`
	DelayMs  int `json:"delayMs"`
}

// UnMarshalCtlFromReader UnMarshalCtlFromReader
func UnMarshalCtlFromReader(r io.Reader) (*Ctl, error) {
	v := &Ctl{}

	decoder := json.NewDecoder(r)
	err := decoder.Decode(v)

	if nil != err {
		return nil, err
	}

	return v, nil
}

// Marshal marshal
func (c *Ctl) Marshal() []byte {
	d, _ := json.Marshal(c)
	return d
}

// CopyFrom copy form
func (c *Ctl) CopyFrom(from *Ctl) {
	c.In.LossRate = from.In.LossRate
	c.In.DelayMs = from.In.DelayMs

	c.Out.LossRate = from.Out.LossRate
	c.Out.DelayMs = from.Out.DelayMs
}

// ProxyUnit proxyUnit
type ProxyUnit struct {
	Src            string `json:"src,omitempty"`
	Target         string `json:"target,omitempty"`
	Desc           string `json:"desc,omitempty"`
	TimeoutConnect int    `json:"timeoutConnect,omitempty"`
	TimeoutWrite   int    `json:"timeoutWrite,omitempty"`
	Ctl            *Ctl   `json:"ctl,omitempty"`
}
