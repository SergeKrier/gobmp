package message

import (
	"github.com/golang/glog"
	"github.com/sbezverk/gobmp/pkg/bgp"
	"github.com/sbezverk/gobmp/pkg/bmp"
)

func (p *producer) processMPUpdate(nlri bgp.MPNLRI, operation int, ph *bmp.PerPeerHeader, update *bgp.Update) {
	labeled := false
	labeledSet := false

	switch nlri.GetAFISAFIType() {
	case 1:
		// MP_REACH_NLRI AFI 1 SAFI 1
		if !labeledSet {
			labeledSet = true
			labeled = false
		}
		fallthrough
	case 2:
		// MP_REACH_NLRI AFI 2 SAFI 1
		if !labeledSet {
			labeledSet = true
			labeled = false
		}
		fallthrough
	case 16:
		// MP_REACH_NLRI AFI 1 SAFI 4
		if !labeledSet {
			labeledSet = true
			labeled = true
		}
		fallthrough
	case 17:
		// MP_REACH_NLRI AFI 2 SAFI 4
		if !labeledSet {
			labeledSet = true
			labeled = true
		}
		msgs, err := p.unicast(nlri, operation, ph, update, labeled)
		if err != nil {
			glog.Errorf("failed to produce Unicast Prefix message message with error: %+v", err)
			return
		}
		// Loop through and publish all collected messages
		for _, m := range msgs {
			if err := p.marshalAndPublish(&m, bmp.UnicastPrefixMsg, []byte(m.RouterHash), false); err != nil {
				glog.Errorf("failed to process Unicast Prefix message with error: %+v", err)
				return
			}
		}
	case 18:
		msg, err := p.l3vpn(nlri, operation, ph, update)
		if err != nil {
			glog.Errorf("failed to produce l3vpn message with error: %+v", err)
			return
		}
		if err := p.marshalAndPublish(&msg, bmp.L3VPNMsg, []byte(msg.RouterHash), false); err != nil {
			glog.Errorf("failed to process L3VPN message with error: %+v", err)
			return
		}
	case 19:
		glog.Infof("2 IP (IP version 6) : 128 MPLS-labeled VPN address, attributes: %+v", update.GetAllAttributeID())
	case 24:
		msgs, err := p.evpn(nlri, operation, ph, update)
		if err != nil {
			glog.Errorf("failed to produce evpn message with error: %+v", err)
			return
		}
		for _, msg := range msgs {
			if err := p.marshalAndPublish(&msg, bmp.EVPNMsg, []byte(msg.RouterHash), false); err != nil {
				glog.Errorf("failed to process EVPNP message with error: %+v", err)
				return
			}
		}
	case 71:
		p.processNLRI71SubTypes(nlri, operation, ph, update)
	}
}

func (p *producer) processNLRI71SubTypes(nlri bgp.MPNLRI, operation int, ph *bmp.PerPeerHeader, update *bgp.Update) {
	// ipv4Flag used to differentiate between IPv4 and IPv6 Prefix NLRI messages
	ipv4Flag := false
	// NLRI 71 carries 6 known sub type
	ls, err := nlri.GetNLRI71()
	if err != nil {
		glog.Errorf("failed to NLRI 71 with error: %+v", err)
		return
	}
	t := ls.GetSubType()
	switch t {
	case 32:
		msg, err := p.lsNode(nlri, operation, ph, update)
		if err != nil {
			glog.Errorf("failed to produce ls_node message with error: %+v", err)
			return
		}
		if err := p.marshalAndPublish(&msg, bmp.LSNodeMsg, []byte(msg.RouterHash), false); err != nil {
			glog.Errorf("failed to process LSNode message with error: %+v", err)
			return
		}
	case 33:
		msg, err := p.lsLink(nlri, operation, ph, update)
		if err != nil {
			glog.Errorf("failed to produce ls_link message with error: %+v", err)
			return
		}
		if err := p.marshalAndPublish(&msg, bmp.LSLinkMsg, []byte(msg.RouterHash), false); err != nil {
			glog.Errorf("failed to process LSLink message with error: %+v", err)
			return
		}
	case 34:
		ipv4Flag = true
		fallthrough
	case 35:
		msg, err := p.lsPrefix(nlri, operation, ph, update, ipv4Flag)
		if err != nil {
			glog.Errorf("failed to produce ls_prefix message with error: %+v", err)
			return
		}
		if err := p.marshalAndPublish(&msg, bmp.LSPrefixMsg, []byte(msg.RouterHash), false); err != nil {
			glog.Errorf("failed to process LSPrefix message with error: %+v", err)
			return
		}
	case 36:
		msg, err := p.lsSRv6SID(nlri, operation, ph, update)
		if err != nil {
			glog.Errorf("failed to produce ls_srv6_sid message with error: %+v", err)
			return
		}
		if err := p.marshalAndPublish(&msg, bmp.LSSRv6SIDMsg, []byte(msg.RouterHash), false); err != nil {
			glog.Errorf("failed to process LSSRv6SID message with error: %+v", err)
			return
		}
	default:
		glog.Warningf("Unknown NLRI 71 Sub type %d", t)
	}
}
