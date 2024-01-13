package internal

import (
	"fmt"
	"testing"
)

func TestConfigLoading(t *testing.T) {
	emailConfig := GetConfigs("../../configs.yaml")
	fmt.Println(emailConfig.Email)
}
