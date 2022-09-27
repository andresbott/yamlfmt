package yamlfmt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
)

var _ = spew.Dump // prevent spew from being removed if unused

const sampleDir = "./sampledata"

func TestName(t *testing.T) {

	tcs := []struct {
		name   string
		in     string
		expect string
	}{
		{
			name:   "assert line breaks",
			in:     "remove_lines.yaml",
			expect: "remove_lines.out.yaml",
		},
		{
			name:   "assert quotes are removed",
			in:     "unquote_text.yaml",
			expect: "unquote_text.out.yaml",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			inBytes, err := os.ReadFile(filepath.Join(sampleDir, tc.in))
			if err != nil {
				t.Fatal(err)
			}

			got, err := Format(inBytes)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}

			expectBytes, err := os.ReadFile(filepath.Join(sampleDir, tc.expect))
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(string(got), string(expectBytes)); diff != "" {
				t.Errorf("unexpected value (-got +want)\n%s", diff)
			}

			//_ = os.WriteFile(filepath.Join(sampleDir, tc.expect), got, 0644)
		})
	}

}
