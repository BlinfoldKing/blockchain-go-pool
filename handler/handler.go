package handler

import (
	"github.com/blinfoldking/blockchain-go-pool/handler/graphql"
)

type Handler struct {
	graphql.GraphQLHandler
}

func New() Handler {
	h := Handler{}

	return h
}
