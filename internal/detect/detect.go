package detect

import "os"

// ProjectType represents the type of project detected.
type ProjectType string

const (
	Laravel  ProjectType = "laravel"
	Python   ProjectType = "python"
	Dolibarr ProjectType = "dolibarr"
	Go       ProjectType = "go"
	Node     ProjectType = "node"
	Generic  ProjectType = "generic"
)

// FromCurrentDir detects the project type by looking at files in the current directory.
func FromCurrentDir() ProjectType {
	return FromDir(".")
}

// FromDir detects the project type in the given directory.
func FromDir(dir string) ProjectType {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Generic
	}
	files := make(map[string]bool, len(entries))
	for _, e := range entries {
		files[e.Name()] = true
	}

	switch {
	case files["artisan"]:
		return Laravel
	case files["main.inc.php"]:
		return Dolibarr
	case files["go.mod"]:
		return Go
	case files["requirements.txt"] || files["pyproject.toml"] || files["setup.py"]:
		return Python
	case files["package.json"]:
		return Node
	default:
		return Generic
	}
}

// String returns a display-friendly name.
func (pt ProjectType) String() string {
	return string(pt)
}
