package sysproxy

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/textproto"
	"regexp"
	"strings"
)

func OnNoProxy(list []string) error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	return set("proxybypassdomains", s, list...)
}

func OffNoProxy() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	return set("proxybypassdomains", s, "")
}

func GetNoProxy() ([]string, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return nil, err
	}
	m, err := get("proxybypassdomains", s)
	if err != nil {
		return nil, err
	}
	m = strings.TrimSpace(m)
	list := strings.Split(m, "\n")
	if len(list) != 0 && list[len(list)-1] == "" {
		list = list[:len(list)-1]
	}
	return list, nil
}

func OnHTTPS(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("securewebproxy", s, host, port)
	if err != nil {
		return err
	}
	err = set("securewebproxystate", s, "on")
	if err != nil {
		return err
	}
	return nil
}

func OffHTTPS() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("securewebproxystate", s, "off")
	if err != nil {
		return err
	}
	return nil
}

func GetHTTPS() (string, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return "", err
	}
	buf, err := get("securewebproxy", s)
	if err != nil {
		return "", err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	header, err := reader.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return "", err
	}
	if header.Get("Enabled") == "Yes" {
		return net.JoinHostPort(header.Get("Server"), header.Get("Port")), nil
	}
	return "", nil
}

func OnHTTP(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("webproxy", s, host, port)
	if err != nil {
		return err
	}
	err = set("webproxystate", s, "on")
	if err != nil {
		return err
	}
	return nil
}

func OffHTTP() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("webproxystate", s, "off")
	if err != nil {
		return err
	}
	return nil
}

func GetHTTP() (string, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return "", err
	}
	buf, err := get("webproxy", s)
	if err != nil {
		return "", err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	header, err := reader.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return "", err
	}
	if header.Get("Enabled") == "Yes" {
		return net.JoinHostPort(header.Get("Server"), header.Get("Port")), nil
	}
	return "", nil
}

func OnPAC(pac string) error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("autoproxyurl", s, pac)
	if err != nil {
		return err
	}
	err = set("autoproxystate", s, "on")
	if err != nil {
		return err
	}
	return nil
}

func OffPAC() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("autoproxystate", s, "off")
	if err != nil {
		return err
	}
	return nil
}

func GetPAC() (string, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return "", err
	}
	buf, err := get("autoproxyurl", s)
	if err != nil {
		return "", err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	header, err := reader.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return "", err
	}
	if header.Get("Enabled") == "Yes" {
		return header.Get("URL"), nil
	}
	return "", nil
}

func set(key string, inter string, value ...string) error {
	_, err := command("networksetup", append([]string{"-set" + key, inter}, value...)...)
	return err
}

func get(key string, inter string) (string, error) {
	return command("networksetup", "-get"+key, inter)
}

func getNetworkInterface() (string, error) {
	buf, err := command("sh", "-c", "networksetup -listnetworkserviceorder | grep -B 1 $(route -n get default | grep interface | awk '{print $2}')")
	if err != nil {
		return "", err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	reg := regexp.MustCompile(`^\(\d+\)\s(.*)$`)
	device, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	match := reg.FindStringSubmatch(device)
	if len(match) <= 1 {
		return "", fmt.Errorf("unable to get network interface")
	}
	return match[1], nil
}
