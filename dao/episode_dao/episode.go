package episode_dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type Episode struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TopicId primitive.ObjectID `json:"top_id,omitempty" bson:"top_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Images  []string           `json:"images,omitempty" bson:"image,omitempty"`
	Url     string             `json:"url,omitempty" bson:"url,omitempty"`
}
