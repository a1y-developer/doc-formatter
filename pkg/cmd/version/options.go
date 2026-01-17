package version

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"

	"github.com/a1y/doc-formatter/pkg/util/gitutil"
	"github.com/a1y/doc-formatter/pkg/version"
)

const jsonOutput = "json"

type VersionOptions struct {
	Output string
}

func NewVersionOptions() *VersionOptions {
	return &VersionOptions{}
}

func (o *VersionOptions) Validate() error {
	if o.Output != "" && o.Output != jsonOutput {
		return errors.New("invalid output type, output must be 'json'")
	}
	return nil
}

func (o *VersionOptions) Run() {
	if strings.ToLower(o.Output) == jsonOutput {
		fmt.Println(version.JSON())
	} else {
		fmt.Println(version.String())
		if msg := checkForUpdate(); msg != "" {
			fmt.Println(msg)
		}
	}
}

// checkForUpdate checks to see if the CLI needs to be updated,
// and if so emits a warning, as well as information as to how it can be upgraded.
func checkForUpdate() string {
	curVer, err := semver.ParseTolerant(version.ReleaseVersion())
	if err != nil {
		logrus.Errorf("error parsing current version: %s", err)
	}

	// We don't care about warning for you to update if you have installed a developer version
	if isDevVersion(curVer) {
		return ""
	}

	latestVer, err := getLatestVersionInfo()
	if err != nil {
		logrus.Errorf("error fetching latest version information: %v", err)
	}

	if latestVer.GT(curVer) {
		return fmt.Sprintf("A new version of doc-formatter is available. To upgrade from version '%s' to '%s'", curVer, latestVer)
	}

	return ""
}

func isDevVersion(s semver.Version) bool {
	if len(s.Build) != 0 {
		return true
	}

	if len(s.Pre) == 0 {
		return false
	}

	devStrings := regexp.MustCompile(`alpha|beta|dev|rc`)
	return !s.Pre[0].IsNum && devStrings.MatchString(s.Pre[0].VersionStr)
}

// getLatestVersionInfo returns information about the latest version of the CLI.
func getLatestVersionInfo() (semver.Version, error) {
	latestTag, err := gitutil.GetLatestTag()
	if err != nil {
		return semver.Version{}, err
	}

	latest, err := semver.ParseTolerant(latestTag)
	if err != nil {
		return semver.Version{}, err
	}

	return latest, nil
}
