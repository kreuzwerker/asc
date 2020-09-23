package transfer

type SKU struct {
	SKU              string                 // Internal ID
	OfferTermCode    string                 // Internal ID
	RateCode         string                 // Internal ID
	TermType         string                 // OnDemand or Reserved
	PriceDescription string                 // Price description as text e.g. $0.025 per GB - first 50 TB / month of storage used
	EffectiveDate    string                 // Price in effect since, e.g. 2020-04-01
	StartingRange    string                 // Tier start in units, e.g. 1024000
	EndingRange      string                 // Tier ends in units, e.g. 512000 or "Inf"
	Unit             string                 // Unit description as text - many expressions such as GB-month, GB-months, inbound-minutes, inbound_minutes
	PricePerUnit     float64                // Unit price
	Currency         string                 // CNY or USD
	ProductFamily    string                 // Always empty
	Other            map[string]interface{} `csv:"-"`
}
