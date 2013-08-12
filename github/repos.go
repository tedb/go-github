// Copyright 2013 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"fmt"
	"net/url"
	"strconv"
)

// RepositoriesService handles communication with the repository related
// methods of the GitHub API.
//
// GitHub API docs: http://developer.github.com/v3/repos/
type RepositoriesService struct {
	client *Client
}

// Repository represents a GitHub repository.
type Repository struct {
	ID          int        `json:"id,omitempty"`
	Owner       *User      `json:"owner,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *Timestamp `json:"created_at,omitempty"`
	PushedAt    *Timestamp `json:"pushed_at,omitempty"`
	UpdatedAt   *Timestamp `json:"updated_at,omitempty"`

	// Additional mutable fields when creating and editing a repository
	HasIssues *bool `json:"has_issues"`
	HasWiki   *bool `json:"has_wiki"`
}

// RepositoryListOptions specifies the optional parameters to the
// RepositoriesService.List method.
type RepositoryListOptions struct {
	// Type of repositories to list.  Possible values are: all, owner, public,
	// private, member.  Default is "all".
	Type string

	// How to sort the repository list.  Possible values are: created, updated,
	// pushed, full_name.  Default is "full_name".
	Sort string

	// Direction in which to sort repositories.  Possible values are: asc, desc.
	// Default is "asc" when sort is "full_name", otherwise default is "desc".
	Direction string

	// For paginated result sets, page of results to retrieve.
	Page int
}

// List the repositories for a user.  Passing the empty string will list
// repositories for the authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/repos/#list-user-repositories
func (s *RepositoriesService) List(user string, opt *RepositoryListOptions) ([]Repository, *Response, error) {
	var u string
	if user != "" {
		u = fmt.Sprintf("users/%v/repos", user)
	} else {
		u = "user/repos"
	}
	if opt != nil {
		params := url.Values{
			"type":      []string{opt.Type},
			"sort":      []string{opt.Sort},
			"direction": []string{opt.Direction},
			"page":      []string{strconv.Itoa(opt.Page)},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repos := new([]Repository)
	resp, err := s.client.Do(req, repos)
	return *repos, resp, err
}

// RepositoryListByOrgOptions specifies the optional parameters to the
// RepositoriesService.ListByOrg method.
type RepositoryListByOrgOptions struct {
	// Type of repositories to list.  Possible values are: all, public, private,
	// forks, sources, member.  Default is "all".
	Type string

	// For paginated result sets, page of results to retrieve.
	Page int
}

// ListByOrg lists the repositories for an organization.
//
// GitHub API docs: http://developer.github.com/v3/repos/#list-organization-repositories
func (s *RepositoriesService) ListByOrg(org string, opt *RepositoryListByOrgOptions) ([]Repository, *Response, error) {
	u := fmt.Sprintf("orgs/%v/repos", org)
	if opt != nil {
		params := url.Values{
			"type": []string{opt.Type},
			"page": []string{strconv.Itoa(opt.Page)},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repos := new([]Repository)
	resp, err := s.client.Do(req, repos)
	return *repos, resp, err
}

// RepositoryListAllOptions specifies the optional parameters to the
// RepositoriesService.ListAll method.
type RepositoryListAllOptions struct {
	// ID of the last repository seen
	Since int
}

// ListAll lists all GitHub repositories in the order that they were created.
//
// GitHub API docs: http://developer.github.com/v3/repos/#list-all-repositories
func (s *RepositoriesService) ListAll(opt *RepositoryListAllOptions) ([]Repository, *Response, error) {
	u := "repositories"
	if opt != nil {
		params := url.Values{
			"since": []string{strconv.Itoa(opt.Since)},
		}
		u += "?" + params.Encode()
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repos := new([]Repository)
	resp, err := s.client.Do(req, repos)
	return *repos, resp, err
}

// Create a new repository.  If an organization is specified, the new
// repository will be created under that org.  If the empty string is
// specified, it will be created for the authenticated user.
//
// GitHub API docs: http://developer.github.com/v3/repos/#create
func (s *RepositoriesService) Create(org string, repo *Repository) (*Repository, *Response, error) {
	var u string
	if org != "" {
		u = fmt.Sprintf("orgs/%v/repos", org)
	} else {
		u = "user/repos"
	}

	req, err := s.client.NewRequest("POST", u, repo)
	if err != nil {
		return nil, nil, err
	}

	r := new(Repository)
	resp, err := s.client.Do(req, r)
	return r, resp, err
}

// Get fetches a repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/#get
func (s *RepositoriesService) Get(owner, repo string) (*Repository, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v", owner, repo)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	repository := new(Repository)
	resp, err := s.client.Do(req, repository)
	return repository, resp, err
}

// Edit updates a repository.
//
// GitHub API docs: http://developer.github.com/v3/repos/#edit
func (s *RepositoriesService) Edit(owner, repo string, repository *Repository) (*Repository, *Response, error) {
	u := fmt.Sprintf("repos/%v/%v", owner, repo)
	req, err := s.client.NewRequest("PATCH", u, repository)
	if err != nil {
		return nil, nil, err
	}
	r := new(Repository)
	resp, err := s.client.Do(req, r)
	return r, resp, err
}

// ListLanguages lists languages for the specified repository. The returned map
// specifies the languages and the number of bytes of code written in that
// language. For example:
//
//     {
//       "C": 78769,
//       "Python": 7769
//     }
//
// GitHub API Docs: http://developer.github.com/v3/repos/#list-languages
func (s *RepositoriesService) ListLanguages(owner string, repository string) (map[string]int, *Response, error) {
	u := fmt.Sprintf("/repos/%v/%v/languages", owner, repository)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	languages := make(map[string]int)
	resp, err := s.client.Do(req, &languages)
	return languages, resp, err
}

// ListTeams lists teams for the specified repository. Returned is a list of Team structs
// GitHub API Docs: http://developer.github.com/v3/repos/#list-teams
func (s *RepositoriesService) ListTeams(owner string, repository string) ([]Team, *Response, error) {
	u := fmt.Sprintf("/repos/%v/%v/teams", owner, repository)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	teams := new([]Team)
	resp, err := s.client.Do(req, teams)
	return *teams, resp, err
}

// this ListTeams method can be invoked against a repository struct
func (repo Repository) ListTeams(s *RepositoriesService) ([]Team, *Response, error) {
	return s.ListTeams(repo.Owner.Name, repo.Name)
}
