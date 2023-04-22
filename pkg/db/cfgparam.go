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

// Node is for single instance objects for which path for each device is unique, e.g. IP
type CfgParamNode struct {
	Dev    DevType
	Path   string
	Params map[string]string
}

func (u *UspDb) GetCfgParams(dev *DevType, path string) (map[string]string, error) {
	if dev == nil {
		return nil, errors.New("Uninitialized dev pointer")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return nil, errors.New("DevType is not set")
	}
	if path == "" {
		return nil, errors.New("Empty Path")
	}
	var elems bson.A
	elems = append(elems, bson.D{{"dev.product_class", dev.ProductClass}})
	elems = append(elems, bson.D{{"dev.manufacturer", dev.Manufacturer}})
	elems = append(elems, bson.D{{"dev.model_name", dev.ModelName}})
	elems = append(elems, bson.D{{"path", path}})

	var filter bson.D
	if len(elems) == 0 {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"$and", elems}}
	}
	//log.Printf("Filter:%+v\n", filter)
	var paramNode CfgParamNode
	err := u.cfgParamColl.FindOne(context.Background(), filter).Decode(&paramNode)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return paramNode.Params, nil
}
func (u *UspDb) GetCfgParamNodes(dev *DevType) ([]*CfgParamNode, error) {
	if dev == nil {
		return nil, errors.New("uninitialized dev pointer")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return nil, errors.New("DevType is not set")
	}

	var elems bson.A
	elems = append(elems, bson.D{{"dev.product_class", dev.ProductClass}})
	elems = append(elems, bson.D{{"dev.manufacturer", dev.Manufacturer}})
	elems = append(elems, bson.D{{"dev.model_name", dev.ModelName}})
	var filter bson.D
	if len(elems) == 0 {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"$and", elems}}
	}
	//log.Printf("Filter:%+v\n", filter)
	cur, err := u.cfgParamColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}

	var paramNodes []*CfgParamNode
	if err := cur.All(context.Background(), &paramNodes); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return paramNodes, nil
}

func (u *UspDb) GetCfgParamsByRegex(dev *DevType, path string) ([]*CfgParamNode, error) {
	if dev == nil {
		return nil, errors.New("uninitialized dev pointer")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return nil, errors.New("DevType is not set")
	}
	if path == "" {
		return nil, errors.New("Empty Path")
	}
	regex := primitive.Regex{}
	regex.Pattern = path

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"dev.product_class", dev.ProductClass}},
			bson.D{{"dev.manufacturer", dev.Manufacturer}},
			bson.D{{"dev.model_name", dev.ModelName}},
			bson.D{{"path", regex}},
		},
		}}
	cur, err := u.cfgParamColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	var cfgNodes []*CfgParamNode
	if err := cur.All(context.Background(), &cfgNodes); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return cfgNodes, nil
}

func (u *UspDb) WriteCfgParamNode(p *CfgParamNode) error {
	if p == nil {
		return errors.New("Uninitialized cfgparam pointer")
	}
	if p.Dev.ProductClass == "" || p.Dev.Manufacturer == "" || p.Dev.ModelName == "" {
		return errors.New("DevType is not set")
	}
	if p.Path == "" {
		return errors.New("Empty Path")
	}
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"dev.product_class", p.Dev.ProductClass}},
			bson.D{{"dev.manufacturer", p.Dev.Manufacturer}},
			bson.D{{"dev.model_name", p.Dev.ModelName}},
			bson.D{{"path", p.Path}},
		},
		}}
	opt := options.FindOneAndReplace().SetUpsert(true)
	err := u.cfgParamColl.FindOneAndReplace(context.TODO(), filter, *p, opt).Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in updating cfg param in DB")
		}
	}
	return nil
}

func (u *UspDb) DeleteCfgParamNode(dev *DevType, path string) error {
	if dev == nil {
		return errors.New("uninitialized dev pointer")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return errors.New("DevType is not set")
	}
	if path == "" {
		return errors.New("Empty Path")
	}
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"dev.product_class", dev.ProductClass}},
			bson.D{{"dev.manufacturer", dev.Manufacturer}},
			bson.D{{"dev.model_name", dev.ModelName}},
			bson.D{{"path", path}},
		},
		}}
	log.Println("Deleting cfg param node with Path:", path)
	_, err := u.cfgParamColl.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting cfg param node from DB")
		}
	}
	return err
}

func (u *UspDb) DeleteCfgParamNodesByRegex(dev *DevType, path string) error {
	if dev == nil {
		return errors.New("uninitialized dev pointer")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return errors.New("DevType is not set")
	}
	if path == "" {
		return errors.New("Empty Path")
	}
	regex := primitive.Regex{}
	regex.Pattern = path

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"dev.product_class", dev.ProductClass}},
			bson.D{{"dev.manufacturer", dev.Manufacturer}},
			bson.D{{"dev.model_name", dev.ModelName}},
			bson.D{{"path", regex}},
		},
		}}
	log.Println("Deleting cfg param nodes with Path:", path)
	_, err := u.cfgParamColl.DeleteMany(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting cfg param nodes from DB")
		}
	}
	return err
}
