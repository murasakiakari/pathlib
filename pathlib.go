package pathlib

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Full control for owner, read and execution for group and public
const DEFAULT_PERM = 0755

var (
	CurrentWorkingDirectory, _ = GetCurrentWorkingDirectory()
	CurrentExecutablePath, _   = GetCurrentExecutablePath()
)

func pathWrapper(path string, err error) (Path, error) {
	return Path(path), err
}

func GetCurrentWorkingDirectory() (Path, error) {
	path, err := pathWrapper(os.Getwd())
	if err == nil {
		path, err = path.EvalSymlinks()
	}
	return path, err
}

func GetCurrentExecutablePath() (Path, error) {
	path, err := pathWrapper(os.Executable())
	if err == nil {
		path, err = path.EvalSymlinks()
	}
	return path, err
}

type Path string

func (p Path) String() string {
	return string(p)
}

// Method that wrap filepath function

// Abs returns an absolute path of the current path.
//
// More please see: https://pkg.go.dev/path/filepath#Abs
func (p Path) Abs() (Path, error) {
	return pathWrapper(filepath.Abs(p.String()))
}

// Base returns the filename with extension (if any) of the path.
//
// More please see: https://pkg.go.dev/path/filepath#Base
func (p Path) Base() string {
	return filepath.Base(p.String())
}

// Clean returns the shortest path equivalent of the path.
//
// More please see: https://pkg.go.dev/path/filepath#Clean
func (p Path) Clean() Path {
	return Path(filepath.Clean(p.String()))
}

// Dir returns the directory of the path.
//
// More please see: https://pkg.go.dev/path/filepath#Dir
func (p Path) Dir() Path {
	return Path(filepath.Dir(p.String()))
}

// EvalSymlinks returns the path after evaluating symbolic links (if any).
//
// More please see: https://pkg.go.dev/path/filepath#EvalSymlinks
func (p Path) EvalSymlinks() (Path, error) {
	return pathWrapper(filepath.EvalSymlinks(p.String()))
}

// Ext returns the extension (if any) of the path.
//
// More please see: https://pkg.go.dev/path/filepath#Ext
func (p Path) Ext() string {
	return filepath.Ext(p.String())
}

// FromSlash returns the result of replacing each slash ('/') character in path with a separator character.
//
// More please see: https://pkg.go.dev/path/filepath#FromSlash
func (p Path) FromSlash() Path {
	return Path(filepath.FromSlash(p.String()))
}

// Glob returns the names of all files matching pattern or nil if there is no matching file.
//
// More please see: https://pkg.go.dev/path/filepath#FromSlash
func (p Path) Glob(patterns ...string) ([]Path, error) {
	realPattern := p.Join(patterns...)
	files, err := filepath.Glob(realPattern.String())
	if err != nil {
		return nil, err
	}
	if files == nil {
		return nil, nil
	}
	paths := make([]Path, len(files))
	for i, file := range files {
		paths[i] = Path(file)
	}
	return paths, nil
}

// Join returns the result of joining any number of elements into the path with an OS specific Separator.
//
// More please see: https://pkg.go.dev/path/filepath#Join
func (p Path) Join(element ...string) Path {
	elements := make([]string, len(element)+1)
	elements[0] = p.String()
	copy(elements[1:], element)
	return Path(filepath.Join(elements...))
}

// Match returns true if it matches to the pathPattern, error is not nil if the pattern is malformed.
//
// More please see: https://pkg.go.dev/path/filepath#Match
func (p Path) Match(pathPattern Path) (bool, error) {
	return filepath.Match(pathPattern.String(), p.String())
}

// Rel returns a relative path of the path to target path,
// error is not nil if the current path pleasenot make relative to the target path.
//
// More please see: https://pkg.go.dev/path/filepath#Rel
func (p Path) Rel(target Path) (Path, error) {
	return pathWrapper(filepath.Rel(p.String(), target.String()))
}

// Split returns the directory and the filename of the path.
//
// More please  see: https://pkg.go.dev/path/filepath#Split
func (p Path) Split() (dir Path, file string) {
	dir_, file_ := filepath.Split(p.String())
	return Path(dir_), file_
}

// ToSlash returns the result of replacing each separator character in path with a slash.
//
// More please see: https://pkg.go.dev/path/filepath#ToSlash
func (p Path) ToSlash() Path {
	return Path(filepath.ToSlash(p.String()))
}

// VolumeName returns the volume name of the path.
//
// More please see: https://pkg.go.dev/path/filepath#VolumeName
func (p Path) VolumeName() Path {
	return Path(filepath.VolumeName(p.String()))
}

// Walk walks the file tree rooted at path and calling walkFunc for each file and directory including the root.
// Walk is less efficient than WalkDir.
//
// More please see: https://pkg.go.dev/path/filepath#Walk
func (p Path) Walk(walkFunc filepath.WalkFunc) error {
	return filepath.Walk(p.String(), walkFunc)
}

// Walk walks the file tree rooted at path and calling walkFunc for each file and directory including the root.
// WalkDir is more efficient than Walk.
//
// More please see: https://pkg.go.dev/path/filepath#Walk
func (p Path) WalkDir(walkDirFunc fs.WalkDirFunc) error {
	return filepath.WalkDir(p.String(), walkDirFunc)
}

// SplitList split the path which joined by the OS-specific ListSeparator into the path list.
//
// More please see: https://pkg.go.dev/path/filepath#SplitList
func SplitList(path string) []Path {
	paths := filepath.SplitList(path)
	resultPaths := make([]Path, len(paths))
	for i, path := range paths {
		resultPaths[i] = Path(path)
	}
	return resultPaths
}

// Helper method using filepath function

// List returns a list of all file and directory root at path.
func (p Path) List() ([]Path, error) {
	return p.Glob(Path("*").Join("**").String())
}

// SplitAll returns the directory, filename without extension and extension of the path.
func (p Path) SplitAll() (dir Path, filename, ext string) {
	dir, filename = p.Split()
	ext = filepath.Ext(filename)
	return dir, strings.TrimSuffix(filename, ext), ext
}

// AddPrefix returns the path with the given prefix to the filename.
func (p Path) AddPrefix(prefix string) Path {
	dir, filename := p.Split()
	return dir.Join(prefix + filename)
}

// AddPostfix returns the path with the given postfix to the filename.
func (p Path) AddPostfix(postfix string) Path {
	dir, filename, ext := p.SplitAll()
	return dir.Join(filename + postfix + ext)
}

// Implementation for os function involving path

// Chdir changes the current working directory to the path.
//
// More please see: https://pkg.go.dev/os#Chdir
func (p Path) Chdir() error {
	return os.Chdir(p.String())
}

// Chmod changes the mode of the file with the path name.
//
// More please see: https://pkg.go.dev/os#Chmod
func (p Path) Chmod(mode os.FileMode) error {
	return os.Chmod(p.String(), mode)
}

// Chown changes the numeric uid and gid of the file with the path name.
//
// More please see: https://pkg.go.dev/os#Chown
func (p Path) Chown(uid, gid int) error {
	return os.Chown(p.String(), uid, gid)
}

// Chtimes changes the access and modification times of the file with the path name.
//
// More please see: https://pkg.go.dev/os#Chtimes
func (p Path) Chtimes(atime, mtime time.Time) error {
	return os.Chtimes(p.String(), atime, mtime)
}

// DirFS returns a file system for the file tree rooted at path.
//
// More please see: https://pkg.go.dev/os#DirFS
func (p Path) DirFS() fs.FS {
	return os.DirFS(p.String())
}

// Lchown changes the numeric uid and gid of the file with the path name.
//
// More please see: https://pkg.go.dev/os#Lchown
func (p Path) Lchown(uid, gid int) error {
	return os.Lchown(p.String(), uid, gid)
}

// Link creates newname as a hard link to the path.
//
// More please see: https://pkg.go.dev/os#Link
func (p Path) Link(newname Path) error {
	return os.Link(p.String(), newname.String())
}

// Mkdir creates a new directory with the path name and specific permission bit.
//
// More please see: https://pkg.go.dev/os#Mkdir
func (p Path) Mkdir(perm os.FileMode) error {
	return os.Mkdir(p.String(), perm)
}

// MkdirAll creates all necessary directory with the path and returns an error if not directory is created.
//
// More please see: https://pkg.go.dev/os#MkdirAll
func (p Path) MkdirAll(perm os.FileMode) error {
	return os.MkdirAll(p.String(), perm)
}

// MkdirTemp create a new temporary directory in the path and returns the path of the temporary directory.
//
// More please see: https://pkg.go.dev/os#MkdirTemp
func (p Path) MkdirTemp(pattern string) (Path, error) {
	return pathWrapper(os.MkdirTemp(p.String(), pattern))
}

// ReadFile reads the path and return the contains.
//
// More please see: https://pkg.go.dev/os#ReadFile
func (p Path) ReadFile() ([]byte, error) {
	return os.ReadFile(p.String())
}

// Readlink returns the destination if the path is a symbolic link.
//
// More please see: https://pkg.go.dev/os#Readlink
func (p Path) ReadLink() (Path, error) {
	return pathWrapper(os.Readlink(p.String()))
}

// Remove removes the file or empty directory of the path name.
//
// More please see: https://pkg.go.dev/os#Remove
func (p Path) Remove() error {
	return os.Remove(p.String())
}

// RemoveAll removes anything of the path.
// RemoveAll return nil if the path is not exist.
//
// More please see: https://pkg.go.dev/os#RemoveAll
func (p Path) RemoveAll() error {
	return os.RemoveAll(p.String())
}

// Rename renames the path to the newpath, it replace the file with newpath if it exists.
//
// More please see: https://pkg.go.dev/os#Rename
func (p Path) Rename(newpath Path) error {
	return os.Rename(p.String(), newpath.String())
}

// Symlink creates newname as a symbolic link to the path.
//
// More please see: https://pkg.go.dev/os#Symlink
func (p Path) Symlink(newname Path) error {
	return os.Symlink(p.String(), newname.String())
}

// Truncate changes the size of the file with the path name.
//
// More please see: https://pkg.go.dev/os#Truncate
func (p Path) Truncate(size int64) error {
	return os.Truncate(p.String(), size)
}

// WriteFile writes data to the file with the path name.
//
// More please see: https://pkg.go.dev/os#WriteFile
func (p Path) WriteFile(data []byte, perm os.FileMode) error {
	return os.WriteFile(p.String(), data, perm)
}

// ReadDir reads the directory with the path name and returns all directory entries sorted by filename.
//
// More please see: https://pkg.go.dev/os#ReadDir
func (p Path) ReadDir() ([]fs.DirEntry, error) {
	return os.ReadDir(p.String())
}

// Create creates the file with the path name if the file does not exist,
// and truncates the file if it exists.
//
// More please see: https://pkg.go.dev/os#Create
func (p Path) Create() (*os.File, error) {
	return os.Create(p.String())
}

// CreateTemp creates a new temporary file in the directory with the path name,
// opens that file for reading and writing and returns the resulting file.
//
// More please see: https://pkg.go.dev/os#CreateTemp
func (p Path) CreateTemp(pattern string) (*os.File, error) {
	return os.CreateTemp(p.String(), pattern)
}

// Open opens the file with the path name for reading.
//
// More please see: https://pkg.go.dev/os#Open
func (p Path) Open() (*os.File, error) {
	return os.Open(p.String())
}

// OpenFile opens the file with the path name with the specified flag (eg os.O_RDONLY).
//
// More please see: https://pkg.go.dev/os#OpenFile
func (p Path) OpenFile(flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(p.String(), flag, perm)
}

// Lstat returns a FileInfo describing the file with the path name and it will not follow the link.
//
// More please see: https://pkg.go.dev/os#Lstat
func (p Path) Lstat() (fs.FileInfo, error) {
	return os.Lstat(p.String())
}

// Lstat returns a FileInfo describing the file with the path name.
//
// More please see: https://pkg.go.dev/os#Stat
func (p Path) Stat() (fs.FileInfo, error) {
	return os.Stat(p.String())
}

// TempDir returns the default directory to use for temporary files.
//
// More please see: https://pkg.go.dev/os#TempDir
func TempDir() Path {
	return Path(os.TempDir())
}

// UserCacheDir returns the default root directory to use for user-specific cached data.
//
// More please see: https://pkg.go.dev/os#UserCacheDir
func UserCacheDir() (Path, error) {
	return pathWrapper(os.UserCacheDir())
}

// UserConfigDir returns the default root directory to use for user-specific configuration data.
//
// More please see: https://pkg.go.dev/os#UserConfigDir
func UserConfigDir() (Path, error) {
	return pathWrapper(os.UserConfigDir())
}

// UserHomeDir returns the ome directory of the current user.
//
// More please see: https://pkg.go.dev/os#UserHomeDir
func UserHomeDir() (Path, error) {
	return pathWrapper(os.UserHomeDir())
}

// Misc boolean function involving path

// IsAbs returns true if the path is absolute.
func (p Path) IsAbs() bool {
	return filepath.IsAbs(p.String())
}

// IsExist returns true if the path is a file or directory.
func (p Path) IsExist() bool {
	_, err := p.Stat()
	return !errors.Is(err, fs.ErrNotExist)
}

// IsDir returns true if the path is a directory, it returns false if the path not exist.
func (p Path) IsDir() bool {
	info, err := p.Stat()
	if !errors.Is(err, fs.ErrExist) {
		return false
	}
	return info.IsDir()
}

// IO method

// CopyToFile copy the data of the file with the path name to the destinationPath.
// If destinationPath is a directory, it will copy the file to that directory with the same filename.
func (p Path) CopyToFile(destinationPath Path, bufferSize uint64) (copiedSize uint64, err error) {
	err = destinationPath.Dir().MkdirAll(DEFAULT_PERM)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return 0, err
	}

	source, err := p.Open()
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := destinationPath.Create()
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	return buffedCopy(source, destination, bufferSize)
}

// CopyToDirectory copy the data of the file with the path name to directoryPath
// and returns the path of the copied file.
func (p Path) CopyToDirectory(directoryPath Path, bufferSize uint64) (destinationPath Path, copiedSize uint64, err error) {
	destinationPath = directoryPath.Join(p.Base())
	copiedSize, err = p.CopyToFile(destinationPath, bufferSize)
	return destinationPath, copiedSize, err
}

// AppendFile appends data to the end of the file with the path name.
func (p Path) AppendFile(data string, bufferSize uint64) (appendedSize uint64, err error) {
	reader := strings.NewReader(data)

	file, err := p.OpenFile(os.O_WRONLY|os.O_CREATE|os.O_APPEND, DEFAULT_PERM)
	if err != nil {
		return
	}
	defer file.Close()

	return buffedCopy(reader, file, bufferSize)
}

// BuffedReadFile reads data from the file with the path name with buffer.
func (p Path) BuffedReadFile(bufferSize uint64) (data []byte, readSize uint64, err error) {
	writer := &bytes.Buffer{}

	file, err := p.OpenFile(os.O_RDONLY|os.O_CREATE, DEFAULT_PERM)
	if err != nil {
		return
	}
	defer file.Close()

	copiedSize, err := buffedCopy(file, writer, bufferSize)
	return writer.Bytes(), copiedSize, err
}

// BuffedWriteFile write data to the file with the path name with buffer.
func (p Path) BuffedWriteFile(data []byte, bufferSize uint64) (writtenSize uint64, err error) {
	reader := bytes.NewReader(data)

	file, err := p.OpenFile(os.O_WRONLY|os.O_CREATE, DEFAULT_PERM)
	if err != nil {
		return
	}
	defer file.Close()

	return buffedCopy(reader, file, bufferSize)
}

func buffedCopy(reader io.Reader, writer io.Writer, bufferSize uint64) (copiedSize uint64, err error) {
	buffer := make([]byte, bufferSize)
	bufferReader := bufio.NewReader(reader)
	bufferWriter := bufio.NewWriter(writer)

	for {
		nr, err := bufferReader.Read(buffer)
		if err != nil && err != io.EOF {
			return copiedSize, err
		}
		if nr == 0 {
			bufferWriter.Flush()
			break
		}
		nw, err := bufferWriter.Write(buffer[:nr])
		if err != nil {
			return copiedSize, err
		}
		if nr != nw {
			return copiedSize, io.ErrShortWrite
		}
		copiedSize += uint64(nw)
	}
	return
}
