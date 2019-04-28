package main

import (
	"encoding/json"
	. "github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestParsePaymentsList(t *testing.T) {
	file, err := ioutil.ReadFile("test_resources/payments.json")
	Nil(t, err)

	var payments PaymentsList
	err = json.Unmarshal([]byte(file), &payments)
	Nil(t, err)

	Equal(t, "https://api.test.form3.tech/v1/payments", payments.Links.Self)
	Equal(t, 14, len(payments.Data))
}

func TestParsePaymentJson(t *testing.T) {
	file, err := ioutil.ReadFile("test_resources/single_payment.json")
	Nil(t, err)

	var payment Payment
	err = json.Unmarshal([]byte(file), &payment)
	Nil(t, err)

	Equal(t, "Payment", payment.Type)
	Equal(t, "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43", payment.ID)
	Equal(t, 0, payment.Version)
	Equal(t, "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb", payment.OrganisationId)

	Equal(t, float64(100.21), payment.Attributes.Amount)
	Equal(t, "GBP", payment.Attributes.Currency)

	Equal(t, "W Owens", payment.Attributes.BeneficiaryParty.AccountName)
	Equal(t, "31926819", payment.Attributes.BeneficiaryParty.AccountNumber)
	Equal(t, "BBAN", payment.Attributes.BeneficiaryParty.AccountNumberCode)
	Equal(t, uint(0), payment.Attributes.BeneficiaryParty.AccountType)
	Equal(t, "1 The Beneficiary Localtown SE2", payment.Attributes.BeneficiaryParty.Address)
	Equal(t, "403000", payment.Attributes.BeneficiaryParty.BankId)
	Equal(t, "GBDSC", payment.Attributes.BeneficiaryParty.BankIdCode)
	Equal(t, "Wilfred Jeremiah Owens", payment.Attributes.BeneficiaryParty.Name)

	Equal(t, "SHAR", payment.Attributes.ChargesInformation.BearerCode)
	Equal(t, float64(5.00), payment.Attributes.ChargesInformation.SenderCharges[0].Amount)
	Equal(t, "GBP", payment.Attributes.ChargesInformation.SenderCharges[0].Currency)
	Equal(t, float64(10.00), payment.Attributes.ChargesInformation.SenderCharges[1].Amount)
	Equal(t, "USD", payment.Attributes.ChargesInformation.SenderCharges[1].Currency)
	Equal(t, float64(1.0), payment.Attributes.ChargesInformation.Amount)
	Equal(t, "USD", payment.Attributes.ChargesInformation.Currency)

	Equal(t, "EJ Brown Black", payment.Attributes.DebtorParty.AccountName)
	Equal(t, "GB29XABC10161234567801", payment.Attributes.DebtorParty.AccountNumber)
	Equal(t, "IBAN", payment.Attributes.DebtorParty.AccountNumberCode)
	Equal(t, "10 Debtor Crescent Sourcetown NE1", payment.Attributes.DebtorParty.Address)
	Equal(t, "203301", payment.Attributes.DebtorParty.BankId)
	Equal(t, "GBDSC", payment.Attributes.DebtorParty.BankIdCode)
	Equal(t, "Emelia Jane Brown", payment.Attributes.DebtorParty.Name)

	Equal(t, "Wil piano Jan", payment.Attributes.EndToEndReference)

	Equal(t, "FX123", payment.Attributes.FX.ContractReference)
	Equal(t, float64(2.00000), payment.Attributes.FX.ExchangeRate)
	Equal(t, float64(200.42), payment.Attributes.FX.OriginalAmount)
	Equal(t, "USD", payment.Attributes.FX.OriginalCurrency)

	Equal(t, 1002001, payment.Attributes.NumericReference)
	Equal(t, "123456789012345678", payment.Attributes.PaymentID)
	Equal(t, "Paying for goods/services", payment.Attributes.PaymentPurpose)
	Equal(t, "FPS", payment.Attributes.PaymentScheme)
	Equal(t, "Credit", payment.Attributes.PaymentType)

	Equal(t, "2017-01-18", payment.Attributes.ProcessingDate)
	Equal(t, "Payment for Em's piano lessons", payment.Attributes.Reference)

	Equal(t, "InternetBanking", payment.Attributes.SchemePaymentSubType)
	Equal(t, "ImmediatePayment", payment.Attributes.SchemePaymentType)

	Equal(t, "56781234", payment.Attributes.SponsorParty.AccountNumber)
	Equal(t, "123123", payment.Attributes.SponsorParty.BankId)
	Equal(t, "GBDSC", payment.Attributes.SponsorParty.BankIdCode)
}
