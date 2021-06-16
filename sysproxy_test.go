package sysproxy

import (
	"reflect"
	"testing"
)

func TestNoProxy(t *testing.T) {
	reset, err := GetNoProxy()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if len(reset) == 0 {
			OffNoProxy()
		} else {
			OnNoProxy(reset)
		}
	})

	list := []string{"a", "b", "c"}
	err = OnNoProxy(list)
	if err != nil {
		t.Fatal(err)
	}
	got, err := GetNoProxy()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(list, got) {
		t.Fatalf("want %q, got %q", list, got)
	}

	err = OffNoProxy()
	if err != nil {
		t.Fatal(err)
	}
	empty, err := GetNoProxy()
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("want empty, got %q", empty)
	}
}

func TestPAC(t *testing.T) {
	reset, err := GetPAC()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if reset == "" {
			OffPAC()
		} else {
			OnPAC(reset)
		}
	})

	pac := "http://127.0.0.1:1080"
	err = OnPAC(pac)
	if err != nil {
		t.Fatal(err)
	}
	got, err := GetPAC()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(pac, got) {
		t.Fatalf("want %q, got %q", pac, got)
	}

	err = OffPAC()
	if err != nil {
		t.Fatal(err)
	}
	ori, err := GetPAC()
	if err != nil {
		t.Fatal(err)
	}
	if ori != "" {
		t.Fatalf("want empty, got %q", ori)
	}
}

func TestHTTP(t *testing.T) {
	reset, err := GetHTTP()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if reset == "" {
			OffHTTP()
		} else {
			OnHTTP(reset)
		}
	})

	address := "127.0.0.1:1080"
	err = OnHTTP(address)
	if err != nil {
		t.Fatal(err)
	}
	got, err := GetHTTP()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(address, got) {
		t.Fatalf("want %q, got %q", address, got)
	}

	err = OffHTTP()
	if err != nil {
		t.Fatal(err)
	}
	ori, err := GetHTTP()
	if err != nil {
		t.Fatal(err)
	}
	if ori != "" {
		t.Fatalf("want empty, got %q", ori)
	}
}

func TestHTTPS(t *testing.T) {
	reset, err := GetHTTPS()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if reset == "" {
			OffHTTPS()
		} else {
			OnHTTPS(reset)
		}
	})

	address := "127.0.0.1:1080"
	err = OnHTTPS(address)
	if err != nil {
		t.Fatal(err)
	}
	got, err := GetHTTPS()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(address, got) {
		t.Fatalf("want %q, got %q", address, got)
	}

	err = OffHTTPS()
	if err != nil {
		t.Fatal(err)
	}
	ori, err := GetHTTPS()
	if err != nil {
		t.Fatal(err)
	}
	if ori != "" {
		t.Fatalf("want empty, got %q", ori)
	}
}
