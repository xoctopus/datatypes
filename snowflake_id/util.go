package snowflake_id

import (
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"time"
)

func WorkerIDFromLocalIP() (uint32, error) {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = os.Getenv("HOSTNAME")
	}

	var ipv4 net.IP
	addresses, err := net.LookupIP(hostname)
	if err != nil {
		return 0, err
	}

	for _, addr := range addresses {
		if ipv4 = addr.To4(); ipv4 != nil {
			break
		}
	}
	return WorkerIDFromIP(ipv4), nil
}

func WorkerIDFromIP(ipv4 net.IP) uint32 {
	if ipv4 == nil {
		return 0
	}
	ip := ipv4.To4()
	return uint32(ip[2])<<8 + uint32(ip[3])
}

var _rand *rand.Rand

func init() {
	ts := uint64(time.Now().UnixNano())
	_rand = rand.New(rand.NewPCG(ts<<32, ts>>32))
}

func RandU32N(n uint32) uint32 {
	return _rand.Uint32N(n)
}

type errBeforeBaseTime [2]time.Time

func (e errBeforeBaseTime) Error() string {
	return fmt.Sprintf("timestamp: %s is before base timestamp %s", e[1], e[0])
}

type errOverMaxTimestamp [2]uint64

func (e errOverMaxTimestamp) Error() string {
	return fmt.Sprintf("elapsed: %d over max timestamp %d", e[1], e[0])
}
