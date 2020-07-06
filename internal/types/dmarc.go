package types

import "github.com/jinzhu/gorm"

type DmarcPOReason struct {
	Reason   string `xml:"type" db:"Reason"`
	Comment  string `xml:"comment" db:"Comment"`
	RecordID uint
}

type DmarcDKIMAuthResult struct {
	Domain      string `xml:"domain" db:"Domain"`
	Selector    string `xml:"selector" db:"Selector"`
	Result      string `xml:"result" db:"Result"`
	HumanResult string `xml:"human_result" db:"HumanResult"`
	RecordID    uint
}

type DmarcSPFAuthResult struct {
	Domain   string `xml:"domain" db:"Domain"`
	Scope    string `xml:"scope" db:"Scope"`
	Result   string `xml:"result" db:"Result"`
	RecordID uint
}

type DmarcReport struct {
	gorm.Model
	MessageId         string        `db:"MessageId" gorm:"type:varchar(255);unique;not null"`
	Organization      string        `xml:"report_metadata>org_name" db:"Organization" gorm:"type:varchar(255);not null"`
	Email             string        `xml:"report_metadata>email" db:"Email" gorm:"type:varchar(255);not null"`
	ExtraContact      string        `xml:"report_metadata>extra_contact_info" db:"ExtraContact" gorm:"type:varchar(100)"` // minOccurs="0"
	ReportID          string        `xml:"report_metadata>report_id" db:"ReportID" gorm:"type:varchar(100);not null"`
	RawDateRangeBegin string        `xml:"report_metadata>date_range>begin" db:"RawDateRangeBegin" gorm:"-"`
	RawDateRangeEnd   string        `xml:"report_metadata>date_range>end" db:"RawDateRangeEnd" gorm:"-"`
	DateRangeBegin    int64         `db:"DateRangeBegin"`
	DateRangeEnd      int64         `db:"DateRangeEnd"`
	Errors            []string      `xml:"report_metadata>error" db:"Errors" gorm:"-"`
	Domain            string        `xml:"policy_published>domain" gorm:"type:varchar(255);not null;index"`
	AlignDKIM         string        `xml:"policy_published>adkim" db:"AlignDKIM"` // minOccurs="0"
	AlignSPF          string        `xml:"policy_published>aspf" db:"AlignSPF"`   // minOccurs="0"
	Policy            string        `xml:"policy_published>p" db:"Policy"`
	SubdomainPolicy   string        `xml:"policy_published>sp" db:"SubdomainPolicy"`
	Percentage        int           `xml:"policy_published>pct" db:"Percentage"`
	FailureReport     string        `xml:"policy_published>fo" db:"FailureReport"`
	Records           []DmarcRecord `xml:"record"`
}
type DmarcRecord struct {
	gorm.Model
	SourceIP     string                `xml:"row>source_ip" gorm:"type:varchar(255);not null"`
	Count        int64                 `xml:"row>count" db:"Count"`
	Disposition  string                `xml:"row>policy_evaluated>disposition" db:"Disposition"` // ignore, quarantine, reject
	EvalDKIM     string                `xml:"row>policy_evaluated>dkim" db:"EvalDKIM"`           // pass, fail
	EvalSPF      string                `xml:"row>policy_evaluated>spf" db:"EvalSPF"`             // pass, fail
	POReason     []DmarcPOReason       `xml:"row>policy_evaluated>reason"`
	HeaderFrom   string                `xml:"identifiers>header_from" db:"HeaderFrom"`
	EnvelopeFrom string                `xml:"identifiers>envelope_from" db:"EnvelopeFrom"`
	EnvelopeTo   string                `xml:"identifiers>envelope_to" db:"EnvelopeTo"` // min 0
	AuthDKIM     []DmarcDKIMAuthResult `xml:"auth_results>dkim"`                       // min 0
	AuthSPF      []DmarcSPFAuthResult  `xml:"auth_results>spf"`
	ReportID     uint
}
