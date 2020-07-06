package utils

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/menta2l/dmarc-parser/internal/types"
)

func ByteReverseIP4(ip net.IP) (revip types.RevIP4) {

	for j := 0; j < len(ip); j++ {
		revip.Byte[len(ip)-j-1] = ip[j]
		revip.String = fmt.Sprintf("%d.%s", ip[j], revip.String)
	}

	revip.String = strings.TrimRight(revip.String, ".")

	return
}
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//ResolveAddrNames returns a struct containging an address and a list of names mapping to it
func ResolveAddrNames(addr string) (addrNames types.AddrNames, err error) {

	addrNames = types.AddrNames{
		Addr: addr,
	}

	names, errLookupAddr := net.LookupAddr(addr)
	if len(names) > 0 {
		addrNames.Names = names
	} else if errLookupAddr != nil {
		err = errLookupAddr
	}

	return
}
