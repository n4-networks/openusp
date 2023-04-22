// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"context"
	"errors"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DmCommand struct {
	Name    string   `json:"name" bson:"name"`
	Inputs  []string `json:"inputs" bson:"inputs"`
	Outputs []string `json:"outputs" bson:"outputs"`
}
type DmEvent struct {
	Name string   `json:"name" bson:"name"`
	Args []string `json:"args" bson:"args"`
}
type DmParam struct {
	Name   string `json:"name" bson:"name"`
	Access string `json:"access" bson:"access"`
}

type DmObject struct {
	EndpointId    string      `json:"endpoint_id" bson:"endpoint_id"`
	Path          string      `json:"path" bson:"path"`
	MultiInstance bool        `json:"multi_instance" bson:"multi_instance"`
	Access        string      `json:"access" bson:"access"`
	Params        []DmParam   `json:"params" bson:"params"`
	Events        []DmEvent   `json:"events" bson:"events"`
	Cmds          []DmCommand `json:"cmds" bson:"cmds"`
}

func (u *UspDb) GetDmByRegex(epId string, path string) ([]*DmObject, error) {
	var elems bson.A
	if epId != "" {
		elems = append(elems, bson.D{{"endpoint_id", epId}})
	}
	if path != "" {
		elems = append(elems, bson.D{{"path", primitive.Regex{path, ""}}})
	}
	var filter bson.D
	if len(elems) == 0 {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"$and", elems}}
	}
	//log.Printf("Filter:%+v\n", filter)
	cur, err := u.dmColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}

	var objs []*DmObject
	if err := cur.All(context.Background(), &objs); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return objs, nil
}

func (u *UspDb) GetDm(epId string, path string) (*DmObject, error) {
	var elems bson.A
	if epId != "" {
		elems = append(elems, bson.D{{"endpoint_id", epId}})
	}
	if path != "" {
		elems = append(elems, bson.D{{"path", path}})
	}
	var filter bson.D
	if len(elems) == 0 {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"$and", elems}}
	}
	//log.Printf("Filter:%+v\n", filter)
	obj := &DmObject{}
	if err := u.dmColl.FindOne(context.Background(), filter).Decode(obj); err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return obj, nil
}

func (u *UspDb) WriteDmObjectToDb(dm *DmObject) error {
	dm.Path = strings.Replace(dm.Path, "{i}.", "", -1)
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", dm.EndpointId}},
			bson.D{{"path", dm.Path}},
		},
		}}
	//log.Println("Adding DM object wth Path ", dm.Path)
	opt := options.FindOneAndReplace().SetUpsert(true)
	err := u.dmColl.FindOneAndReplace(context.TODO(), filter, dm, opt).Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in creating datamodel object in DB")
		}
	}
	return nil
}

func (u *UspDb) DeleteDmObjectManyFromDb(epId string, path string) error {
	regex := primitive.Regex{}
	regex.Pattern = path

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", epId}},
			bson.D{{"path", regex}},
		},
		}}
	log.Printf("Deleting DM Object for agentId: %s and Path:%s\n", epId, path)
	result, err := u.dmColl.DeleteMany(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting datamodels from DB")
		}
	}
	log.Println("Deleted object count:", result.DeletedCount)
	return nil
}
