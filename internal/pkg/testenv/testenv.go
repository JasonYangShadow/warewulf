// Package testenv provides functions and data structures for
// constructing and manipulating a temporary Warewulf environment for
// use during automated testing.
//
// The testenv package should only be used in tests.
package testenv

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/warewulf/warewulf/internal/pkg/config"
)

const initWarewulfConf = ``
const initNodesConf = `nodeprofiles:
  default: {}
nodes:
  node1: {}
`

type TestEnv struct {
	t       *testing.T
	b       *testing.B
	BaseDir string
}

const Sysconfdir = "etc"
const Bindir = "bin"
const Datadir = "usr/share"
const Localstatedir = "var/local"
const Srvdir = "srv"
const Tftpdir = "srv/tftp"
const Firewallddir = "usr/lib/firewalld/services"
const Systemddir = "usr/lib/systemd/system"
const WWOverlaydir = "var/lib/warewulf/overlays"
const WWChrootdir = "var/lib/warewulf/chroots"
const WWProvisiondir = "srv/warewulf"
const Cachedir = "var/cache"

// New creates a test environment in a temporary directory and configures
// Warewulf to use it.
//
// Caller is responsible to delete env.BaseDir by calling
// env.RemoveAll. Note that this does not restore Warewulf to its
// previous state.
//
// Asserts no errors occur.
func New(t *testing.T) (env *TestEnv) {
	env = new(TestEnv)
	env.t = t
	env.init()
	return env
}

// NewBenchmark creates a benchmark environment in a temporary directory and configures
// Warewulf to use it.
//
// Caller is responsible to delete env.BaseDir by calling
// env.RemoveAll. Note that this does not restore Warewulf to its
// previous state.
//
// Asserts no errors occur.
func NewBenchmark(b *testing.B) (env *TestEnv) {
	env = new(TestEnv)
	env.b = b
	env.init()
	return env
}

// init creates a temporary directory and configured Warewulf to use it on behalf of
// New() or NewBenchmark().
//
// Asserts no errors occur.
func (env *TestEnv) init() {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "ww4test-*")
	env.assertNoError(err)
	env.BaseDir = tmpDir

	env.WriteFile(path.Join(Sysconfdir, "warewulf/nodes.conf"), initNodesConf)
	env.WriteFile(path.Join(Sysconfdir, "warewulf/warewulf.conf"), initWarewulfConf)

	// re-read warewulf.conf
	_ = env.Configure()
}

func (env *TestEnv) Configure() *config.WarewulfYaml {
	conf := config.New()
	err := conf.Read(env.GetPath(path.Join(Sysconfdir, "warewulf/warewulf.conf")), false)
	env.assertNoError(err)
	conf.Paths.Sysconfdir = env.GetPath(Sysconfdir)
	conf.Paths.Bindir = env.GetPath(Bindir)
	conf.Paths.Datadir = env.GetPath(Datadir)
	conf.Paths.Localstatedir = env.GetPath(Localstatedir)
	conf.Paths.Srvdir = env.GetPath(Srvdir)
	conf.TFTP.TftpRoot = env.GetPath(Tftpdir)
	conf.Paths.Firewallddir = env.GetPath(Firewallddir)
	conf.Paths.Systemddir = env.GetPath(Systemddir)
	conf.Paths.WWOverlaydir = env.GetPath(WWOverlaydir)
	conf.Paths.WWChrootdir = env.GetPath(WWChrootdir)
	conf.Paths.WWProvisiondir = env.GetPath(WWProvisiondir)
	conf.Paths.Cachedir = env.GetPath(Cachedir)
	conf.Paths.WWClientdir = "/warewulf"

	for _, confPath := range []string{
		conf.Paths.Sysconfdir,
		conf.Paths.Bindir,
		conf.Paths.Datadir,
		conf.Paths.Localstatedir,
		conf.Paths.Srvdir,
		conf.TFTP.TftpRoot,
		conf.Paths.Firewallddir,
		conf.Paths.Systemddir,
		conf.Paths.WWOverlaydir,
		conf.Paths.WWChrootdir,
		conf.Paths.WWProvisiondir,
	} {
		env.MkdirAll(confPath)
	}
	return conf
}

// assertNoError handles error conditions generated by env,
// using the semantics of either a testing or a benchmark environment.
func (env *TestEnv) assertNoError(err error, msgAndArgs ...interface{}) {
	if env.t != nil {
		assert.NoError(env.t, err, msgAndArgs...)
	}

	if env.b != nil && err != nil {
		if len(msgAndArgs) == 0 {
			env.b.Errorf("%s", err)
		} else {
			env.b.Errorf(msgAndArgs[0].(string), (msgAndArgs[1:])...)
		}
	}
}

// GetPath returns the absolute path name for fileName specified
// relative to the test environment.
func (env *TestEnv) GetPath(fileName string) string {
	return path.Join(env.BaseDir, fileName)
}

// MkdirAll creates dirName and any intermediate directories relative
// to the test environment.
//
// Asserts no errors occur.
func (env *TestEnv) MkdirAll(dirName string) {
	err := os.MkdirAll(env.GetPath(dirName), 0755)
	env.assertNoError(err)
}

// Chmod changes the mode of fileName in the test environment, which must already exist.
//
// Asserts no errors occur.
func (env *TestEnv) Chmod(fileName string, mode int) {
	err := os.Chmod(env.GetPath(fileName), os.FileMode(mode))
	env.assertNoError(err)
}

// WriteFile writes content to fileName, creating any necessary
// intermediate directories relative to the test environment.
//
// Asserts no errors occur.
func (env *TestEnv) WriteFile(fileName string, content string) {
	dirName := filepath.Dir(fileName)
	env.MkdirAll(dirName)

	f, err := os.Create(env.GetPath(fileName))
	env.assertNoError(err)
	defer f.Close()
	_, err = f.WriteString(content)
	env.assertNoError(err)
	err = os.Chtimes(env.GetPath(fileName),
		time.Date(2006, time.February, 1, 3, 4, 5, 0, time.UTC),
		time.Date(2006, time.February, 1, 3, 4, 5, 0, time.UTC))
	env.assertNoError(err)
}

// ImportFile writes the contents of inputFileName to fileName,
// creating any necessary intermediate directories relative to the
// test environment.
func (env *TestEnv) ImportFile(fileName string, inputFileName string) {
	buffer, err := os.ReadFile(inputFileName)
	env.assertNoError(err)
	env.WriteFile(fileName, string(buffer))
}

func (env *TestEnv) ImportDir(dirName string, inputDirName string) {
	env.MkdirAll(path.Dir(dirName))
	cmd := exec.Command("cp", "--recursive", inputDirName, env.GetPath(dirName))
	output, err := cmd.CombinedOutput()
	env.assertNoError(err, string(output))
}

// CreateFile creates an empty file at fileName, creating any necessary intermediate directories
// relative to the test environment.
func (env *TestEnv) CreateFile(fileName string) {
	env.WriteFile(fileName, "")
}

// Symlink creates a symlink at fileName to target, creating any necessary intermediate directories
// relative to the test environment.
func (env *TestEnv) Symlink(target string, fileName string) {
	dirName := filepath.Dir(fileName)
	env.MkdirAll(dirName)

	err := os.Symlink(target, env.GetPath(fileName))
	env.assertNoError(err)
}

// ReadFile returns the content of fileName as converted to a
// string.
//
// Asserts no errors occur.
func (env *TestEnv) ReadFile(fileName string) string {
	buffer, err := os.ReadFile(env.GetPath(fileName))
	env.assertNoError(err)
	return string(buffer)
}

// ReadDir returns the content of dirName as converted to a
// slice of strings.
//
// Asserts no errors occur.
func (env *TestEnv) ReadDir(dirName string) []string {
	entries, err := os.ReadDir(env.GetPath(dirName))
	env.assertNoError(err)
	var entryStrs []string
	for _, entry := range entries {
		entryStrs = append(entryStrs, entry.Name())
	}
	return entryStrs
}

// RemoveAll deletes the temporary directory, and all its contents,
// for the test environment.
//
// Asserts no errors occur.
func (env *TestEnv) RemoveAll() {
	err := os.RemoveAll(env.BaseDir)
	env.assertNoError(err)
}
