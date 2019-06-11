package commands

import (
	"fmt"

	"os"

	"reflect"

	"code.cloudfoundry.org/credhub-cli/errors"
	"code.cloudfoundry.org/credhub-cli/models"
)

type ImportCommand struct {
	File string `short:"f" long:"file" description:"File containing credentials to import" required:"true"`
	ClientCommand
}

func (c *ImportCommand) Execute([]string) error {
	var bulkImport models.CredentialBulkImport
	err := bulkImport.ReadFile(c.File)

	if err != nil {
		return err
	}

	err = c.setCredentials(bulkImport)

	return err
}

func (c *ImportCommand) setCredentials(bulkImport models.CredentialBulkImport) error {
	var (
		name       string
		successful int
		failed     int
	)
	importErrors := make([]string, 0)

	for i, credential := range bulkImport.Credentials {
		switch credentialName := credential["name"].(type) {
		case string:
			name = credentialName
		default:
			name = ""
		}

		switch credential["type"].(string) {
		case "ssh":
			if _, ok := credential["value"].(map[string]interface{})["public_key_fingerprint"]; ok {
				delete(credential["value"].(map[string]interface{}), "public_key_fingerprint")
			}
		case "user":
			if _, ok := credential["value"].(map[string]interface{})["password_hash"]; ok {
				delete(credential["value"].(map[string]interface{}), "password_hash")
			}
		}

		_, err := c.client.SetCredential(name, credential["type"].(string), credential["value"])

		if err != nil {
			if isAuthenticationError(err) {
				return err
			}
			failure := fmt.Sprintf("Credential '%s' at index %d could not be set: %v", name, i, err)
			fmt.Println(failure + "\n")
			importErrors = append(importErrors, " - "+failure)
			failed++
			continue
		} else {
			successful++
		}
	}

	fmt.Println("Import complete.")
	fmt.Fprintf(os.Stdout, "Successfully set: %d\n", successful)
	fmt.Fprintf(os.Stdout, "Failed to set: %d\n", failed)
	for _, v := range importErrors {
		fmt.Println(v)
	}

	if failed > 0 {
		return errors.NewFailedToImportError()
	}

	return nil
}

func isAuthenticationError(err error) bool {
	return reflect.DeepEqual(err, errors.NewNoApiUrlSetError()) ||
		reflect.DeepEqual(err, errors.NewRevokedTokenError()) ||
		reflect.DeepEqual(err, errors.NewRefreshError())
}
