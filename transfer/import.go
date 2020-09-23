package transfer

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/pkg/errors"
)

var headerPattern = regexp.MustCompile(`,"(.+)"$`)

type Header struct {
	Disclaimer      string
	FormatVersion   string
	OfferCode       string
	PublicationDate time.Time
	Version         string
}

func Import(r1 io.Reader, cb func(*Header, *SKU) error) error {

	var (
		header = new(Header)
		i      = 0
	)

	r2, err := gzip.NewReader(r1)

	if err != nil {
		return errors.Wrapf(err, "failed to parse zip")
	}

	r3 := bufio.NewReader(r2)

	for ; i < 5; i++ {

		meta, _, _ := r3.ReadLine()
		line := string(meta)

		result := headerPattern.FindStringSubmatch(line)

		if len(result) != 2 {
			return fmt.Errorf("unexpected header pattern %q", line)
		}

		line = result[1]

		if i == 0 {
			header.FormatVersion = line
		} else if i == 1 {
			header.Disclaimer = line
		} else if i == 2 {

			date, err := time.Parse(time.RFC3339, line)

			if err != nil {
				return errors.Wrapf(err, "failed to parse date for header from line %q", line)
			}

			header.PublicationDate = date

		} else if i == 3 {
			header.Version = line
		} else if i == 4 {
			header.OfferCode = line
		}

	}

	r4 := csv.NewReader(r3)

	dec, err := csvutil.NewDecoder(r4)

	if err != nil {
		return errors.Wrapf(err, "failed to parse CSV")
	}

	columns := dec.Header()

	for {

		i = i + 1

		sku := SKU{
			Other: make(map[string]interface{}),
		}

		if err := dec.Decode(&sku); err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrapf(err, "failed to decode line %d", i)
		}

		for _, i := range dec.Unused() {
			sku.Other[columns[i]] = dec.Record()[i]
		}

		if err := cb(header, &sku); err != nil {
			return errors.Wrapf(err, "failed to handle callback in line %d", i)
		}

	}

	return nil

}
