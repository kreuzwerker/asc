package transfer

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Index struct {
	FormatVersion   string           `json:"formatVersion"`
	Disclaimer      string           `json:"disclaimer"`
	PublicationDate time.Time        `json:"publicationDate"`
	Offers          map[string]Offer `json:"offers"`
}

type Offer struct {
	CurrentRegionIndexUrl      string `json:"currentRegionIndexUrl"`
	CurrentSavingsPlanIndexUrl string `json:"currentSavingsPlanIndexUrl"`
	CurrentVersionUrl          string `json:"currentVersionUrl"`
	OfferCode                  string `json:"offerCode"`
	SavingsPlanVersionIndexUrl string `json:"savingsPlanVersionIndexUrl"`
	VersionIndexUrl            string `json:"versionIndexUrl"`
}

func Export(dir string) error {

	var (
		i    Index
		size int
	)

	res, err := http.Get("https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/index.json")

	if err != nil {
		log.Fatal(err)
	}

	out, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(out, &i); err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	size = len(i.Offers)

	for _, o := range i.Offers {

		url := strings.Replace(o.CurrentVersionUrl, ".json", ".csv", 1)
		url = fmt.Sprintf("https://pricing.us-east-1.amazonaws.com%s", url)

		log.Printf("START (%3d left) %q downloading current version: %s", size, o.OfferCode, url)

		res, err := http.Get(url)

		if err != nil {
			return errors.Wrapf(err, "failed to retrieve url %q", url)
		}

		file := filepath.Join(dir, fmt.Sprintf("%s.csv.gz", o.OfferCode))

		w1, err := os.Create(file)

		if err != nil {
			return errors.Wrapf(err, "failed to create file %q", file)
		}

		w2, err := gzip.NewWriterLevel(w1, gzip.BestCompression)

		if err != nil {
			return errors.Wrapf(err, "failed to create zip writer")
		}

		if _, err := io.Copy(w2, res.Body); err != nil {
			return errors.Wrapf(err, "failed to copy CSV contents")
		}

		for _, e := range []io.Closer{res.Body, w2, w1} {

			if err := e.Close(); err != nil {
				return errors.Wrapf(err, "failed to close writer")
			}

		}

		size = size - 1

		log.Printf(" DONE (%3d left) %q downloading current version: %s", size, o.OfferCode, url)

	}

	return nil

}
