package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {

	// Create inputs
	var gh_URL, gh_Token, gh_rel, gl_URL, gl_Token string
	fmt.Println("Write your Github URL(http): ")
	fmt.Scanln(&gh_URL)
	fmt.Println("Write your Github Token: ")
	fmt.Scanln(&gh_Token)
	fmt.Println("Write Github releas name ('repo name' + '-' + 'releas name'): ")
	fmt.Scanln(&gh_rel)
	fmt.Println("Write your Gitlab URL(http): ")
	fmt.Scanln(&gl_URL)
	fmt.Println("Write your Gitlab Token: ")
	fmt.Scanln(&gl_Token)

	// Cloning from the Github

	_, err := git.PlainClone("./repo", false, &git.CloneOptions{
		URL: gh_URL,
		Auth: &githttp.BasicAuth{
			Username: "empty",
			Password: gh_Token,
		},
	})

	if err != nil {
		panic(err)
	}
	// Clonning from the Github (releases)
	token := gh_Token

	url := "https://api.github.com/repos/ayxank/test/zipball/assad"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("GitHub error: " + resp.Status)
	}

	out, err := os.Create("release.zip")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
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
			When: time.Now(),
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
		Auth: &githttp.BasicAuth{
			Username: "empty",
			Password: gl_Token,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully pushed.")

	// Push to the remote (releases)

	// Removing folders from the repository

	err = os.RemoveAll("./repo")
	if err != nil {
		log.Fatalf("Error removing directory: %v", err)
	}

}
