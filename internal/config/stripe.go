package config

import (
	"os"

	"github.com/stripe/stripe-go/v78"
)

// InitStripe inicializa a configuração da API do Stripe.
func InitStripe() {
	apiKey := os.Getenv("STRIPE_API_KEY")
	stripe.Key = apiKey
}
