// +build !windows,!darwin,!linux

package sysproxy

func OnNoProxy(list []string) error {
	return nil
}

func OffNoProxy() error {
	return nil
}

func GetNoProxy() ([]string, error) {
	return nil, nil
}

func OnHTTPS(address string) error {
	return nil
}

func OffHTTPS() error {
	return nil
}

func GetHTTPS() (string, error) {
	return "", nil
}

func OnHTTP(address string) error {
	return nil
}

func OffHTTP() error {
	return nil
}

func GetHTTP() (string, error) {
	return "", nil
}

func OnPAC(pac string) error {
	return nil
}

func OffPAC() error {
	return nil
}

func GetPAC() (string, error) {
	return "", nil
}
