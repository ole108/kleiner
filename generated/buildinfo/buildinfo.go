package buildinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/can3p/kleiner/generated/internal/version"
	"github.com/pkg/errors"
)

var (
	// set during init
	cachedVersion   version.Version
	cachedBuildTime time.Time
	cachedName      string
	branchName      string
)

func init() {
	if err := loadMeta(); err != nil {
		panic(err)
	}
}

func loadMeta() error {
	if err := loadCachedName(); err != nil {
		return errors.Wrap(err, "error loading executable name")
	}

	if err := loadBuildTime(); err != nil {
		return errors.Wrap(err, "error loading build date from embedded flag")
	}

	if err := loadVersion(); err != nil {
		return errors.Wrap(err, "error loading build version from embedded flag")
	}

	return nil
}

func loadCachedName() error {
	execName, err := os.Executable()
	if err != nil {
		return err
	}
	cachedName = filepath.Base(execName)
	return nil
}

// Name returns the executable that started the current process.
func Name() string {
	return cachedName
}

type info struct {
	Name         string
	Version      version.Version
	Commit       string
	BranchName   string
	BuildDate    time.Time
	OS           string
	Architecture string
	Environment  string
	GithubRepo   string
}

func (i info) String() string {
	res := fmt.Sprintf("%s v%s %s/%s Commit: %s BuildDate: %s",
		i.Name,
		i.Version,
		i.OS,
		i.Architecture,
		i.Commit,
		i.BuildDate.Format(time.RFC3339))
	if i.BranchName != "" {
		res += fmt.Sprintf(" BranchName: %s", i.BranchName)
	}

	res += fmt.Sprintf(" Github Repo: https://github.com/%s", i.GithubRepo)

	return res
}

func Info() info {
	return info{
		Name:         Name(),
		Version:      Version(),
		Commit:       Commit(),
		BranchName:   BranchName(),
		BuildDate:    BuildTime(),
		OS:           OS(),
		Architecture: Arch(),
		Environment:  Environment(),
		GithubRepo:   GithubRepo(),
	}
}

func OS() string {
	return runtime.GOOS
}

func Arch() string {
	return runtime.GOARCH
}

func BranchName() string {
	return branchName
}

func Version() version.Version {
	return cachedVersion
}

func BuildTime() time.Time {
	return cachedBuildTime
}

// it's generated only once, no need to have it
// as a parameter
func GithubRepo() string {
	return "can3p/kleiner"
}

func Commit() string {
	info, _ := debug.ReadBuildInfo()
	var rev string = "<none>"
	var dirty string = ""
	for _, v := range info.Settings {
		if v.Key == "vcs.revision" {
			rev = v.Value
		}
		if v.Key == "vcs.modified" {
			if v.Value == "true" {
				dirty = "-dirty"
			} else {
				dirty = ""
			}
		}
	}
	return rev + dirty
}
