package dto

type HelloWorldReq struct {
}

type HelloWorldResp struct {
	Hello string `json:"hello"`
	World string `json:"world"`
}
