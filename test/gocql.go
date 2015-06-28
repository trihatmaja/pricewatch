/* Before you execute the program, Launch `cqlsh` and execute:
create keyspace example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
create table example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));
create index on example.tweet(timeline);
*/
package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
)

func main() {
	// connect to the cluster
	cluster := gocql.NewCluster("192.168.1.54")
	cluster.Keyspace = "example"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// insert a tweet
	// if err := session.Query(`INSERT INTO tweet (timeline, id, text, name, time) VALUES (?, ?, ?, ?, ?)`,
	// 	"me", gocql.TimeUUID(), "hello world", "Ninoy", time.Now()).Exec(); err != nil {
	// 	log.Fatal(err)
	// }

	var id gocql.UUID
	var text string

	/* Search for a specific set of records whose 'timeline' column matches
	 * the value 'me'. The secondary index that we created earlier will be
	 * used for optimizing the search */
	// if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
	// 	"me").Scan(&id, &text); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Tweet:", id, text)

	// list all tweets
	iter := session.Query(`SELECT id, text FROM tweet WHERE time < ? allow filtering`, time.Now()).Iter()
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:", id, text)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}
