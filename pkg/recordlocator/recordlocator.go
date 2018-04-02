package recordlocator

import (
	"crypto/md5"
	"math/rand"
	"strings"
)

// This package creates and validates Record Locators.
// A Record Locator is a six character alphanumeric string.
// The sixth character is a check digit, if any single digit and/or transposition errors are made
// the checksum will be off and it will be known that the Record Locator has been miscommunicated.

// the characters 0/O, and 1/I have been removed to prevent common reading errors
// this leaves 32 possible characters, with 5 non-check digits,
// so there are 32 ^ 5 == 33,500,000 valid Record Locators using this scheme.

// Our checksum scheme is simple. We take an md5 of the first 5 digits, this returns
// an array of bytes, we take the first byte mod 32 and that's the check digit.

// DANGER WILL ROBINSON
// If any of these constants are changed, all extant record locators will become invalid!
const locatorCharacterSet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
const finalLocatorLength = 6

// end DANGER

func checksumForPrefix(locator string) string {
	checksum := md5.Sum([]byte(locator))
	characterIndex := checksum[0] % 32 // it's a 16 byte number, we just mod the first byte
	checkDigit := string(locatorCharacterSet[characterIndex])
	return checkDigit
}

// addCheckDigit adds a check digit to the five digit prefix
func addCheckDigit(prefix string) string {
	checkDigit := checksumForPrefix(prefix)

	return prefix + checkDigit
}

// NewRecordLocator generates a random record locator (NOT guaranteed to be unique in the db.)
func NewRecordLocator() string {
	prefix := ""
	for i := 0; i < finalLocatorLength-1; i++ {
		newChar := string(locatorCharacterSet[rand.Intn(len(locatorCharacterSet))])
		prefix = prefix + newChar
	}

	return addCheckDigit(prefix)
}

// CheckRecordLocator returns true if the checksum is correct for the given string.
// If this returns false, it means that the record locator was not communicated correctly.
func CheckRecordLocator(locator string) bool {
	normalizedLocator := strings.ToUpper(locator)

	if len(normalizedLocator) != finalLocatorLength {
		return false
	}

	// Make sure all the characters are in the valid character set
	for _, char := range normalizedLocator {
		digit := string(char)
		if strings.Index(locatorCharacterSet, digit) == -1 {
			return false
		}
	}

	prefix := string(normalizedLocator[0 : finalLocatorLength-1])
	checkDigit := string(normalizedLocator[finalLocatorLength-1])

	checksum := checksumForPrefix(prefix)
	return checksum == checkDigit
}
