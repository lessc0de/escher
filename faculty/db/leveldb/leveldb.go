// Written in 2014 by Petar Maymounkov.
//
// It helps future understanding of past knowledge to save
// this notice, so peers of other times and backgrounds can
// see history clearly.

package basic

import (
	// "fmt"
	"strconv"

	"github.com/gocircuit/escher/faculty"
	. "github.com/gocircuit/escher/image"
	"github.com/gocircuit/escher/be"
	"github.com/gocircuit/escher/kit/plumb"

	"github.com/gocircuit/escher/github.com/syndtr/goleveldb/leveldb"
	"github.com/gocircuit/escher/github.com/syndtr/goleveldb/leveldb/util"
)

func init() {
	ns := faculty.Root.Refine("db").Refine("leveldb")
	ns.AddTerminal("File", File{})
}

// File
type File struct{}

//	File string
//	Put {Key []byte, Value []byte}
//	Query {Name interface{}, Start []byte, Limit []byte}
//	Result {Name interface{}, Result Image}
func (File) Materialize(*be.Matter) be.Reflex {
	reflex, eye := plumb.NewEye("File", "Put", "Query", "Result")
	go func() { // dispatch
		var err error
		var db *leveldb.DB
		connected, put, query := make(chan struct{}), make(chan Image, 5), make(chan Image, 5)
		go func() { // Put loop
			<-connected // wait for db connection
			for {
				p := <-put
				if err := db.Put(p["Key"].([]byte), p["Value"].([]byte), nil); err != nil {
					panic(err)
				}
			}
		}()
		go func() { // Query loop
			<-connected // wait for db connection
			for {
				q := <-query
				iter := db.NewIterator(
					&util.Range{
						Start: q["Start"].([]byte), 
						Limit: q["Limit"].([]byte),
					},
					nil,
				)
				slice := Make()
				if iter.First() {
					for i := 0; ; i++ {
						if err := iter.Error(); err != nil {
							panic(err)
						}
						slice.Grow(
							strconv.Itoa(i),
							Image{
								"Key": iter.Key(),
								"Value": iter.Value(),
							},
						)
						if !iter.Next() {
							break
						}
					}
				}
				eye.Show("Result", Make().Grow("Name", q["Name"]).Grow("Slice", slice))
			}
		}()
		for {
			valve, value := eye.See()
			switch valve {
			case "File":
				if db, err = leveldb.OpenFile(value.(string), nil); err != nil {
					panic(err)
				}
				close(connected)
			case "Put":
				put <- value.(Image)
			case "Query":
				query <- value.(Image)
 			}
		}
	}()
	return reflex
}
