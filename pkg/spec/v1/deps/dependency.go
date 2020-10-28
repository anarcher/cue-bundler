package deps

type Dependency struct {
	Source  Source `json:"source"`
	Version string `json:"version"`
	Sum     string `json:"sum,omitempty"`
	Single  bool   `json:"single,omitempty"`
}

func Parse(dir, uri string) *Dependency {
	if uri == "" {
		return nil
	}

	if d := parseGit(uri); d != nil {
		return d
	}

	return parseLocal(dir, uri)
}

func (d Dependency) Name() string {
	return d.Source.Name()
}
