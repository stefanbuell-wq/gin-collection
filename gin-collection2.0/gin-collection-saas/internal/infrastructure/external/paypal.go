package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/gin-collection-saas/pkg/logger"
)

// PayPalClient handles PayPal API interactions
type PayPalClient struct {
	clientID     string
	clientSecret string
	baseURL      string
	mode         string // "sandbox" or "live"
	httpClient   *http.Client
	accessToken  string
	tokenExpiry  time.Time
}

// PayPalConfig holds PayPal configuration
type PayPalConfig struct {
	ClientID     string
	ClientSecret string
	Mode         string // "sandbox" or "live"
}

// PayPalSubscriptionRequest represents a subscription creation request
type PayPalSubscriptionRequest struct {
	PlanID             string                        `json:"plan_id"`
	StartTime          string                        `json:"start_time,omitempty"`
	Subscriber         *PayPalSubscriber             `json:"subscriber,omitempty"`
	ApplicationContext *PayPalApplicationContext     `json:"application_context,omitempty"`
}

// PayPalSubscriber represents subscriber information
type PayPalSubscriber struct {
	EmailAddress string             `json:"email_address,omitempty"`
	Name         *PayPalName        `json:"name,omitempty"`
}

// PayPalName represents name information
type PayPalName struct {
	GivenName string `json:"given_name,omitempty"`
	Surname   string `json:"surname,omitempty"`
}

// PayPalApplicationContext represents application context
type PayPalApplicationContext struct {
	BrandName          string `json:"brand_name,omitempty"`
	ReturnURL          string `json:"return_url,omitempty"`
	CancelURL          string `json:"cancel_url,omitempty"`
	UserAction         string `json:"user_action,omitempty"` // "SUBSCRIBE_NOW" or "CONTINUE"
	PaymentMethod      *PayPalPaymentMethod `json:"payment_method,omitempty"`
}

// PayPalPaymentMethod represents payment method preferences
type PayPalPaymentMethod struct {
	PayerSelected  string `json:"payer_selected,omitempty"`
	PayeePreferred string `json:"payee_preferred,omitempty"` // "IMMEDIATE_PAYMENT_REQUIRED"
}

// PayPalSubscriptionResponse represents a subscription response
type PayPalSubscriptionResponse struct {
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	StatusUpdateTime string              `json:"status_update_time"`
	PlanID        string                 `json:"plan_id"`
	StartTime     string                 `json:"start_time"`
	Quantity      string                 `json:"quantity"`
	Links         []PayPalLink           `json:"links"`
	BillingInfo   *PayPalBillingInfo     `json:"billing_info,omitempty"`
}

// PayPalLink represents a HATEOAS link
type PayPalLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// PayPalBillingInfo represents billing information
type PayPalBillingInfo struct {
	OutstandingBalance    *PayPalMoney          `json:"outstanding_balance,omitempty"`
	CycleExecutions       []PayPalCycleExecution `json:"cycle_executions,omitempty"`
	NextBillingTime       string                `json:"next_billing_time,omitempty"`
	FailedPaymentsCount   int                   `json:"failed_payments_count,omitempty"`
}

// PayPalCycleExecution represents a billing cycle execution
type PayPalCycleExecution struct {
	TenureType          string `json:"tenure_type"`
	Sequence            int    `json:"sequence"`
	CyclesCompleted     int    `json:"cycles_completed"`
	CyclesRemaining     int    `json:"cycles_remaining"`
	TotalCycles         int    `json:"total_cycles"`
}

// PayPalMoney represents money
type PayPalMoney struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

// PayPalAuthResponse represents OAuth2 token response
type PayPalAuthResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int    `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

// PayPalError represents PayPal API error
type PayPalError struct {
	Name    string               `json:"name"`
	Message string               `json:"message"`
	DebugID string               `json:"debug_id"`
	Details []PayPalErrorDetail  `json:"details,omitempty"`
}

// PayPalErrorDetail represents error detail
type PayPalErrorDetail struct {
	Field       string `json:"field"`
	Value       string `json:"value"`
	Location    string `json:"location"`
	Issue       string `json:"issue"`
	Description string `json:"description"`
}

// NewPayPalClient creates a new PayPal client
func NewPayPalClient(cfg *PayPalConfig) *PayPalClient {
	baseURL := "https://api-m.sandbox.paypal.com"
	if cfg.Mode == "live" {
		baseURL = "https://api-m.paypal.com"
	}

	return &PayPalClient{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		baseURL:      baseURL,
		mode:         cfg.Mode,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getAccessToken retrieves or refreshes the access token
func (c *PayPalClient) getAccessToken() (string, error) {
	// Return cached token if still valid
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	// Request new token
	url := c.baseURL + "/v1/oauth2/token"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.SetBasicAuth(c.clientID, c.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp PayPalAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	// Cache token (subtract 60s for safety margin)
	c.accessToken = authResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn-60) * time.Second)

	logger.Debug("PayPal access token refreshed", "expires_in", authResp.ExpiresIn)

	return c.accessToken, nil
}

// doRequest performs an authenticated API request
func (c *PayPalClient) doRequest(method, path string, body interface{}, result interface{}) error {
	token, err := c.getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
		logger.Debug("PayPal request", "method", method, "path", path, "body", string(jsonBody))
	}

	url := c.baseURL + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	logger.Debug("PayPal response", "status", resp.StatusCode, "body", string(bodyBytes))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var paypalErr PayPalError
		if err := json.Unmarshal(bodyBytes, &paypalErr); err == nil {
			return fmt.Errorf("PayPal API error: %s - %s (debug_id: %s)", paypalErr.Name, paypalErr.Message, paypalErr.DebugID)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// CreateSubscription creates a new PayPal subscription
func (c *PayPalClient) CreateSubscription(req *PayPalSubscriptionRequest) (*PayPalSubscriptionResponse, error) {
	var resp PayPalSubscriptionResponse
	err := c.doRequest("POST", "/v1/billing/subscriptions", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	logger.Info("PayPal subscription created", "subscription_id", resp.ID, "status", resp.Status)

	return &resp, nil
}

// GetSubscription retrieves a subscription by ID
func (c *PayPalClient) GetSubscription(subscriptionID string) (*PayPalSubscriptionResponse, error) {
	var resp PayPalSubscriptionResponse
	path := fmt.Sprintf("/v1/billing/subscriptions/%s", subscriptionID)
	err := c.doRequest("GET", path, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &resp, nil
}

// CancelSubscription cancels a subscription
func (c *PayPalClient) CancelSubscription(subscriptionID, reason string) error {
	path := fmt.Sprintf("/v1/billing/subscriptions/%s/cancel", subscriptionID)

	reqBody := map[string]string{
		"reason": reason,
	}

	err := c.doRequest("POST", path, reqBody, nil)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	logger.Info("PayPal subscription cancelled", "subscription_id", subscriptionID, "reason", reason)

	return nil
}

// SuspendSubscription suspends a subscription
func (c *PayPalClient) SuspendSubscription(subscriptionID, reason string) error {
	path := fmt.Sprintf("/v1/billing/subscriptions/%s/suspend", subscriptionID)

	reqBody := map[string]string{
		"reason": reason,
	}

	err := c.doRequest("POST", path, reqBody, nil)
	if err != nil {
		return fmt.Errorf("failed to suspend subscription: %w", err)
	}

	logger.Info("PayPal subscription suspended", "subscription_id", subscriptionID, "reason", reason)

	return nil
}

// ActivateSubscription activates a suspended subscription
func (c *PayPalClient) ActivateSubscription(subscriptionID, reason string) error {
	path := fmt.Sprintf("/v1/billing/subscriptions/%s/activate", subscriptionID)

	reqBody := map[string]string{
		"reason": reason,
	}

	err := c.doRequest("POST", path, reqBody, nil)
	if err != nil {
		return fmt.Errorf("failed to activate subscription: %w", err)
	}

	logger.Info("PayPal subscription activated", "subscription_id", subscriptionID, "reason", reason)

	return nil
}

// GetApprovalURL extracts the approval URL from subscription response links
func (r *PayPalSubscriptionResponse) GetApprovalURL() string {
	for _, link := range r.Links {
		if link.Rel == "approve" {
			return link.Href
		}
	}
	return ""
}
