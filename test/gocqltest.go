package main

import (
	"fmt"
	"github.com/elvtechnology/gocqltable"
	"github.com/elvtechnology/gocqltable/recipes"
	"github.com/gocql/gocql"
	"log"
	"time"
)

func main() {

	// Generic initialization of gocql
	c := gocql.NewCluster("192.168.1.54")
	s, err := c.CreateSession()
	if err != nil {
		log.Fatalln("Unable to open up a session with the Cassandra database (err=" + err.Error() + ")")
	}

	// Tell gocqltable to use this session object as the default for new objects
	gocqltable.SetDefaultSession(s)

	// Now we're ready to create our first keyspace. We start by getting a keyspace object
	keyspace := gocqltable.NewKeyspace("gocqltable_test")
	// Now lets create that in the database using the simple strategy and durable writes (true)
	// err = keyspace.Create(map[string]interface{}{
	// 	"class":              "SimpleStrategy",
	// 	"replication_factor": 1,
	// }, true)
	// if err != nil { // If something went wrong we print the error and quit.
	// 	log.Fatalln(err)
	// }

	// Now that we have a very own keyspace to play with, lets create our first table.

	// First we need a Row-object to base the table on. It will later be passed to the table wrapper
	// to be used for returning row-objects as the answer to fetch requests.
	type User struct {
		Email    string // Our primary key
		Password string `password`     // Use Tags to rename fields
		Active   bool   `cql:"active"` // If there are multiple tags, use `cql:""` to specify what the table column will be
		Created  time.Time
	}

	// Let's define and instantiate a table object for our user table
	userTable := struct {
		recipes.CRUD // If you looked at the base example first, notice we replaced this line with the recipe
	}{
		recipes.CRUD{ // Here we didn't replace, but rather wrapped the table object in our recipe, effectively adding more methods to the end API
			keyspace.NewTable(
				"users",           // The table name
				[]string{"email"}, // Row keys
				nil,               // Range keys
				User{},            // We pass an instance of the user struct that will be used as a type template during fetches.
			),
		},
	}

	// // Lets create this table in our cassandra database
	// err = userTable.Create()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// Now that we have a keyspace with a table in it: lets make a few rows! In the base example we had to write out the CQL manually, this time
	// around, however, we can insert entire User objects.

	// Lets instantiate a user object, set its values and insert it
	user1 := User{
		Email:    "1@example.com",
		Password: "123456",
		Active:   true,
		Created:  time.Now().UTC(),
	}
	err = userTable.Insert(user1)
	if err != nil {
		log.Fatalln(err)
	}

	// With our database filled up with users, lets query it and print out the results (containing all users in the database).
	rowset, err := userTable.List() // Our rowset variable is a "interface{}", and here we type assert it to a slice of pointers to "User"
	for _, user := range rowset.([]*User) {
		fmt.Println(user)
	}
	if err != nil {
		log.Fatalln(err)
	}

	// You can also fetch a single row, obviously
	row, err := userTable.Get("1@example.com")
	if err != nil {
		log.Fatalln(err)
	}
	user := row.(*User)

	// Lets update this user by changing his password
	user.Password = "654321"
	err = userTable.Update(user)
	if err != nil {
		log.Fatalln(err)
	}

	// Lets delete user 1@example.com
	err = userTable.Delete(user)
	if err != nil {
		log.Fatalln(err)
	}

	// Lets clean up after ourselves by dropping the keyspace.
	keyspace.Drop()
}
