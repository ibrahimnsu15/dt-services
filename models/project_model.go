package models

import (
	"time"

	"dt-services/stores"

	"github.com/globalsign/mgo/bson"
)

type Project struct {
	Id        bson.ObjectId   `json:"id" bson:"_id"`
	CreatedAt time.Time       `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" bson:"updated_at"`
	Name      string          `json:"name" bson:"name"`
	RepoIds   []bson.ObjectId `json:"repo_ids" bson:"repo_ids"`
	Repos     []Repo          `json:"repos"`
}

func NewProject(name string) *Project {
	return &Project{
		Id:        bson.NewObjectId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}
}

func (p *Project) Save() error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	p.Repos = nil
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_PROJECTS)
	if err := collection.Insert(p); err != nil {
		return err
	}

	return nil
}

func (p *Project) Update() error {
	p.UpdatedAt = time.Now()

	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	p.Repos = nil
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_PROJECTS)
	if err := collection.UpdateId(p.Id, p); err != nil {
		return err
	}

	return nil
}

func GetProjects() ([]Project, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var projects []Project
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_PROJECTS)
	if err := collection.Find(nil).All(&projects); err != nil {
		return nil, err
	}

	//sort.Slice(jobs, func(i, j int) bool {
	//	return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	//})

	return projects, nil
}

func GetProjectById(id string, withRepos bool) (*Project, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var project *Project
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_PROJECTS)
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&project); err != nil {
		return nil, err
	}

	// todo: waitgroup

	if withRepos {
		//for _, repoId := range project.RepoIds {
		//	repo, err := GetRepoById(repoId.Hex(), false)
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	project.Repos = append(project.Repos, *repo)
		//}
		repos, err := GetReposByIds(project.RepoIds, true)
		if err != nil {
			return nil, err
		}

		project.Repos = repos
	}

	return project, nil
}

func DeleteProjectById(id string) error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_PROJECTS)
	if err := collection.Remove(bson.M{"_id": bson.ObjectIdHex(id)}); err != nil {
		return err
	}

	return nil
}
