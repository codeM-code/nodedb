package nodedb

type Flag int64

const (
	First         Flag = 1 << (62 - iota) // first node in a sequence of nodes (the head of a double-linked-list)
	Bidirectional                         // in relational node this points in both directions, not the usual left->right
	Content                               // node has a Content, this needs to be loaded separately, needs to be set if one of the following is set
	// Jsondata                              // node is a (regular) data node, read its type node for schema and additional info
	// Document                              // node is a document, its type node Text-Field contains its Mime-Type String. Textfield contains document name
	// String                                // node is just a string for reference, not intended as a document, use document node for that purpose
	// Binary                                // node contains a byte array,
)

func (f Flag) Check(flag Flag) bool {
	return (f & flag) != 0
}

func (f *Flag) Set(flag Flag, value bool) {
	if value {
		*f = *f | flag // sets the flag
	} else {
		*f = *f &^ flag // clears the flag (AND with NOT(flag))
	}
}

func (f Flag) IsFirst() bool {
	return f.Check(First)
}
func (f Flag) IsBidirectional() bool {
	return f.Check(Bidirectional)
}

// func (f Flag) IsContent() bool {
// 	return f.Check(Content)
// }
// func (f Flag) IsJsondata() bool {
// 	return f.Check(Jsondata)
// }
// func (f Flag) IsDocument() bool {
// 	return f.Check(Document)
// }
// func (f Flag) IsString() bool {
// 	return f.Check(String)
// }
// func (f Flag) IsBinary() bool {
// 	return f.Check(Binary)
// }
