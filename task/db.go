package db

import (
	"encoding/binary"
	"time"

	"go.etcd.io/bbolt"
)

type Task struct {
	ID    int
	Value string
}

var taskBucket = []byte("tasks")
var db *bbolt.DB

func Init(dbPath string) error {
	var err error

	db, err = bbolt.Open(dbPath, 0606, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				ID:    btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	return tasks, err
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}
