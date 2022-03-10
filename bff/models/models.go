package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token for use of plaid public token retrieval
type Token struct {
	ID          primitive.ObjectID `bson:"_id"`
	Username    string             `bson:"username"`
	Value       string             `bson:"value"`
	ItemId      string             `bson:"item_id"`
	Institution string             `bson:"institution"`
}

type LiabilitiesResponse struct {
	Liabilities []CreditCardLiability `json:"liabilities"`
}

type TransactionsResponse struct {
	Accounts     []AccountBase `json:"accounts,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

// CreditCardLiability An object representing a credit card account.
type CreditCardLiability struct {
	// The ID of the account that this liability belongs to.
	AccountId string `json:"account_id"`
	// The various interest rates that apply to the account.
	Aprs []struct {
		AprPercentage        float64 `json:"apr_percentage,omitempty"`
		AprType              string  `json:"apr_type,omitempty"`
		BalanceSubjectToApr  float64 `json:"balance_subject_to_apr,omitempty"`
		InterestChargeAmount float64 `json:"interest_charge_amount,omitempty"`
	} `json:"aprs"`
	// true if a payment is currently overdue. Availability for this field is limited.
	IsOverdue bool `json:"is_overdue,omitempty"`
	// The amount of the last payment.
	LastPaymentAmount float32 `json:"last_payment_amount,omitempty"`
	// The date of the last payment. Dates are returned in an [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format (YYYY-MM-DD). Availability for this field is limited.
	LastPaymentDate string `json:"last_payment_date,omitempty"`
	// The date of the last statement. Dates are returned in an [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format (YYYY-MM-DD).
	LastStatementIssueDate string `json:"last_statement_issue_date,omitempty"`
	// The total amount owed as of the last statement issued
	LastStatementBalance float32 `json:"last_statement_balance,omitempty"`
	// The minimum payment due for the next billing cycle.
	MinimumPaymentAmount float32 `json:"minimum_payment_amount,omitempty"`
	// The due date for the next payment. The due date is `null` if a payment is not expected. Dates are returned in an [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format (YYYY-MM-DD).
	NextPaymentDueDate string `json:"next_payment_due_date,omitempty"`
}

// AccountBase A single account at a financial institution.
type AccountBase struct {
	// Plaid’s unique identifier for the account. This value will not change unless Plaid can't reconcile the account with the data returned by the financial institution. This may occur, for example, when the name of the account changes. If this happens a new `account_id` will be assigned to the account.  The `account_id` can also change if the `access_token` is deleted and the same credentials that were used to generate that `access_token` are used to generate a new `access_token` on a later date. In that case, the new `account_id` will be different from the old `account_id`.  If an account with a specific `account_id` disappears instead of changing, the account is likely closed. Closed accounts are not returned by the Plaid API.  Like all Plaid identifiers, the `account_id` is case sensitive.
	AccountId string         `json:"account_id"`
	Balances  AccountBalance `json:"balances"`
	// The last 2-4 alphanumeric characters of an account's official account number. Note that the mask may be non-unique between an Item's accounts, and it may also not match the mask that the bank displays to the user.
	Mask string `json:"mask"`
	// The name of the account, either assigned by the user or by the financial institution itself
	Name string `json:"name"`
	// The official name of the account as given by the financial institution
	OfficialName string `json:"official_name"`
	Type         string `json:"type"`
	Subtype      string `json:"subtype"`
	// The current verification status of an Auth Item initiated through Automated or Manual micro-deposits.  Returned for Auth Items only.  `pending_automatic_verification`: The Item is pending automatic verification  `pending_manual_verification`: The Item is pending manual micro-deposit verification. Items remain in this state until the user successfully verifies the two amounts.  `automatically_verified`: The Item has successfully been automatically verified   `manually_verified`: The Item has successfully been manually verified  `verification_expired`: Plaid was unable to automatically verify the deposit within 7 calendar days and will no longer attempt to validate the Item. Users may retry by submitting their information again through Link.  `verification_failed`: The Item failed manual micro-deposit verification because the user exhausted all 3 verification attempts. Users may retry by submitting their information again through Link.
	VerificationStatus *string `json:"verification_status,omitempty"`
}

// Transaction A representation of a transaction
type Transaction struct {
	// Please use the `payment_channel` field, `transaction_type` will be deprecated in the future.  `digital:` transactions that took place online.  `place:` transactions that were made at a physical location.  `special:` transactions that relate to banks, e.g. fees or deposits.  `unresolved:` transactions that do not fit into the other three types.
	TransactionType *string `json:"transaction_type,omitempty"`
	// The ID of a posted transaction's associated pending transaction, where applicable.
	PendingTransactionId string `json:"pending_transaction_id"`
	// The ID of the category to which this transaction belongs. See [Categories](https://plaid.com/docs/#category-overview).  If the `transactions` object was returned by an Assets endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
	CategoryId string `json:"category_id"`
	// A hierarchical array of the categories to which this transaction belongs. See [Categories](https://plaid.com/docs/#category-overview).  If the `transactions` object was returned by an Assets endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
	Category    []string    `json:"category"`
	Location    Location    `json:"location"`
	PaymentMeta PaymentMeta `json:"payment_meta"`
	// The name of the account owner. This field is not typically populated and only relevant when dealing with sub-accounts.
	AccountOwner string `json:"account_owner"`
	// The merchant name or transaction description.  If the `transactions` object was returned by a Transactions endpoint such as `/transactions/get`, this field will always appear. If the `transactions` object was returned by an Assets endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
	Name string `json:"name"`
	// The string returned by the financial institution to describe the transaction. For transactions returned by `/transactions/get`, this field is in beta and will be omitted unless the client is both enrolled in the closed beta program and has set `options.include_original_description` to `true`.
	OriginalDescription string `json:"original_description,omitempty"`
	// The ID of the account in which this transaction occurred.
	AccountId string `json:"account_id"`
	// The settled value of the transaction, denominated in the account's currency, as stated in `iso_currency_code` or `unofficial_currency_code`. Positive values when money moves out of the account; negative values when money moves in. For example, debit card purchases are positive; credit card payments, direct deposits, and refunds are negative.
	Amount float32 `json:"amount"`
	// The ISO-4217 currency code of the transaction. Always `null` if `unofficial_currency_code` is non-null.
	IsoCurrencyCode string `json:"iso_currency_code"`
	// The unofficial currency code associated with the transaction. Always `null` if `iso_currency_code` is non-`null`. Unofficial currency codes are used for currencies that do not have official ISO currency codes, such as cryptocurrencies and the currencies of certain countries.  See the [currency code schema](https://plaid.com/docs/api/accounts#currency-code-schema) for a full listing of supported `iso_currency_code`s.
	UnofficialCurrencyCode string `json:"unofficial_currency_code"`
	// For pending transactions, the date that the transaction occurred; for posted transactions, the date that the transaction posted. Both dates are returned in an [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format ( `YYYY-MM-DD` ).
	Date string `json:"date"`
	// When `true`, identifies the transaction as pending or unsettled. Pending transaction details (name, type, amount, category ID) may change before they are settled.
	Pending bool `json:"pending"`
	// The unique ID of the transaction. Like all Plaid identifiers, the `transaction_id` is case sensitive.
	TransactionId string `json:"transaction_id"`
	// The merchant name, as extracted by Plaid from the `name` field.
	MerchantName string `json:"merchant_name,omitempty"`
	// The check number of the transaction. This field is only populated for check transactions.
	CheckNumber string `json:"check_number,omitempty"`
	// The channel used to make a payment. `online:` transactions that took place online.  `in store:` transactions that were made at a physical location.  `other:` transactions that relate to banks, e.g. fees or deposits.  This field replaces the `transaction_type` field.
	PaymentChannel string `json:"payment_channel"`
	// The date that the transaction was authorized. Dates are returned in an [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format ( `YYYY-MM-DD` ).
	AuthorizedDate string `json:"authorized_date"`
	// Date and time when a transaction was authorized in [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format ( `YYYY-MM-DDTHH:mm:ssZ` ).  This field is only populated for UK institutions. For institutions in other countries, will be `null`.
	AuthorizedDatetime *time.Time `json:"authorized_datetime"`
	// Date and time when a transaction was posted in [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format ( `YYYY-MM-DDTHH:mm:ssZ` ).  This field is only populated for UK institutions. For institutions in other countries, will be `null`.
	Datetime                *time.Time `json:"datetime"`
	TransactionCode         string     `json:"transaction_code"`
	PersonalFinanceCategory struct {
		// A high level category that communicates the broad category of the transaction.
		Primary string `json:"primary"`
		// Provides additional granularity to the primary categorization.
		Detailed string `json:"detailed"`
	} `json:"personal_finance_category,omitempty"`
}

// Location A representation of where a transaction took place
type Location struct {
	// The street address where the transaction occurred.
	Address string `json:"address"`
	// The city where the transaction occurred.
	City string `json:"city"`
	// The region or state where the transaction occurred. In API versions 2018-05-22 and earlier, this field is called `state`.
	Region string `json:"region"`
	// The postal code where the transaction occurred. In API versions 2018-05-22 and earlier, this field is called `zip`.
	PostalCode string `json:"postal_code"`
	// The ISO 3166-1 alpha-2 country code where the transaction occurred.
	Country string `json:"country"`
	// The latitude where the transaction occurred.
	Lat float32 `json:"lat"`
	// The longitude where the transaction occurred.
	Lon float32 `json:"lon"`
	// The merchant defined store number where the transaction occurred.
	StoreNumber string `json:"store_number"`
}

// PaymentMeta Transaction information specific to inter-bank transfers. If the transaction was not an inter-bank transfer, all fields will be `null`.  If the `transactions` object was returned by a Transactions endpoint such as `/transactions/get`, the `payment_meta` key will always appear, but no data elements are guaranteed. If the `transactions` object was returned by an Assets endpoint such as `/asset_report/get/` or `/asset_report/pdf/get`, this field will only appear in an Asset Report with Insights.
type PaymentMeta struct {
	// The transaction reference number supplied by the financial institution.
	ReferenceNumber string `json:"reference_number"`
	// The ACH PPD ID for the payer.
	PpdId string `json:"ppd_id"`
	// For transfers, the party that is receiving the transaction.
	Payee string `json:"payee"`
	// The party initiating a wire transfer. Will be `null` if the transaction is not a wire transfer.
	ByOrderOf string `json:"by_order_of"`
	// For transfers, the party that is paying the transaction.
	Payer string `json:"payer"`
	// The type of transfer, e.g. 'ACH'
	PaymentMethod string `json:"payment_method"`
	// The name of the payment processor
	PaymentProcessor string `json:"payment_processor"`
	// The payer-supplied description of the transfer.
	Reason string `json:"reason"`
}

// AccountBalance A set of fields describing the balance for an account. Balance information may be cached unless the balance object was returned by `/accounts/balance/get`.
type AccountBalance struct {
	// The amount of funds available to be withdrawn from the account, as determined by the financial institution.  For `credit`-type accounts, the `available` balance typically equals the `limit` less the `current` balance, less any pending outflows plus any pending inflows.  For `depository`-type accounts, the `available` balance typically equals the `current` balance less any pending outflows plus any pending inflows. For `depository`-type accounts, the `available` balance does not include the overdraft limit.  For `investment`-type accounts (or `brokerage`-type accounts for API versions 2018-05-22 and earlier), the `available` balance is the total cash available to withdraw as presented by the institution.  Note that not all institutions calculate the `available`  balance. In the event that `available` balance is unavailable, Plaid will return an `available` balance value of `null`.  Available balance may be cached and is not guaranteed to be up-to-date in realtime unless the value was returned by `/accounts/balance/get`.  If `current` is `null` this field is guaranteed not to be `null`.
	Available float32 `json:"available"`
	// The total amount of funds in or owed by the account.  For `credit`-type accounts, a positive balance indicates the amount owed; a negative amount indicates the lender owing the account holder.  For `loan`-type accounts, the current balance is the principal remaining on the loan, except in the case of student loan accounts at Sallie Mae (`ins_116944`). For Sallie Mae student loans, the account's balance includes both principal and any outstanding interest.  For `investment`-type accounts (or `brokerage`-type accounts for API versions 2018-05-22 and earlier), the current balance is the total value of assets as presented by the institution.  Note that balance information may be cached unless the value was returned by `/accounts/balance/get`; if the Item is enabled for Transactions, the balance will be at least as recent as the most recent Transaction update. If you require realtime balance information, use the `available` balance as provided by `/accounts/balance/get`.  When returned by `/accounts/balance/get`, this field may be `null`. When this happens, `available` is guaranteed not to be `null`.
	Current float32 `json:"current"`
	// For `credit`-type accounts, this represents the credit limit.  For `depository`-type accounts, this represents the pre-arranged overdraft limit, which is common for current (checking) accounts in Europe.  In North America, this field is typically only available for `credit`-type accounts.
	Limit float32 `json:"limit"`
	// The ISO-4217 currency code of the balance. Always null if `unofficial_currency_code` is non-null.
	IsoCurrencyCode string `json:"iso_currency_code"`
	// The unofficial currency code associated with the balance. Always null if `iso_currency_code` is non-null. Unofficial currency codes are used for currencies that do not have official ISO currency codes, such as cryptocurrencies and the currencies of certain countries.  See the [currency code schema](https://plaid.com/docs/api/accounts#currency-code-schema) for a full listing of supported `unofficial_currency_code`s.
	UnofficialCurrencyCode string `json:"unofficial_currency_code"`
	// Timestamp in [ISO 8601](https://wikipedia.org/wiki/ISO_8601) format (`YYYY-MM-DDTHH:mm:ssZ`) indicating the last time that the balance for the given account has been updated  This is currently only provided when the `min_last_updated_datetime` is passed when calling `/accounts/balance/get` for `ins_128026` (Capital One).
	LastUpdatedDatetime *time.Time `json:"last_updated_datetime,omitempty"`
}
