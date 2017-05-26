package pretty

import (
	"encoding/json"
	"strings"

	"github.com/aviaryan/abc/function"
	"github.com/aviaryan/abc/log"
	"github.com/aviaryan/abc/message"
	"github.com/compose/mejson"
)

const (
	DefaultIndent = 2
)

var (
	DefaultPrettifier = &Prettify{Spaces: DefaultIndent}
)

func init() {
	function.Add(
		"pretty",
		func() function.Function {
			return DefaultPrettifier
		},
	)
}

type Prettify struct {
	Spaces int `json:"spaces"`
}

func (p *Prettify) Apply(msg message.Msg) (message.Msg, error) {
	d, _ := mejson.Unmarshal(msg.Data())
	b, _ := json.Marshal(d)
	if p.Spaces > 0 {
		b, _ = json.MarshalIndent(d, "", strings.Repeat(" ", p.Spaces))
	}
	log.Infof("\n%s", string(b))
	return msg, nil
}
