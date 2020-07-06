package dmarc

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/menta2l/dmarc-parser/internal/archive"
	"github.com/menta2l/dmarc-parser/internal/types"
	"github.com/menta2l/dmarc-parser/internal/utils"
	"golang.org/x/net/html/charset"
)

func Parse(input string, db *gorm.DB) error {
	var reader io.Reader
	if input == "stdin" {
		reader = os.Stdin
	} else {
		if !utils.FileExists(input) {
			return fmt.Errorf("File %s not exist", input)
		}
		file, err := os.Open(input)
		if err != nil {
			return err
		}
		reader = file
		defer file.Close()
	}
	msg, err := utils.ReadMail(reader)
	if err != nil {
		return err
	}
	r, err := DmarcReportPrepareAttachment(msg)
	if err != nil {
		return err
	}
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charset.NewReaderLabel
	//fb := &types.AggregateReport{}
	fb := &types.DmarcReport{}

	if err := decoder.Decode(fb); err != nil {
		return err
	}
	timestamp1, _ := strconv.Atoi(strings.TrimSpace(fb.RawDateRangeBegin))
	fb.DateRangeBegin = int64(timestamp1)
	timestamp2, _ := strconv.Atoi(strings.TrimSpace(fb.RawDateRangeEnd))
	fb.DateRangeEnd = int64(timestamp2)

	fb.MessageId = msg.Header.Get("Message-Id")
	if err = db.Create(&fb).Error; err != nil {
		return err
	}
	//	chanARR, wg := ParseDmarcARRParallel(50, 4, *fb)
	//	for k, _ := range fb.Records {
	//		fb.Records[k].RecordNumber = int64(k)
	//		chanARR <- &fb.Records[k]
	//	}
	//	close(chanARR)
	//	wg.Wait()
	return nil
}
func DmarcReportPrepareAttachment(m *mail.Message) (io.Reader, error) {

	header := m.Header

	mediaType, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("PrepareAttachment: error parsing media type")
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(m.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				return nil, fmt.Errorf("PrepareAttachment: EOF before valid attachment")
			}
			if err != nil {
				return nil, err
			}

			// need to add checks to ensure base64
			partType, _, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			if err != nil {
				return nil, fmt.Errorf("PrepareAttachment: error parsing media type of part")
			}

			// if gzip
			if strings.HasPrefix(partType, "application/gzip") ||
				strings.HasPrefix(partType, "application/x-gzip") ||
				strings.HasPrefix(partType, "application/gzip-compressed") ||
				strings.HasPrefix(partType, "application/gzipped") ||
				strings.HasPrefix(partType, "application/x-gunzip") ||
				strings.HasPrefix(partType, "application/x-gzip-compressed") ||
				strings.HasPrefix(partType, "gzip/document") {

				decodedBase64 := base64.NewDecoder(base64.StdEncoding, p)
				decompressed, err := gzip.NewReader(decodedBase64)
				if err != nil {
					return nil, err
				}

				return decompressed, nil
			}

			// if zip
			if strings.HasPrefix(partType, "application/zip") || // google style
				strings.HasPrefix(partType, "application/x-zip-compressed") { // yahoo style

				decodedBase64 := base64.NewDecoder(base64.StdEncoding, p)
				decompressed, err := archive.ExtractZipFile(decodedBase64)
				if err != nil {
					return nil, err
				}

				return decompressed, nil
			}

			// if xml
			if strings.HasPrefix(partType, "text/xml") {
				return p, nil
			}

			// if application/octetstream, check filename
			if strings.HasPrefix(partType, "application/octet-stream") {

				if strings.HasSuffix(p.FileName(), ".zip") {
					decodedBase64 := base64.NewDecoder(base64.StdEncoding, p)
					decompressed, err := archive.ExtractZipFile(decodedBase64)
					if err != nil {
						return nil, err
					}

					return decompressed, nil
				}
				if strings.HasSuffix(p.FileName(), ".gz") {
					decodedBase64 := base64.NewDecoder(base64.StdEncoding, p)
					decompressed, _ := gzip.NewReader(decodedBase64)

					return decompressed, nil
				}
			}
		}

	}

	// if gzip
	if strings.HasPrefix(mediaType, "application/gzip") || // proper :)
		strings.HasPrefix(mediaType, "application/x-gzip") || // gmail attachment
		strings.HasPrefix(mediaType, "application/gzip-compressed") ||
		strings.HasPrefix(mediaType, "application/gzipped") ||
		strings.HasPrefix(mediaType, "application/x-gunzip") ||
		strings.HasPrefix(mediaType, "application/x-gzip-compressed") ||
		strings.HasPrefix(mediaType, "gzip/document") {

		decodedBase64 := base64.NewDecoder(base64.StdEncoding, m.Body)
		decompressed, _ := gzip.NewReader(decodedBase64)

		return decompressed, nil

	}

	// if zip
	if strings.HasPrefix(mediaType, "application/zip") || // google style
		strings.HasPrefix(mediaType, "application/x-zip-compressed") { // yahoo style
		decodedBase64 := base64.NewDecoder(base64.StdEncoding, m.Body)
		decompressed, err := archive.ExtractZipFile(decodedBase64)
		if err != nil {
			return nil, err
		}

		return decompressed, nil
	}

	// if xml
	if strings.HasPrefix(mediaType, "text/xml") {
		return m.Body, nil
	}

	return nil, fmt.Errorf("PrepareAttachment: reached the end, no attachment found.")
}
