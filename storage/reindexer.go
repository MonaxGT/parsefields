package storage

import (
	"github.com/restream/reindexer"
)

const insertOption = "id=serial()"

// Reindexer represents github.com/restream/reindexer database
type Reindexer struct {
	DB             *reindexer.Reindexer
	NamespaceField string
	NamespaceEvent string
	DBName         string
}

// Open opens connection to a remote reindexer server
func (r *Reindexer) Open(url string) error {
	r.DB = reindexer.NewReindex(url)
	err := r.DB.OpenNamespace(r.NamespaceField, reindexer.DefaultNamespaceOptions(), Fields{})
	if err != nil {
		return err
	}
	err = r.DB.OpenNamespace(r.NamespaceEvent, reindexer.DefaultNamespaceOptions(), Events{})
	if err != nil {
		return err
	}
	return nil
}

func (r *Reindexer) InsertFields(f *Fields) error {
	err := r.DB.Upsert(r.NamespaceField, f, insertOption)
	return err
}

func (r *Reindexer) InsertEvents(e *Events) error {
	err := r.DB.Upsert(r.NamespaceEvent, e, insertOption)
	return err
}

func (r *Reindexer) RestoreFields() ([]*Fields, error) {
	iter := r.DB.Query(r.NamespaceField).
		Exec()
	defer iter.Close()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	var F []*Fields
	for iter.Next() {
		elem := iter.Object().(*Fields)
		F = append(F, elem)
	}
	return F, nil

}

func (r *Reindexer) RestoreEvents() ([]*Events, error) {
	iter := r.DB.Query(r.NamespaceField).
		Exec()
	defer iter.Close()
	if err := iter.Error(); err != nil {
		return nil, err
	}
	var F []*Events
	for iter.Next() {
		elem := iter.Object().(*Events)
		F = append(F, elem)
	}
	return F, nil

}

func (r *Reindexer) DeleteFields(field string) error {
	query := r.DB.Query(r.NamespaceField).
		WhereString("field", reindexer.EQ, field).
		Limit(10).
		Offset(0).
		Exec()
	defer query.Close()
	if err := query.Error(); err != nil {
		return err
	}

	var err error
	for query.Next() {
		elem := query.Object().(*Fields)
		err = r.DB.Delete(r.NamespaceField, elem)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *Reindexer) DeleteEvents(logname string, eventid int32) error {
	query := r.DB.Query(r.NamespaceEvent).
		WhereString("LogName", reindexer.EQ, logname).
		WhereInt32("EventID", reindexer.EQ, eventid).
		Limit(10).
		Offset(0).
		Exec()
	defer query.Close()
	if err := query.Error(); err != nil {
		return err
	}
	var err error
	for query.Next() {
		elem := query.Object().(*Events)
		err = r.DB.Delete(r.NamespaceEvent, elem)
		if err != nil {
			return err
		}
	}
	return err
}
