package recordlocator

import (
	"fmt"
	"testing"
)

func TestCheckDigitAdder(t *testing.T) {
	goodOnes := []string{
		"3URJK9",
		"K9KRZZ",
		"ZFYMVN",
		"AJA78P",
		"E94UBL",
		"Y5445U",
		"UD65BT",
		"G9M4RM",
		"JEY76F",
		"E534PK",
		"C9QBAN",
		"9FUCBV",
		"Z72SBD",
		"7E7UEN",
		"HN4XVR",
		"8Z362E",
		"RTGNCY",
	}

	for _, goodOne := range goodOnes {
		prefix := string(goodOne[0:5])
		stillGood := addCheckDigit(prefix)
		if stillGood != goodOne {
			t.Error("These should always be the same!", goodOne, stillGood)
		}
	}

	newguy := addCheckDigit("ABCDE")
	fmt.Println(newguy)
	if newguy != "ABCDEQ" {
		t.Fail()
	}
}

func TestGenerateLocators(t *testing.T) {
	for i := 0; i < 100; i++ {
		newLocator := NewRecordLocator()
		if !CheckRecordLocator(newLocator) {
			t.Error("Created invalid locator: ", newLocator)
		}
	}

}

func TestValidator(t *testing.T) {

	goodOnes := []string{
		"ABCDEQ",
		"abcdeq",
		"3H9V3M",
		"8TEA6P",
		"8H4KSE",
		"UDP74F",
		"HUALZ4",
		"DHSQR7",
		"PSVH7C",
		"EBRXL5",
		"FLC5K9",
		"MY6ZTZ",
		"M7XDKF",
		"K8NCAT",
		"S64PX8",
		"E69UAF",
		"EHXFWP",
		"AHXXBC",
		"6F7H9W",
		"Q28W2S",
		"NPVV7N",
		"8YP82Z",
		"C2ADGT",
		"7NF8JE",
		"AWG5PW",
		"93URJK",
		"K9KRZZ",
		"ZFYMVN",
	}

	badOnes := []string{
		"ABDCEQ",
		"3H8V3M",
		"8TE6P",
		"8I4KSE",
		"UDPA4F",
		"HLAUZ4",
		"D😈SQR7",
		"PSVH7CA",
		"EBRXM5",
		"AAAAAA",
		"MACRAE",
		"7MXDKF",
		"DOGCAT",
		"S64XP8",
		"F69UAF",
		"EHRFWP",
		"AHXBXC",
		"W6F7H9",
		"O28W2S",
		"NYVV7N",
		"8YP81Z",
		"C2AOGT",
		"7NF8EE",
		"AWN6PW",
		"93EOJK",
		"KXKRZZ",
		"ZFPMVN",
		"BBBBBB",
	}

	for _, goodOne := range goodOnes {
		if !(CheckRecordLocator(goodOne)) {
			t.Error("This is supposed to be a good one!", goodOne)
		}
	}

	for _, badOne := range badOnes {
		if CheckRecordLocator(badOne) {
			t.Error("This is supposed to be a bad one!", badOne)
		}
	}
}
