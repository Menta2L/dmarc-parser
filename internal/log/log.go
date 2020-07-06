package log

import (
	"log/syslog"
	"os"

	"github.com/jedisct1/dlog"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func init() {
	dlog.Init("dmarc-parser", dlog.SeverityNotice, "DAEMON")
	dlog.UseSyslog(true)

	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)
	syslogOutput, err := lSyslog.NewSyslogHook("", "",
		syslog.LOG_INFO|syslog.LOG_DAEMON, "")
	if err != nil {
		log.Fatal("main: unable to setup syslog output")
	}
	log.AddHook(syslogOutput)

}

// GormLogger is a custom logger for Gorm, making it use logrus.
type GormLogger struct{}

// Print handles log events from Gorm for the custom logger.
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		log.WithFields(
			log.Fields{
				"module":  "gorm",
				"type":    "sql",
				"rows":    v[5],
				"src_ref": v[1],
				"values":  v[4],
			},
		).Debug(v[3])
	case "log":
		log.WithFields(log.Fields{"module": "gorm", "type": "log"}).Print(v[2])
	}
}
