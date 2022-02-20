package data

import (
	"github.com/jalexanderII/zero_fintech/gen/Go/common"
)

type FakePaymentTask struct {
	ID        string  `faker:"_"`
	UserId    string  `faker:"oneof: 6212a0101fca9390a37a32d2, 6212a0101fca9390a37a32d3, 6212a0101fca9390a37a32d4"`
	AccountId string  `faker:"oneof: 6212a29794c88ffb3de9d76b, 6212a29794c88ffb3de9d76a, 6212a29794c88ffb3de9d769, 6212a29794c88ffb3de9d768, 6212a29794c88ffb3de9d767, 6212a29794c88ffb3de9d766, 6212a29794c88ffb3de9d765, 6212a29794c88ffb3de9d764, 6212a29794c88ffb3de9d763, 6212a29794c88ffb3de9d762"`
	Amount    float64 `faker:"oneof: 340.0, 530.0, 250.0, 684.0"`
}

type FakeMetaData struct {
	PreferredPlanType         common.PlanType         `faker:"preferred_plan_type"`
	PreferredTimelineInMonths float64                 `faker:"preferred_timeline_in_months" `
	PreferredPaymentFreq      common.PaymentFrequency `faker:"preferred_payment_freq" `
}

type FakeAnnualPercentageRates struct {
	AprPercentage        float64 `faker:"oneof: 20.0, 22.0, 40.0"`
	AprType              string  `faker:"oneof: cash_apr, purchase_apr"`
	BalanceSubjectToApr  float64 `faker:"oneof: 1000.0, 1500.0"`
	InterestChargeAmount float64 `faker:"-"`
}

type FakeAccount struct {
	ID                     string  `faker:"-"`
	PlaidAccountId         string  `faker:"oneof: 1, 2, 3"`
	UserId                 string  `faker:"oneof: 6212a0101fca9390a37a32d2, 6212a0101fca9390a37a32d3, 6212a0101fca9390a37a32d4"`
	Name                   string  `faker:"oneof: Chase, Amex, Citi, CreditOne, CapitalOne, Barclay, BestBuy, Crate-and-Barrel"`
	OfficialName           string  `faker:"oneof: Chase, Amex, Citi, CreditOne, CapitalOne, Barclay, BestBuy, Crate-and-Barrel"`
	AvailableBalance       float64 `faker:"oneof: 0.0, 50.0, 500.0, 1000.0, 7000.0"`
	CurrentBalance         float64 `faker:"oneof: 0.0, 50.0, 500.0, 1000.0"`
	CreditLimit            float64 `faker:"oneof: 8000.0, 9000.0, 10000.0, 15000.0"`
	IsoCurrencyCode        string  `faker:"currency"`
	LastPaymentAmount      float64 `faker:"oneof: 0.0, 50.0, 500.0"`
	LastStatementIssueDate string  `faker:"date"`
	LastStatementBalance   float64 `faker:"oneof: 0.0, 50.0, 500.0, 1000.0"`
	MinimumPaymentAmount   float64 `faker:"oneof: 100.0, 150.0, 250.0"`
	NextPaymentDueDate     string  `faker:"date"`
	Type                   string  `faker:"word"`
	Subtype                string  `faker:"word"`
	AnnualPercentageRate   []FakeAnnualPercentageRates
	IsOverdue              bool
}

type FakeTransaction struct {
	ID                   string   `faker:"-"`
	PlaidTransactionId   string   `faker:"oneof: 1, 2, 3"`
	AccountId            string   `faker:"oneof: 6212a29794c88ffb3de9d76b, 6212a29794c88ffb3de9d76a, 6212a29794c88ffb3de9d769, 6212a29794c88ffb3de9d768, 6212a29794c88ffb3de9d767, 6212a29794c88ffb3de9d766, 6212a29794c88ffb3de9d765, 6212a29794c88ffb3de9d764, 6212a29794c88ffb3de9d763, 6212a29794c88ffb3de9d762"`
	PlaidAccountId       string   `faker:"oneof: 1, 2, 3"`
	UserId               string   `faker:"oneof: 6212a0101fca9390a37a32d2, 6212a0101fca9390a37a32d3, 6212a0101fca9390a37a32d4"`
	OriginalDescription  string   `faker:"sentence"`
	Amount               float64  `faker:"oneof: 340.0, 530.0, 250.0, 684.0"`
	IsoCurrencyCode      string   `faker:"currency"`
	Date                 string   `faker:"date"`
	MerchantName         string   `faker:"name"`
	PaymentChannel       string   `faker:"oneof: online, in-store"`
	AuthorizedDate       string   `faker:"date"`
	Category             []string `faker:"slice_len=2"`
	TransactionType      string   `faker:"word"`
	PendingTransactionId string   `faker:"uuid_digit"`
	CategoryId           string   `faker:"uuid_digit"`
	Name                 string   `faker:"word"`
	PrimaryCategory      string   `faker:"word"`
	DetailedCategory     string   `faker:"word"`
	TransactionDetails   FakeTransactionDetails
	Pending              bool
}

type FakeTransactionDetails struct {
	Address         string `faker:"word"`
	City            string `faker:"word"`
	State           string `faker:"word"`
	Zipcode         string `faker:"oneof: 10012, 10002, 30083"`
	Country         string `faker:"word"`
	StoreNumber     string `faker:"uuid_digit"`
	ReferenceNumber string `faker:"uuid_digit"`
}

type FakeUser struct {
	ID               string            `faker:"-"`
	Username         string            `faker:"name"`
	Email            string            `faker:"email"`
	Password         string            `faker:"password"`
	AccountIdToToken map[string]string `faker:"account_id_to_token"`
}
