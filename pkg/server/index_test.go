package server

import (
	"os"
	"sort"
	"testing"
)

func TestGetCommaSeparatedKVPairEnv(t *testing.T) {
	t.Run("With single kv pair", func(t *testing.T) {
		os.Setenv("TEST_KV_PAIRS", "bar=baz")

		result := getCommaSeperatedKVPairEnv("TEST_KV_PAIRS")

		if val, contained := result["bar"]; !contained || val != "baz" {
			t.Errorf("expect result['bar'] to be 'baz' but %s", val)
		}

		os.Unsetenv("TEST_KV_PAIRS")
	})
	t.Run("With multiple kv pairs", func(t *testing.T) {
		os.Setenv("TEST_KV_PAIRS", "foo=bar,qux=quux")

		result := getCommaSeperatedKVPairEnv("TEST_KV_PAIRS")

		if val, contained := result["foo"]; !contained || val != "bar" {
			t.Errorf("expect result['foo'] to be 'bar' but %s", val)
		}
		if val, contained := result["qux"]; !contained || val != "quux" {
			t.Errorf("expect result['qux'] to be 'quux' but %s", val)
		}
		os.Unsetenv("TEST_KV_PAIRS")
	})
}

func TestGenerateGaMetaTags(t *testing.T) {
	input := map[string]string{
		"foo": "bar",
		"qux": "quux",
	}

	result := generateGaMetaTags(input)
	sort.Strings(result)

	expect0 := `<meta id="ga-meta-foo" content="bar">`
	expect1 := `<meta id="ga-meta-qux" content="quux">`
	if len(result) != 2 {
		t.Errorf("expect len(result) = 2 but %d", len(result))
	}
	if result[0] != expect0 {
		t.Errorf("expect result[0] to be %s but %s", expect0, result[0])
	}
	if result[1] != expect1 {
		t.Errorf("expect result[1] to be %s but %s", expect1, result[1])
	}
}

func TestGetServerBasePathMetaTag(t *testing.T) {
	testCases := []struct {
		name       string
		before     func()
		after      func()
		wantResult string
	}{
		{
			name: "With server base path env",
			before: func() {
				os.Setenv(
					"KHI_SERVER_BASE_PATH",
					"/foo/bar/",
				)
			},
			after: func() {
				os.Unsetenv("KHI_SERVER_BASE_PATH")
			},
			wantResult: `<meta id="server-base-path" content="/foo/bar/">`,
		},
		{
			name: "without the training slash",
			before: func() {
				os.Setenv(
					"KHI_SERVER_BASE_PATH",
					"/foo/bar",
				)
			},
			after: func() {
				os.Unsetenv("KHI_SERVER_BASE_PATH")
			},
			wantResult: `<meta id="server-base-path" content="/foo/bar/">`,
		},
		{
			name: "Without server base path env",
			before: func() {
				os.Unsetenv("KHI_SERVER_BASE_PATH")
			},
			after: func() {
				os.Unsetenv("KHI_SERVER_BASE_PATH")
			},
			wantResult: `<meta id="server-base-path" content="/">`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			got := getServerBasePathMetaTag()
			if got != tc.wantResult {
				t.Errorf("got %s, want %s", got, tc.wantResult)
			}
		})
	}
}

func TestGetBaseTag(t *testing.T) {
	testCases := []struct {
		name       string
		before     func()
		after      func()
		wantResult string
	}{
		{
			name: "With frontend resource base path env",
			before: func() {
				os.Setenv(
					"KHI_FRONTEND_RESOURCE_BASE_PATH",
					"/foo/bar/",
				)
			},
			after: func() {
				os.Unsetenv("KHI_FRONTEND_RESOURCE_BASE_PATH")
			},
			wantResult: `<base href="/foo/bar/">`,
		},
		{
			name: "With server base path env",
			before: func() {
				os.Setenv(
					"KHI_SERVER_BASE_PATH",
					"/foo/bar/",
				)
			},
			after: func() {
				os.Unsetenv("KHI_SERVER_BASE_PATH")
			},
			wantResult: `<base href="/foo/bar/">`,
		},
		{
			name: "complements the training slash",
			before: func() {
				os.Setenv(
					"KHI_FRONTEND_RESOURCE_BASE_PATH",
					"/foo/bar",
				)
			},
			after: func() {
				os.Unsetenv("KHI_FRONTEND_RESOURCE_BASE_PATH")
			},
			wantResult: `<base href="/foo/bar/">`,
		},
		{
			name: "prioritize KHI_FRONTEND_RESOURCE_BASE_PATH env if both envs are set",
			before: func() {
				os.Setenv(
					"KHI_SERVER_BASE_PATH",
					"/foo/bar/",
				)
				os.Setenv(
					"KHI_FRONTEND_RESOURCE_BASE_PATH",
					"/qux/",
				)
			},
			after: func() {
				os.Unsetenv("KHI_SERVER_BASE_PATH")
				os.Unsetenv("KHI_FRONTEND_RESOURCE_BASE_PATH")
			},
			wantResult: `<base href="/qux/">`,
		},
		{
			name: "Without any env",
			before: func() {
				os.Unsetenv("KHI_FRONTEND_RESOURCE_BASE_PATH")
				os.Unsetenv("KHI_SERVER_BASE_PATH")
			},
			after: func() {
				os.Unsetenv("KHI_FRONTEND_RESOURCE_BASE_PATH")
				os.Unsetenv("KHI_SERVER_BASE_PATH")
			},
			wantResult: `<base href="/">`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			got := getBaseTag()
			if got != tc.wantResult {
				t.Errorf("got %s, want %s", got, tc.wantResult)
			}
		})
	}
}
