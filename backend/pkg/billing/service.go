package billing

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"web-openssl-backend/internal/models"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/invoice"
	"github.com/stripe/stripe-go/v75/paymentmethod"
	"github.com/stripe/stripe-go/v75/price"
	"github.com/stripe/stripe-go/v75/product"
	"github.com/stripe/stripe-go/v75/subscription"
	"github.com/stripe/stripe-go/v75/webhook"
)

type Service struct {
	secretKey     string
	webhookSecret string
}

type Plan struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           int64   `json:"price"` // in cents
	Currency        string  `json:"currency"`
	Interval        string  `json:"interval"`
	Features        []string `json:"features"`
	OperationsLimit int     `json:"operationsLimit"`
	Popular         bool    `json:"popular"`
}

type CreateSubscriptionRequest struct {
	PriceID       string `json:"priceId"`
	PaymentMethod string `json:"paymentMethod"`
}

type SubscriptionResponse struct {
	ID                   string `json:"id"`
	Status               string `json:"status"`
	CurrentPeriodStart   int64  `json:"currentPeriodStart"`
	CurrentPeriodEnd     int64  `json:"currentPeriodEnd"`
	ClientSecret         string `json:"clientSecret,omitempty"`
	LatestInvoice        string `json:"latestInvoice,omitempty"`
}

type UsageResponse struct {
	CurrentUsage   int    `json:"currentUsage"`
	MonthlyLimit   int    `json:"monthlyLimit"`
	ResetDate      string `json:"resetDate"`
	Plan           string `json:"plan"`
	OverageAllowed bool   `json:"overageAllowed"`
}

func NewService(secretKey string) *Service {
	stripe.Key = secretKey
	return &Service{
		secretKey: secretKey,
	}
}

func (s *Service) GetPlans() []Plan {
	return []Plan{
		{
			ID:          "free",
			Name:        "Free",
			Description: "Perfect for trying out OpenSSL operations",
			Price:       0,
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"50 operations per month",
				"Basic certificate operations",
				"Community support",
				"Web interface access",
			},
			OperationsLimit: 50,
			Popular:         false,
		},
		{
			ID:          "pro",
			Name:        "Pro",
			Description: "For professionals and small teams",
			Price:       1900, // $19.00
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"5,000 operations per month",
				"All OpenSSL operations",
				"Priority support",
				"API access",
				"Advanced features",
				"Team collaboration",
			},
			OperationsLimit: 5000,
			Popular:         true,
		},
		{
			ID:          "enterprise",
			Name:        "Enterprise",
			Description: "For large organizations with advanced needs",
			Price:       9900, // $99.00
			Currency:    "usd",
			Interval:    "month",
			Features: []string{
				"Unlimited operations",
				"White-label solution",
				"SSO integration",
				"24/7 priority support",
				"SLA guarantee",
				"Custom integrations",
				"On-premise deployment",
			},
			OperationsLimit: -1, // unlimited
			Popular:         false,
		},
	}
}

func (s *Service) CreateCustomer(user *models.User) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.FirstName + " " + user.LastName),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", user.ID),
		},
	}

	return customer.New(params)
}

func (s *Service) CreateSubscription(customerID, priceID, paymentMethodID string) (*stripe.Subscription, error) {
	// Attach payment method to customer
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	_, err := paymentmethod.Attach(paymentMethodID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to attach payment method: %w", err)
	}

	// Set as default payment method
	customerParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethodID),
		},
	}
	_, err = customer.Update(customerID, customerParams)
	if err != nil {
		return nil, fmt.Errorf("failed to set default payment method: %w", err)
	}

	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(priceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")

	return subscription.New(subscriptionParams)
}

func (s *Service) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{}
	params.AddExpand("latest_invoice.payment_intent")

	return subscription.Get(subscriptionID, params)
}

func (s *Service) CancelSubscription(subscriptionID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionCancelParams{
		Prorate: stripe.Bool(true),
	}

	return subscription.Cancel(subscriptionID, params)
}

func (s *Service) GetInvoices(customerID string, limit int64) *invoice.Iter {
	params := &stripe.InvoiceListParams{
		Customer: stripe.String(customerID),
	}
	params.Filters.AddFilter("limit", "", fmt.Sprintf("%d", limit))

	return invoice.List(params)
}

func (s *Service) CreateProducts() error {
	plans := s.GetPlans()

	for _, plan := range plans {
		if plan.ID == "free" {
			continue // Skip free plan for Stripe
		}

		// Create product
		productParams := &stripe.ProductParams{
			Name:        stripe.String(plan.Name),
			Description: stripe.String(plan.Description),
			Metadata: map[string]string{
				"plan_id":          plan.ID,
				"operations_limit": fmt.Sprintf("%d", plan.OperationsLimit),
			},
		}

		prod, err := product.New(productParams)
		if err != nil {
			log.Printf("Failed to create product %s: %v", plan.Name, err)
			continue
		}

		// Create price
		priceParams := &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			UnitAmount: stripe.Int64(plan.Price),
			Currency:   stripe.String(plan.Currency),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String(plan.Interval),
			},
			Metadata: map[string]string{
				"plan_id": plan.ID,
			},
		}

		_, err = price.New(priceParams)
		if err != nil {
			log.Printf("Failed to create price for %s: %v", plan.Name, err)
		}
	}

	return nil
}

func (s *Service) HandleWebhook(payload []byte, sigHeader string) error {
	event, err := webhook.ConstructEvent(payload, sigHeader, s.webhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	switch event.Type {
	case "customer.subscription.created":
		return s.handleSubscriptionCreated(event.Data.Object)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(event.Data.Object)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(event.Data.Object)
	case "invoice.payment_succeeded":
		return s.handlePaymentSucceeded(event.Data.Object)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(event.Data.Object)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	return nil
}

func (s *Service) handleSubscriptionCreated(obj map[string]interface{}) error {
	// TODO: Update user subscription in database
	log.Printf("Subscription created: %+v", obj)
	return nil
}

func (s *Service) handleSubscriptionUpdated(obj map[string]interface{}) error {
	// TODO: Update user subscription in database
	log.Printf("Subscription updated: %+v", obj)
	return nil
}

func (s *Service) handleSubscriptionDeleted(obj map[string]interface{}) error {
	// TODO: Cancel user subscription in database
	log.Printf("Subscription deleted: %+v", obj)
	return nil
}

func (s *Service) handlePaymentSucceeded(obj map[string]interface{}) error {
	// TODO: Handle successful payment
	log.Printf("Payment succeeded: %+v", obj)
	return nil
}

func (s *Service) handlePaymentFailed(obj map[string]interface{}) error {
	// TODO: Handle failed payment
	log.Printf("Payment failed: %+v", obj)
	return nil
}

func (s *Service) ValidateWebhookSignature(r *http.Request) ([]byte, error) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(nil, r.Body, MaxBodyBytes)
	payload, err := json.Marshal(r.Body)
	if err != nil {
		return nil, err
	}

	return payload, nil
}