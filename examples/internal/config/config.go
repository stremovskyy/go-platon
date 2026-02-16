package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config aggregates credentials that examples require.
type Config struct {
	MerchantName         string
	MerchantNameWithdraw string
	Login                string
	MerchantID           string
	MerchantIDWithdraw   string
	MerchantKey          string
	SecretKey            string
	SuccessRedirect      string
	FailRedirect         string

	PayerEmail string
	CardToken  string
	WebhookURL string

	AppleContainer string
	GoogleToken    string

	CardNumber string
	CardMonth  string
	CardYear   string
	CardCVV    string
}

var defaultEnvPaths = []string{
	".env.local",
	".env",
	"examples/.env.local",
	"examples/.env",
}

// MustLoad loads configuration and exits the process if required values are missing.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	return cfg
}

// Load populates configuration from environment variables. It optionally pulls
// values from a .env-compatible file to simplify local development.
func Load() (*Config, error) {
	if err := hydrateEnv(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	var err error

	if cfg.MerchantName, err = requireString("PLATON_MERCHANT_NAME"); err != nil {
		return nil, err
	}
	if cfg.MerchantNameWithdraw, err = requireString("PLATON_MERCHANT_NAME_WITHDRAW"); err != nil {
		return nil, err
	}
	if cfg.Login, err = requireString("PLATON_LOGIN"); err != nil {
		return nil, err
	}
	if cfg.MerchantID, err = requireString("PLATON_MERCHANT_ID"); err != nil {
		return nil, err
	}
	if cfg.MerchantIDWithdraw, err = requireString("PLATON_MERCHANT_ID_WITHDRAW"); err != nil {
		return nil, err
	}
	if cfg.MerchantKey, err = requireString("PLATON_MERCHANT_KEY"); err != nil {
		return nil, err
	}
	if cfg.SecretKey, err = requireString("PLATON_SECRET_KEY"); err != nil {
		return nil, err
	}
	if cfg.SuccessRedirect, err = requireString("PLATON_SUCCESS_REDIRECT"); err != nil {
		return nil, err
	}
	if cfg.FailRedirect, err = requireString("PLATON_FAIL_REDIRECT"); err != nil {
		return nil, err
	}
	if cfg.PayerEmail, err = requireString("PLATON_PAYER_EMAIL"); err != nil {
		return nil, err
	}
	if cfg.CardToken, err = requireString("PLATON_CARD_TOKEN"); err != nil {
		return nil, err
	}
	if cfg.WebhookURL, err = requireString("PLATON_WEBHOOK_URL"); err != nil {
		return nil, err
	}
	if cfg.AppleContainer, err = requireString("PLATON_APPLE_CONTAINER"); err != nil {
		return nil, err
	}
	if cfg.GoogleToken, err = requireString("PLATON_GOOGLE_TOKEN"); err != nil {
		return nil, err
	}
	if cfg.CardNumber, err = requireString("PLATON_CARD_NUMBER"); err != nil {
		return nil, err
	}
	if cfg.CardMonth, err = requireString("PLATON_CARD_MONTH"); err != nil {
		return nil, err
	}
	if cfg.CardYear, err = requireString("PLATON_CARD_YEAR"); err != nil {
		return nil, err
	}
	if cfg.CardCVV, err = requireString("PLATON_CARD_CVV"); err != nil {
		return nil, err
	}
	return cfg, nil
}

func hydrateEnv() error {
	if custom := strings.TrimSpace(os.Getenv("PLATON_EXAMPLES_ENV_FILE")); custom != "" {
		if err := loadEnvFile(custom); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("env file %q (referenced by PLATON_EXAMPLES_ENV_FILE) not found", custom)
			}
			return err
		}
		return nil
	}

	for _, path := range defaultEnvPaths {
		if err := loadEnvFile(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}
		// stop after successfully loading the first available file
		return nil
	}

	// none of the default files existed; that's acceptable if env vars are pre-set
	return nil
}

func loadEnvFile(path string) error {
	// #nosec G304 -- configuration files are explicitly chosen by the developer.
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if idx := strings.Index(line, "="); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			val = strings.Trim(val, `"'`)

			if key == "" {
				continue
			}

			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
	}

	return scanner.Err()
}

func requireString(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return "", fmt.Errorf("environment variable %s is empty", key)
	}
	return val, nil
}
