package vmess

import (
	"io"
	"os"
	"testing"
)

func TestDialer(t *testing.T) {
	link := os.Getenv("link")

	client := NewClient(link)

	resp, err := client.Get("https://ipinfo.io/")
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)

}
