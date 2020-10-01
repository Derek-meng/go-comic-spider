package topic

import "go.mongodb.org/mongo-driver/bson/primitive"

type Topic struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	WebId primitive.ObjectID `json:"web_id,omitempty" bson:"web_id,omitempty"`
	Title string             `json:"title,omitempty" bson:"title,omitempty"`
	Url   string             `json:"url,omitempty" bson:"url,omitempty"`
}
