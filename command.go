// +build !windows

package sysproxy

import (
	"fmt"
	"os/exec"
	"strings"
)

func command(name string, arg ...string) (string, error) {
	c := exec.Command(name, arg...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%q: %w: %q", strings.Join(append([]string{name}, arg...), " "), err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
