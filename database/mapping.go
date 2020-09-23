package database

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
)

const (
	None = 1 << iota
	IncludeInAll
	IncludeTermVectors
	Index
	Store
)

type Mapping struct {
	// TODO: also track fields to detect unmapped ones
	Mapping *mapping.DocumentMapping
}

func NewMapping(m *mapping.DocumentMapping) *Mapping {

	return &Mapping{
		Mapping: m,
	}

}

func (m *Mapping) Other() *Mapping {

	other := bleve.NewDocumentMapping()
	m.Mapping.AddSubDocumentMapping("Other", other)

	return &Mapping{
		Mapping: other,
	}

}

func (m *Mapping) AddKeyword(name string, flags int) {

	field := bleve.NewTextFieldMapping()
	field.Analyzer = keyword.Name
	field.IncludeInAll = (IncludeInAll & flags) == IncludeInAll
	field.IncludeTermVectors = (IncludeTermVectors & flags) == IncludeTermVectors
	field.Index = (Index & flags) == Index
	field.Store = (Store & flags) == Store

	m.Mapping.AddFieldMappingsAt(name, field)

}

func (m *Mapping) Ignore(names ...string) {

	for _, name := range names {
		m.AddKeyword(name, Store)
	}

}

func mapKeyword(d *mapping.DocumentMapping, name string, flags int) {

	Mapping := bleve.NewTextFieldMapping()
	Mapping.Analyzer = keyword.Name
	Mapping.IncludeInAll = (IncludeInAll & flags) == IncludeInAll
	Mapping.IncludeTermVectors = (IncludeTermVectors & flags) == IncludeTermVectors
	Mapping.Index = (Index & flags) == Index
	Mapping.Store = (Store & flags) == Store

	d.AddFieldMappingsAt(name, Mapping)

}
