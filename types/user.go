package types

type User struct {
	ID        int    `bson:"_id" json:"id", omitempty`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
