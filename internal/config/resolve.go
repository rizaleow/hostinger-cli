package config

import "os"

// Resolved is the effective runtime configuration for one CLI invocation.
type Resolved struct {
	Token   string
	BaseURL string
	Profile string
	Source  string // where the token came from: flag|env|keyring|file
}

// ResolveOptions controls token resolution precedence.
type ResolveOptions struct {
	FlagToken   string
	FlagBaseURL string
	Profile     string
	UseKeyring  bool
}

// Resolve applies the precedence chain: flag → env → keyring → file.
// It never returns an error for a missing token; callers that need auth
// must check Resolved.Token themselves.
func Resolve(file *File, opts ResolveOptions) Resolved {
	r := Resolved{Profile: opts.Profile}
	if r.Profile == "" {
		r.Profile = file.CurrentProfile
	}
	if r.Profile == "" {
		r.Profile = "default"
	}
	p := file.Profile(r.Profile)

	switch {
	case opts.FlagToken != "":
		r.Token, r.Source = opts.FlagToken, "flag"
	case os.Getenv(EnvToken) != "":
		r.Token, r.Source = os.Getenv(EnvToken), "env"
	case opts.UseKeyring || file.UseKeyring:
		if tok, _ := KeyringGet(r.Profile); tok != "" {
			r.Token, r.Source = tok, "keyring"
		}
	}
	if r.Token == "" && p.Token != "" {
		r.Token, r.Source = p.Token, "file"
	}

	switch {
	case opts.FlagBaseURL != "":
		r.BaseURL = opts.FlagBaseURL
	case os.Getenv(EnvBaseURL) != "":
		r.BaseURL = os.Getenv(EnvBaseURL)
	case p.BaseURL != "":
		r.BaseURL = p.BaseURL
	}
	return r
}
