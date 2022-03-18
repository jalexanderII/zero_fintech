package models

import (
	"github.com/plaid/plaid-go/plaid"
	"github.com/stripe/stripe-go/v72"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token for use of plaid public token retrieval
type Token struct {
	ID            primitive.ObjectID `bson:"_id"`
	User          *User              `bson:"user"`
	Value         string             `bson:"value"`
	ItemId        string             `bson:"item_id"`
	Institution   string             `bson:"institution"`
	InstitutionID string             `bson:"institution_id"`
}

type StripeToken struct {
	ID               primitive.ObjectID                        `bson:"_id"`
	User             *User                                     `bson:"user"`
	CustomerId       string                                    `bson:"customer_id"`
	CheckingAccounts map[PlaidAccountId]*StripeCustomerAccount `bson:"checking_accounts"`
	CreditAccounts   map[PlaidAccountId]*StripeCustomerAccount `bson:"credit_accounts"`
}

type PlaidAccountId struct {
	Value string `bson:"account_id"`
}

type StripeCustomerAccount struct {
	StripeToken   string `bson:"token"`
	AccountId     string `bson:"account_id"`
	AccountName   string `bson:"account_name"`
	AccountStatus string `bson:"account_status"`
}

// StripeChargeParams To charge a credit card or other payment source, you create a Charge object.
// If your API key is in test mode, the supplied payment source (e.g., card) won't actually be charged,
// although everything else will occur as if in live mode.
type StripeChargeParams struct {
	// Housing source and destination stripe account tokens
	StripeToken *StripeToken `json:"stripe_token"`
	// Amount intended to be collected by this payment. A positive integer representing how much to charge in the
	// [smallest currency unit](https://stripe.com/docs/currencies#zero-decimal) (e.g., 100 cents to charge $1.00 or
	// 100 to charge ¥100, a zero-decimal currency). The minimum amount is $0.50 US or
	// [equivalent in charge currency](https://stripe.com/docs/currencies#minimum-and-maximum-charge-amounts).
	// The amount value supports up to eight digits (e.g., a value of 99999999 for a USD charge of $999,999.99).
	Amount int64 `json:"amount"`
	// Three-letter [ISO currency code](https://www.iso.org/iso-4217-currency-codes.html), in lowercase.
	// Must be a [supported currency](https://stripe.com/docs/currencies).
	Currency string `json:"currency"`
	// The ID of an existing customer that will be associated with this request. This field may only be updated
	// if there is no existing associated customer with this charge.
	Customer string `json:"customer"`
	// An arbitrary string which you can attach to a charge object. It is displayed when in the web interface alongside
	// the charge. Note that if you use Stripe to send automatic email receipts to your customers, your
	// receipt emails will include the `description` of the charge(s) that they are describing.
	Description string             `json:"description"`
	Destination *DestinationParams `json:"destination"`
	// This is the email address that the receipt for this charge will be sent to.
	// If this field is updated, then a new email receipt will be sent to the updated address.
	ReceiptEmail string `json:"receipt_email"`
}

type StripeChargeResponse struct {
	stripe.APIResource
	// Amount intended to be collected by this payment. A positive integer representing how much to charge in the [smallest currency unit](https://stripe.com/docs/currencies#zero-decimal) (e.g., 100 cents to charge $1.00 or 100 to charge ¥100, a zero-decimal currency). The minimum amount is $0.50 US or [equivalent in charge currency](https://stripe.com/docs/currencies#minimum-and-maximum-charge-amounts). The amount value supports up to eight digits (e.g., a value of 99999999 for a USD charge of $999,999.99).
	Amount int64 `json:"amount"`
	// Authorization code on the charge.
	AuthorizationCode string `json:"authorization_code"`
	// Time at which the object was created. Measured in seconds since the Unix epoch.
	Created int64 `json:"created"`
	// ID of the customer this charge is for if one exists.
	Customer *stripe.Customer `json:"customer"`
	// Details about the dispute if the charge has been disputed.
	Dispute *stripe.Dispute `json:"dispute"`
	// Whether the charge has been disputed.
	Disputed bool `json:"disputed"`
	// Error code explaining reason for charge failure if available (see [the errors section](https://stripe.com/docs/api#errors) for a list of codes).
	FailureCode string `json:"failure_code"`
	// Message to user further explaining reason for charge failure if available.
	FailureMessage string `json:"failure_message"`
	// Information on fraud assessments for the charge.
	FraudDetails *stripe.FraudDetails `json:"fraud_details"`
	// Unique identifier for the object.
	ID string `json:"id"`
	// Details about whether the payment was accepted, and why. See [understanding declines](https://stripe.com/docs/declines) for details.
	Outcome *stripe.ChargeOutcome `json:"outcome"`
	// `true` if the charge succeeded, or was successfully authorized for later capture.
	Paid bool `json:"paid"`
	// Whether the charge has been fully refunded. If the charge is only partially refunded, this attribute will still be false.
	Refunded bool `json:"refunded"`
	// A list of refunds that have been applied to the charge.
	Refunds *stripe.RefundList `json:"refunds"`
	// ID of the review associated with this charge if one exists.
	Review *stripe.Review `json:"review"`
	// The status of the payment is either `succeeded`, `pending`, or `failed`.
	Status stripe.ChargeStatus `json:"status"`
}

type DestinationParams struct {
	// ID of an existing, connected Stripe account.
	Account string `json:"account"`
	// The amount to transfer to the destination account without creating an `Application Fee` object.
	// Cannot be combined with the `application_fee` parameter. Must be less than or equal to the charge amount.
	Amount int64 `json:"amount"`
}

type CreateLinkTokenResponse struct {
	UserId string
	Token  string
}

// User is a DB Serialization of Proto User
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
}

type LiabilitiesResponse struct {
	Liabilities []plaid.CreditCardLiability `json:"liabilities"`
}

type TransactionsResponse struct {
	Accounts     []plaid.AccountBase `json:"accounts,omitempty"`
	Transactions []plaid.Transaction `json:"transactions,omitempty"`
}

type PlaidMetaData struct {
	Institution struct {
		Name          string `json:"name"`
		InstitutionId string `json:"institution_id"`
	} `json:"institution"`
	Accounts []struct {
		Id                 string `json:"id"`
		Name               string `json:"name"`
		Mask               string `json:"mask"`
		Type               string `json:"type"`
		Subtype            string `json:"subtype"`
		VerificationStatus string `json:"verification_status,omitempty"`
	} `json:"accounts"`
	LinkSessionId string `json:"link_session_id"`
}
