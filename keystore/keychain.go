package keystore

import (
	"fmt"

	"github.com/keybase/go-keychain"
)

const (
	service = "Secure MCP"
)

type AccountType string

type KeyChain struct {
}

func (k *KeyChain) Store(accountKey AccountType, secret string) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(service)
	item.SetAccount(string(accountKey))
	item.SetData([]byte(secret))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	err := keychain.AddItem(item)
	if err == keychain.ErrorDuplicateItem {
		return k.update(accountKey, secret)
	}
	if err != nil {
		return fmt.Errorf("storing API key: %w", err)
	}
	return nil
}

func (*KeyChain) update(accountKey AccountType, secret string) error {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(string(accountKey))

	update := keychain.NewItem()
	update.SetData([]byte(secret))

	err := keychain.UpdateItem(query, update)
	if err != nil {
		return fmt.Errorf("failed to update API key: %w", err)
	}
	fmt.Println("API key updated successfully")
	return nil
}

func (*KeyChain) Retrieve(accountKey AccountType) (string, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(string(accountKey))
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return "", fmt.Errorf("retrieve API key: %w", err)
	}
	if len(results) == 0 {
		return "", fmt.Errorf("key not found")
	}
	return string(results[0].Data), nil
}
