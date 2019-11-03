package googleCredentials

type GoogleCredentials struct {
	ClientId                string   `json:"client_id"`
	ProjectId               string   `json:"project_id"`
	AuthUri                 string   `json:"auth_uri"`
	TokenUri                string   `json:"token_uri"`
	AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
	ClientSecret            string   `json:"client_secret"`
	RedirectUris            []string `json:"redirect_uris"`
}

type GoogleCredentialsWrapper struct {
	Installed GoogleCredentials `json:"installed"`
}

func Get() *GoogleCredentialsWrapper {
	return &GoogleCredentialsWrapper{
		Installed: GoogleCredentials{
			"687759170793-sgl4fsv331ipq5639dm25maj0vqodge1.apps.googleusercontent.com",
			"quickstart-1572131343445",
			"https://accounts.google.com/o/oauth2/auth",
			"https://oauth2.googleapis.com/token",
			"https://www.googleapis.com/oauth2/v1/certs",
			"dHzrSU3ZqyQFPDyNajBrqAfb",
			[]string{"urn:ietf:wg:oauth:2.0:oob", "http://localhost"},
		},
	}
}
