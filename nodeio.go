package nodedb

import (
	"database/sql"
	"fmt"
	"time"
)

func newNode() *Node {
	newNode := Node{}
	newNode.Created = time.Now()
	newNode.Start = time.Now()
	newNode.End = time.Now()
	newNode.Deleted = time.Now()
	newNode.TypeID = 0
	newNode.Flags = 0
	newNode.OwnerID = 0
	newNode.Up, newNode.Left, newNode.Right = -1, -2, -3  // sentry values, if they show up, then something is wrong!
	newNode.Text = "New Node, not integrated into graph!" // sentry value, if this show up, then something is wrong!
	newNode.Content = nil
	return &newNode
}

func loadNode(coll *Collection, nodeID int64, withContent bool) (*Node, error) {
	result := Node{}

	sql := SQLString(SQLLoadNode).Collection(coll.Name).WithContent(withContent).String()

	coll.Conn.log(sql, "nodeID", nodeID)

	rows, err := coll.Conn.DB.Query(sql, nodeID)
	if err == nil {
		if rows.Next() { // only pick the first row
			scanNode(rows, &result, withContent)
		}
	}
	rows.Close()

	return &result, nil
}

func loadNodes(coll *Collection, nodeIDs []int64, withContent bool) ([]*Node, error) {

	result := []*Node{}

	sql := SQLString(SQLLoadNodes).Collection(coll.Name).NodeIDs(nodeIDs).WithContent(withContent).String()

	coll.Conn.log(sql)

	rows, err := coll.Conn.DB.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		node := newNode()
		scanNode(rows, node, withContent)

		result = append(result, node)
	}

	return result, nil
}

func insertNode(coll *Collection, newNode *Node) (*Node, error) {

	tx, err := coll.Conn.DB.Begin()
	if err != nil {
		return newNode, err
	}

	sql := SQLString(SQLInsertNode).Collection(coll.Name).String()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return newNode, err
	}
	defer stmt.Close()

	res, err := stmtExec(coll, sql, stmt,
		newNode.TypeID, newNode.OwnerID, newNode.Flags, newNode.Created, newNode.Start, newNode.End, newNode.Deleted,
		newNode.Up, newNode.Left, newNode.Right, newNode.Text, newNode.Content)
	// res, err := stmt.Exec(
	// 	newNode.TypeID, newNode.OwnerID, newNode.Flags, newNode.Created, newNode.Start, newNode.End, newNode.Deleted,
	// 	newNode.Up, newNode.Left, newNode.Right, newNode.Text, newNode.Content)

	if err != nil {
		tx.Rollback()
		return newNode, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return newNode, err
	}

	tx.Commit()
	stmt.Close()
	if err != nil {
		newNode.ID = -1
		return newNode, err
	}

	newNode.ID = id

	return newNode, nil
}

func insertNodes(coll *Collection, newNodes []*Node) ([]*Node, error) {

	tx, err := coll.Conn.DB.Begin()
	if err != nil {
		return newNodes, err
	}

	sql := SQLString(SQLInsertNode).Collection(coll.Name).String()

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return newNodes, err
	}
	defer stmt.Close()

	var _err error

	for _, n := range newNodes {

		res, err := stmtExec(coll, sql, stmt,
			n.TypeID, n.OwnerID, n.Flags, n.Created, n.Start, n.End, n.Deleted,
			n.Up, n.Left, n.Right, n.Text, n.Content)

		// res, err := stmt.Exec(
		// 	n.TypeID, n.OwnerID, n.Flags, n.Created, n.Start, n.End, n.Deleted,
		// 	n.Up, n.Left, n.Right, n.Text, n.Content)

		if err != nil {
			_err = err
			break
		}
		id, err := res.LastInsertId()
		if err != nil {
			_err = err
			break
		}
		n.ID = id
	}

	if _err != nil {
		tx.Rollback()
		return newNodes, err

	}

	tx.Commit()
	err = stmt.Close()
	if err != nil {
		return newNodes, err
	}

	return newNodes, nil
}

func updateNode(coll *Collection, updateNode *Node) error {
	tx, err := coll.Conn.DB.Begin()
	if err != nil {
		return err
	}

	sql := SQLString(SQLUpdateNode).Collection(coll.Name).WithContent(updateNode.contentLoaded).String()

	coll.Conn.log(sql)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if updateNode.contentLoaded {
		_, err = stmtExec(coll, sql, stmt,
			// _, err = stmt.Exec(
			updateNode.TypeID, updateNode.OwnerID, updateNode.Flags, updateNode.Created, updateNode.Start, updateNode.End, updateNode.Deleted,
			updateNode.Up, updateNode.Left, updateNode.Right, updateNode.Text, updateNode.Content,
			updateNode.ID,
		)
	} else {
		_, err = stmtExec(coll, sql, stmt,
			// _, err = stmt.Exec(
			updateNode.TypeID, updateNode.OwnerID, updateNode.Flags, updateNode.Created, updateNode.Start, updateNode.End, updateNode.Deleted,
			updateNode.Up, updateNode.Left, updateNode.Right, updateNode.Text,
			updateNode.ID,
		)
	}

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func updateNodes(coll *Collection, updateNodes []*Node) error {
	tx, err := coll.Conn.DB.Begin()
	if err != nil {
		return err
	}

	sqlwContent := SQLString(SQLUpdateNode).Collection(coll.Name).WithContent(true).String()
	sqlwoContent := SQLString(SQLUpdateNode).Collection(coll.Name).WithContent(false).String()

	stmtwContent, err := tx.Prepare(sqlwContent)
	if err != nil {
		return err
	}
	stmtwoContent, err := tx.Prepare(sqlwoContent)
	if err != nil {
		return err
	}

	defer stmtwContent.Close()
	defer stmtwoContent.Close()

	var _err error

	for _, n := range updateNodes {

		if n.contentLoaded {
			_, err = stmtExec(coll, sqlwContent, stmtwContent,
				// _, err := stmtwContent.Exec(
				n.TypeID, n.OwnerID, n.Flags, n.Created, n.Start, n.End, n.Deleted,
				n.Up, n.Left, n.Right, n.Text, n.Content,
				n.ID,
			)
			if err != nil {
				_err = err
				break
			}

		} else {

			_, err = stmtExec(coll, sqlwoContent, stmtwoContent,
				// _, err := stmtwoContent.Exec(
				n.TypeID, n.OwnerID, n.Flags, n.Created, n.Start, n.End, n.Deleted,
				n.Up, n.Left, n.Right, n.Text,
				n.ID,
			)
			if err != nil {
				_err = err
				break
			}
		}

	}

	if _err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func deleteNode(coll *Collection, deleteNode *Node) error {

	deleteNode.Deleted = time.Now()
	return updateNode(coll, deleteNode)

}

func deleteNodes(coll *Collection, deleteNodes []*Node) error {

	for _, dn := range deleteNodes {
		dn.Deleted = time.Now()
	}
	return updateNodes(coll, deleteNodes)

}

func scanNode(rows *sql.Rows, node *Node, withContent bool) error {

	text := sql.NullString{}
	scanSlice := []interface{}{&node.ID, &node.TypeID, &node.OwnerID, &node.Flags, &node.Created, &node.Start, &node.End, &node.Deleted,
		&node.Up, &node.Left, &node.Right, &text} // can be null, validity test required

	if withContent {
		scanSlice = append(scanSlice, &node.Content)
	}

	err := rows.Scan(scanSlice...)

	if text.Valid {
		node.Text = text.String
	}

	node.contentLoaded = withContent

	return err
}

func stmtExec(coll *Collection, sql string, stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {

	coll.Conn.log(sql, "Parameters:", fmt.Sprint(args...))

	return stmt.Exec(args...)

}
