package main

import (
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
	bitbucket "github.com/gfleury/go-bitbucket-v1"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

var args struct {
	Debug        *bool  `arg:"-D"`
	Username     string `arg:"-u,--username,env:BITBUCKET_USERNAME"`
	Password     string `arg:"-p,--password,env:BITBUCKET_PASSWORD"`
	Endpoint     string `arg:"-e,--endpoint,env:BITBUCKET_ENDPOINT"`
	AuthorFilter string `arg:"-a,--author-filter,env:BITBUCKET_AUTHOR_FILTER"`
	AddComment   *bool  `arg:"-c,--add-comment,env:BITBUCKET_ADD_COMMENT" default:"true" help:"\"true\" to add a comment in addition to approving a PR, \"false\" to not add a comment."`
	Version      *bool  `arg:"-v"`
}

var version = "unknown"

type Client struct {
	bitbucketClient *bitbucket.APIClient
	logger          *logrus.Logger
}

func main() {
	logger := logrus.New()
	arg.MustParse(&args)

	if args.Debug != nil && *args.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	if args.Version != nil && *args.Version {
		fmt.Printf("approve-bot %s\n", version)
		return
	}

	if args.Endpoint == "" {
		logger.Fatalf("provide and endpoint")
	}

	if args.Username == "" {
		logger.Fatalf("username cannot be empty")
	}

	if args.Password == "" {
		logger.Fatalf("password cannot be empty")
	}

	_, err := url.Parse(args.Endpoint)
	if err != nil {
		logger.Fatalf("invalid BitBucket endpoint provided: %v", err)
	}

	basicAuth := bitbucket.BasicAuth{UserName: args.Username, Password: args.Password}
	ctx := context.WithValue(context.Background(), bitbucket.ContextBasicAuth, basicAuth)
	bitbucketClient := bitbucket.NewAPIClient(ctx, bitbucket.NewConfiguration(args.Endpoint))

	c := Client{
		bitbucketClient: bitbucketClient,
		logger:          logger,
	}

	// https://docs.atlassian.com/bitbucket-server/rest/7.20.0/bitbucket-rest.html#idp96

	prsChannel := make(chan bitbucket.PullRequest)
	start := 0
	size := 50
	go fetchAllPrs(c, start, size, logger, prsChannel)
	for v := range prsChannel {
		if args.AuthorFilter != "" {
			if v.Author.User.Slug != args.AuthorFilter {
				logger.Infof("skipping %d because %s != %s", v.ID, v.Author.User.Slug, args.AuthorFilter)
			}
		}

		c.logger.Infof("Auto-approving %s - %s - %s",
			v.Title,
			v.Author.User.DisplayName,
			v.Links.Self[0].Href,
		)

		if args.AddComment != nil && *args.AddComment {
			err = c.addComment(&v)
			if err != nil {
				logger.Warnf("unable to add comment: %v. Skipping PR.", err)
				continue
			}
		}
		err = c.approvePr(&v)
		if err != nil {
			logger.Warnf("unable to approve PR: %v", err)
		}
	}
}

func (c *Client) addComment(v *bitbucket.PullRequest) error {
	res, err := c.bitbucketClient.DefaultApi.CreatePullRequestComment(
		v.ToRef.Repository.Project.Key,
		v.ToRef.Repository.Slug,
		v.ID,
		bitbucket.Comment{
			Text: "Auto approving this PR - please make sure the tests are passing!",
		},
		[]string{"application/json"},
	)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("invalid status code %d: 201 expected", res.StatusCode)
	}
	return nil
}

func (c *Client) approvePr(v *bitbucket.PullRequest) error {
	res, err := c.bitbucketClient.DefaultApi.Approve(
		v.ToRef.Repository.Project.Key,
		v.ToRef.Repository.Slug,
		int64(v.ID),
	)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code %d: 200 expected", res.StatusCode)
	}
	return nil
}

func fetchAllPrs(c Client, start int, size int, logger *logrus.Logger, prsChannel chan bitbucket.PullRequest) {
	for {
		prs, isLastPage, nextPageStart, err := c.GetPRs(start, map[string]interface{}{
			"state":             "OPEN",
			"role":              "REVIEWER",
			"participantStatus": "UNAPPROVED",
			"limit":             size,
		})

		if err != nil {
			logger.Fatalf("unable to get prs: %v", err)
		}

		for _, v := range prs {
			prsChannel <- v
		}

		if isLastPage {
			break
		}

		start = nextPageStart
	}
	close(prsChannel)
}

func (c *Client) GetPRs(start int, params map[string]interface{}) (prs []bitbucket.PullRequest, isLastPage bool, nextPageStart int, err error) {
	if start != 0 {
		params["start"] = start
	}
	res, err := c.bitbucketClient.DefaultApi.GetPullRequests(params)

	if err != nil {
		c.logger.Fatalf("unable to get pull requests: %v", err)
	}

	isLastPage, nextPageStart = bitbucket.HasNextPage(res)

	prs, err = bitbucket.GetPullRequestsResponse(res)
	if err != nil {
		c.logger.Fatalf("unable to parse pull requests: %v", err)
	}
	return prs, isLastPage, nextPageStart, nil
}
