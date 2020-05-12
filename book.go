package main

// Book struct
type Book struct {
	RemoteId        string
	ImagePath       string
	Title           string
	AuthorName      string
	AuthorId        string
	Rating          string
	NumberRating    string
	Price           string
	Length          string
	Genre           string
	FileSize        string
	Country         string
	DatePublication string
	Publisher       string
	Resume          string
	FilePath        string
}

/* Exemple JSON
{
    "ImagePath": "camus_la_peste.png",
    "Title": "camus la peste2",
    "Author": "John Doe",
    "Rating": "4",
    "NumberRating": "35",
    "Price": "0",
    "Length": "324",
    "Genre": "Roman",
    "FileSize": "0.85",
    "Country": "France",
    "DatePublication": "10/09/2015",
    "Publisher": "Publish Inc.",
    "Resume": "Super livre de fou",
    "FilePath": "camus_la_peste.epub"
}
*/
