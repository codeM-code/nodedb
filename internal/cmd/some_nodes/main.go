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

	"github.com/codeM-code/nodedb/internal/cli"
	"github.com/davecgh/go-spew/spew"
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
		panic(err)
	}
	conn.SetLogging(true)
	defer conn.Close()

	coll, err := conn.Collection("node")
	if err != nil {
		log.Fatal(err)
	}

	node1, err := coll.NewNode()
	if err != nil {
		log.Fatal(err)
	}

	id := node1.ID
	node1.Text = "Hi, I am the first node in this collection"

	coll.UpdateNode(node1)
	if err != nil {
		log.Fatal(err)
	}

	loadedNode, err := coll.LoadNode(id, false)
	if err != nil {
		log.Fatal(err)
	}

	loadedNode.Text += " and I was reloaded and updated."

	coll.UpdateNode(loadedNode)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(loadedNode)

	node2, err := coll.NewNode()
	if err != nil {
		log.Fatal(err)
	}
	node2.Text = "Hi I'm number two, have a bad mood since nobody cares about number two"
	id2 := node2.ID

	coll.UpdateNode(node2)
	if err != nil {
		log.Fatal(err)
	}

	coll.DeleteNode(node2)
	if err != nil {
		log.Fatal(err)
	}

	loadedNode, err = coll.LoadNode(id2, false)
	if err != nil {
		log.Fatal(err)
	}

	// currently a deleted node can still be loaded, this will change, once the design and implentation of NodeDD
	// has been completed. Note that the deletion date has been set to time.Now().
	spew.Dump(loadedNode)

	node3, err := coll.NewNode()
	if err != nil {
		log.Fatal(err)
	}

	node3.Text = "Number 3 has some content, finally"
	node3.Content = []byte("could be anything, that can be expressed as an array of bytes, but this is just a string")
	node3.SetContentLoaded(true) // this is required for content to be saved!
	err = coll.UpdateNode(node3)
	if err != nil {
		log.Fatal(err)
	}

	loadedNode, err = coll.LoadNode(node3.ID, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("node 3 loaded without content")
	spew.Dump(loadedNode)

	loadedNode, err = coll.LoadNode(node3.ID, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("node 3 loaded with content")
	spew.Dump(loadedNode)

}
