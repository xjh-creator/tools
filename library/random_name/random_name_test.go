package random_name

import "testing"

func TestGetFullName(t *testing.T) {

	for i := 0; i <= 99; i++ {
		fn := GetFullName()
		t.Log("full_name:", fn)
	}
}
