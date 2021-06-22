package cryptoId

import (
	"encoding/json"
	"fmt"
	"testing"

	pb "github.com/lights-T/lib-go/crypto/cryptoId/pb"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	//id := 10
	type User struct {
		Id   int64  `json:"id,cv_id"`
		Name string `json:"name"`
	}

	t.Run("marshal", func(t *testing.T) {
		u := &User{
			Id:   10,
			Name: "Lucy",
		}
		by, err := MarshalJSON(u)
		fmt.Println("data", string(by), err)
	})
	t.Run("unmarshal", func(t *testing.T) {
		var u User
		s := `{"id":400040,"name":"Lucy"}`
		err1 := json.Unmarshal([]byte(s), &u)
		fmt.Println("err", err1, "u", u)
		err := UnmarshalJSON([]byte(s), &u)
		fmt.Println(u, "err", err)
	})

	u := pb.User{
		Id:           10,
		Nickname:     "Lily",
		FollowStatus: 1,
	}
	t.Run("pb Marshal", func(t *testing.T) {
		by, err := MarshalJSON(u)
		assert.NoError(t, err)
		fmt.Println("data", string(by))
	})

	t.Run("pb unmarshl", func(t *testing.T) {
		var u pb.User
		s := `{"id":400040,"nickname":"Lily","followStatus":1}`
		err1 := UnmarshalJSON([]byte(s), &u)
		assert.NoError(t, err1)
		assert.Equal(t, int64(10), u.Id)
		fmt.Println("u", u)
	})

	t.Run("array marshal", func(t *testing.T) {
		l := pb.Label{
			Id:      10,
			Name:    "Weather",
			PostIds: []int64{19, 20, 300},
		}
		by, err := MarshalJSON(&l)
		assert.NoError(t, err)
		t.Logf("res:%s", string(by))
	})
	t.Run("array nil marshal ", func(t *testing.T) {
		l := pb.Label{
			Id:      10,
			Name:    "Weather",
			PostIds: nil,
		}
		by, err := MarshalJSON(&l)
		assert.NoError(t, err)
		t.Logf("res:%s", string(by))
	})

	t.Run("array unmarshal", func(t *testing.T) {

		t.Run("unmarshal", func(t *testing.T) {
			s := `{"id":10,"name":"Weather","postIds":[400043,400020,400900]}`
			var l pb.Label
			err := UnmarshalJSON([]byte(s), &l)
			assert.NoError(t, err)
			t.Logf("postIds:%v", l.PostIds)
		})

		t.Run("null array", func(t *testing.T) {
			s := `{"postIds":null}`
			var l pb.Label
			err := UnmarshalJSON([]byte(s), &l)
			assert.NoError(t, err)
			t.Logf("postIds:%v", l.PostIds)
		})

		t.Run("empty elem array", func(t *testing.T) {
			s := `{"postIds":[]}`
			var l pb.Label
			err := UnmarshalJSON([]byte(s), &l)
			assert.NoError(t, err)
			t.Logf("postIds:%v", l.PostIds)
		})

		t.Run("empty array with space chart", func(t *testing.T) {
			s := `{"postIds":[  ]}`
			var l pb.Label
			err := UnmarshalJSON([]byte(s), &l)
			assert.NoError(t, err)
			t.Logf("postIds:%v", l.PostIds)
		})
	})
	t.Run("no covert id ", func(t *testing.T) {
		l := pb.Label{
			Id:   10,
			Name: "Weather",
		}
		by, err := MarshalJSON(&l)
		assert.NoError(t, err)
		fmt.Println("data", string(by))

	})

}
