package pipeline

import (
	"testing"

	"github.com/boxproject/boxwallet/bccore"
)

var pipe = &Pipeline{}

func TestPipeline(t *testing.T) {
	t.Log("tilen1", len(p.ti))
	t.Log("kvlen1", len(p.kv))
	t.Log("utilen1", len(p.uti))
	defer t.Log("tilenDefer", len(p.ti))
	defer t.Log("kvlenDefer", len(p.kv))
	defer t.Log("utilenDefer", len(p.uti))
	pipe.Send(bccore.BC_BTC, "uuid-1", []string{"1", "2", "3"})
	for k, v := range p.uti {
		t.Log("uti1:", k)
		t.Log("uti1:", v)
	}
	for k, v := range p.ti {
		t.Log("ti1:", k)
		t.Log("ti1:", v)
	}
	for k, v := range p.kv {
		t.Log("kv1:", k)
		t.Log("kv1:", v)
	}

	t.Log("read", pipe.AddressExist(bccore.BC_BTC, []string{"1"}))
	t.Log("tilen2", len(p.ti))
	t.Log("kvlen2", len(p.kv))
	t.Log("utilen2", len(p.uti))
	ch, _ := pipe.CheckSend(bccore.BC_BTC, "txId-1", "uuid-1")
	t.Log("tilen3", len(p.ti))
	t.Log("kvlen3", len(p.kv))
	t.Log("utilen3", len(p.uti))
	go func() {
		for c := range ch {
			t.Log(c)
		}
	}()

	pipe.TxOver(bccore.BC_BTC, "txId-1", true)
	for k, v := range p.uti {
		t.Log("uti2:", k)
		t.Log("uti2:", v)
	}
	for k, v := range p.ti {
		t.Log("ti2:", k)
		t.Log("ti2:", v)
	}
	for k, v := range p.kv {
		t.Log("kv2:", k)
		t.Log("kv2:", v)
	}

}
