package resolver

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type JSON map[string]interface{}

func (JSON) ImplementsGraphQLType(name string) bool { return name == "JSON" }
func (j *JSON) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case JSON:
		*j = input
		return nil
	default:
		logrus.Error("wrong type")
		return errors.New("wrong type")
	}
}
