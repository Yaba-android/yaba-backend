package main

// Author struct
type Author struct {
	RemoteId  string
	ImagePath string
	Name      string
	Desc      string
	BooksId   []string
}

/* Exemple JSON
{
    "ImagePath": "john_doe.png",
    "Name": "John Doe",
	"Desc": "John doe is a famous french author.",
	"BooksId": [
		"id1", "id2", "id3"
	]
}
*/
