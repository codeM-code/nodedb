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

// this creates a double linked list with a head (root) node. All nodes in the list have the same root node in their up property.
// The first node has left == 0, the last node has right == 0

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
	head.Text = "List-One"
	head.Up = 0
	coll.UpdateNode(head)
	if err != nil {
		log.Fatal(err)
	}

	count := 100

	list1, err := createList(head.ID, count, coll, 3)
	if err != nil {
		log.Fatal(err)
	}

	ids := make([]int64, count, count)

	for i, n := range list1 {
		// fmt.Println("Text:", n.Text, "   id:", n.ID, "    left", n.Left, "    right", n.Right)
		ids[i] = n.ID
	}

	fmt.Println("loaded by IDs in sequence")
	list2, err := coll.LoadNodesInSequence(ids, false)
	if err != nil {
		log.Fatal(err)
	}

	for i, n := range list2 {
		fmt.Println("Text:", n.Text, "   id:", n.ID, "    type", n.TypeID, "    left", n.Left, "    right", n.Right)
		ids[i] = n.ID
	}

	fmt.Println("loaded by QueryList")
	list3, err := nodedb.QueryList(coll, head.ID, false, 3)
	if err != nil {
		log.Fatal(err)
	}

	for i, n := range list3 {
		fmt.Println("Text:", n.Text, "   id:", n.ID, "    type", n.TypeID, "    left", n.Left, "    right", n.Right)
		ids[i] = n.ID
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

func createList(headID int64, count int, coll *nodedb.Collection, typeID int64) ([]*listnode, error) {
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

	// connecting left-right-axis
	for i, n := range list {
		n.Text = fmt.Sprintf("I am number %d", i+1)
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
