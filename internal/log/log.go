package log

import (
	"github.com/jedisct1/dlog"
)

func init() {
	dlog.Init("dmarc-parser", dlog.SeverityNotice, "DAEMON")
	dlog.UseSyslog(true)

}
