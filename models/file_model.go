package models

import (
	"dt-services/stores"
	"github.com/globalsign/mgo/bson"
	"time"
)

type File struct {
	Id           bson.ObjectId `json:"id" bson:"_id"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	Path         string        `json:"path" bson:"path"`
	Sha          string        `json:"sha" bson:"sha"`
	Additions    int           `json:"additions" bson:"additions"`
	Subtractions int           `json:"subtractions" bson:"subtractions"`
	Changes      int           `json:"changes" bson:"changes"`
}

func NewFile() *File {
	return &File{
		Id:        bson.NewObjectId(),
		CreatedAt: time.Now(),
	}
}

func (f *File) Save() error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_FILES)
	if err := collection.Insert(f); err != nil {
		return err
	}

	return nil
}

func (f *File) Update() error {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_FILES)
	if err := collection.UpdateId(f.Id, f); err != nil {
		return err
	}

	return nil
}

func GetFiles() ([]File, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var files []File
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_FILES)
	if err := collection.Find(nil).All(&files); err != nil {
		return nil, err
	}

	//sort.Slice(jobs, func(i, j int) bool {
	//	return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	//})

	return files, nil
}

func GetFilesByCommit(commit *Commit) ([]File, error) {
	//session := stores.DB.Mongo.Session.Clone()
	//defer session.Close()

	var files []File

	//collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_FILES)
	//query := bson.M{
	//	"_id": bson.M{
	//		"in:": commit.FileIds,
	//	},
	//}
	//if err := collection.Find(query).All(&files); err != nil {
	//	return nil, err
	//}

	for _, fileId := range commit.FileIds {
		file, err := GetFileById(fileId.Hex())
		if err != nil {
			return nil, err
		}

		files = append(files, *file)
	}

	//sort.Slice(jobs, func(i, j int) bool {
	//	return jobs[i].CreatedAt.After(jobs[j].CreatedAt)
	//})

	return files, nil
}

func GetFileById(id string) (*File, error) {
	session := stores.DB.Mongo.Session.Clone()
	defer session.Close()

	var file *File
	collection := session.DB(stores.DB_NAME).C(stores.DB_COLLECTION_FILES)
	if err := collection.FindId(bson.ObjectIdHex(id)).One(&file); err != nil {
		return nil, err
	}

	return file, nil
}
