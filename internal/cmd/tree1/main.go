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

	coll, err := conn.Collection("data")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Collection:", coll.Name)

	// create a 3 level tree
	root1, nodeCount := treeInMemory("The First Tree", 1, 3) // 3 levels
	err = attachNodesStart(coll, nodeCount, root1)           // create and save nodes
	if err != nil {
		log.Fatal(err)
	}
	root2, nodeCount := treeInMemory("The Second Tree", 1, 4) // 4 levels
	err = attachNodesStart(coll, nodeCount, root2)            // create and save nodes
	if err != nil {
		log.Fatal(err)
	}
	root3, nodeCount := treeInMemory("The Third Tree", 1, 2) // 2 levels
	err = attachNodesStart(coll, nodeCount, root3)           // create and save nodes
	if err != nil {
		log.Fatal(err)
	}
	// root4, nodeCount := treeInMemory("The Forth Tree", 1, 5) // 5 levels
	// err = attachNodesStart(coll, nodeCount, root4)           // create and save nodes
	// if err != nil {
	// 	log.Fatal(err)
	// }

	res, err := nodedb.QueryAxis(coll, 1112, 2, 0) // load nodes down to level 2 (skip level 3 and 4) from second tree

	for i, n := range res {
		fmt.Println(i+1, "id:", n.ID, "parent:", n.Up, "Text:", n.Text)
	}

}

type Treenode struct {
	*nodedb.Node
	Text    string
	Childen []*Treenode
}

// ---- Tree Phase 1
func treeInMemory(text string, typeid int64, depth int) (*Treenode, int) {
	var root = Treenode{Text: text}

	count := addChildren(1, depth, &root, text, 10, 1)

	// fmt.Println("count", count)
	return &root, count

}

func addChildren(level int, maxlevel int, node *Treenode, text string, n int, count int) int {

	if level > maxlevel {
		return count
	}
	for i := 0; i < n; i++ {
		newNode := Treenode{Text: fmt.Sprintf("%s -- level %d, child %d", text, level, i)}
		count++
		node.Childen = append(node.Childen, &newNode)
		count = addChildren(level+1, maxlevel, &newNode, text, n, count)

	}

	return count
}

// ---- Tree Phase 2 attache nodes breadth first traversal

func attachNodesStart(col *nodedb.Collection, numberOfNodes int, root *Treenode) error {

	nodes, err := col.NewNodes(numberOfNodes)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return err
	}

	pos := new(int)
	*pos = 0

	attachNodes(root, nodes, 0, pos)

	return col.UpdateNodes(nodes)

}

func attachNodes(treeNode *Treenode, nodes []*nodedb.Node, up int64, pos *int) error {

	treeNode.Node = nodes[*pos]
	*pos++
	treeNode.Up = up
	treeNode.Node.Text = treeNode.Text

	for _, childNode := range treeNode.Childen {
		attachNodes(childNode, nodes, treeNode.ID, pos)
	}
	return nil
}
