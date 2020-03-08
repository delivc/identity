package api

import "net/http"

// ProviderSettings an JSON struct to determine of a provider is enabled
type ProviderSettings struct {
	Bitbucket bool `json:"bitbucket"`
	GitHub    bool `json:"github"`
	GitLab    bool `json:"gitlab"`
	Google    bool `json:"google"`
	Facebook  bool `json:"facebook"`
	Email     bool `json:"email"`
	SAML      bool `json:"saml"`
}

// ProviderLabels the labels of the SAML providers
type ProviderLabels struct {
	SAML string `json:"saml,omitempty"`
}

// Settings holds the settings for the external proiders
type Settings struct {
	ExternalProviders ProviderSettings `json:"external"`
	ExternalLabels    ProviderLabels   `json:"external_labels"`
	DisableSignup     bool             `json:"disable_signup"`
	Autoconfirm       bool             `json:"autoconfirm"`
}

// Settings returns the current IDentity provider settings
func (a *API) Settings(w http.ResponseWriter, r *http.Request) error {
	config := a.getConfig(r.Context())

	return sendJSON(w, http.StatusOK, &Settings{
		ExternalProviders: ProviderSettings{
			Bitbucket: config.External.Bitbucket.Enabled,
			GitHub:    config.External.Github.Enabled,
			GitLab:    config.External.Gitlab.Enabled,
			Google:    config.External.Google.Enabled,
			Facebook:  config.External.Facebook.Enabled,
			Email:     !config.External.Email.Disabled,
			SAML:      config.External.Saml.Enabled,
		},
		ExternalLabels: ProviderLabels{
			SAML: config.External.Saml.Name,
		},
		DisableSignup: config.DisableSignup,
		Autoconfirm:   config.Mailer.Autoconfirm,
	})
}
