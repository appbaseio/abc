package function

import (
	"github.com/aviaryan/abc/log"
	"github.com/aviaryan/abc/message"
)

var (
	_ Function = &Mock{}
)

type Mock struct {
	ApplyCount int
	Err        error
}

func (m *Mock) Apply(msg message.Msg) (message.Msg, error) {
	m.ApplyCount++
	log.With("apply_count", m.ApplyCount).With("err", m.Err).Debugln("applying...")
	return msg, m.Err
}
