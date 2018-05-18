package gh

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ListStarredRepos will fetch and return all of the starred repos by a user
func ListStarredRepos(token string) ([]string, error) {
	starredRepos := []string{}
	operationErrors := []error{}
	reposPerPage := 50

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(oauthClient)

	firstPage := 1
	repos, res, err := ghClient.Activity.ListStarred(
		context.Background(),
		"",
		&github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page:    firstPage,
				PerPage: reposPerPage,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		starredRepos = append(starredRepos, repo.GetRepository().GetFullName())
	}

	if firstPage == res.LastPage {
		return starredRepos, nil
	}

	wg := sync.WaitGroup{}
	allRepos := make(chan []string)
	errs := make(chan error)
	for i := firstPage + 1; i <= res.LastPage; i++ {
		wg.Add(1)
		go getStarredReposByPage(
			&wg,
			ghClient,
			i,
			reposPerPage,
			allRepos,
		)
	}

	// Handle successful return of starred repositories
	go func() {
		for pageRepos := range allRepos {
			starredRepos = append(starredRepos, pageRepos...)
		}
	}()

	// Handle any errors returned from the goroutines
	go func() {
		for err := range errs {
			operationErrors = append(operationErrors, err)
		}
	}()

	wg.Wait()
	close(allRepos)
	close(errs)

	if len(operationErrors) > 0 {
		// If there are any errors, fail out returning the first error
		//
		// This is because they are likely related, and the assumption
		// is that it is ok to assume this here
		return nil, operationErrors[0]
	}

	return starredRepos, nil
}

func getStarredReposByPage(wg *sync.WaitGroup, client *github.Client, pageNumber, reposPerPage int, reposByPage chan<- []string) {
	defer wg.Done()

	starredReposCurrentPage := []string{}

	repos, _, err := client.Activity.ListStarred(
		context.Background(),
		"",
		&github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				Page:    pageNumber,
				PerPage: reposPerPage,
			},
		},
	)
	if err != nil {
		fmt.Printf("[page %d] err is not nil\n", pageNumber)
	}

	for _, repo := range repos {
		starredReposCurrentPage = append(starredReposCurrentPage, repo.GetRepository().GetFullName())
	}

	reposByPage <- starredReposCurrentPage
}
