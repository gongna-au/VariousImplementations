package SingleChannelSchedue

import (
	"github.com/VariousImplementations/SingleChannelSchedue/pkg"
	"testing"
)

func TestHRRF(t *testing.T) {
	ClientHRRF()
}
func TestFCFS(t *testing.T) {
	ClientFCFS()
}

func TestSJF(t *testing.T) {
	ClientSJF()
}

func TestHPF(t *testing.T) {
	ClientHPF()
}
