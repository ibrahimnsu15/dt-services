package services

import (
	"dt-services/models"

	"github.com/labstack/echo"
)

func Ingest(c echo.Context) error {
	if err := ingestRepos(); err != nil {
		return err
	}

	repos, err := models.GetRepos()
	if err != nil {
		return nil
	}

	for _, repo := range repos {
		if err := ingestRepoCommits(&repo); err != nil {
			return err
		}
	}

	repos, err = models.GetRepos()
	if err != nil {
		return err
	}

	for _, repo := range repos {
		commits, err := models.GetCommitsByRepo(&repo, nil, nil)
		if err != nil {
			return err
		}

		for _, commit := range commits {
			if err := ingestCommitFiles(&repo, &commit); err != nil {
				return err
			}
		}
	}

	return nil
}

func ingestRepos() error {
	ghRepos, err := fetchRepos(GITHUB_ORG)
	if err != nil {
		return err
	}

	repos, err := models.GetRepos()
	if err != nil {
		return err
	}

	for _, ghRepo := range ghRepos {
		found := false
		for _, repo := range repos {
			if int(*ghRepo.ID) == repo.GithubId {
				found = true
				break
			}
		}

		if !found {
			repo := models.NewRepo(int(*ghRepo.ID), *ghRepo.Name)
			if err := repo.Save(); err != nil {
				return err
			}
		}
	}

	// todo: handle repo deleted

	return nil
}

func ingestRepoCommits(repo *models.Repo) error {
	commits, err := models.GetCommitsByRepo(repo, nil, nil)
	if err != nil {
		return err
	}

	ghCommits, err := fetchRepoCommitsByName(GITHUB_ORG, repo.Name)
	if err != nil {
		return err
	}

	for _, ghCommit := range ghCommits {
		found := false
		for _, commit := range commits {
			if *ghCommit.SHA == commit.Sha {
				found = true
				break
			}
		}

		if !found {
			commit := models.NewCommit()

			if ghCommit.SHA != nil {
				commit.Sha = *ghCommit.SHA
			}

			if ghCommit.Author != nil {
				if ghCommit.Commit.Author.Date != nil {
					commit.CommittedAt = *ghCommit.Commit.Author.Date
				}
				if ghCommit.Author.Name != nil {
					commit.Username = *ghCommit.Author.Name
				}
			}

			if err := commit.Save(); err != nil {
				return err
			}

			repo.CommitIds = append(repo.CommitIds, commit.Id)
			if err := repo.Update(); err != nil {
				return err
			}
		}
	}

	// todo: handle commit deleted

	return nil
}

func ingestCommitFiles(repo *models.Repo, commit *models.Commit) error {
	files, err := models.GetFilesByCommit(commit)
	if err != nil {
		return err
	}

	ghCommit, err := fetchRepoCommitBySha(GITHUB_ORG, repo.Name, commit.Sha)
	if err != nil {
		return err
	}

	if ghCommit != nil {
		for _, ghFile := range ghCommit.Files {
			found := false
			for _, file := range files {
				if ghFile.SHA != nil {
					if *ghFile.SHA == file.Sha {
						found = true
						break
					}
				}
			}

			if !found {
				file := models.NewFile()
				if ghFile.Filename != nil {
					file.Path = *ghFile.Filename
				}
				if ghFile.SHA != nil {
					file.Sha = *ghFile.SHA
				}
				if ghFile.Additions != nil {
					file.Additions = *ghFile.Additions
				}
				if ghFile.Deletions != nil {
					file.Subtractions = *ghFile.Deletions
				}
				if ghFile.Changes != nil {
					file.Changes = *ghFile.Changes
				}

				if err := file.Save(); err != nil {
					return err
				}

				commit.FileIds = append(commit.FileIds, file.Id)
				if err := commit.Update(); err != nil {
					return err
				}
			}
		}
	}

	// todo: handle file deleted maybe?

	return nil
}
