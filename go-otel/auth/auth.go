package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var AZURE_TENANT_ID = "33e01921-4d64-4f8c-a055-5bdaffd5e33d"

var AZURE_CLIENT_ID = "9c7ae59d-9323-4423-a0da-38ddce774875"

type AADToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

type Token struct {
	Token string
	ExpiresOn time.Time
}

var cachedToken = Token {
	Token: "",
	ExpiresOn: time.Now(),
}

// Get token from Azure managed environment, AML or other Azure VMs
func GetTokenInAzure(
	ctx context.Context,
	resource string,
	clientId string,
) (string, error) {

	println("Try to get token in Azure.")
	println(resource)
	fromAML := false
	if os.Getenv("USING_AML") == "true" {
		fromAML = true
	}

	println("Not from AML")

	logPrefix := fmt.Sprintf(
		"GetMSIToken(fromAML=%v, resource=%s, UAIdentityClientId=%s): ",
		fromAML, resource, clientId)

	if resource == "" || clientId == "" {
		fmt.Printf(logPrefix + "one or more inputs is blank\n")
		return "", fmt.Errorf("one or more inputs is blank")
	}

	var queryStr string
	var secret string
	if fromAML {
		queryStr = fmt.Sprintf(
			"%v?version=2019-08-01",
			os.Getenv("MSI_ENDPOINT"),
		)
		secret = os.Getenv("MSI_SECRET")
	} else {
		queryStr = "http://169.254.169.254/metadata/" +
			"identity/oauth2/token?api-version=2018-02-01"
	}
	queryUrl, _ := url.Parse(queryStr)
	msiParams := queryUrl.Query()
	msiParams.Add("resource", resource)
	msiParams.Add("client_id", clientId)
	queryUrl.RawQuery = msiParams.Encode()
	req, err := http.NewRequestWithContext(
		ctx, "GET", queryUrl.String(), nil)
	if err != nil {
		fmt.Printf(logPrefix+
			"Error creating http request to IMDS: %s\n",
			err)
		return "", err
	}

	client := &http.Client{}
	req.Header.Add("Metadata", "true")
	if fromAML {
		req.Header.Add("secret", secret)
	}
	resp, err := client.Do(req)
	if ctx.Err() != nil {
		return "", ctx.Err()
	} else if err != nil {
		fmt.Printf(logPrefix+
			"Error calling IMDS token endpoint: HTTP GET %v: %s\n",
			queryUrl, err)
		return "", err
	}

	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf(logPrefix+
			"Error reading response body after HTTP GET %v: %s\n",
			queryUrl, err)
		return "", err
	}

	if resp.StatusCode != 200 {
		fmt.Printf(logPrefix+
			"Failed to get token. Details: %s %s\n",
			resp.Status, string(respBytes))
		return "", fmt.Errorf("IMDS returned: %s", resp.Status)
	}

	var token AADToken
	err = json.Unmarshal(respBytes, &token)
	if err != nil {
		fmt.Printf(logPrefix+
			"Error unmarshalling response: %s", err)
		return "", err
	}

	fmt.Printf("Successfully get token in Azure Env: %s\n", token.AccessToken)
	// set the cached token.
	cachedToken.Token = token.AccessToken

	timestamp, err := strconv.ParseInt(token.ExpiresOn, 10, 64)  
	expiresOn := time.Unix(timestamp, 0)
	if err != nil {  
		fmt.Println("Error parsing AADToken.ExpiresOn", err)  
		return "", err
	}
	cachedToken.ExpiresOn = expiresOn
	return token.AccessToken, nil
}

// Get token using app id, should be use in local development or in public environment.
func GetTokenWithAppID(ctx context.Context, resource string) (string, error) {
	//os.Setenv("AZURE_TENANT_ID", AZURE_TENANT_ID)
	//os.Setenv("AZURE_CLIENT_ID", AZURE_CLIENT_ID)

	// Make sure secret is present
	OTELP_CLIENT_SECRET := os.Getenv("OTLP_CLIENT_SECRET")
	if OTELP_CLIENT_SECRET == "" {
		fmt.Printf("OTLP_CLIENT_SECRET is missing.\n")
	}

	logPrefix := fmt.Sprintf(
		"GetTokenWithAppID(tenantId=%s, clientId=%s, clientSecret=%s): ",
		os.Getenv("AZURE_TENANT_ID"), os.Getenv("AZURE_CLIENT_ID"), os.Getenv("AZURE_CLIENT_SECRET"))

	//credential, err := azidentity.NewDefaultAzureCredential(nil)
	print("Use clientid and secret\n")
	credential, err := azidentity.NewClientSecretCredential(  
		AZURE_TENANT_ID,  
		AZURE_CLIENT_ID,  
		OTELP_CLIENT_SECRET,  
		nil)

	if err != nil {
		fmt.Printf(logPrefix + "Error creating default azure credential: %s\n", err)
		return "", err
	}
	options := policy.TokenRequestOptions{
		Scopes: []string{resource},
	}

	token, err := credential.GetToken(ctx, options)
	if err != nil {
		fmt.Printf(logPrefix + "Error GetToken: %s\n", err)
		return "", err
	}

	cachedToken.Token = token.Token
	cachedToken.ExpiresOn = token.ExpiresOn
	// print(token.ExpiresOn.Format(time.RFC3339))

	fmt.Printf("Successfully get token with AppId and Secret: %s\n", token.Token)
	return token.Token, nil
}

func GetTokenWithDefaultUAI(ctx context.Context, resource string) (string, error) {
	
	logPrefix := fmt.Sprintf(
		"GetTokenWithDefaultUAI %s", os.Getenv("AZURE_CLIENT_ID"))
	credential, _ := azidentity.NewDefaultAzureCredential(nil)
	options := policy.TokenRequestOptions{
		Scopes: []string{resource},
	}

	UAI_CLIENT_ID := "1baa67a6-59c1-4c0f-a675-ee2682793b42"
	fmt.Printf("Use hard coded id %s\n", UAI_CLIENT_ID)
	//clientID := azidentity.ClientID(UAI_CLIENT_ID)
	//opts := azidentity.ManagedIdentityCredentialOptions{ID: clientID}
	//credential, _ := azidentity.NewManagedIdentityCredential(&opts)

	token, err := credential.GetToken(ctx, options)
	if err != nil {
		fmt.Printf(logPrefix + "Error GetToken: %s\n", err)
		return "", err
	}

	cachedToken.Token = token.Token
	cachedToken.ExpiresOn = token.ExpiresOn
	// print(token.ExpiresOn.Format(time.RFC3339))

	fmt.Printf("Successfully get token with Default UAI %s\n", token.Token)
	return token.Token, nil
}

// Main api to call to get token.
func GetToken(ctx context.Context, resource string) (string, error) {
	// Try to get cached token.
	currentTime := time.Now()
	fiveMinsfromNow := currentTime.Add(5 * time.Minute)

	if cachedToken.ExpiresOn.After(fiveMinsfromNow) {
		fmt.Printf("Cached token is still valid. \n")
		return cachedToken.Token, nil
	}

	// Get token from UAI_CLIENT_ID
	UAI_CLIENT_ID := os.Getenv("UAI_CLIENT_ID")
	fmt.Printf("Getting UAI_CLIENT_ID: %s\n", UAI_CLIENT_ID)

	if UAI_CLIENT_ID != "" {
		// First try Azure way
		
		for i := 0; i < 10; i++ {
			token, err := GetTokenInAzure(ctx, resource, UAI_CLIENT_ID)
			// token, err := GetTokenWithDefaultUAI(ctx, resource)

			if err == nil {
				return token, nil
			}

			// This is to give me enough time to check container instance log.
			print("Going to try again in 60s")
			time.Sleep(60 * time.Second)
		}
	}

	// Fall back to use clientId and secret
	token, err := GetTokenWithAppID(ctx, resource)

	if err == nil {
		return token, nil
	}

	return "", err
}