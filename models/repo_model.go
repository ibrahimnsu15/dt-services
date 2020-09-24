package models

import (
	"time"

	"dt-services/stores"

	"github.com/globalsign/mgo/bson"
)

type Repo struct {
	Id        bson.ObjectId   `json:"id" bson:"_id"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
	GithubId  int             `json:"github_id" bson:"github_id"`
	Name      string          `json:"name" bson:"name"`
	ProjectId *bson.ObjectId  `json:"project_id,omitempty" bson:"project_id"`
	CommitIds []bson.ObjectId `json:"commit_ids" bson:"commit_ids"`
	Commits   []Commit        `json:"commits"`
}

func NewRepo(githubId int, name string) *Repo {
	return &Repo{
		Id:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		GithubId:  githubId,
		Name:      name,
	}
}

func (r *Repo) Save() error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_REPOS)
	if err := collection.Insert(r); err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update() error {
	r.UpdatedAt = time.Now()

	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_REPOS)
	if err := collection.UpdateId(r.Id, r); err != nil {
		return err
	}

	return nil
}

func GetRepos() ([]Repo, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var repos []Repo
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_REPOS)
	if err := collection.Find(nil).All(&repos); err != nil {
		return nil, err
	}

	//sort.Slice(jobs, func(i, j int) bool {
	//	return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	//})

	return repos, nil
}

func GetReposByIds(ids []bson.ObjectId, withCommits bool) ([]Repo, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var repos []Repo
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_REPOS)
	query := bson.M{
		"_id": bson.M{"$in": ids},
	}
	if err := collection.Find(query).All(&repos); err != nil {
		return nil, err
	}

	if withCommits {
		for _, r := range repos {
			commits, err := GetCommitsByRepo(&r, nil, nil)
			if err != nil {
				return nil, err
			}

			r.Commits = commits
		}
	}

	return repos, nil
}

func GetRepoById(id string, withCommits bool) (*Repo, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var repo *Repo
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_REPOS)
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&repo); err != nil {
		return nil, err
	}

	// todo: waitgroup

	if withCommits {
		for _, commitId := range repo.CommitIds {
			commit, err := GetCommitById(commitId.Hex(), false)
			if err != nil {
				return nil, err
			}

			repo.Commits = append(repo.Commits, *commit)
		}
	}

	return repo, nil
}
