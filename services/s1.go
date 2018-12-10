package services

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
Should connect to a local mongo database with a collection “records”
Should have one endpoint that take the raw data (look at csv structure) through
this endpoint https://localhost:5555/api/records
Should check that the float value and the date given is valid or return an error
Should provide an endpoint to get a record by uuid
Should provide an endpoint to update a record by incrementing the value with
an int given in parameter
Should provide an endpoint to delete a record Only S2 can talk with this service
*/

func NewS1() (h http.Handler, err error) {
	s1 := S1{}

	if s1.sess, err = mgo.Dial("localhost"); err != nil {
		return
	}

	s1.collection = s1.sess.DB("spark").C("records")

	h = newRouter([]route{
		{"PUT", "/api/records", s1.putRecordsHandler},
		{"GET", "/api/records/:uuid", s1.getRecordHandler},
		{"POST", "/api/records/:uuid", s1.updateRecordHandler},
		{"DELETE", "/api/records/:uuid", s1.deleteRecordHandler},
	})
	return
}

type S1 struct {
	sess       *mgo.Session
	collection *mgo.Collection
}

func (s *S1) putRecordsHandler(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-type"); ct != "text/csv" {
		Error(w, ErrBadRequest)
	} else if err := s.putRecords(r.Body); err != nil {
		Error(w, ErrBadRequest, err)
	}
}

func (s *S1) putRecords(r io.Reader) (err error) {
	var data []string

	csvr := csv.NewReader(r)
	for {
		//TODO: introduce concurency to insert record faster here
		//Could probably spawn multiple goroutine. Would have to throttle
		//to not overload mongo maybe....
		if data, err = csvr.Read(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		r := &Record{}

		if err = r.Parse(data); err != nil {
			//With more time, I would make sure the user of the api
			//gets a good error message...but right now they'll get http bad request
			log.Println(err)
			continue
			//Could return here if we want to stop the whole batch...
		}

		if err = s.collection.Insert(r); err != nil {
			return
		}
	}

	return
}

func (s *S1) getRecordHandler(w http.ResponseWriter, r *http.Request) {
	uuid := urlParam(r, "uuid")
	if !isValidUUID(uuid) {
		Error(w, ErrBadRequest)
	} else if rec, err := s.findRecord(uuid); err != nil {
		Error(w, ErrNotFound, err)
	} else {
		JSONResponse(w, rec)
	}
}

func (s *S1) findRecord(uuid string) (r *Record, err error) {
	r = &Record{}
	//TODO: differentiate between not found and other error...n00b in mongo right now
	err = s.collection.Find(bson.M{"uuid": uuid}).One(&r)
	return
}

type updateRequest struct {
	Increment int
}

func (s *S1) updateRecordHandler(w http.ResponseWriter, r *http.Request) {
	var (
		uuid = urlParam(r, "uuid")
		ur   = updateRequest{}
	)

	if !isValidUUID(uuid) {
		Error(w, ErrBadRequest)
	} else if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		Error(w, ErrBadRequest, err)
	} else if rec, err := s.updateRecord(uuid, ur.Increment); err != nil {
		Error(w, err)
	} else {
		JSONResponse(w, rec)
	}
}

func (s *S1) updateRecord(uuid string, increment int) (rec *Record, err error) {
	rec = &Record{}

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"num": increment}},
		ReturnNew: true,
	}
	_, err = s.collection.Find(bson.M{"uuid": uuid}).Apply(change, &rec)
	return
}

func (s *S1) deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	uuid := urlParam(r, "uuid")
	if !isValidUUID(uuid) {
		Error(w, ErrBadRequest)
	} else if err := s.deleteRecord(uuid); err != nil {
		Error(w, err)
	}
}

func (s *S1) deleteRecord(uuid string) error {
	return s.collection.Remove(bson.M{"uuid": uuid})
}
