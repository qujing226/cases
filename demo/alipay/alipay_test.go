package alipay

import "testing"

func TestPay(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			Pay()
		})
	}
}
