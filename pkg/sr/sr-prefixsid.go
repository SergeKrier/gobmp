package sr

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/tools"
)

// PrefixSIDTLV defines Prefix SID TLV Object
// https://tools.ietf.org/html/draft-ietf-idr-bgp-ls-segment-routing-ext-08#section-2.3.1
type PrefixSIDTLV struct {
	Flags     uint8
	Algorithm uint8
	SID       []byte
}

func (psid *PrefixSIDTLV) String() string {
	var s string
	s += fmt.Sprintf("   Flags: %02x\n", psid.Flags)
	s += fmt.Sprintf("   Algorithm: %d\n", psid.Algorithm)
	s += fmt.Sprintf("   SID: %s\n", tools.MessageHex(psid.SID))

	return s
}

// UnmarshalPrefixSIDTLV builds Prefix SID TLV Object
func UnmarshalPrefixSIDTLV(b []byte) (*PrefixSIDTLV, error) {
	glog.V(6).Infof("Prefix SID TLV Raw: %s", tools.MessageHex(b))
	psid := PrefixSIDTLV{}
	p := 0
	psid.Flags = b[p]
	p++
	psid.Algorithm = b[p]
	p++
	// SID length would be Length of b - Flags 1 byte - Algorithm 1 byte - 2 bytes Reserved
	sl := len(b) - 4
	psid.SID = make([]byte, len(b)-4)
	p += 2
	copy(psid.SID, b[p:p+sl])

	return &psid, nil
}
