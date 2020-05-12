package nodedb

import (
	"database/sql"
	"fmt"
	"log"
)

func QueryAxis(collection *Collection, nodeID int64, maxLevel int, limit int64) ([]*Node, error) {

	result := []*Node{}
	templateParams := map[string]string{}

	withContent := true // TODO make this a parameter to the function

	templateParams["collection"] = collection.Name
	templateParams["axis"] = "up" // TODO make this a parameter to the function
	templateParams["commaContent"] = ""
	templateParams["commaPContent"] = ""
	templateParams["limit"] = ""
	templateParams["where"] = ""

	if withContent {
		templateParams["commaContent"] = ", content"
		templateParams["commaPContent"] = ", p.content"
	}

	if limit > 0 {
		templateParams["limit"] = fmt.Sprint("limit ", limit)
	}

	if maxLevel > 0 {
		templateParams["where"] = fmt.Sprint("where level < ", maxLevel) // TODO handling mehrerer where clauses überlegen
	}

	rows, err := collection.Conn.DB.Query(processString(SQLRecursiveAxis, templateParams), nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		newNode := new(Node)
		// level, err := scanQueryAxisNode(rows, newNode, withContent)
		_, err := scanQueryAxisNode(rows, newNode, withContent) //TODO den level mit zurück übergeben, ggf neues struct mit embedded node
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("level:", level, newNode)
		count++
		result = append(result, newNode)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(count)
	return result, nil
}

// doubled code !!! TODO replace with scanNode
func scanQueryAxisNode(rows *sql.Rows, node *Node, withContent bool) (int, error) {

	text := sql.NullString{}
	level := 0
	var err error

	if withContent {
		err = rows.Scan(&level,
			&node.ID, &node.TypeID, &node.OwnerID, &node.Flags, &node.Created, &node.Start, &node.End, &node.Deleted,
			&node.Up, &node.Left, &node.Right, &text, &node.Content)
	} else {
		err = rows.Scan(&level,
			&node.ID, &node.TypeID, &node.OwnerID, &node.Flags, &node.Created, &node.Start, &node.End, &node.Deleted,
			&node.Up, &node.Left, &node.Right, &text)
	}

	if text.Valid {
		node.Text = text.String
	}

	node.contentLoaded = withContent

	return level, err
}
