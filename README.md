# nodedb
NodeDB - a simple yet powerful persistency layer.

## Why NodeDB?

NodeDB was created to serve my need to have a simple yet powerful persistency layer. After decades of writing software and using ORMs, I grew tired of designing data base schemas, object hierarchies and doing all the plumbing and wiring up to overcome the object relational-impedance-mismatch-problem, again and again, ...

After many years of working with relational data bases, I got interested in exploring different data base concepts, for example key-value stores, document databases and in particular graph databases. They all have their strengths and draw backs. Wouldn't it be great to combine all strengths and avoid the weaknesses? Somehow?

PostgreSQL, Oracle, MySQL, MongoDB, DynamoDB, Couchbase, Neo4j, Cassandra, ... these databases are all big and heavy players and I wanted something really small and simple, light-weight but capable to help me getting my job done.

As with all software, data needs to be written into file(s) eventually and I didn't want to write the io-code to do so myself. So I picked a relational data base as the lowest layer of abstraction. In principal it could be any SQL data base but it should provide Common Table Expressions. My pick is SQLite! It is fast, reliable, modern and extremely easy to use/administer. And it should be fairly easy to move to a "bigger" database like PostgreSQL, should the need arise later.

In order to overcome the object-relational gap, there are simplications and compromises to be made on both sides:

### Relational side
There is only on data schema: the "Node". This node is the unit of data exchange, all data is expressed as nodes, the data base excepts only nodes and delivers only nodes.

#### Schema of a node:
	
|   Field  | Type                  | Description                                                           |
|--------- |-----------------------|---------------------------------------------------------------------- |
| id       | Integer, Primary Key  | unique, automatically assigned                                        |
| typeid   | Integer, Indexed      | schema node axis                                                      |
| ownerid  | Integer, Indexed      | owner axis                                                            |
| up       | Integer, Indexed      | structural axis: tree, list, composite, master-detail, ...            |
| left     | Integer, Indexed      | edge axis to the left  - uni- or bidirectional links between nodes    |
| right    | Integer, Indexed      | edge axis to the right - uni- or bidirectional links between nodes    |
| flags    | Integer               | top 31 bits reserved, low 32 bits for custom use                      |
| text     | String,  Indexed      | any text, node name, etc. could be used with like operator            |
| created  | DateTime              | creation date/time of this node, automatically set to current date/time |
| start    | DateTime              | custom date/time, initialized with created date/time                  |
| end      | DateTime              | custom date/time, initialized with deleted date/time                  |
| deleted  | DateTime              | deletion data/time initialized to a date far into the future and set to the current date/time on deletion |
| content  | Blob/[]Byte           | payload of this node, can be anything, i.e. an image, a document, a serialized object, ... |

This simple schema can provide for all my needs in designing data structures, lists, documents, trees, composite, graphs, ... more on that later.

With only one schema in the data base there is only one set of SQL-Statements needed to handle all data updates and queries. And everything could be done with only one table. Using only one table felt like a loss of design "ressources". So I came up with the idea to "call" tables collections (of nodes). The structur  Database - Tables - Records becomes Connection - Collections - Nodes. The node-ID is unique only within its collection, copying or moving nodes between collections requires additional adjustment work of source and target ids.

### API:
 
* Connection:
	* Open(filename), Close(), Collections(), Collection(name), RemoveCollection(name), CloneCollection(source_name, target_name), Ping()
* Collection:
	* NewNode(), NewNodes(number), LoadNode(id), LoadNodes([]id), UpdateNode(node), UpdateNodes([]node), DeleteNode(node), DeleteNodes([]node)
* Node:
	* Public Fields + IsContentLoaded(), SetContentLoaded(bool), GetFlag(), SetFlag(), ...
* Query:
	* QueryAxes recursively, QueryList, QueryNeighbour, ... work in progress

* Graph-Algorithms:
	* Shortest Path, Minimum Spanning Tree, ... future work
* ...

## The Object-Oriented side

There is no help from an ORM!
Put up with the work of serializing your objects into a single unit of data exchange.
The objects don't get filled with data automagically, but the process of storing and loading is fairly straightforward
(the implemention language is Go, this makes it possible to directly embed a node as an anonymous property -> easy access to its data fields and content).


--------------

The current work represents only the lowest possible layer of abstractions and even this is still a long way away from completion.

The current design principles are:

- keep it as simple as possible, but not simpler, keep as many future options open as possible.
- use every part to its best performance and do the work where it fits best, i.e. let the underlying database do as much work as possible.
- provide the means to explore what can be done with a single unit of data exchange.
- abstract away the database layer, stay in the realm of writing code for the most of time.
- create and persist arbitrary object structures
- deletion should not be scary and be undone easily
- every information has a lifecycle: a creation and deletion interval and a custom start-end interval to i.e. indicate when an information becomes valid or invalid or a duration for a task, project, etc.
- there are three hierarchical axes: up, type and owner, to build type and ownership hierarchies and any tree-like structures with the up axis.
- the is a left-right edge axis to build i.e double-linked lists or relationships between nodes like links in graph-databases.
- put anything into the content field: object properties, json-schema, json-data, images, text documents, html, markdown, ...
- create type nodes that provide schema information and use their id as the typeid in the corresponding data nodes.
- it is easy to create queries and all kinds of algorithms on top of NodeDB, it is just SQL and the underlying database is readily available.



## Intended uses:

A quick and simple persistency layer for a toy project or prototype.
Simple Data-Management for a web site, backend for ... see JAMStack (Javascript - API - Markup)
A Data-Layer for multiple client web apps, where each client gets their own single collection, this should provide a good data storage solution and great separation between clients.

Ideas:
- Put a GraphQL-Server-Layer on top of NodeDB, that would be really cool!
- Use JSON-Schema to provide Schema, Validation (and UI-Information for the frontend) in schema-nodes. Data nodes have corresponding JSON-data in their content field.
- Create Composite Objects, i.e. Person with Adress, Email, Phone- Information embedded. Adress, Email and Phone-Information are individual nodes themselves, or projects with master-detail-task information.


This is just the start of my journey, I want to create a second-brain-personal-information-project-management-system on top of NodeDB! Care to join?

Cheers, Thomas
