package credentials

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const (
	CredentialsFileContentTemplate = `# smtp:
#   user_name: my-smtp-user
#   password: my-smtp-password
#
# aws:
#   access_key_id: 123
#   secret_access_key: 345

# Used as the base secret for all MessageVerifiers in Rails, including the one protecting cookies.
secret_key_base: %s
`
)

func NewCredentialsFileContent() (string, error) {
	// generate secret_key_base
	// https://github.com/rails/rails/blob/04df9bc3d120b51447bde54caa56e9237cb8da0e/railties/lib/rails/generators/rails/credentials/credentials_generator.rb#L42
	r := make([]byte, 32)
	_, err := rand.Read(r)
	if err != nil {
		return "", fmt.Errorf("unable to generate randomness: %w", err)
	}

	return fmt.Sprintf(CredentialsFileContentTemplate, hex.EncodeToString(r)), nil
}
