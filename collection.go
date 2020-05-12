package nodedb

import (
	"time"
)

type Collection struct {
	Conn *Connection
	Name string
}

// NewNode creates and inserts a new node into the collection (runs SQL insert command), this way a node has allways a valid ID
// The ID may never be zero. A node that you use in your code must allways have a valid representation in the database.
func (c *Collection) NewNode() (*Node, error) {
	return insertNode(c, setupNode(c))
}

func (c *Collection) NewNodes(count int) ([]*Node, error) {

	result := make([]*Node, count, count)
	for i := 0; i < count; i++ {
		result[i] = setupNode(c)
	}
	return insertNodes(c, result)
}

// LoadNode loads a node by its ID, provide withContent = true, if the payload (content) should
// be loaded with the node.
func (c *Collection) LoadNode(nodeID int64, withContent bool) (*Node, error) {
	return loadNode(c, nodeID, withContent)
}

func (c *Collection) LoadNodes(nodeIDs []int64, withContent bool) ([]*Node, error) {
	return loadNodes(c, nodeIDs, withContent)
}

func (c *Collection) LoadNodesInSequence(nodeIDs []int64, withContent bool) ([]*Node, error) {

	nodes, err := loadNodes(c, nodeIDs, withContent)
	if err != nil {
		return nil, err
	}

	nmap := make(map[int64]*Node)
	for _, n := range nodes {
		nmap[n.ID] = n
	}
	for i, id := range nodeIDs {
		nodes[i] = nmap[id]
	}

	return nodes, nil
}

func (c *Collection) UpdateNode(update *Node) error {
	return updateNode(c, update)
}

func (c *Collection) UpdateNodes(updates []*Node) error {
	return updateNodes(c, updates)
}

func (c *Collection) DeleteNode(update *Node) error {
	return deleteNode(c, update)
}

func (c *Collection) DeleteNodes(updates []*Node) error {
	return deleteNodes(c, updates)
}

func setupNode(c *Collection) *Node {
	result := newNode()
	result.Created = time.Now()
	result.Start = time.Now()
	result.End = c.Conn.End
	result.Deleted = c.Conn.Deleted
	result.OwnerID = c.Conn.RootID
	result.Up = c.Conn.RootID
	result.Left = 0
	result.Right = 0
	return result
}
