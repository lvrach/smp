package prompt

import (
	"fmt"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lvrach/smp/internal/config"
	"github.com/lvrach/smp/internal/state"
	"github.com/lvrach/smp/keystore"
)

// PromptEnvironmentVariables prompts the user for environment variable values
func PromptEnvironmentVariables(mcpConfig *config.MCPConfig, mcpState *state.MCPServer) error {
	// TODO decouple this from the mcpConfig and mcpState

	// Ask about keychain storage on macOS if there are any secrets
	var useKeychain bool
	hasSecrets := false
	for _, envVar := range mcpConfig.EnvironmentVars {
		if envVar.Type == "secret" {
			hasSecrets = true
			break
		}
	}

	if runtime.GOOS == "darwin" && hasSecrets {
		prompt := &survey.Confirm{
			Message: "Do you want to store secrets in the macOS keychain?",
			Default: true,
		}
		if err := survey.AskOne(prompt, &useKeychain); err != nil {
			return fmt.Errorf("failed to get keychain preference: %w", err)
		}
	}

	for _, envVar := range mcpConfig.EnvironmentVars {
		// Skip if already set
		if _, exists := mcpState.GetEnvironmentVariable(envVar.Name); exists {
			continue
		}

		var value string
		var prompt survey.Prompt

		// Create appropriate prompt based on variable type
		switch envVar.Type {
		case "secret":
			prompt = &survey.Password{
				Message: fmt.Sprintf("Enter %s (%s):", envVar.Name, envVar.Description),
			}
		default:
			prompt = &survey.Input{
				Message: fmt.Sprintf("Enter %s (%s):", envVar.Name, envVar.Description),
			}
		}

		// Show the prompt
		if err := survey.AskOne(prompt, &value, survey.WithValidator(survey.Required)); err != nil {
			return fmt.Errorf("failed to get value for %s: %w", envVar.Name, err)
		}

		if envVar.Type == "secret" && useKeychain {
			kc := keystore.KeyChain{}
			accountKey := mcpConfig.Name + "_" + envVar.Name
			kc.Store(keystore.AccountType(accountKey), value)

			if mcpState.KeyChainEnvVars == nil {
				mcpState.KeyChainEnvVars = make(map[string]string)
			}
			mcpState.KeyChainEnvVars[accountKey] = envVar.Name
		} else {
			mcpState.SetEnvironmentVariable(envVar.Name, value)
		}
	}

	return nil
}

// MultiSelect prompts the user to select multiple options from a list
func MultiSelect(message string, options []string, defaultSelections map[string]bool) ([]string, error) {
	// Convert defaultSelections to a slice of indices
	var defaultIndices []int
	for i, option := range options {
		if defaultSelections[option] {
			defaultIndices = append(defaultIndices, i)
		}
	}

	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Default: defaultIndices,
	}

	if err := survey.AskOne(prompt, &selectedIndices); err != nil {
		return nil, fmt.Errorf("failed to get selection: %w", err)
	}

	// Convert selected indices back to options
	selected := make([]string, len(selectedIndices))
	for i, idx := range selectedIndices {
		selected[i] = options[idx]
	}

	return selected, nil
}
