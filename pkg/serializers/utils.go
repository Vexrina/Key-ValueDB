package serializers

type Name struct {
	FirstName  string
	SecondName string
	ThirdName  string
}

type Person struct {
	Name   Name 
	Age    int  
	Active bool 
}

type User struct {
	UUID   int    
	Person Person 
}