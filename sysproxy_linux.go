package sysproxy

import (
	"bytes"
	"net"
	"strings"
)

const (
	scheme = "org.gnome.system.proxy"
)

func OnNoProxy(list []string) error {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("[ ")
	for i, item := range list {
		if item == "" {
			continue
		}
		buf.WriteByte('\'')
		buf.WriteString(item)
		buf.WriteByte('\'')
		if len(list)-1 != i {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" ]")
	return set("", "ignore-hosts", buf.String())
}

func OffNoProxy() error {
	return set("", "ignore-hosts", "[]")
}

func GetNoProxy() ([]string, error) {
	data, err := get("", "ignore-hosts")
	if err != nil {
		return nil, err
	}
	data = strings.TrimPrefix(data, "@as")
	data = strings.TrimSpace(data)
	data = strings.TrimPrefix(data, "[")
	data = strings.TrimSuffix(data, "]")
	data = strings.TrimSpace(data)
	if data == "" {
		return []string{}, nil
	}
	list := strings.Split(data, ",")
	for i := range list {
		item := list[i]
		item = strings.TrimSpace(item)
		item = strings.Trim(item, "'")
		list[i] = item
	}
	return list, nil
}

func OnHTTPS(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	err = set("https", "host", host)
	if err != nil {
		return err
	}
	err = set("https", "port", port)
	if err != nil {
		return err
	}
	err = set("", "mode", "manual")
	if err != nil {
		return err
	}
	return nil
}

func OffHTTPS() error {
	err := reset("", "mode")
	if err != nil {
		return err
	}
	err = reset("https", "host")
	if err != nil {
		return err
	}
	err = reset("https", "port")
	if err != nil {
		return err
	}
	return nil
}

func GetHTTPS() (string, error) {
	mode, err := get("", "mode")
	if err != nil {
		return "", err
	}
	mode = strings.Trim(mode, "'")
	if mode != "manual" {
		return "", nil
	}
	host, err := get("https", "host")
	if err != nil {
		return "", err
	}
	port, err := get("https", "port")
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(host, port), nil
}

func OnHTTP(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	err = set("http", "host", host)
	if err != nil {
		return err
	}
	err = set("http", "port", port)
	if err != nil {
		return err
	}
	err = set("", "mode", "manual")
	if err != nil {
		return err
	}
	return nil
}

func OffHTTP() error {
	err := reset("", "mode")
	if err != nil {
		return err
	}
	err = reset("http", "host")
	if err != nil {
		return err
	}
	err = reset("http", "port")
	if err != nil {
		return err
	}
	return nil
}

func GetHTTP() (string, error) {
	mode, err := get("", "mode")
	if err != nil {
		return "", err
	}
	if mode != "manual" {
		return "", nil
	}
	host, err := get("http", "host")
	if err != nil {
		return "", err
	}
	port, err := get("http", "port")
	if err != nil {
		return "", err
	}
	return net.JoinHostPort(host, port), nil
}

func OnPAC(pac string) error {
	err := set("", "autoconfig-url", pac)
	if err != nil {
		return err
	}
	err = set("", "mode", "auto")
	if err != nil {
		return err
	}
	return nil
}

func OffPAC() error {
	err := reset("", "mode")
	if err != nil {
		return err
	}
	err = reset("", "autoconfig-url")
	if err != nil {
		return err
	}
	return nil
}

func GetPAC() (string, error) {
	mode, err := get("", "mode")
	if err != nil {
		return "", err
	}
	mode = strings.Trim(mode, "'")
	if mode != "auto" {
		return "", err
	}
	pac, err := get("", "autoconfig-url")
	if err != nil {
		return "", err
	}
	return pac, nil
}

func reset(sub, key string) error {
	scheme := scheme
	if sub != "" {
		scheme = scheme + "." + sub
	}
	_, err := command("gsettings", "reset", scheme, key)
	return err
}

func get(sub, key string) (string, error) {
	scheme := scheme
	if sub != "" {
		scheme = scheme + "." + sub
	}
	out, err := command("gsettings", "get", scheme, key)
	if err != nil {
		return "", err
	}
	out = strings.Trim(out, "'")
	return out, nil
}

func set(sub, key string, val string) error {
	scheme := scheme
	if sub != "" {
		scheme = scheme + "." + sub
	}
	_, err := command("gsettings", "set", scheme, key, val)
	return err
}
