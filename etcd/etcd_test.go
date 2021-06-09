package etcd

import (
	"testing"
)

func Test(t *testing.T) {
	s, err := NewKV([]string{"10.4.61.145:31169"})
	if err != nil {
		t.Fatal(err)
	}

	type MsgLimitConfig struct {
		FollowLimit       int32 `json:"followLimit"`
		PostPraiseLimit   int32 `json:"postPraiseLimit"`
		PostCommentLimit  int32 `json:"postCommentLimit"`
		PostFavoriteLimit int32 `json:"postFavoriteLimit"`
		NotifyTimeLimit   int32 `json:"notifyTimeLimit"`
	}
	no := &MsgLimitConfig{
		FollowLimit:       2,
		PostPraiseLimit:   2,
		PostCommentLimit:  2,
		PostFavoriteLimit: 2,
		NotifyTimeLimit:   30,
	}

	t.Run("setPathObj", func(t *testing.T) {
		if err = s.SetKV("/micro/config/test", no); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("setKey", func(t *testing.T) {
		if err = s.SetKV("/micro/config/test", 100, "followLimit"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("getPathObj", func(t *testing.T) {
		v, err := s.GetKV("/micro/config/test")
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("===>:%s", string(v))
	})

	t.Run("getKey", func(t *testing.T) {
		v, err := s.GetKV("/micro/config/test", "followLimit1")
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("===>:%s", string(v))
	})

	t.Run("delPath", func(t *testing.T) {
		if err = s.DelKV("/micro/config/test1"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("delKey", func(t *testing.T) {
		if err = s.DelKV("/micro/config/test1", "followLimit"); err != nil {
			t.Fatal(err)
		}
	})
}
