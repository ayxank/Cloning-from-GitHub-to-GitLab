package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {

	// Create inputs
	var gh_URL, gh_Token, gl_URL, gl_Token, gl_Usr, gl_Ml string
	fmt.Println("Write your Github URL(http): ")
	fmt.Scanln(&gh_URL)
	fmt.Println("Write your Github Token: ")
	fmt.Scanln(&gh_Token)
	fmt.Println("Write your Gitlab URL(http): ")
	fmt.Scanln(&gl_URL)
	fmt.Println("Write your Gitlab Token: ")
	fmt.Scanln(&gl_Token)
	fmt.Println("Write your Gitlab username: ")
	fmt.Scanln(&gl_Usr)
	fmt.Println("Write your Gitlab mail: ")
	fmt.Scanln(&gl_Ml)
	// Cloning from the Github

	_, err := git.PlainClone("./repo", false, &git.CloneOptions{
		URL: gh_URL,
		Auth: &http.BasicAuth{
			Username: "empty",
			Password: gh_Token,
		},
	})

	if err != nil {
		panic(err)
	}

	// Transfer to Gitlab
	// Open repo

	repo, err := git.PlainOpen("./repo")
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	// Get the worktree

	w, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Error getting worktree: %v", err)
	}

	// Add all changes to the staging area

	_, err = w.Add(".")
	if err != nil {
		log.Fatalf("Error adding changes: %v", err)
	}

	// Commit the changes

	_, err = w.Commit("My commit message", &git.CommitOptions{
		All:               true,
		AllowEmptyCommits: true,

		Author: &object.Signature{
			Name:  gl_Usr,
			Email: gl_Ml,
			When:  time.Now(),
		},
	})
	if err != nil {
		log.Fatalf("Error committing changes: %v", err)
	}

	// Update

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origi",
		URLs: []string{gl_URL},
	})

	if err != nil {
		log.Fatalf("Error creating remote: %v", err)
	}

	// Push to the remote

	err = repo.Push(&git.PushOptions{
		RemoteURL: gl_URL,
		Force:     true,
		Auth: &http.BasicAuth{
			Username: "empty",
			Password: gl_Token,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully pushed.")

	// Removing folders from the repository

	err = os.RemoveAll("./repo")
	if err != nil {
		log.Fatalf("Error removing directory: %v", err)
	}
}
