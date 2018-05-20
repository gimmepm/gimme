package gh

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ListStarredRepos will fetch and return all of the starred repos by a user
func ListStarredRepos(token string) ([]*github.StarredRepository, error) {
	starredRepos := []*github.StarredRepository{}
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
		starredRepos = append(starredRepos, repo)
	}

	if firstPage == res.LastPage {
		return starredRepos, nil
	}

	wg := sync.WaitGroup{}
	allRepos := make(chan []*github.StarredRepository)
	errs := make(chan error)
	for i := firstPage + 1; i <= res.LastPage; i++ {
		wg.Add(1)
		go getStarredReposByPage(
			&wg,
			ghClient,
			i,
			reposPerPage,
			allRepos,
			errs,
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

func getStarredReposByPage(wg *sync.WaitGroup, client *github.Client, pageNumber, reposPerPage int, reposByPage chan<- []*github.StarredRepository, errs chan<- error) {
	defer wg.Done()

	starredReposCurrentPage := []*github.StarredRepository{}

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
		errs <- err
		return
	}

	for _, repo := range repos {
		starredReposCurrentPage = append(starredReposCurrentPage, repo)
	}

	reposByPage <- starredReposCurrentPage
}

// ListStarredReposLatestReleases fetches all latest releases for starred repos by a user
func ListStarredReposLatestReleases(token string) (map[*github.Repository]*github.RepositoryRelease, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(oauthClient)

	starredRepos, err := ListStarredRepos(token)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	latestReleases := make(chan map[*github.Repository]*github.RepositoryRelease)
	errs := make(chan error)
	for _, starredRepo := range starredRepos {
		wg.Add(1)
		go getLatestReleaseForRepo(
			&wg,
			ghClient,
			starredRepo.GetRepository(),
			latestReleases,
			errs,
		)
	}

	allLatestReleases := map[*github.Repository]*github.RepositoryRelease{}
	go func() {
		for latestRelease := range latestReleases {
			for repo, release := range latestRelease {
				allLatestReleases[repo] = release
			}
		}
	}()

	allErrs := []error{}
	go func() {
		for err := range errs {
			allErrs = append(allErrs, err)
		}
	}()

	wg.Wait()
	close(latestReleases)
	close(errs)

	if len(allErrs) > 0 {
		return nil, allErrs[0]
	}

	return allLatestReleases, nil
}

func getLatestReleaseForRepo(wg *sync.WaitGroup, client *github.Client, repo *github.Repository, latestRelease chan<- map[*github.Repository]*github.RepositoryRelease, errs chan<- error) {
	defer wg.Done()

	releasesChan := make(chan []*github.RepositoryRelease)
	errsChan := make(chan error)
	for i := 0; i < 5; i++ {
		go func() {
			releases, _, err := client.Repositories.ListReleases(
				context.Background(),
				repo.GetOwner().GetLogin(),
				repo.GetName(),
				&github.ListOptions{},
			)
			if err != nil {
				errsChan <- err
				return
			}
			releasesChan <- releases
		}()

		select {
		case releases := <-releasesChan:
			if len(releases) > 0 {
				latestRelease <- map[*github.Repository]*github.RepositoryRelease{repo: releases[0]}
			}
			return
		case err := <-errsChan:
			errs <- err
			return
		case <-time.After(5 * time.Second):
		}
	}

	errs <- fmt.Errorf("Too many failed attempts for %s\n", repo.GetFullName())
}
