package pathlib

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var CurrentWorkingDirectory = Path(os.Args[0]).Dir()

type Path string

func (path Path) Abs() Path {
	pathString, _ := filepath.Abs(string(path))
	return Path(pathString)
}

func (path Path) Base() string {
	return strings.ReplaceAll(filepath.Base(string(path)), path.Ext(), "")
}

func (path Path) Dir() Path {
	return Path(filepath.Dir(string(path)))
}

func (path Path) Ext() string {
	return filepath.Ext(string(path))
}

func (path Path) IsExist() bool {
	_, err := path.open()
	return !os.IsNotExist(err)
}

func (path Path) IsDir() bool {
	if path.IsExist() {
		info, _ := path.open()
		return info.IsDir()
	} else {
		return false
	}
}

func (path Path) Join(element ...string) Path {
	tempPath := make([]string, len(element)+1)
	tempPath[0] = string(path)
	copy(tempPath[1:], element)
	return Path(filepath.Join(tempPath...))
}

func (path Path) ReadFile() ([]byte, error) {
	return os.ReadFile(string(path))
}

func (path Path) open() (fs.FileInfo, error) {
	return os.Stat(string(path))
}

func Glob(pattern Path) ([]Path, error) {
	matches, err := filepath.Glob(string(pattern))
	if err != nil {
		return nil, err
	}
	var paths []Path
	for _, m := range matches {
		paths = append(paths, Path(m))
	}
	return paths, nil
}
