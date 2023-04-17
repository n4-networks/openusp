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

type DevType struct {
	ProductClass string `bson:"product_class"`
	Manufacturer string `bson:"manufacturer"`
	ModelName    string `bson:"model_name"`
}

type CfgInstance struct {
	Dev    DevType
	Path   string
	Level  int
	Key    string
	Params map[string]string
}

func (u *UspDb) GetCfgInstances(dev *DevType) ([]*CfgInstance, error) {
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
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"level", 1}})
	cur, err := u.cfgInstanceColl.Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}

	var instances []*CfgInstance
	if err := cur.All(context.Background(), &instances); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return instances, nil
}

func (u *UspDb) GetCfgInstancesByPath(dev *DevType, path string) ([]*CfgInstance, error) {
	if dev == nil {
		return nil, errors.New("uninitialized dev pointer")
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
	cur, err := u.cfgInstanceColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}

	var instances []*CfgInstance
	if err := cur.All(context.Background(), &instances); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return instances, nil
}

func (u *UspDb) GetCfgInstancesByRegex(dev *DevType, path string) ([]*CfgInstance, error) {
	if dev == nil {
		return nil, errors.New("uninitialized dev pointer")
	}
	if path == "" {
		return nil, errors.New("Empty Path")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return nil, errors.New("DevType is not set")
	}

	var elems bson.A
	elems = append(elems, bson.D{{"dev.product_class", dev.ProductClass}})
	elems = append(elems, bson.D{{"dev.manufacturer", dev.Manufacturer}})
	elems = append(elems, bson.D{{"dev.model_name", dev.ModelName}})
	elems = append(elems, bson.D{{"path", primitive.Regex{path, ""}}})
	var filter bson.D
	if len(elems) == 0 {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"$and", elems}}
	}
	//log.Printf("Filter:%+v\n", filter)
	cur, err := u.cfgInstanceColl.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	if cur.RemainingBatchLength() == 0 {
		return nil, errors.New("No documents found")
	}

	var instances []*CfgInstance
	if err := cur.All(context.Background(), &instances); err != nil {
		log.Println("Error in decoding:", err)
		return nil, err
	}
	return instances, nil
}

func (u *UspDb) WriteCfgInstance(inst *CfgInstance) error {
	if inst == nil {
		return errors.New("uninitialized cfgInstance pointer")
	}
	if inst.Path == "" {
		return errors.New("Empty Path")
	}
	if inst.Dev.ProductClass == "" || inst.Dev.Manufacturer == "" || inst.Dev.ModelName == "" {
		return errors.New("DevType is not set")
	}
	if inst.Key == "" {
		return errors.New("key is not provided")
	}
	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"dev.product_class", inst.Dev.ProductClass}},
			bson.D{{"dev.manufacturer", inst.Dev.Manufacturer}},
			bson.D{{"dev.model_name", inst.Dev.ModelName}},
			bson.D{{"path", inst.Path}},
			bson.D{{"key", inst.Key}},
		},
		}}
	//log.Println("Adding Instance object with  Path:", inst.Path)
	opt := options.FindOneAndReplace().SetUpsert(true)
	err := u.cfgInstanceColl.FindOneAndReplace(context.TODO(), filter, inst, opt).Err()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in creating Cfg Instance in DB")
		}
	}
	return nil
}

func (u *UspDb) DeleteCfgInstancesByRegex(dev *DevType, path string) error {
	if dev == nil {
		return errors.New("Uninitialized dev pointer")
	}
	if path == "" {
		return errors.New("Empty Path")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return errors.New("DevType is not set")
	}
	var elems bson.A
	elems = append(elems, bson.D{{"dev.product_class", dev.ProductClass}})
	elems = append(elems, bson.D{{"dev.manufacturer", dev.Manufacturer}})
	elems = append(elems, bson.D{{"dev.model_name", dev.ModelName}})
	elems = append(elems, bson.D{{"path", primitive.Regex{path, ""}}})

	var filter bson.D
	if len(elems) == 0 {
		filter = bson.D{{}}
	} else {
		filter = bson.D{{"$and", elems}}
	}
	log.Println("Deleting Cfg Instances with Path:", path)
	_, err := u.cfgInstanceColl.DeleteMany(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting Cfg Instance from DB")
		}
	}
	return err
}

func (u *UspDb) DeleteCfgInstance(dev *DevType, path string, key string) error {
	if dev == nil {
		return errors.New("Uninitialized dev pointer")
	}
	if path == "" {
		return errors.New("Empty Path")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return errors.New("DevType is not set")
	}
	//var filter bson.D
	var elems bson.A
	elems = append(elems, bson.D{{"dev.product_class", dev.ProductClass}})
	elems = append(elems, bson.D{{"dev.manufacturer", dev.Manufacturer}})
	elems = append(elems, bson.D{{"dev.model_name", dev.ModelName}})
	elems = append(elems, bson.D{{"path", path}})
	elems = append(elems, bson.D{{"key", key}})

	filter := bson.D{{"$and", elems}}

	//log.Printf("Filter:%+v\n", filter)
	log.Println("Deleting Cfg Instance with Path:", path)

	result, err := u.cfgInstanceColl.DeleteOne(context.TODO(), filter, nil)
	log.Println("Delete count:", result.DeletedCount)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting Cfg Instance from DB")
		}
	}
	return nil
}

func (u *UspDb) DeleteCfgInstancesByDevType(dev *DevType) error {
	if dev == nil {
		return errors.New("Unitialized dev pointer")
	}
	if dev.ProductClass == "" || dev.Manufacturer == "" || dev.ModelName == "" {
		return errors.New("DevType is not set")
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

	log.Printf("Deleting Cfg Instance with ProductClass: %v Manufacturer: %v ModelName: %v\n", dev.ProductClass, dev.Manufacturer, dev.ModelName)
	_, err := u.cfgInstanceColl.DeleteMany(context.TODO(), filter, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New("Error in deleting Cfg instances from DB")
		}
	}
	return err
}
