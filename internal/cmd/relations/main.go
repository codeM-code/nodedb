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

	"github.com/codeM-code/nodedb"
	"github.com/codeM-code/nodedb/internal/cli"
)

// Select link.text, right.text from relation as link join relation as right on link.right = right.id where link.left in (1,2,3);

// !!! nicht fertig !!!
// !! this is my current area of work !!!

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
		panic(err)
	}
	defer conn.Close()

	coll, err := conn.Collection("relation")
	if err != nil {
		log.Fatal(err)
	}

	rootID, err := createNodesAndLinks(coll)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("the root ID is", rootID)

}

func createNodesAndLinks(coll *nodedb.Collection) (int64, error) {

	pos := 0

	nodes, err := coll.NewNodes(17)
	if err != nil {
		return 0, err
	}

	// nodes/vertices
	thomasID := fillNode(nodes[pos], "Thomas", 1, 0, 0) // person
	pos++
	andreasID := fillNode(nodes[pos], "Andreas", 1, 0, 0) // person
	pos++
	ralphID := fillNode(nodes[pos], "Ralph", 1, 0, 0) // person
	pos++

	codemID := fillNode(nodes[pos], "codeM GmbH", 2, 0, 0) // company
	pos++
	acmeID := fillNode(nodes[pos], "Acme Corporation", 2, 0, 0) // company
	pos++

	berichteID := fillNode(nodes[pos], "World Peace", 3, 0, 0) // project
	pos++
	stundenID := fillNode(nodes[pos], "World Dominance", 3, 0, 0) // project
	pos++

	// links/edges
	fillNode(nodes[pos], "has client", 22, codemID, acmeID)
	pos++

	fillNode(nodes[pos], "employed by", 12, thomasID, codemID)
	pos++
	fillNode(nodes[pos], "employed by", 12, andreasID, acmeID)
	pos++
	fillNode(nodes[pos], "employed by", 12, ralphID, acmeID)
	pos++

	fillNode(nodes[pos], "has project", 23, acmeID, berichteID)
	pos++
	fillNode(nodes[pos], "has project", 23, acmeID, stundenID)
	pos++

	fillNode(nodes[pos], "works on", 13, thomasID, stundenID)
	pos++
	fillNode(nodes[pos], "works on", 13, thomasID, berichteID)
	pos++
	fillNode(nodes[pos], "works on", 13, andreasID, stundenID)
	pos++
	fillNode(nodes[pos], "works on", 13, ralphID, berichteID)
	pos++

	err = coll.UpdateNodes(nodes)
	if err != nil {
		return 0, err
	}
	return thomasID, nil
}

func fillNode(node *nodedb.Node, text string, typeid int64, left int64, right int64) int64 {
	node.Up = 0
	node.Text = text
	node.TypeID = typeid
	node.Left = left
	node.Right = right
	return node.ID
}
