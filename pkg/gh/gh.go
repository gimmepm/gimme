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

	totalProcessed := 0
	allLatestReleases := map[*github.Repository]*github.RepositoryRelease{}
	go func() {
		for latestRelease := range latestReleases {
			for repo, release := range latestRelease {
				allLatestReleases[repo] = release
			}
			totalProcessed++
		}
	}()

	allErrs := []error{}
	go func() {
		for err := range errs {
			allErrs = append(allErrs, err)
			totalProcessed++
		}
	}()

	wg.Wait()
	for totalProcessed < len(starredRepos) {
		// We have to block until all repos are completely processed. It is not enough
		// to wait until all syncs are done, as there can be a race condition where all
		// operations in the wait group have been Done() but the consuming goroutine to
		// process hasn't looped through the item. This means that the channel can be at
		// zero (and therefore gets closed prematurely) under certain conditions.
		//
		// Scenario:
		//	The last item to process gets pushed into the channel, so WaitGroup.Wait() no
		//	longer blocks. And then the channel close() succeeds because we have looped
		//	through the last element, making len(releasesChannel) to zero and closing it
		//	resulting in a return of this current function without adding the last item
		//	to the map.
		//
		//	With having this additional check to see if we've processed all items (instead
		//	of just having zero items left in the channel) we can block on the right metric
		//	and workaround the race condition.
	}
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
			} else {
				latestRelease <- map[*github.Repository]*github.RepositoryRelease{repo: nil}
			}
			return
		case err := <-errsChan:
			errs <- err
			return
		case <-time.After(5 * time.Second):
		}
	}

	errs <- fmt.Errorf("Too many failed attempts for %s", repo.GetFullName())
}
