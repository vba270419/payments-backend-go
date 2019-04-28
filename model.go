package main

type PaymentsList struct {
	Data  []Payment `json:"data,omitempty"`
	Links Links     `json:"links,omitempty"`
}

type Links struct {
	Self string `json:"self,omitempty"`
}

type Payment struct {
	Type           string     `json:"type,omitempty" bson:"type,omitempty"`
	ID             string     `json:"id,omitempty" bson:"_id"`
	Version        int        `json:"version,omitempty" bson:"version"`
	OrganisationId string     `json:"organisation_id,omitempty" bson:"organisation_id,omitempty"`
	Attributes     Attributes `json:"attributes,omitempty" bson:"attributes,omitempty"`
}

type Attributes struct {
	Amount               float64            `json:"amount,string,omitempty" bson:"amount,omitempty"`
	BeneficiaryParty     BeneficiaryParty   `json:"beneficiary_party,omitempty" bson:"beneficiary_party,omitempty"`
	ChargesInformation   ChargesInformation `json:"charges_information,omitempty" bson:"charges_information,omitempty"`
	Currency             string             `json:"currency,omitempty" bson:"currency,omitempty"`
	DebtorParty          DebtorParty        `json:"debtor_party,omitempty" bson:"debtor_party,omitempty"`
	EndToEndReference    string             `json:"end_to_end_reference,omitempty" bson:"end_to_end_reference,omitempty"`
	FX                   FX                 `json:"fx,omitempty" bson:"fx,omitempty"`
	NumericReference     int                `json:"numeric_reference,string,omitempty" bson:"numeric_reference,omitempty"`
	PaymentID            string             `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	PaymentPurpose       string             `json:"payment_purpose,omitempty" bson:"payment_purpose,omitempty"`
	PaymentScheme        string             `json:"payment_scheme,omitempty" bson:"payment_scheme,omitempty"`
	PaymentType          string             `json:"payment_type,omitempty" bson:"payment_type,omitempty"`
	ProcessingDate       string             `json:"processing_date,omitempty" bson:"processing_date,omitempty"`
	Reference            string             `json:"reference,omitempty" bson:"reference,omitempty"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type,omitempty" bson:"scheme_payment_sub_type,omitempty"`
	SchemePaymentType    string             `json:"scheme_payment_type,omitempty" bson:"scheme_payment_type,omitempty"`
	SponsorParty         SponsorParty       `json:"sponsor_party,omitempty" bson:"sponsor_party,omitempty"`
}

type SponsorParty struct {
	AccountNumber string `json:"account_number,omitempty" bson:"account_number,omitempty"`
	BankId        string `json:"bank_id,omitempty" bson:"bank_id,omitempty"`
	BankIdCode    string `json:"bank_id_code,omitempty" bson:"bank_id_code,omitempty"`
}

type DebtorParty struct {
	*SponsorParty
	AccountName       string `json:"account_name,omitempty" bson:"account_name,omitempty"`
	AccountNumberCode string `json:"account_number_code,omitempty" bson:"account_number_code,omitempty"`
	Address           string `json:"address,omitempty" bson:"address,omitempty"`
	Name              string `json:"name,omitempty" bson:"name,omitempty"`
}

type BeneficiaryParty struct {
	*DebtorParty
	AccountType uint `json:"account_type,omitempty" bson:"account_type,omitempty"`
}

type ChargesInformation struct {
	BearerCode    string          `json:"bearer_code,omitempty" bson:"bearer_code,omitempty"`
	SenderCharges []SenderCharges `json:"sender_charges,omitempty" bson:"sender_charges,omitempty"`
	Amount        float64         `json:"receiver_charges_amount,string,omitempty" bson:"receiver_charges_amount,omitempty"`
	Currency      string          `json:"receiver_charges_currency,omitempty" bson:"receiver_charges_currency,omitempty"`
}

type SenderCharges struct {
	Amount   float64 `json:"amount,string,omitempty" bson:"amount,omitempty"`
	Currency string  `json:"currency,omitempty" bson:"currency,omitempty"`
}

type FX struct {
	ContractReference string  `json:"contract_reference,omitempty" bson:"contract_reference,omitempty"`
	ExchangeRate      float64 `json:"exchange_rate,string,omitempty" bson:"exchange_rate,omitempty"`
	OriginalAmount    float64 `json:"original_amount,string,omitempty" bson:"original_amount,omitempty"`
	OriginalCurrency  string  `json:"original_currency,omitempty" bson:"original_currency,omitempty"`
}
