package types

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
