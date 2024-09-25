package change

import "testing"

func TestParseGithubURL(t *testing.T) {
	urls := []struct {
		url   string
		owner string
		repo  string
	}{
		{"https://github.com/plum330/kratos.git", "plum330", "kratos"},
		{"https://github.com/plum330/kratos", "plum330", "kratos"},
		{"git@github.com:plum330/kratos.git", "plum330", "kratos"},
		{"https://github.com/plum330/go-kratos.dev.git", "plum330", "go-kratos.dev"},
	}
	for _, url := range urls {
		owner, repo := ParseGithubURL(url.url)
		if owner != url.owner {
			t.Fatalf("owner want: %s, got: %s", owner, url.owner)
		}
		if repo != url.repo {
			t.Fatalf("repo want: %s, got: %s", repo, url.repo)
		}
	}
}
