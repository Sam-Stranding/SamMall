package rpc

import "testing"

func TestParseLarkUserInfoResponse(t *testing.T) {
	body := []byte(`{"code":0,"msg":"success","data":{"open_id":"ou_test","name":"Sam"}}`)

	userInfo, err := parseLarkUserInfoResponse(body)
	if err != nil {
		t.Fatalf("parseLarkUserInfoResponse returned error: %v", err)
	}
	if userInfo.OpenID != "ou_test" {
		t.Fatalf("expected open_id ou_test, got %q", userInfo.OpenID)
	}
	if userInfo.Name != "Sam" {
		t.Fatalf("expected name Sam, got %q", userInfo.Name)
	}
}
