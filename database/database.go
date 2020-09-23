package database

import (
	"fmt"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/scorch"
	"github.com/kreuzwerker/asc/transfer"
	"github.com/pkg/errors"
)

type Database struct {
	current *bleve.Batch
	file    string
	index   bleve.Index
}

type Document struct {
	OfferCode string
	transfer.SKU
}

func (d *Document) BleveType() string {
	return d.OfferCode
}

func New(file string) (*Database, error) {

	var (
		index = bleve.NewIndexMapping()
		root  = NewMapping(bleve.NewDocumentStaticMapping())
		other = NewMapping(bleve.NewDocumentStaticMapping())
		// ignored       = bleve.NewDocumentDisabledMapping()
	)

	index.DefaultMapping = root.Mapping

	root.AddKeyword("SKU", Store)
	root.AddKeyword("OfferTermCode", Store)
	root.AddKeyword("RateCode", Store)
	root.AddKeyword("TermType", IncludeTermVectors|Index|Store)

	root.Mapping.AddSubDocumentMapping("Other", other.Mapping)

	other.AddKeyword("Location", IncludeTermVectors|Index|Store)

	// mapKeyword(root, "SKU", Store)
	// mapKeyword(root, "OfferTermCode", Store)
	// mapKeyword(root, "RateCode", Store)
	// mapKeyword(root, "TermType", IncludeTermVectors|Index|Store)
	//
	// root.AddSubDocumentMapping("Other", other)
	//
	// mapKeyword(other, "Location", IncludeTermVectors|Index|Store)

	// {
	//
	// 	var (
	// 		s3      = bleve.NewDocumentMapping()
	// 		other   = bleve.NewDocumentMapping()
	// 		mapping = NewMapping(other)
	// 	)
	//
	// 	mapping.AddKeyword("Location", IncludeTermVectors|Index|Store)
	//
	// 	s3.AddSubDocumentMapping("Other", other)
	// 	index.AddDocumentMapping("S3", s3)
	//
	// }

	// other := bleve.NewDocumentMapping()
	// mapKeyword(other, "Location", IncludeTermVectors|Index|Store)
	//
	// root.AddSubDocumentMapping("Other", other)

	// TODO: build a document hierarchy per service with individual mappings for others (full ignore for all fields) with panic / errors logs for unmapped fields - organize it so that the root / default document contains all the base fields and the Other fields are document type specific

	// TODO: add document mappings per service, dont index ids etc

	// TODO: add proper mappings and correct analyzers (e.g. keyword or english)
	// TODO: ignore fields we don't really have a use for, enumerate them as a helper function

	// TODO: cmdline helper for Other headers and example values / CSV

	main, err := bleve.New(file, index)

	if err != nil {
		return nil, err
	}

	return &Database{
		index: main,
	}, nil

}

func Open(file string) (*Database, error) {

	db, err := bleve.Open(file)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database %q", file)
	}

	return &Database{
		index: db,
	}, nil

}

func (d *Database) Terms(field string) ([]string, error) {

	var terms []string

	f, err := d.index.FieldDict(field)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	for {

		d, err := f.Next()

		if d == nil || err != nil {
			break
		}

		terms = append(terms, d.Term)

	}

	return terms, nil

}

func (d *Database) Add(header *transfer.Header, sku *transfer.SKU) error {

	if d.current == nil {
		d.current = d.index.NewBatch()
	}

	// sku.Other = nil

	offerCode := header.OfferCode
	offerCode = strings.TrimPrefix(offerCode, "Amazon")
	offerCode = strings.TrimPrefix(offerCode, "AWS")
	offerCode = strings.TrimPrefix(offerCode, "aws")

	document := Document{
		OfferCode: offerCode,
		SKU:       *sku,
	}

	if err := d.current.Index(sku.SKU, document); err != nil {
		return err
	}

	if d.current.Size() > 50000 {
		return d.Commit()
	}

	return nil

}

func (d *Database) Commit() error {

	if d.current == nil {
		return nil
	}

	defer d.current.Reset()

	return d.index.Batch(d.current)

}

func (d *Database) Close() error {

	if err := d.Commit(); err != nil {
		return err
	}

	return d.index.Close()

}

func (d *Database) Search(region string, service string, query string) error {

	// TODO: proper query merging

	query = fmt.Sprintf(`+Other.Location:"%s" +OfferCode:%s %s`, region, service, query)

	qs := bleve.NewQueryStringQuery(query)

	if err := qs.Validate(); err != nil {
		return errors.Wrapf(err, "failed to parse query %q", query)
	}

	req := bleve.NewSearchRequest(qs)

	req.Fields = []string{
		"PriceDescription",
		"PricePerUnit", // TODO: fixme / missing (during import?)
		"TermType",
		"Unit",
		"Other.Location",
	}

	req.Fields = []string{
		"*",
	}

	// fCurrency := bleve.NewFacetRequest("Currency", 2)
	// req.AddFacet(fCurrency.Field, fCurrency)
	//
	// fTermType := bleve.NewFacetRequest("TermType", 2)
	// req.AddFacet(fTermType.Field, fTermType)
	//
	// fUnit := bleve.NewFacetRequest("Unit", 10)
	// req.AddFacet(fUnit.Field, fUnit)

	res, err := d.index.Search(req)

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil

}

func init() {
	bleve.Config.DefaultIndexType = scorch.Name
}
