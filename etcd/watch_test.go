package etcd

import (
	"testing"
)

func TestConf(t *testing.T) {
	c, err := New([]string{"10.4.60.128:2379"}, "/micro/config/app")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get", func(t *testing.T) {
		var url string
		if err := c.Get("mediaUrl").Scan(&url); err != nil {
			t.Fatal(err)
		}
		t.Log(url)
	})

	t.Run("watch", func(t *testing.T) {
		var url string
		if err := c.Watch(&url, "mediaUrl"); err != nil {
			t.Fatal(err)
		}
		t.Log(url)
	})
}
