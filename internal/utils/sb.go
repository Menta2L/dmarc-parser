package utils

import (
	"fmt"
	"log"
	"net"
	"sort"
	"strings"

	"github.com/menta2l/dmarc-parser/internal/types"
)

//SenderbaseIPData query the senderbase to find out the org name of ip
func SenderbaseIPData(sip string) (sbGeo types.SBGeo, err error) {

	// convert from string input to net.IP:
	ip := net.ParseIP(sip).To4()
	if ip == nil {
		log.Println("ip6 address")
		return
	}

	// reverse the byte-order of IP:
	srevip := ByteReverseIP4(ip)

	// senderbase ip-specific domain to query:
	domain := fmt.Sprintf("%s.query.senderbase.org", srevip.String)

	// perform the lookup:
	log.Println("SB:  lookupTXT  ", sip)
	txtRecords, errLookupTXT := net.LookupTXT(domain)
	log.Println("SB:  lookupTXT2  ", sip)
	if errLookupTXT != nil {
		err = fmt.Errorf("SBIPD - errLookupTXT:  %s\n%s", domain, errLookupTXT)
		log.Println(err)
		return
	}
	if len(txtRecords) < 1 {
		err = fmt.Errorf("no TXT records found for IP %s\n%s", domain, sip)
		log.Println(err)
		return
	}

	rr := txtRecords[0]
	log.Println("SB:  TXT proc  ", rr)

	// handle multiple TXT records:
	sbStr := rr
	if len(txtRecords) > 1 {
		sort.Slice(txtRecords, func(i, j int) bool { return txtRecords[i][0] < txtRecords[j][0] })
		sbStr = ""
		for j := range txtRecords {
			// each TXT leads with '[0-9]-'...strip away this 2-char prefix
			txtRecords[j] = txtRecords[j][2:len(txtRecords[j])] // this could use regex improvement
			sbStr = fmt.Sprintf("%s%s", sbStr, txtRecords[j])
		}
	}
	sbFields := strings.Split(sbStr, "|")
	sbMap := map[string]string{}
	for j := range sbFields {
		sbm := strings.Split(sbFields[j], "=")
		sbMap[sbm[0]] = sbm[1]
	}

	log.Println("SB:  struct  ", sip)
	sbGeo.OrgName = sbMap["1"]
	sbGeo.OrgID = sbMap["4"]
	sbGeo.OrgCategory = sbMap["5"]
	sbGeo.Hostname = strings.ToLower(sbMap["20"])
	sbGeo.DomainName = strings.ToLower(sbMap["21"])
	sbGeo.HostnameMatchesIP = sbMap["22"]
	sbGeo.City = sbMap["50"]
	sbGeo.State = sbMap["51"]
	sbGeo.Country = sbMap["53"]
	sbGeo.Longitude = sbMap["54"]
	sbGeo.Latitude = sbMap["55"]

	return
}
