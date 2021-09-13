package omit

import (
	"github.com/appbaseio/abc/importer/function"
	"github.com/appbaseio/abc/importer/message"
)

func init() {
	function.Add(
		"omit",
		func() function.Function {
			return &Omitter{}
		},
	)
}

type Omitter struct {
	Fields []string `json:"fields"`
}

func (o *Omitter) Apply(msg message.Msg) (message.Msg, error) {
	for _, k := range o.Fields {
		msg.Data().Delete(k)
	}
	return msg, nil
}
