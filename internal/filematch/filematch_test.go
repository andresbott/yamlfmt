package filematch

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseGlob(t *testing.T) {

	tcs := []struct {
		name   string
		in     string
		expect []string
	}{
		{
			name: "assert alias ./",
			in:   "./",
			expect: []string{
				"main.yaml",
				"other.yml",
			},
		},
		{
			name: "assert default to  ./",
			in:   "",
			expect: []string{
				"main.yaml",
				"other.yml",
			},
		},
		{
			name: "assert alias ./..",
			in:   "./..",
			expect: []string{
				"main.yaml",
				"fruits/banana/banana.yaml",
				"other.yml",
				"sub1/sub2/sub.yml",
			},
		},
		{
			name: "assert single file",
			in:   "fruits/banana/banana.yaml",
			expect: []string{
				"fruits/banana/banana.yaml",
			},
		},
		{
			name: "assert directory",
			in:   "fruits/banana/",
			expect: []string{
				"fruits/banana/banana.yaml",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			os.Chdir("./sampledata")
			got, err := FindFiles(tc.in)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			os.Chdir("..")
			if diff := cmp.Diff(got, tc.expect); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}
		})
	}
}
