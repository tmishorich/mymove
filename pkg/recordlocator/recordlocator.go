package recordlocator

import (
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

// Our checksum scheme is simple. Given a record locator with six digits d1, d2, d3...
// The checksum is valid if and only if:
//        d1 * 11 + d2 * 9 + d3 * 7 + d4 * 5 + d5 * 3 + d6 % 32 == 0

// It is important that the coefficients in this equation are all relatively prime to 32.
// I picked ascending prime numbers (leaving out 2) for simplicity.

// DANGER WILL ROBINSON
// If any of these constants are changed, all extant record locators will become invalid!
const locatorCharacterSet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
const finalLocatorLength = 6

var relativelyPrimeCoefficients = [...]int{11, 9, 7, 5, 3, 1}

// end DANGER

// checksumForString returns the checksum for this string. It returns the checksum so far,
// so it is safe to use this on both 5 and 6 character strings.
// It expects the strings to only contain valid characters, though.
func checksumForString(locator string) int {
	checksum := 0
	for i := 0; i < len(locator); i++ {
		digit := string(locator[i])
		intValue := strings.Index(locatorCharacterSet, digit)
		checksum += intValue * relativelyPrimeCoefficients[i]
	}
	return checksum % len(locatorCharacterSet)
}

// addCheckDigit adds a check digit to the five digit prefix
func addCheckDigit(prefix string) string {
	incompleteChecksum := checksumForString(prefix)
	inverse := len(locatorCharacterSet) - ((incompleteChecksum) % len(locatorCharacterSet))
	checkDigit := string(locatorCharacterSet[inverse])

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

	checksum := checksumForString(strings.ToUpper(normalizedLocator))
	return checksum == 0
}
