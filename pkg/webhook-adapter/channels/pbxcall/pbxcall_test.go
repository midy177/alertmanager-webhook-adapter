package pbxcall

import "testing"

func TestName(t *testing.T) {
	s := NewSender("pbxcall-test.yeastardigital.com:8181", "http", "smartpbxsmartpbx", "6866", "6866", []string{"L0-alarm-zh-female"})
	s.Send(nil)
}
