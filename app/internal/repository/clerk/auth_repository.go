package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/lestrrat-go/jwx/jwk"
)

type ClerkAuthRepository interface {
	FetchJWKS() (jwk.Set, error)
}

type clerkAuthRepository struct{}

func NewClerkAuthRepository() ClerkAuthRepository {
	return &clerkAuthRepository{}
}

func (r *clerkAuthRepository) FetchJWKS() (jwk.Set, error) {
	clerkJWKSURL := os.Getenv("CLERK_JWKS_URL")
	if clerkJWKSURL == "" {
		return nil, fmt.Errorf("CLERK_JWKS_URL environment variable is not set")
	}

	keySet, err := jwk.Fetch(context.Background(), clerkJWKSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	return keySet, nil
}
