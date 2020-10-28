package deps

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	GitSchemeSSH   = "ssh://git@"
	GitSchemeHTTPS = "https://"
)

var gitProtoFmts = map[string]string{
	GitSchemeSSH:   GitSchemeSSH + "%s/%s/%s.git",
	GitSchemeHTTPS: GitSchemeHTTPS + "%s/%s/%s.git",
}

type Git struct {
	Scheme string

	Host string
	// User (example.com/<user>)
	User string
	// Repo (example.com/<user>/<repo>)
	Repo string
	// Subdir (example.com/<user></<repo>/<subdir>)
	Subdir string
}

func (g *Git) Name() string {
	return fmt.Sprintf("%s/%s/%s%s", g.Host, g.User, strings.TrimSuffix(g.Repo, ".git"), g.Subdir)
}

func (g *Git) Remote() string {
	return fmt.Sprintf(gitProtoFmts[g.Scheme],
		g.Host, g.User, g.Repo,
	)
}

// regular expressions for matching package uris
const (
	gitSSHExp = `ssh://git@(?P<host>.+)/(?P<user>.+)/(?P<repo>.+).git`
	gitSCPExp = `^git@(?P<host>.+):(?P<user>.+)/(?P<repo>.+).git`
	// The long ugly pattern for ${host} here is a generic pattern for "valid URL with zero or more subdomains and a valid TLD"
	gitHTTPSSubgroup = `(?P<host>[a-zA-Z0-9][a-zA-Z0-9-\.]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,})/(?P<user>[-_a-zA-Z0-9/\.]+)/(?P<repo>[-_a-zA-Z0-9\.]+)\.git`
	gitHTTPSExp      = `(?P<host>[a-zA-Z0-9][a-zA-Z0-9-\.]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,})/(?P<user>[-_a-zA-Z0-9\.]+)/(?P<repo>[-_a-zA-Z0-9\.]+)`
)

var (
	VersionRegex        = `@(?P<version>.*)`
	PathRegex           = `/(?P<subdir>.*)`
	PathAndVersionRegex = `/(?P<subdir>.*)@(?P<version>.*)`
)

func parseGit(uri string) *Dependency {
	var d = Dependency{
		Version: "master",
		Source:  Source{},
	}
	var g *Git
	var version string

	switch {
	case reMatch(gitSSHExp, uri):
		g, version = match(uri, gitSSHExp)
		g.Scheme = GitSchemeSSH
	case reMatch(gitSCPExp, uri):
		g, version = match(uri, gitSCPExp)
		g.Scheme = GitSchemeSSH
	case reMatch(gitHTTPSSubgroup, uri):
		g, version = match(uri, gitHTTPSSubgroup)
		g.Scheme = GitSchemeHTTPS
	case reMatch(gitHTTPSExp, uri):
		g, version = match(uri, gitHTTPSExp)
		g.Scheme = GitSchemeHTTPS
	default:
		return nil
	}

	if g.Subdir != "" {
		g.Subdir = "/" + g.Subdir
	}

	d.Source.GitSource = g
	if version != "" {
		d.Version = version
	}
	return &d
}

func match(p string, exp string) (git *Git, version string) {
	git = &Git{}
	exps := []*regexp.Regexp{
		regexp.MustCompile(exp + PathAndVersionRegex),
		regexp.MustCompile(exp + PathRegex),
		regexp.MustCompile(exp + VersionRegex),
		regexp.MustCompile(exp),
	}

	for _, e := range exps {
		if !e.MatchString(p) {
			continue
		}

		matches := reSubMatchMap(e, p)
		git.Host = matches["host"]
		git.User = matches["user"]
		git.Repo = matches["repo"]

		if sd, ok := matches["subdir"]; ok {
			git.Subdir = sd
		}

		return git, matches["version"]
	}
	return git, ""
}

func reMatch(exp string, str string) bool {
	return regexp.MustCompile(exp).MatchString(str)
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}
