package types

type SBGeo struct {
	OrgName           string `json:"org_name" db:"org_name"`
	OrgID             string `json:"org_id" db:"org_id"`
	OrgCategory       string `json:"org_category" db:"org_category"`
	Hostname          string `json:"hostname" db:"hostname"`
	DomainName        string `json:"domain_name" db:"domain_name"`
	HostnameMatchesIP string `json:"hostname_matches_ip" db:"hostname_matches_ip"`
	City              string `json:"city" db:"city"`
	State             string `json:"state" db:"state"`
	Country           string `json:"country" db:"country"`
	Longitude         string `json:"longitude" db:"longitude"`
	Latitude          string `json:"latitude" db:"latitude"`
}
