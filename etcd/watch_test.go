package etcd

import (
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"testing"
)

var (
	ShowCompleteLimit int
	MediaUrl          string
)

func Test(t *testing.T) {
	address := []string{"10.4.61.147:31979"}
	t.Run("Watch:", func(t *testing.T) {
		c, err := NewWatch(address, "/micro/config/app")
		if err != nil {
			t.Error(err)
		}
		if err := c.Watch(&MediaUrl, "mediaUrl"); err != nil {
			t.Error(err)
		}
		t.Logf("MediaUrl: %s", MediaUrl)
	})
	t.Run("WatchCallback:", func(t *testing.T) {
		c, err := NewWatch(address, fmt.Sprintf("%s/%s", "/micro/config", "USER_COMMON"))
		if err != nil {
			t.Error(err)
		}
		if err = c.WatchCallback(initUserCommon); err != nil {
			t.Error(err)
		}
		t.Logf("ShowCompleteLimit: %d", ShowCompleteLimit)
	})
}

type UserCommon struct {
	ShowCompleteLimit int `json:"showCompleteLimit"`
}

func initUserCommon(d []byte) error {
	//fmt.Println(d) //[123 34 115 104 111 119 67 111 109 112 108 101 116 101 76 105 109 105 116 34 58 55 48 125]
	var UserCommon *UserCommon
	err := ffjson.Unmarshal(d, &UserCommon)
	if err != nil {
		return err
	}
	if UserCommon == nil {
		return nil
	}
	if UserCommon.ShowCompleteLimit > 0 {
		ShowCompleteLimit = UserCommon.ShowCompleteLimit
	}

	return nil
}
