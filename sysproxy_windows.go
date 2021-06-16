package sysproxy

import (
	"bufio"
	"bytes"
	"io"
	"net/textproto"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

const (
	settingPath = `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`
)

func OnNoProxy(list []string) error {
	return set("ProxyOverride", "REG_SZ", strings.Join(list, ";"))
}

func OffNoProxy() error {
	return set("ProxyOverride", "REG_SZ", "")
}

func GetNoProxy() ([]string, error) {
	m, err := get("ProxyOverride")
	if err != nil {
		return nil, err
	}
	proxyOverride := m["ProxyOverride"]
	if proxyOverride == "" {
		return []string{}, nil
	}
	return strings.Split(proxyOverride, ";"), nil
}

func OnHTTPS(address string) error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	m["https"] = address
	return setAllProxy(m)
}

func OffHTTPS() error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	delete(m, "https")
	return setAllProxy(m)
}

func GetHTTPS() (string, error) {
	b, err := getProxy()
	if err != nil {
		return "", err
	}
	if !b {
		return "", nil
	}
	m, err := getAllProxy()
	if err != nil {
		return "", err
	}
	return m["https"], nil
}

func OnHTTP(address string) error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	m["http"] = address
	return setAllProxy(m)
}

func OffHTTP() error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	delete(m, "http")
	return setAllProxy(m)
}

func GetHTTP() (string, error) {
	b, err := getProxy()
	if err != nil {
		return "", err
	}
	if !b {
		return "", nil
	}
	m, err := getAllProxy()
	if err != nil {
		return "", err
	}
	return m["http"], nil
}

func setAllProxy(m map[string]string) error {
	list := make([]string, 0, len(m))
	for key, item := range m {
		if item == "" {
			continue
		}
		list = append(list, strings.Join([]string{key, item}, "="))
	}
	sort.Strings(list)
	err := set("ProxyServer", "REG_SZ", strings.Join(list, ";"))
	if err != nil {
		return err
	}
	return useProxy(len(list) != 0)
}

func getAllProxy() (map[string]string, error) {
	m, err := get("ProxyServer")
	if err != nil {
		return nil, err
	}
	list := strings.Split(m["ProxyServer"], ";")
	proxy := map[string]string{}
	for _, item := range list {
		n := strings.SplitN(item, "=", 2)
		if len(n) == 1 {
			proxy["http"] = item
			proxy["https"] = item
			proxy["ftp"] = item
			proxy["socks"] = item
			break
		}
		proxy[n[0]] = n[1]
	}
	return proxy, nil
}

func OnPAC(pac string) error {
	err := set("AutoConfigURL", "REG_SZ", pac)
	if err != nil {
		return err
	}
	err = useProxy(false)
	if err != nil {
		return err
	}
	return reloadWinProxy()
}

func OffPAC() error {
	err := del("AutoConfigURL")
	if err != nil {
		return err
	}
	return reloadWinProxy()
}

func GetPAC() (string, error) {
	m, err := get("AutoConfigURL", "ProxyEnable")
	if err != nil {
		return "", err
	}
	autoConfigURL := m["AutoConfigURL"]
	proxyEnable, _ := strconv.ParseInt(m["ProxyEnable"], 0, 0)

	if autoConfigURL != "" && proxyEnable == 0 {
		return autoConfigURL, nil
	}
	return "", nil
}

func set(key string, typ string, value string) error {
	_, err := command(`reg`, `add`, settingPath, `/v`, key, `/t`, typ, `/d`, value, `/f`)
	return err
}

func get(keys ...string) (map[string]string, error) {
	buf, err := command(`reg`, `query`, settingPath)
	if err != nil {
		return nil, err
	}
	return getFrom(buf, settingPath, keys...)
}

func del(key string) error {
	_, err := command(`reg`, `delete`, settingPath, `/v`, key, `/f`)
	return err
}

func strBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func useProxy(b bool) error {
	return set("ProxyEnable", "REG_DWORD", strBool(b))
}

func getProxy() (bool, error) {
	m, err := get("ProxyEnable", "REG_DWORD")
	if err != nil {
		return false, err
	}
	i, err := strconv.ParseInt(m["ProxyEnable"], 0, 0)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

func reloadWinProxy() error {
	h, err := syscall.LoadLibrary("wininet.dll")
	if err != nil {
		return err
	}
	f, err := syscall.GetProcAddress(h, "InternetSetOptionW")
	if err != nil {
		return err
	}
	ret, _, errno := syscall.Syscall6(uintptr(f), 4, 0, 39, 0, 0, 0, 0)
	if ret != 1 {
		return errno
	}
	ret, _, errno = syscall.Syscall6(uintptr(f), 4, 0, 37, 0, 0, 0, 0)
	if ret != 1 {
		return errno
	}
	return nil
}

func getFrom(data string, path string, keys ...string) (map[string]string, error) {
	m := map[string]string{}
	index := strings.Index(data, path)
	if index == -1 {
		return m, nil
	}
	data = data[index+len(path):]
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(data)))
	reader.ReadLine()
	for len(m) != len(keys) {
		row, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if row == "" {
			break
		}
		row = strings.TrimSpace(row)
		s := strings.SplitN(row, "    ", 3)
		key := s[0]
		skip := true
		for _, k := range keys {
			if k == key {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		val := ""
		if len(s) == 3 {
			val = s[2]
		}
		m[key] = val
	}
	return m, nil
}
