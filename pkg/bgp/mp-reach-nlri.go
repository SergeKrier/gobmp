package bgp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/ls"
	"github.com/sbezverk/gobmp/pkg/tools"
)

// MPReachNLRI defines an MP Reach NLRI object
type MPReachNLRI struct {
	AddressFamilyID      uint16
	SubAddressFamilyID   uint8
	NextHopAddressLength uint8
	NextHopAddress       []byte
	NLRI                 []byte
}

func (mp *MPReachNLRI) String() string {
	var s string
	s += fmt.Sprintf("Address Family ID: %d\n", mp.AddressFamilyID)
	s += fmt.Sprintf("Subsequent Address Family ID: %d\n", mp.SubAddressFamilyID)
	if mp.NextHopAddressLength == 4 {
		s += fmt.Sprintf("Next Hop Network Address: %s\n", net.IP(mp.NextHopAddress).To4().String())
	} else if mp.NextHopAddressLength == 16 {
		s += fmt.Sprintf("Next Hop Network Address: %s\n", net.IP(mp.NextHopAddress).To16().String())
	}
	switch mp.SubAddressFamilyID {
	case 71:
		nlri, err := ls.UnmarshalLSNLRI71(mp.NLRI)
		if err != nil {
			s += err.Error()
		} else {
			s += nlri.String()
		}
	default:
		s += fmt.Sprintf("NLRI: %s\n", tools.MessageHex(mp.NLRI))
	}

	return s
}

// GetNextHop return a string representation of the next hop ip address.
func (mp *MPReachNLRI) GetNextHop() string {
	if mp.NextHopAddressLength == 4 {
		return net.IP(mp.NextHopAddress).To4().String()
	} else if mp.NextHopAddressLength == 16 {
		return net.IP(mp.NextHopAddress).To16().String()
	}
	return "invalid"
}

// GetNLRI71 check for presense of NLRI 71 in the NLRI 14 NLRI data and if exists, instantiate NLRI71 object
func (mp *MPReachNLRI) GetNLRI71() (*ls.NLRI71, error) {
	if mp.SubAddressFamilyID == 71 {
		nlri71, err := ls.UnmarshalLSNLRI71(mp.NLRI)
		if err != nil {
			return nil, err
		}
		return nlri71, nil
	}

	// TODO return new type of errors to be able to check for the code
	return nil, fmt.Errorf("not found")
}

// MarshalJSON defines a custom method to convert MP REACH NLRI object into JSON object
func (mp *MPReachNLRI) MarshalJSON() ([]byte, error) {
	var jsonData []byte

	jsonData = append(jsonData, []byte("{\"AddressFamilyID\":")...)
	jsonData = append(jsonData, []byte(fmt.Sprintf("%d,", mp.AddressFamilyID))...)
	jsonData = append(jsonData, []byte("\"SubAddressFamilyID\":")...)
	jsonData = append(jsonData, []byte(fmt.Sprintf("%d,", mp.SubAddressFamilyID))...)
	jsonData = append(jsonData, []byte("\"NextHopAddressLength\":")...)
	jsonData = append(jsonData, []byte(fmt.Sprintf("%d,", mp.NextHopAddressLength))...)
	jsonData = append(jsonData, []byte("\"NextHopAddress\":")...)
	jsonData = append(jsonData, tools.RawBytesToJSON(mp.NextHopAddress)...)
	jsonData = append(jsonData, ',')
	jsonData = append(jsonData, []byte("\"NLRI\":")...)
	switch mp.SubAddressFamilyID {
	case 71:
		nlri, err := ls.UnmarshalLSNLRI71(mp.NLRI)
		if err != nil {
			return nil, err
		}
		b, err := json.Marshal(&nlri)
		if err != nil {
			return nil, err
		}
		jsonData = append(jsonData, b...)
	default:
		jsonData = append(jsonData, tools.RawBytesToJSON(mp.NLRI)...)
	}
	jsonData = append(jsonData, '}')

	return jsonData, nil
}

// UnmarshalMPReachNLRI builds MP Reach NLRI attributes
func UnmarshalMPReachNLRI(b []byte) (*MPReachNLRI, error) {
	glog.V(6).Infof("MPReachNLRI Raw: %s", tools.MessageHex(b))
	mp := MPReachNLRI{}
	p := 0
	mp.AddressFamilyID = binary.BigEndian.Uint16(b[p : p+2])
	p += 2
	mp.SubAddressFamilyID = uint8(b[p])
	p++
	mp.NextHopAddressLength = uint8(b[p])
	p++
	mp.NextHopAddress = b[p : p+int(mp.NextHopAddressLength)]
	p += int(mp.NextHopAddressLength)
	// Skip reserved byte
	p++
	mp.NLRI = make([]byte, len(b)-p)
	copy(mp.NLRI, b[p:len(b)])

	return &mp, nil
}
