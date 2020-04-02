package message

import (
	"net"

	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/bgp"
	"github.com/sbezverk/gobmp/pkg/bmp"
)

func (p *producer) lsNode(operation string, ph *bmp.PerPeerHeader, update *bgp.Update) (*LSNode, error) {
	nlri14, err := update.GetNLRI14()
	if err != nil {
		return nil, err
	}
	nlri71, err := nlri14.GetNLRI71()
	if err != nil {
		return nil, err
	}
	msg := LSNode{
		Action:       operation,
		RouterHash:   p.speakerHash,
		RouterIP:     p.speakerIP,
		BaseAttrHash: update.GetBaseAttrHash(),
		PeerHash:     ph.GetPeerHash(),
		PeerASN:      ph.PeerAS,
		Timestamp:    ph.PeerTimestamp,
	}
	msg.Nexthop = nlri14.GetNextHop()
	if ph.FlagV {
		// IPv6 specific conversions
		msg.PeerIP = net.IP(ph.PeerAddress).To16().String()
	} else {
		// IPv4 specific conversions
		msg.PeerIP = net.IP(ph.PeerAddress[12:]).To4().String()
	}
	// Processing other nlri and attributes, since they are optional, processing only if they exist
	node, err := nlri71.GetNodeNLRI()
	if err == nil {
		msg.Protocol = node.GetProtocolID()
		msg.IGPRouterID = node.GetIGPRouterID()
		msg.LSID = node.GetLSID()
		msg.ASN = node.GetASN()
		msg.OSPFAreaID = node.GetOSPFAreaID()
	}

	lsnode, err := update.GetNLRI29()
	if err == nil {
		msg.Flags = lsnode.GetNodeFlags()
		msg.Name = lsnode.GetNodeName()
		msg.MTID = lsnode.GetMTID()
		msg.ISISAreaID = lsnode.GetISISAreaID()
		if ph.FlagV {
			msg.RouterID = lsnode.GetNodeIPv6RouterID()
		} else {
			msg.RouterID = lsnode.GetNodeIPv4RouterID()
		}
		msg.NodeMSD = lsnode.GetNodeMSD()
		msg.SRCapabilities = lsnode.GetNodeSRCapabilities()
		msg.SRAlgorithm = lsnode.GetSRAlgorithm()
		msg.SRLocalBlock = lsnode.GetNodeSRLocalBlock()
		msg.SRv6CapabilitiesTLV = lsnode.GetNodeSRv6CapabilitiesTLV()
	}

	if count, path := update.GetAttrASPathString(p.as4Capable); count != 0 {
		msg.ASPath = path
	}
	if med := update.GetAttrMED(); med != nil {
		msg.MED = *med
	}
	glog.V(6).Infof("LS Node messages: %+v", msg)

	return &msg, nil
}
