package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UspDbName                = "usp"
	UspParamCollection       = "params"
	UspDmCollection          = "datamodel"
	UspInstanceCollection    = "instances"
	UspCfgInstanceCollection = "cfginstances" // DefaultConfig Instance collection
	UspCfgParamCollection    = "cfgparams"    //  DefaultConfig Param collection
)

type UspDb struct {
	paramColl       *mongo.Collection
	dmColl          *mongo.Collection
	instanceColl    *mongo.Collection
	cfgInstanceColl *mongo.Collection
	cfgParamColl    *mongo.Collection
}

func (u *UspDb) Init(client *mongo.Client, dbName string) error {
	if client == nil {
		err := errors.New("DB is not connected, please try again...")
		return err
	}

	u.paramColl = client.Database(dbName).Collection(UspParamCollection)
	u.dmColl = client.Database(dbName).Collection(UspDmCollection)
	u.instanceColl = client.Database(dbName).Collection(UspInstanceCollection)
	u.cfgInstanceColl = client.Database(dbName).Collection(UspCfgInstanceCollection)
	u.cfgParamColl = client.Database(dbName).Collection(UspCfgParamCollection)

	return nil
}

func (u *UspDb) DeleteCollection(collName string) error {
	var err error
	switch collName {
	case UspParamCollection:
		err = u.paramColl.Drop(context.Background())
	case UspDmCollection:
		err = u.dmColl.Drop(context.Background())
	case UspInstanceCollection:
		err = u.instanceColl.Drop(context.Background())
	case UspCfgInstanceCollection:
		err = u.cfgInstanceColl.Drop(context.Background())
	case UspCfgParamCollection:
		err = u.cfgParamColl.Drop(context.Background())
	default:
		err = errors.New("Invalid collection name:" + collName)
	}
	return err
}
