package auth

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/transcom/mymove/pkg/app"
	"go.uber.org/zap"
)

const myProviderName = "myProvider"
const officeProviderName = "officeProvider"

func getLoginGovProviderForRequest(r *http.Request) (*openidConnect.Provider, error) {
	providerName := myProviderName
	if app.IsOfficeApp(r) {
		providerName = officeProviderName
	}
	gothProvider, err := goth.GetProvider(providerName)
	if err != nil {
		return nil, err
	}
	return gothProvider.(*openidConnect.Provider), nil
}

// LoginGovProvider facilitates generating URLs and parameters for interfacing with Login.gov
type LoginGovProvider struct {
	hostname  string
	secretKey string
	logger    *zap.Logger
}

// NewLoginGovProvider returns a new LoginGovProvider
func NewLoginGovProvider(hostname string, secretKey string, logger *zap.Logger) LoginGovProvider {
	return LoginGovProvider{
		hostname:  hostname,
		secretKey: secretKey,
		logger:    logger,
	}
}

func (p LoginGovProvider) getOpenIDProvider(hostname string, clientID string, callbackProtocol string, callbackPort string) (goth.Provider, error) {
	return openidConnect.New(
		clientID,
		p.secretKey,
		fmt.Sprintf("%s%s:%s/auth/login-gov/callback", callbackProtocol, hostname, callbackPort),
		fmt.Sprintf("https://%s/.well-known/openid-configuration", p.hostname),
	)
}

// RegisterProvider registers Login.gov with Goth, which uses
// auto-discovery to get the OpenID configuration
func (p LoginGovProvider) RegisterProvider(myHostname string, myClientID string, officeHostname string, officeClientID string, callbackProtocol string, callbackPort string) error {

	myProvider, err := p.getOpenIDProvider(myHostname, myClientID, callbackProtocol, callbackPort)
	if err != nil {
		p.logger.Error("getting open_id provider", zap.String("host", myHostname), zap.Error(err))
		return err
	}
	myProvider.SetName(myProviderName)
	officeProvider, err := p.getOpenIDProvider(officeHostname, officeClientID, callbackProtocol, callbackPort)
	if err != nil {
		p.logger.Error("getting open_id provider", zap.String("host", officeHostname), zap.Error(err))
		return err
	}
	officeProvider.SetName(officeProviderName)
	goth.UseProviders(myProvider, officeProvider)
	return nil
}

func generateNonce() string {
	nonceBytes := make([]byte, 64)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 64; i++ {
		nonceBytes[i] = byte(random.Int63() % 256)
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

// AuthorizationURL returns a URL for login.gov authorization with required params
func (p LoginGovProvider) AuthorizationURL(r *http.Request) (string, error) {
	provider, err := getLoginGovProviderForRequest(r)
	if err != nil {
		p.logger.Error("Get Goth provider", zap.Error(err))
		return "", err
	}
	state := generateNonce()
	sess, err := provider.BeginAuth(state)
	if err != nil {
		p.logger.Error("Goth begin auth", zap.Error(err))
		return "", err
	}

	baseURL, err := sess.GetAuthURL()
	if err != nil {
		p.logger.Error("Goth get auth URL", zap.Error(err))
		return "", err
	}

	authURL, err := url.Parse(baseURL)
	if err != nil {
		p.logger.Error("Parse auth URL", zap.Error(err))
		return "", err
	}

	params := authURL.Query()
	params.Add("acr_values", "http://idmanagement.gov/ns/assurance/loa/1")
	params.Add("nonce", state)
	params.Set("scope", "openid email")

	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}

// LogoutURL returns a full URL to log out of login.gov with required params
func (p LoginGovProvider) LogoutURL(redirectURL string, idToken string) string {
	logoutPath, _ := url.Parse(fmt.Sprintf("https://%s/openid_connect/logout", p.hostname))
	// Parameters taken from https://developers.login.gov/oidc/#logout
	params := url.Values{
		"id_token_hint":            {idToken},
		"post_logout_redirect_uri": {redirectURL},
		"state":                    {generateNonce()},
	}

	logoutPath.RawQuery = params.Encode()
	return logoutPath.String()
}

// TokenURL returns a full URL to retrieve a user token from login.gov
func (p LoginGovProvider) TokenURL() string {
	// TODO: Get the token endpoint URL from Goth instead when
	// https://github.com/markbates/goth/pull/207 is resolved
	return fmt.Sprintf("https://%s/api/openid_connect/token", p.hostname)
}

// TokenParams creates query params for use in the token endpoint
func (p LoginGovProvider) TokenParams(code string, clientID string, expiry time.Time) (url.Values, error) {
	clientAssertion, err := p.createClientAssertionJWT(clientID, expiry)
	params := url.Values{
		"client_assertion":      {clientAssertion},
		"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
		"code":                  {code},
		"grant_type":            {"authorization_code"},
	}

	return params, err
}

func (p LoginGovProvider) createClientAssertionJWT(clientID string, expiry time.Time) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    clientID,
		Subject:   clientID,
		Audience:  p.TokenURL(),
		Id:        generateNonce(),
		ExpiresAt: expiry.Unix(),
	}

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(p.secretKey))
	if err != nil {
		p.logger.Error("JWT parse private key from PEM", zap.Error(err))
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwt, err := token.SignedString(rsaKey)
	if err != nil {
		p.logger.Error("Signing JWT", zap.Error(err))
	}
	return jwt, err
}

// LoginGovTokenResponse is a struct for parsing responses from the token endpoint
type LoginGovTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}
