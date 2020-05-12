Experiments are provides as examples for how to use NodeDB's basic functions. NodeDB is designed to provide a performant and reliable way to store
application data.

The experiments follow the same pattern:
    - if a database file exists? delete it
    - open a fresh database
    - create collection(s)
    - fill collection(s)
    - do some experiments, load nodes, modify nodes, run queries, etc.
    - close database

this leaves a sqlite database file in folder to be there for further inspection using i.e. sqlitebrowser.

the .log, .db are all automatically named after the executable. They can be deleted comfortably by i.e. rm tree1* 
! Don't name your source files according to this scheme!

To experiment with different parts just comment or uncomment function calls in the main function.

The experiments are not intended as unit tests, but as a way to explore NodeDB while it is still under heavy and possible code breaking development. 
I currently use them to implement NodeDB.

A test suite will follow, once the dust has settled a bit.