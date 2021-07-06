package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type engine struct {
	db         *mongo.Database
	resumeID   []byte
	lock       sync.Mutex
	ctx        context.Context
	cancelFunc func()
}

func NewEngine(db *mongo.Database) (*engine, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	e := &engine{
		db:         db,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	return e, e.Close
}

func (e *engine) CDC(ctx context.Context) {
	e.loadResumeID()

	changeStreams, err := e.db.Watch(ctx, bson.D{}, e.Options())
	if err != nil {
		panic(err)
	}

	defer func() {
		err := changeStreams.Close(context.TODO())
		if err != nil {
			log.Printf("close change stream error: %v\n", err)
		}
	}()

	for changeStreams.Next(ctx) {
		ce := &map[string]interface{}{}
		err := changeStreams.Decode(ce)
		if err != nil {
			panic(err)
		}
		PrintJSON(ce)

		e.lock.Lock()
		_id, err := json.Marshal((*ce)["_id"])
		if err != nil {
			panic(err)
		}
		e.resumeID = _id
		e.lock.Unlock()
	}

	e.cancelFunc()
}

func (e *engine) loadResumeID() {
	data, err := os.ReadFile("resume_id")
	if err == os.ErrNotExist {
		return
	}

	if err != nil {
		panic(err)
	}

	e.resumeID = data
}

// Options to set watch change options: SetStartAfter,...
func (e *engine) Options() *options.ChangeStreamOptions {
	opts := options.ChangeStream()
	if e.resumeID != nil {
		r := bson.M{}
		err := json.Unmarshal(e.resumeID, &r)
		if err != nil {
			panic(err)
		}
		opts.SetStartAfter(r)
	}
	return opts
}

// Close to graceful shutdown
func (e *engine) Close() {
	<-e.ctx.Done()
	if e.resumeID == nil {
		return
	}

	rFile, err := os.OpenFile("resume_id", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Printf("open resume file by error: %v\n", err)
		return
	}

	_, err = rFile.Write(e.resumeID)
	if err != nil {
		log.Printf("write resume file by error: %v\n", err)
	}
}
