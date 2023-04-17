package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Param struct {
	EndpointId string `json:"endpoint_id" bson:"endpoint_id"`
	Path       string `json:"path" bson:"path"`
	Value      string `json:"value" bson:"value"`
}

func (u *UspDb) GetParams(epId string, path string) ([]*Param, error) {
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
	cur, err := u.paramColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}
	var params []*Param
	if err := cur.All(context.Background(), &params); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return params, nil
}

func (u *UspDb) GetParamsByRegex(epId string, path string) ([]*Param, error) {
	regex := primitive.Regex{}
	regex.Pattern = path

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", epId}},
			bson.D{{"path", regex}},
		},
		}}
	cur, err := u.paramColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	var params []*Param
	if err := cur.All(context.Background(), &params); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return params, nil
}

func (u *UspDb) GetParamObjByValue(epId string, path string, name string, value string) ([]*Param, error) {
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
	cur, err := u.paramColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}
	var params []*Param
	if err := cur.All(context.Background(), &params); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return params, nil
}

func (u *UspDb) GetAllEndpoints() ([]string, error) {

	filter := bson.D{}
	//log.Printf("Filter:%+v\n", filter)
	//var epIds []string
	//var err error
	epIdIntfs, err := u.paramColl.Distinct(context.Background(), "endpoint_id", filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	epIdIntfsLen := len(epIdIntfs)

	if epIdIntfsLen == 0 {
		return nil, errors.New("No documents found")
	}
	/*
		var epIds []string
		if err := cur.All(context.Background(), &epIds); err != nil {
			log.Println("Error in decoding:", err)
			return nil, err
		}
	*/
	epIds := make([]string, epIdIntfsLen)
	for i, epIdIntf := range epIdIntfs {
		epIds[i] = fmt.Sprintf("%v", epIdIntf)
	}
	return epIds, nil
}

func (u *UspDb) WriteParamToDb(p *Param) error {
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", p.EndpointId}},
			bson.D{{"path", p.Path}},
		},
		}}
	opt := options.FindOneAndReplace().SetUpsert(true)
	err := u.paramColl.FindOneAndReplace(context.TODO(), filter, *p, opt).Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in updating param in DB")
		}
	}
	return nil
}

func (u *UspDb) DeleteParamManyFromDb(epId string, path string) error {
	regex := primitive.Regex{}
	regex.Pattern = path

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"endpoint_id", epId}},
			bson.D{{"path", regex}},
		},
		}}
	log.Println("Deleting param with  Path:", path)
	_, err := u.paramColl.DeleteMany(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting params from DB")
		}
	}
	return err
}
