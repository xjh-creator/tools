package util

import (
	"fmt"
	"testing"
)

func TestHideStar(t *testing.T) {
	tel := "12345671234"
	tel_hide := HideStar(tel)
	t.Logf("src: %s , hide: %s", tel, tel_hide)
}

func TestTelephoneHide(t *testing.T) {
	tel := "12345671234"
	tel_hide := TelephoneHide(tel)
	t.Logf("src: %s , hide: %s", tel, tel_hide)
}

func TestTelephoneHideMiddle(t *testing.T) {
	tel := "17817447510"
	tel_hide := TelephoneHideMiddle(tel)
	fmt.Printf("src: %s , hide: %s", tel, tel_hide)
}
