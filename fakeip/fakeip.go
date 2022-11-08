package fakeip

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/xjasonlyu/tun2socks/v2/log"
)

type entry struct {
	IP   net.IP
	Port uint16
}

var records map[string]entry = make(map[string]entry)

// Init fake ip mapping
func Init(m map[string]string) {
	records = make(map[string]entry)
	for k, v := range m {
		segs := strings.Split(v, ":")
		if len(segs) != 2 {
			log.Warnf("Failed to split %s to IP and port, discard it", v)
			continue
		}

		ip := net.ParseIP(segs[0])
		if ip == nil {
			log.Warnf("Failed to parse %s to IP, discard %s", segs[0], v)
			continue
		}

		port, err := strconv.Atoi(segs[1])
		if err != nil {
			log.Warnf("Failed to parse %s to port, discard %s: %v", segs[1], v, err)
			continue
		}

		if port <= 0 || port > 65535 {
			log.Warnf("Invalid port number(%s), discard %s", segs[1], v)
			continue
		}

		records[k] = entry{
			IP:   ip,
			Port: uint16(port),
		}
	}
}

// Resolve the real IP and port by fake IP and port
func Resolve(fakeIP net.IP, fakePort uint16) (realIP net.IP, realPort uint16) {
	// TODO: performance issues
	fakeAddr := fmt.Sprintf("%s:%d", fakeIP.String(), fakePort)
	realAddr, ok := records[fakeAddr]
	if !ok {
		return fakeIP, fakePort
	}

	return realAddr.IP, realAddr.Port
}
