package redir

import (
	"errors"
	"net"
	"net/netip"

	"github.com/metacubex/mihomo/common/nnip"
	"github.com/metacubex/mihomo/transport/socks5"
)

func ParserPacket(conn net.Conn) (netip.Addr, uint16, error) {
	target, err := parserPacket(conn)
	if err != nil {
		return netip.Addr{}, 0, err
	}
	return parseSocksAddr(target)
}

func parseSocksAddr(target socks5.Addr) (netip.Addr, uint16, error) {
	ip, port := netip.Addr{}, uint16(0)

	switch target[0] {
	// there is no need to handle socks5.domain name
	/*
		 case socks5.AtypDomainName:
			// trim for FQDN
			metadata.Host = strings.TrimRight(string(target[2:2+target[1]]), ".")
			metadata.DstPort = uint16((int(target[2+target[1]]) << 8) | int(target[2+target[1]+1]))
	*/
	case socks5.AtypIPv4:
		ip = nnip.IpToAddr(net.IP(target[1 : 1+net.IPv4len]))
		port = uint16((int(target[1+net.IPv4len]) << 8) | int(target[1+net.IPv4len+1]))
	case socks5.AtypIPv6:
		ip6, _ := netip.AddrFromSlice(target[1 : 1+net.IPv6len])
		ip = ip6.Unmap()
		port = uint16((int(target[1+net.IPv6len]) << 8) | int(target[1+net.IPv6len+1]))
	default:
		return ip, port, errors.New("unsupported address type")
	}

	return ip, port, nil
}
