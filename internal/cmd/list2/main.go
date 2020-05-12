// Copyright 2020 codeM GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/codeM-code/nodedb"
	"github.com/codeM-code/nodedb/internal/cli"
)

// this creates 3 double linked list under a single head node, list entries are in random order, this random order is preserved
// by linking up the left-right-axes. The 3 lists are separated by a different typeid.

func main() {

	name, err := cli.ExecutableName()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("running", name)

	logFile, err := cli.StartLogFile()
	defer logFile.Close()

	conn, err := cli.OpenFreshDB()
	if err != nil {
		log.Fatal(err)
	}

	// conn.SetLogging(true)

	defer conn.Close()

	coll, err := conn.Collection("list")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Collection:", coll.Name)

	head, err := coll.NewNode()
	if err != nil {
		log.Fatal(err)
	}
	head.Text = "List-Two, child nodes are in random order"
	head.Up = 0
	coll.UpdateNode(head)
	if err != nil {
		log.Fatal(err)
	}

	count := 25

	list1, err := createList(head.ID, count, coll, 3, "First")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("First with %d entries created\n", len(list1))

	list2, err := createList(head.ID, count, coll, 4, "Second")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Second with %d entries created\n", len(list2))
	list3, err := createList(head.ID, count, coll, 5, "Third")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Third with %d entries created\n", len(list3))

	fmt.Println("loaded by QueryList")

	// listResult, err := nodedb.QueryList(coll, head.ID, false) // all types
	// listResult, err := nodedb.QueryList(coll, head.ID, false, 3) // just type 3
	// listResult, err := nodedb.QueryList(coll, head.ID, false, 4, 5) // types 4 and 5
	listResult, err := nodedb.QueryList(coll, head.ID, false, 3, 4, 5) // types 3, 4, 5 explicitely

	if err != nil {
		log.Fatal(err)
	}

	for _, n := range listResult {
		fmt.Println("Text:", n.Text, "   id:", n.ID, "    type", n.TypeID, "    left", n.Left, "    right", n.Right)
	}

}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

type sortedList []*listnode

func (e sortedList) Len() int      { return len(e) }
func (e sortedList) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e sortedList) Less(i, j int) bool {
	return e[i].randomInt < e[j].randomInt
}

type listnode struct {
	*nodedb.Node
	randomInt int // to produce a random order that wil be recreated after load
}

func createList(headID int64, count int, coll *nodedb.Collection, typeID int64, text string) ([]*listnode, error) {
	// make the list
	list := make([]*listnode, count, count)
	nodes, err := coll.NewNodes(count)
	if err != nil {
		return nil, err
	}

	// attach nodedb Nodes
	for i, n := range nodes {
		list[i] = &listnode{n, rnd.Int()}
		n.TypeID = typeID
	}

	// order by random value
	sort.Sort(sortedList(list))

	// connecting left-right-axis to preserve this particular order
	for i, n := range list {
		n.Text = fmt.Sprintf("%s. I am number %d", text, i+1)
		if i > 0 {
			n.Left = list[i-1].ID
		}
		if i < count-1 {
			n.Right = list[i+1].ID
		}
	}

	err = coll.UpdateNodes(nodes)
	if err != nil {
		return nil, err
	}

	return list, nil

}
