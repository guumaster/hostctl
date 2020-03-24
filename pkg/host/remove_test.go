package host

import (
	"testing"
)

func TestRemoveFromProfile(t *testing.T) {
	f := func(lines, remove, want []string) {
		t.Helper()
		data := removeFromProfile(lines, remove)
		if len(data) != len(want) {
			t.Fatalf("bad count of records; got %d; want %d", len(data), len(want))
		}
		for _, d := range data {
			found := len(want) == 0
			for _, w := range want {
				if d == w {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("unexpected record; got %s", d)
			}
		}
		for _, w := range want {
			found := len(data) == 0
			for _, d := range data {
				if d == w {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("unexpected record; want %s", w)
			}
		}
	}
	f([]string{},
		[]string{},
		[]string{})
	f([]string{"127.0.0.1 test1.loc"},
		[]string{},
		[]string{"127.0.0.1 test1.loc"})
	f([]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"},
		[]string{},
		[]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"})
	f([]string{"127.0.0.1 test1.loc"},
		[]string{"test1.loc"},
		[]string{})
	f([]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"},
		[]string{"test1.loc"},
		[]string{"127.0.0.1 test2.loc"})
	f([]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"},
		[]string{"test2.loc"},
		[]string{"127.0.0.1 test1.loc"})
	f([]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"},
		[]string{"test1.loc", "test2.loc"},
		[]string{})

	f([]string{"127.0.0.1 test1.loc"},
		[]string{"test2.loc"},
		[]string{"127.0.0.1 test1.loc"})
	f([]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"},
		[]string{"test3.loc"},
		[]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"})
	f([]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"},
		[]string{"test3.loc", "test4.loc", "test5.loc"},
		[]string{"127.0.0.1 test1.loc", "127.0.0.1 test2.loc"})
	f([]string{},
		[]string{"test1.loc"},
		[]string{})

	f([]string{"# 127.0.0.1 test1.loc"},
		[]string{},
		[]string{"# 127.0.0.1 test1.loc"})
	f([]string{"# 127.0.0.1 test1.loc"},
		[]string{"test1.loc"},
		[]string{})
	f([]string{"# 127.0.0.1 test1.loc", "127.0.0.1 test2.loc"},
		[]string{"test1.loc"},
		[]string{"127.0.0.1 test2.loc"})
	f([]string{"127.0.0.1 test1.loc", "# 127.0.0.1 test2.loc"},
		[]string{"test2.loc"},
		[]string{"127.0.0.1 test1.loc"})
}
