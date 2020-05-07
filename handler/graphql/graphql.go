package graphql

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/blinfoldking/blockchain-go-pool/resolver"
	"github.com/graph-gophers/graphql-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type GraphQLHandler struct{}

// GET server graphql playground
func (handler *GraphQLHandler) Playground(c echo.Context) error {
	c.HTML(http.StatusOK, page)
	return nil
}

// POST graphql query
func (handler *GraphQLHandler) Query(c echo.Context) error {
	s, err := getSchema("schema/schema.gql")
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	resolver.ResolverConnecetion = resolver.Init()
	err = GqlResponse(c, graphql.MustParseSchema(s, resolver.ResolverConnecetion))
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func GqlResponse(ctx echo.Context, Schema *graphql.Schema) error {

	r := ctx.Request()

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	// var graphqlContext context.Context
	// authorization := ctx.Request().Header.Get("Authorization")
	// if authorization != "" {
	// 	logrus.Println(authorization)
	// 	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
	// 		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 			return nil, fmt.Errorf("Signing method invalid")
	// 		} else if method != jwt.SigningMethodHS256 {
	// 			return nil, fmt.Errorf("Signing method invalid")
	// 		}

	// 		return []byte("secret"), nil
	// 	})

	// 	if err != nil {
	// 		logrus.Println(">>>>>>>>>>>>>>>>>>>")
	// 		return err
	// 	}

	// 	claims, ok := token.Claims.(jwt.MapClaims)
	// 	if !ok {
	// 		logrus.Error("not ok when when converting token")
	// 		return errors.New("not ok when when converting token")
	// 	}

	// 	logrus.Println(">>>>>>>>>>>>>>>>>>>", claims)
	// 	logrus.Println(">>>>>", claims["user"].(string))
	// 	if claims["user"] != nil {
	// 		graphqlContext = context.WithValue(context.Background(), "user", claims["user"].(string))
	// 	}
	// } else {
	// 	graphqlContext = r.Context()
	// }
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	graphqlContext := r.Context()

	response := Schema.Exec(graphqlContext, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	ctx.Response().Write(responseJSON)

	return nil
}

func getSchema(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
