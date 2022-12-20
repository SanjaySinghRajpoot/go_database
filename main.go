package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jcelliott/lumber"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opencensus.io/resource"
	"gopkg.in/launchdarkly/go-server-sdk.v5/interfaces"
)

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger = lumber.NewConsoleLogger((lumber.INFO)

	driver := Driver{
		dir: dir,
		mutexes: make(map[string]*sync.Mutex),
		log: opts.Logger
	}

	if _, err := os.Stat(dit); err == nil{
		opts.Logger.Debug("Using '%s' (database already exists) \n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the database at '%s'.....\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interfaces{}) error {
    if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!!") 
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record!!")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil {   // 0755 is the permission to Make Dir
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n')) 

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmppath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
     
	if collection == "" {
		return fmt.Errorf("Missing collection - unable to read")
	}

	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")

	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string)([]string, error) {
   
	if collection == "" {
		return nil, fmt.Errorf("Missing collection from the project - unable to read")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)

	var records []string
    
	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))

		if err != nil {
			return nil,err
		}

		records = append(records, string(b))
	}
    
	return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
    
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi==nil, err != nil:
		return fmt.Errorf("unable to get the file with desired file name %v \n", path)
	case fi.mode().IsDir():
		return os.RemoveAll(dir)
	}
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
    
	d.mutex.Lock
	defer d.mutex.Unlock()
	m, err := d.mutexes[collection]
    
	if !err {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
}

func stat(path string)(fi os.FileInfo, err error){
	  if fi, err = os.Stat(path); os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	  }

	  return 
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, er := New(dir, nil)
	if er {
		fmt.Println("Error", er)
	}

	employees := []User{
		{"Tim", "23", "8899665544", "My tech", Address{"indore", "Madhya Pradesh", "India", "998855663322"}},
		{"Sam", "24", "9966345544", "techBro", Address{"Chennai", "Chennai", "India", "998812663322"}},
		{"Tom", "28", "9924345544", "techBro2", Address{"City", "State", "India", "997612663322"}},
	}

	for _, value := range employees {
		db.write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, er := db.ReadAll("users")

	if er != nil {
		fmt.Println("Error", er)
	}

	fmt.Println(records)

	allusers := []User{} // array to save the data from the JSON response

	for _, f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil { // destructuring
			fmt.Println("Error", err)
		}

		allusers = append(allusers, employeeFound)
	}

	fmt.Println((allusers))

	if err := db.Delete("user", "john"); err != nil {
		fmt.Println("Error", err)
	}

	if err := db.Delete("user", ""); err != nil {
		fmt.Println("Error", err)
	}

}
