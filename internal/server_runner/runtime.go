package serverrunner

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type Runtime struct {
	Name    string
	Version string
}

func ParseRuntime(spec string) (*Runtime, error) {
	// Split the spec into runtime and version ('@' is the separator)
	// If the version is not specified, it is assumed to be the the tag 'latest'.
	// If the version is specified, it must be a valid semver version.
	// If the version is not a valid semver version, return an error.
	// If the runtime is not a valid runtime, return an error. We support `node`
	// and `python` runtimes.

	specParts := strings.SplitN(spec, "@", 2)

	switch specParts[0] {
	case "node", "python":
	default:
		return nil, fmt.Errorf("unsupported runtime: %s", specParts[0])
	}

	if len(specParts) == 1 {
		return &Runtime{
			Name:    specParts[0],
			Version: "",
		}, nil
	}

	version, err := semver.StrictNewVersion(specParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid semver: %s", specParts[1])
	}

	return &Runtime{
		Name:    specParts[0],
		Version: version.String(),
	}, nil
}
