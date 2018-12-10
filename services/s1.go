package services

import (
	"encoding/csv"
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
		{"PUT", "/api/records", s1.putRecords},
		{"GET", "/api/records/:uuid", s1.getRecord},
		{"POST", "/api/records/:uuid", s1.updateRecord},
		{"DELETE", "/api/records/:uuid", s1.deleteRecord},
	})
	return
}

type S1 struct {
	sess       *mgo.Session
	collection *mgo.Collection
}

func (s *S1) putRecords(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-type"); ct != "text/csv" {
		Error(w, ErrBadRequest)
	}

	if err := s.putRecordsReader(r.Body); err != nil {
		Error(w, err)
	}
}

func (s *S1) putRecordsReader(r io.Reader) (err error) {
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
			return ErrBadRequest
		}

		if err = s.collection.Insert(r); err != nil {
			return
		}
	}

	return
}

func (s *S1) getRecord(w http.ResponseWriter, r *http.Request) {
	uuid := urlParam(r, "uuid")
	if !isValidUUID(uuid) {
		Error(w, ErrBadRequest)
	}

	result := Record{}
	if err := s.collection.Find(bson.M{"uuid": uuid}).One(&result); err != nil {
		//TODO: differentiate between not found and other error...n00b in mongo right now
		log.Println(err)
		Error(w, ErrNotFound)
	} else {
		JSONResponse(w, result)
	}
}

func (s *S1) updateRecord(w http.ResponseWriter, r *http.Request) {
}

func (s *S1) deleteRecord(w http.ResponseWriter, r *http.Request) {
}
