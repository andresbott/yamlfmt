package filematch

import (
	"fmt"
	"github.com/yargevad/filepathx"
	"os"
	"strings"
)

func aggregate(glob ...string) ([]string, error) {
	ret := []string{}
	for _, g := range glob {
		got, err := find(g)
		if err != nil {
			return []string{}, err
		}
		ret = append(ret, got...)
	}
	return ret, nil
}

func find(glob string) ([]string, error) {
	matches, err := filepathx.Glob(glob)
	if err != nil {
		return []string{}, fmt.Errorf("unable to glob files: %v", err)
	}
	return matches, nil
}

func FindFiles(glob string) ([]string, error) {

	switch glob {

	case "./", "":
		return aggregate("./*.yaml", "./*.yml")
	case "./..":
		return aggregate("./**/*.yaml", "./**/*.yml")
	default:
		if finfo, err := os.Stat(glob); err == nil {
			if finfo.IsDir() {
				glob = strings.TrimSuffix(glob, "/")
				if !strings.HasPrefix(glob, "./") {
					glob = "./" + glob
				}
				return aggregate(glob+"/*.yaml", glob+"/*.yml")
			}
		}
		return find(glob)
	}

}
