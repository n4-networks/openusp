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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Instance struct {
	EndpointId string            `json:"endpoint_id" bson:"endpoint_id"`
	Path       string            `json:"path" bson:"path"`
	UniqueKeys map[string]string `json:"unique_keys" bson:"unique_keys"`
}

func (u *UspDb) GetInstances(epId string, path string) ([]*Instance, error) {

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
	cur, err := u.instanceColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}

	var instances []*Instance
	if err := cur.All(context.Background(), &instances); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return instances, nil
}

func (u *UspDb) GetInstancesByRegex(epId string, path string) ([]*Instance, error) {

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
	cur, err := u.instanceColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}

	var instances []*Instance
	if err := cur.All(context.Background(), &instances); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return instances, nil
}

func (u *UspDb) GetInstancesByUniqueKeys(epId string, key string, value string) (*Instance, error) {

	var elems bson.A

	if epId == "" {
		return nil, errors.New("Empty endpoint_id")
	}
	if key == "" {
		return nil, errors.New("Empty key")
	}
	elems = append(elems, bson.D{{"endpoint_id", epId}})
	key = "uniquekeys." + key
	elems = append(elems, bson.D{{key, value}})

	var filter bson.D
	if len(elems) == 0 {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"$and", elems}}
	}
	log.Printf("Filter:%+v\n", filter)
	instance := &Instance{}
	err := u.instanceColl.FindOne(context.Background(), filter).Decode(instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (u *UspDb) WriteInstanceToDb(inst Instance) error {
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", inst.EndpointId}},
			bson.D{{"path", inst.Path}},
		},
		}}
	//log.Println("Adding Instance object with  Path:", inst.Path)
	opt := options.FindOneAndReplace().SetUpsert(true)
	err := u.instanceColl.FindOneAndReplace(context.TODO(), filter, inst, opt).Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in creating instance in DB")
		}
	}
	return nil
}

func (u *UspDb) DeleteInstancesByRegex(epId string, path string) error {
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
	log.Println("Deleting Instances with  Path:", path)
	_, err := u.instanceColl.DeleteMany(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting instance from DB")
		}
	}
	return err
}

func (u *UspDb) DeleteInstanceFromDb(epId string, path string) error {
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", epId}},
			bson.D{{"path", path}},
		},
		}}
	log.Println("Deleting Instance with  Path:", path)
	_, err := u.instanceColl.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting instance from DB")
		}
	}
	return err
}

func (u *UspDb) DeleteInstanceByUniqueKey(epId string, key string, value string) error {
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", epId}},
			bson.D{{"uniquekeys." + key, value}},
		},
		}}
	log.Printf("Deleting Instance with uniquekey: %s and value: %s\n", key, value)
	_, err := u.instanceColl.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting instance from DB")
		}
	}
	return err
}
