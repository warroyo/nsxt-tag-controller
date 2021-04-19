package nsxt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	sdkclient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	// +kubebuilder:scaffold:imports
)

type NsxtClient struct {
	PolicySecurityContext *core.SecurityContextImpl
	PolicyHTTPClient      *http.Client
	Host                  string
	User                  string
	Pass                  string
}

type overwriteHeaderProcessor struct {
}

func newOverwriteHeaderProcessor() *overwriteHeaderProcessor {
	return &overwriteHeaderProcessor{}
}

func (processor overwriteHeaderProcessor) Process(req *http.Request) error {
	req.Header.Set("X-Allow-Overwrite", "true")
	return nil
}

func GetPolicyConnector(clients interface{}) *sdkclient.RestConnector {
	c := clients.(NsxtClient)
	connector := sdkclient.NewRestConnector(c.Host, *c.PolicyHTTPClient)
	connector.AddRequestProcessor(newOverwriteHeaderProcessor())
	if c.PolicySecurityContext != nil {
		connector.SetSecurityContext(c.PolicySecurityContext)
	}

	return connector
}

func getConnectorTLSConfig() (*tls.Config, error) {

	insecure := true
	caCert := ""
	tlsConfig := tls.Config{InsecureSkipVerify: insecure}

	if len(caCert) > 0 {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))

		tlsConfig.RootCAs = caCertPool
	}

	return &tlsConfig, nil
}

func ConfigurePolicyConnectorData(clients *NsxtClient) error {

	host := clients.Host
	username := clients.User
	password := clients.Pass
	if host == "" {
		return fmt.Errorf("host must be provided")
	}

	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}

	securityCtx := core.NewSecurityContextImpl()
	securityContextNeeded := true

	if securityContextNeeded {

		if username == "" {
			return fmt.Errorf("username must be provided")
		}

		if password == "" {
			return fmt.Errorf("password must be provided")
		}

		securityCtx.SetProperty(security.AUTHENTICATION_SCHEME_ID, security.USER_PASSWORD_SCHEME_ID)
		securityCtx.SetProperty(security.USER_KEY, username)
		securityCtx.SetProperty(security.PASSWORD_KEY, password)
	}

	tlsConfig, err := getConnectorTLSConfig()
	if err != nil {
		return err
	}

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}

	httpClient := http.Client{Transport: tr}
	clients.PolicyHTTPClient = &httpClient
	if securityContextNeeded {
		clients.PolicySecurityContext = securityCtx
	}
	clients.Host = host

	return nil
}
