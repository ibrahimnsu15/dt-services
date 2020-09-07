package models

import (
	"time"

	"dt-services/stores"

	"github.com/globalsign/mgo/bson"
)

type Developer struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
	Email      string        `json:"email" bson:"email"`
	GithubUser string        `json:"github_user" bson:"github_user"`
}

func NewDeveloper(email, githubUser string) *Developer {
	return &Developer{
		Id:         bson.NewObjectId(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		GithubUser: githubUser,
		Email:      email,
	}
}

func (d *Developer) Save() error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_DEVELOPERS)
	if err := collection.Insert(d); err != nil {
		return err
	}

	return nil
}

func (d *Developer) Update() error {
	d.UpdatedAt = time.Now()

	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_DEVELOPERS)
	if err := collection.UpdateId(d.Id, d); err != nil {
		return err
	}

	return nil
}

func GetDevelopers() ([]Developer, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var developers []Developer
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_DEVELOPERS)
	if err := collection.Find(nil).All(&developers); err != nil {
		return nil, err
	}

	//sort.Slice(jobs, func(i, j int) bool {
	//	return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	//})

	return developers, nil
}

func GetDeveloperById(id string) (*Developer, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var developer *Developer
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_DEVELOPERS)
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&developer); err != nil {
		return nil, err
	}

	return developer, nil
}
