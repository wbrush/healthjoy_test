//  +build unit

package configuration

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	cfg := InitConfig("", "")
	if cfg.IsLoaded() != true {
		t.Error("loading failed to be set correctly")
	}
	if cfg.Environment != "local" {
		t.Errorf("error initing development space expected %s got %s", "local", cfg.Environment)
	}
}
