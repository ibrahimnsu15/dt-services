package models

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type GlobalStats struct {
	NumRepos     int           `json:"num_repos"`
	NumCommits   int           `json:"num_commits"`
	ProjectStats []ProjectStat `json:"project_stats"`
}

func ParseGlobalStats() (*GlobalStats, error) {
	gs := new(GlobalStats)
	projects, err := GetProjects()
	if err != nil {
		return nil, err
	}

	var projectsStats []ProjectStat
	for _, project := range projects {
		projectStat, err := ParseProjectStat(&project, nil, nil)
		if err != nil {
			return nil, err
		}

		projectsStats = append(projectsStats, *projectStat)
	}

	gs.ProjectStats = projectsStats
	numCommits := 0
	for _, ps := range gs.ProjectStats {
		numCommits = numCommits + ps.NumCommits
	}
	gs.NumCommits = numCommits

	return gs, nil
}

type ProjectStat struct {
	ProjectId  bson.ObjectId `json:"project_id"`
	NumRepos   int           `json:"num_repos"`
	NumCommits int           `json:"num_commits"`
	RepoStats  []RepoStat    `json:"repo_stats"`
}

func ParseProjectStat(project *Project, start, end *time.Time) (*ProjectStat, error) {
	ps := new(ProjectStat)
	project, err := GetProjectById(project.Id.Hex(), true)
	if err != nil {
		return nil, err
	}

	var repoStats []RepoStat
	for _, repo := range project.Repos {
		repoStat, err := ParseRepoStat(&repo, start, end)
		if err != nil {
			return nil, err
		}

		repoStats = append(repoStats, *repoStat)
	}

	ps.RepoStats = repoStats
	ps.NumRepos = len(repoStats)
	numCommits := 0
	for _, rs := range ps.RepoStats {
		numCommits = numCommits + rs.NumCommits
	}
	ps.NumCommits = numCommits
	ps.ProjectId = project.Id

	return ps, nil
}

type RepoStat struct {
	RepoId      bson.ObjectId `json:"repo_id"`
	Name        string        `json:"name"`
	NumCommits  int           `json:"num_commits"`
	CommitStats []CommitStat  `json:"commit_stats"`
}

func ParseRepoStat(repo *Repo, start, end *time.Time) (*RepoStat, error) {
	commits, err := GetCommitsByRepo(repo, start, end)
	if err != nil {
		return nil, err
	}

	var commitStats []CommitStat
	for _, commit := range commits {
		commitStat, err := ParseCommitStat(&commit)
		if err != nil {
			return nil, err
		}

		commitStats = append(commitStats, *commitStat)
	}

	rs := new(RepoStat)
	rs.Name = repo.Name
	rs.CommitStats = commitStats
	rs.NumCommits = len(commitStats)
	rs.RepoId = repo.Id

	return rs, nil
}

type CommitStat struct {
	NumFiles    int        `json:"num_files"`
	Username    string     `json:"username"`
	SHA         string     `json:"sha"`
	CommittedAt time.Time  `json:"commited_at"`
	FileStats   []FileStat `json:"file_stats"`
}

func ParseCommitStat(commit *Commit) (*CommitStat, error) {
	cs := new(CommitStat)
	cs.Username = commit.Username
	cs.SHA = commit.Sha
	cs.CommittedAt = commit.CommittedAt

	return cs, nil
}

type FileStat struct {
	Additions    int `json:"additions"`
	Subtractions int `json:"subtractions"`
	Changes      int `json:"changes"`
}
