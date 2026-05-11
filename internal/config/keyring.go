package config

import "github.com/zalando/go-keyring"

const keyringService = "hostinger-cli"

// KeyringGet returns the stored token for a profile, or "" if absent.
func KeyringGet(profile string) (string, error) {
	tok, err := keyring.Get(keyringService, profile)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return tok, nil
}

// KeyringSet stores a token in the OS keychain.
func KeyringSet(profile, token string) error {
	return keyring.Set(keyringService, profile, token)
}

// KeyringDelete removes a token from the OS keychain.
func KeyringDelete(profile string) error {
	err := keyring.Delete(keyringService, profile)
	if err == keyring.ErrNotFound {
		return nil
	}
	return err
}
