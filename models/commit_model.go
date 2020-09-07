package models

import (
	"sort"
	"time"

	"dt-services/stores"

	"github.com/globalsign/mgo/bson"
)

type Commit struct {
	Id          bson.ObjectId   `json:"id" bson:"_id"`
	CreatedAt   time.Time       `json:"created_at" bson:"created_at"`
	Sha         string          `json:"sha" bson:"sha"`
	CommittedAt time.Time       `json:"committed_at" bson:"committed_at"`
	Username    string          `json:"username" bson:"username"`
	FileIds     []bson.ObjectId `json:"file_ids" bson:"file_ids"`
	Files       []File          `json:"files"`
}

func NewCommit() *Commit {
	return &Commit{
		Id:        bson.NewObjectId(),
		CreatedAt: time.Now(),
	}
}

func (c *Commit) Save() error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_COMMITS)
	if err := collection.Insert(c); err != nil {
		return err
	}

	return nil
}

func (c *Commit) Update() error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_COMMITS)
	if err := collection.UpdateId(c.Id, c); err != nil {
		return err
	}

	return nil
}

func GetCommits() ([]Commit, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var commits []Commit
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_COMMITS)
	if err := collection.Find(nil).All(&commits); err != nil {
		return nil, err
	}

	//sort.Slice(jobs, func(i, j int) bool {
	//	return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	//})

	return commits, nil
}

func GetCommitsByRepo(repo *Repo, start, end *time.Time) ([]Commit, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var query bson.M
	if start != nil && end != nil {
		query = bson.M{
			"$and": []bson.M{
				{"_id": bson.M{"$in": repo.CommitIds}},
				{"committed_at": bson.M{"$gt": start}},
				{"committed_at": bson.M{"$lt": end}},
			},
		}
	} else {
		query = bson.M{
			"_id": bson.M{"$in": repo.CommitIds},
		}
	}

	var commits []Commit
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_COMMITS)
	if err := collection.Find(query).All(&commits); err != nil {
		return nil, err
	}

	//for _, commitId := range repo.CommitIds {
	//	commit, err := GetCommitById(commitId.Hex(), false)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	commits = append(commits, *commit)
	//}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].CommittedAt.Before(commits[j].CommittedAt)
	})

	return commits, nil
}

func GetCommitById(id string, withFiles bool) (*Commit, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var commit *Commit
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_COMMITS)
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&commit); err != nil {
		return nil, err
	}

	// todo: waitgroup

	if withFiles {
		for _, fileId := range commit.FileIds {
			file, err := GetFileById(fileId.Hex())
			if err != nil {
				return nil, err
			}

			commit.Files = append(commit.Files, *file)
		}
	}

	return commit, nil
}
