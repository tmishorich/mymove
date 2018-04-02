package recordlocator

import (
	"fmt"
	"testing"
)

func TestCheckDigitAdder(t *testing.T) {
	goodOnes := []string{
		"ABCDEG",
		"3H9V3H",
		"8TEA6K",
		"8H4KSB",
		"UDP74G",
		"HUALZM",
		"DHSQRG",
		"PSVH7C",
		"EBRXLB",
		"FLC5K3",
		"MY6ZT3",
		"PBPRAJ",
		"KAU9DE",
		"5YXWC6",
		"6B4Z4U",
		"DACQ9M",
		"8B9WN7",
		"ZT3XBN",
		"DBM5L8",
		"YS4M6H",
		"9S5FYS",
		"39VX9Y",
		"ULPYJA",
		"CBZ8A2",
		"35TRQH",
		"YCFH4D",
		"QXGV7U",
		"6RPGVD",
		"PSL2HH",
		"S5EW2Y",
		"HEUDWF",
		"A9WWKW",
		"HKU44V",
		"RZCUGP",
		"9YJZ2P",
		"GG9ASL",
		"C3HVZ9",
		"8BNHJ5",
		"HQSRU6",
		"PLXJ2V",
		"NHCA24",
		"P6HBGD",
		"A9URNA",
		"J5H2G7",
		"E5HFWD",
		"M7XDK6",
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
	if newguy != "ABCDEG" {
		t.Fail()
	}
}

func TestGenerateLocators(t *testing.T) {
	for i := 0; i < 100; i++ {
		newLocator := NewRecordLocator()
		if !CheckRecordLocator(newLocator) {
			t.Error("Created invalid locator: ", newLocator)
		}
		// fmt.Println(newLocator)
	}
	// t.Fail()

}

func TestValidator(t *testing.T) {

	goodOnes := []string{
		"ABCDEG",
		"abcdeg",
		"3H9V3H",
		"8TEA6K",
		"8H4KSB",
		"UDP74G",
		"HUALZM",
		"DHSQRG",
		"PSVH7C",
		"EBRXLB",
		"FLC5K3",
		"MY6ZT3",
		"PBPRAJ",
		"KAU9DE",
		"5YXWC6",
		"6B4Z4U",
		"DACQ9M",
		"8B9WN7",
		"ZT3XBN",
		"DBM5L8",
		"YS4M6H",
		"9S5FYS",
		"39VX9Y",
		"ULPYJA",
		"CBZ8A2",
		"35TRQH",
		"YCFH4D",
		"QXGV7U",
		"6RPGVD",
		"PSL2HH",
		"S5EW2Y",
		"HEUDWF",
		"A9WWKW",
		"HKU44V",
		"RZCUGP",
		"9YJZ2P",
		"GG9ASL",
		"C3HVZ9",
		"8BNHJ5",
		"HQSRU6",
		"PLXJ2V",
		"NHCA24",
		"P6HBGD",
		"A9URNA",
		"J5H2G7",
		"E5HFWD",
		"M7XDK6",
	}

	badOnes := []string{
		"ABCDEF",
		"3I9V3H",
		"8TAE6K",
		"84KSB",
		"GUDP74",
		"HUAALZM",
		"DHSIRG",
		"PSVI7C",
		"EBEXLB",
		"FL5CK3",
		"M6YZT3",
		"PRPBAJ",
		"KAU8DE",
		"5YXXC6",
		"6B4😈4U",
		"DAOQ9M",
		"8C9WN7",
		"ZI3XBN",
		"DBN5L8",
		"YS46MH",
		"9S5FSY",
		"39VXY9",
		"ULYJAP",
		"88Z8A2",
		"35TRQX",
		"YCFHXX",
		"QXXX7U",
		"MACRAE",
		"AAAAAA",
		"BBBBBB",
		"A9WWWK",
		"HKA44V",
		"RZIUGP",
		"9JZY2P",
		"GGGASL",
		"C3VZ9H",
		"8BNHJ6",
		"HQSRU5",
		"PUXJ2V",
		"KHCA24",
		"P6HBKD",
		"A9URNALLLL",
		"J5H2G7PLPL",
		"E5HFWDKKLOI",
		"M7XDK6P",
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
