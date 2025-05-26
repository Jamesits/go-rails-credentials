package credentials

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCredPairs = []struct {
	MasterKey     string
	EncryptedData string
	PlainTextData string
}{
	{
		MasterKey:     "a2683380db86af7597f33561b5f11755",
		EncryptedData: "9IOTUFFLHgvM0u0i0sbzTqEn7W3PW45+bkQ6Ce+TSy+JZnjXHFibrKsCL6FRVQT11h9YoqQ4lI1QlEMrrickp3EtOZaFkMPZM9NCv4O+UtTawBmav11dGgYG4Ye+XGrROoiG/8xFalGo6312N+sQwk4BGgD6DANHLBm/RL3jDaR8wSg1IzrlfwSR7E8PMKFEb4cZSUr1OzgV00M6kBrfIK4KX4nLjS6xYyHEoUdJ1zmaE8dw5ppAr41QWQsPkyD5lL/0jjsNUit0J//g/lc5MMpmHePfB9wbwXy+ZTECq+aj6mhYDK7p2l2Z/iT8kpA7HvvqRg2t+qofW6GZR8rnR0HDmH4inWTzaMMoaF4WbVZarXXu8a2w0nkthwzHHbrrNvJHQB0Co27F/vsCyokVKj6rb0BdwJap5TdhThPfwM2zh5s4QRu0YgA+tGUM6cxzWa3a1iqaUqE48Aqd2y+3fLJ6TsMI0rtoB+q0aaUwjSspyeGe8F6c3LFn--c3FaoQhEm/BmM4dM--+Y8philyFTu9MsWmDWPASA==",
		PlainTextData: `# smtp:
#   user_name: my-smtp-user
#   password: my-smtp-password
#
# aws:
#   access_key_id: 123
#   secret_access_key: 345

# Used as the base secret for all MessageVerifiers in Rails, including the one protecting cookies.
secret_key_base: 63ff3f018cc8bf80c23cee7367342f1701462d916b3cbb81144e495bac51c9c22c3af28f7d8af68567b56c21554874c510ab646ddfbbd1da398edb8ceeaf1964
`,
	},
}

func TestDecryption(t *testing.T) {
	for _, p := range testCredPairs {
		dec, err := Decrypt(p.MasterKey, p.EncryptedData)
		assert.NoError(t, err)

		des, err := UnmarshalSingleString(dec)
		assert.NoError(t, err)

		assert.Equal(t, p.PlainTextData, des)
	}
}

func TestLoop(t *testing.T) {
	for _, p := range testCredPairs {
		mk, err := RandomMasterKey()
		assert.NoError(t, err)

		ser, err := MarshalSingleString(p.PlainTextData)
		assert.NoError(t, err)

		enc, err := Encrypt(mk, ser)
		assert.NoError(t, err)

		dec, err := Decrypt(mk, enc)
		assert.NoError(t, err)

		des, err := UnmarshalSingleString(dec)
		assert.NoError(t, err)

		assert.Equal(t, p.PlainTextData, des)
	}
}
