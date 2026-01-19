package models

// BillingCycle define os ciclos de faturamento disponíveis.
type BillingCycle string

const (
	Free         BillingCycle = "free"
	OneTime      BillingCycle = "onetime"
	Monthly      BillingCycle = "monthly"
	Quarterly    BillingCycle = "quarterly"
	Semiannually BillingCycle = "semiannually"
	Annually     BillingCycle = "annually"
	Biennially   BillingCycle = "biennially"
	Triennially  BillingCycle = "triennially"
)
