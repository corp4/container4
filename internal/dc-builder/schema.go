package builder

import "github.com/corp4/container4/pkg/git"

type Workspace struct {
	Name      string `json:"name"`
	UserId    string `json:"user_id"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type Task struct {
	TaskName  string       `json:"name"`
	Workspace Workspace    `json:"workspace"`
	Provider  git.Provider `json:"provider"`
	Repo      git.Repo     `json:"repo"`
}
