package nodedb

// all SQL-Statements are here:

// SQLCollections lists all tables in SQLite database
const SQLCollections = `select tbl_name from sqlite_master where type = "table" and tbl_name != "sqlite_sequence" order by tbl_name;`

// SQLCollection find a single table name
const SQLCollection = `select tbl_name, type from sqlite_master where tbl_name = ?;`

// const SQLVacuumInto = `vacuum into {{.backupFilename}};` // !!! not (yet) working

// SQLNewCollection creates the scheme and the indexes for the node
const SQLNewCollection = `
	PRAGMA page_size = 4096; -- sind die Pragmas hier überhaupt effektiv???
	PRAGMA foreign_keys = "1";

	DROP TABLE IF EXISTS "{{.collection}}";

	CREATE TABLE "{{.collection}}"(
	"id" integer not null primary key autoincrement,
	"typeid" integer not null,
	"ownerid" integer not null,
	"flags" integer not null,
	"created" datetime not null,
	"deleted" datetime not null,
	"start" datetime not null,
	"end" datetime not null,
	"up" integer not null,
	"left" integer not null,
	"right" integer not null,
	"text" text,
	"content" blob
	);

	CREATE INDEX {{.collection}}type_id ON {{.collection}}(typeid);
	CREATE INDEX {{.collection}}owner_id ON {{.collection}}(ownerid);
	CREATE INDEX {{.collection}}up_id ON {{.collection}}(up);
	CREATE INDEX {{.collection}}left_id ON {{.collection}}(left);
	CREATE INDEX {{.collection}}right_id ON {{.collection}}(right);
	CREATE INDEX {{.collection}}text ON {{.collection}}(text);
	`

// SQLDropCollection deletes a table permanently - no undoing this
const SQLDropCollection = `DROP TABLE IF EXISTS "{{.collection}}";`

// SQLCloneCollection perform a record by record from source table into (empty) target table
const SQLCloneCollection = "insert into {{.targetCollection}} select * from {{.sourceCollection}};"

// SQLRecursiveAxis reads records recursively to find all descendents of a single node under a specified axis. It keeps a level info.
// This query can be stopped at a maximum level and/or at a limit of returned records count
const SQLRecursiveAxis = `
	WITH RECURSIVE nodes(level, id, typeid, ownerid, flags, created, start, end, deleted, up, left, right, text{{.commaContent}}) AS (
    SELECT
		0, id, typeid, ownerid, flags, created, start, end, deleted,
		up, left, right,
		text{{.commaContent}}
	FROM
    {{.collection}}
	WHERE
		id = ?
	UNION ALL
	SELECT
		level+1,
		p.id, p.typeid, p.ownerid, p.flags, p.created, p.start, p.end, p.deleted, 
		p.up, p.left, p.right, 
		p.text{{.commaPContent}}
	FROM
		{{.collection}} As p
		INNER JOIN nodes AS c ON (p.{{.axis}} = c.id) {{.where}} {{.limit}}
	)
	SELECT
		level,
		id, typeid, ownerid, flags, created, start, end, deleted,
		up, left, right, text{{.commaContent}}
	FROM
		nodes
	ORDER BY {{.axis}}, id;
`

// SQLLoadNode read a single record by id from selected table (collection)
const SQLLoadNode = `select id, typeid, ownerid, flags, created, start, end, deleted, up, left, right, text{{.commaContent}} 
from {{.collection}} where id=?;`

// SQLLoadNodes read a set of record by ids from selected table (collection)
const SQLLoadNodes = `select id, typeid, ownerid, flags, created, start, end, deleted, up, left, right, text{{.commaContent}}
	from {{.collection}} where id in ({{.nodeids}});`

// SQLInsertNode insert a new node into the collection
const SQLInsertNode = `insert into {{.collection}} (typeid, ownerid,flags,created,start,end,deleted,up,left,right,text,content) 
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?);`

// SQLUpdateNode update a single node in the collection
const SQLUpdateNode = `update {{.collection}} set (typeid, ownerid,flags,created,start,end,deleted,up,left,right,text{{.commaContent}}) = 
	(?,?,?,?,?,?,?,?,?,?,?{{.commaQuestionMark}}) where id = ?;`

//SQLListQuery read all nodes under a provided up-ID and optionally narrow down certain type-IDs
const SQLListQuery = `select id, typeid, ownerid, flags, created, start, end, deleted, up, left, right, text{{.commaContent}}
		from {{.collection}} where up = ?{{.typeids}};`
