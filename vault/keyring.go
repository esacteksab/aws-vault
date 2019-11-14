package vault

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/99designs/keyring"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func NewKeyringCredentials(k keyring.Keyring, credentialsName string) *credentials.Credentials {
	return credentials.NewCredentials(NewKeyringProvider(k, credentialsName))
}

func NewKeyringProvider(k keyring.Keyring, credentialsName string) *KeyringProvider {
	return &KeyringProvider{k, credentialsName}
}

type KeyringProvider struct {
	keyring         keyring.Keyring
	credentialsName string
}

func (p *KeyringProvider) IsExpired() bool {
	return false
}

func (p *KeyringProvider) Retrieve() (val credentials.Value, err error) {
	log.Printf("Looking up keyring for %s", p.credentialsName)
	item, err := p.keyring.Get(p.credentialsName)
	if err != nil {
		log.Println("Error from keyring", err)
		return val, err
	}
	if err = json.Unmarshal(item.Data, &val); err != nil {
		return val, fmt.Errorf("Invalid data in keyring: %v", err)
	}
	return val, err
}

func (p *KeyringProvider) Store(val credentials.Value) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return p.keyring.Set(keyring.Item{
		Key:   p.credentialsName,
		Label: fmt.Sprintf("aws-vault (%s)", p.credentialsName),
		Data:  bytes,

		// specific Keychain settings
		KeychainNotTrustApplication: true,
	})
}

func (p *KeyringProvider) Delete() error {
	return p.keyring.Remove(p.credentialsName)
}