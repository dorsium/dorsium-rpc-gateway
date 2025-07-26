package utils

import "regexp"

var (
	hexRegex    = regexp.MustCompile(`^(0x)?[0-9a-fA-F]{40}$`)
	bech32Regex = regexp.MustCompile(`^[a-z0-9]{1,83}1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}$`)
)

// IsValidAddress validates bech32 or hex formatted addresses.
func IsValidAddress(addr string) bool {
	return hexRegex.MatchString(addr) || bech32Regex.MatchString(addr)
}
