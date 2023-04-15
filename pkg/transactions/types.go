package transactions

type Transaction struct {
	// ID is the id of the transaction.
	ID int64
	// Account is the account number of the account under view.
	Account string
	// BookingDate is the date of the transaction being triggered.
	BookingDate string
	// ValutaDate is the date of the transaction being completed.
	ValutaDate string
	// BookingText is a text that tries to set a type for the transaction.
	BookingText string
	// Purpose is a text that describes the purpose of the transaction.
	Purpose string
	// CreditorID is the creditor identifier of the creditor.
	CreditorID string
	// MandateRef is the id that identifies the mandate that allows the creditor to collect the amount.
	MandateRef string
	// CustomerRef ???
	CustomerRef string
	// CollectorRef is some id that seems to be specific to the creditor.
	CollectorRef string
	// OrigAmount ???
	OrigAmount float64
	// ChargebackFee is the fee charged by the bank for a chargeback.
	ChargebackFee float64
	// Beneficiary is the name of the creditor.
	Beneficiary string
	// AccountNumber is the IBAN of the beneficiary.
	AccountNumber string
	// BIC is the Bank Identifier Code of the beneficiary.
	BIC string
	// Amount is the amount of the transaction.
	Amount float64
	// Currency is the currency of the transaction.
	Currency string
	// AdditionalDetails describes the current state of the transaction.
	AdditionalDetails string
}
