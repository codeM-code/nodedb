package nodedb

import (
	"fmt"
	"log"
)

// QueryList finds all Nodes with provided parentID (up) and and optional list of typeids.
// It first queries the database and then builds a sequence ordered by the left-right-axis.
// There may be more then one list in the result set, check for start nodes with left==0 or end nodes with right == 0
// It reports an error if the number of sequenced nodes differs from the number of nodes return by the data base query.
// The result set may still be usable.
func QueryList(collection *Collection, upID int64, withContent bool, typeids ...int64) ([]*Node, error) {

	processing := 0
	length := 0
	listHeads := []*Node{}
	result := []*Node{}

	templateParams := map[string]string{}
	templateParams["collection"] = collection.Name
	templateParams["typeids"] = ""
	templateParams["commaContent"] = ""

	if withContent {
		templateParams["commaContent"] = ", content"
	}
	if len(typeids) == 1 {
		templateParams["typeids"] = fmt.Sprint(" and typeid =", typeids[0])
	}
	if len(typeids) > 1 {
		templateParams["typeids"] = fmt.Sprintf("and typeid in (%s)", intSliceToString(typeids))
	}

	rows, err := collection.Conn.DB.Query(processString(SQLListQuery, templateParams), upID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		newNode := new(Node)
		err := scanNode(rows, newNode, withContent)
		if err != nil {
			log.Fatal(err)
		}
		length++
		result = append(result, newNode)
		if newNode.Left == 0 { // this starts a sequence
			listHeads = append(listHeads, newNode)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(listHeads) == length || len(listHeads) == 0 {
		return result, fmt.Errorf("query result of length %d is not a list", length)
	}

	// for _, n := range result {
	// 	fmt.Println("QUERY!!! Text:", n.Text, "   id:", n.ID, "    left", n.Left, "    right", n.Right)
	// }

	// build the list(s)
	nodeMap := map[int64]*Node{}
	for _, n := range result {
		nodeMap[n.ID] = n
	}

	for _, h := range listHeads {
		result[processing] = h
		processing++

		for n := h; n.Right != 0 && processing < length; processing++ {
			n = nodeMap[n.Right]
			result[processing] = n
		}
	}

	return result, nil

}
