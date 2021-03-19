package uid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUid(t *testing.T) {
	id := 10
	type User struct {
		Id   UID    `json:"id"`
		Name string `json:"name"`
	}
	u := &User{
		Id:   UID(id),
		Name: "Lucy",
	}
	t.Run("encode", func(t *testing.T) {
		by, err := json.Marshal(u)
		assert.NoError(t, err)
		assert.Contains(t, string(by), "400040")
	})

	t.Run("decode", func(t *testing.T) {
		s := `{"id":400040,"name":"Lucy"}`
		var u User
		err := json.Unmarshal([]byte(s), &u)
		assert.NoError(t, err)
		assert.Equal(t, UID(10), u.Id)
	})
}
