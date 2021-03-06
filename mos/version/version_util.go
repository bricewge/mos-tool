package version

//go:generate sh -c "../../common/tools/fw_meta.py gen_build_info --version=`[ -f version ] && cat version` --tag_as_version=true --id=`[ -f build_id ] && cat build_id` --go_output=version.go  --json_output=version.json"

import (
	"regexp"
	"strings"

	"cesanta.com/common/go/ourutil"
	moscommon "cesanta.com/mos/common"
)

type VersionJson struct {
	BuildId        string `json:"build_id"`
	BuildTimestamp string `json:"build_timestamp"`
	BuildVersion   string `json:"build_version"`
}

const (
	LatestVersionName = "latest"
)

var (
	regexpVersionNumber = regexp.MustCompile(`^\d+\.[0-9.]*$`)
	regexpBuildIdDistr  = regexp.MustCompile(
		`^(?P<version>[^+]+)\+(?P<hash>[^~]+)\~(?P<distr>.+)$`,
	)

	ubuntuDistrNames = []string{"xenial", "bionic", "cosmic"}
)

// GetMosVersion returns this binary's version, or "latest" if it's not a release build.
func GetMosVersion() string {
	if LooksLikeVersionNumber(Version) {
		return Version
	}
	return LatestVersionName
}

// GetMosVersionSuffix returns an empty string if mos version is "latest";
// otherwise returns the mos version prepended with a dash, like "-1.6".
func GetMosVersionSuffix() string {
	return moscommon.GetVersionSuffix(GetMosVersion())
}

func LooksLikeVersionNumber(s string) bool {
	return regexpVersionNumber.MatchString(s)
}

// Returns whether the build id looks like the mos was built in some distro
// environment (like, ubuntu or brew), and thus it shouldn't update itself.
func LooksLikeDistrBuildId(s string) bool {
	return ourutil.FindNamedSubmatches(regexpBuildIdDistr, s) != nil
}

func LooksLikeUbuntuBuildId(s string) bool {
	return GetUbuntuPackageName(s) != ""
}

func LooksLikeBrewBuildId(s string) bool {
	return strings.HasSuffix(s, "~brew")
}

// GetUbuntuPackageName parses given build id string, and if it looks like a
// debian build id, returns either "mos-latest" or "mos". Otherwise, returns
// an empty string.
func GetUbuntuPackageName(buildId string) string {
	matches := ourutil.FindNamedSubmatches(regexpBuildIdDistr, buildId)
	if matches != nil {
		for _, v := range ubuntuDistrNames {
			if strings.HasPrefix(matches["distr"], v) {
				if LooksLikeVersionNumber(matches["version"]) {
					return "mos"
				} else {
					return "mos-latest"
				}
			}
		}

		// Some distro other than Ubuntu
		return ""
	} else {
		// Doesn't look like distro build id
		return ""
	}
}
