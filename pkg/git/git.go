package git

import (
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Repo struct {
	Url    string `json:"url"`
	Branch string `json:"branch"`
	Commit string `json:"commit"`
}

type Provider struct {
	Name     string `json:"name"`
	Token    string `json:"token"`
	Username string `json:"username"`
}

func CloneRepo(repo Repo, provider Provider, workspacePath string) error {
	// Clone the given repository to the given directory
	_, err := git.PlainClone(workspacePath, false, &git.CloneOptions{
		URL:        repo.Url,
		Progress:   os.Stdout,
		RemoteName: "origin/" + repo.Branch,
		Depth:      1,
		Auth: &http.BasicAuth{
			Username: provider.Username,
			Password: provider.Token,
		},
	})

	return err
}
