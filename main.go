package main

/*
	This program runs an unencrypted http server on a local port.

	The server uses http.ServeMux to make it easy to add more endpoints in the future.

	Each endpoint should require POST requests with a JSON body. The request
	includes many rows of data, and the response should include result of
	running some computation on each row of data..

	Every endpoint will receive a JSON object structured like this:
	{
		"data": [
			[ ROWNUM, VALUE, VALUE, ... ],
			[ ROWNUM, VALUE, VALUE, ... ],
			...
		]
	}

	Every endpoint will respond with a JSON object structured like this:
	{
		"data": [
			[ ROWNUM, RESPONSE ],
			...
		]
	}

	If an endpoint wants to return a JSON object as the RESPONSE, it must
	encode it as a string. This is because SingleStore external functions
	doesn't yet support nested JSON. This limitation will be resolved in a
	future release.
*/

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/sugarme/tokenizer/pretrained"
)

type ReqResp struct {
	Data [][]interface{} `json:"data"`
}

var (
	BertModel = pretrained.BertBaseUncased()
)

func TokenizeText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReqResp
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := make([][]interface{}, len(req.Data))
	for i, row := range req.Data {
		valueString, ok := row[1].(string)
		if !ok {
			http.Error(w, "value must be a string", http.StatusBadRequest)
			return
		}

		encoded, err := BertModel.EncodeSingle(valueString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// SingleStore external functions doesn't yet support nested JSON
		// objects in a response so we need to encode each value as a string.
		marshalled, err := json.Marshal(map[string]interface{}{
			"tokens":  encoded.Tokens,
			"offsets": encoded.Offsets,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		out[i] = []interface{}{row[0], string(marshalled)}
	}

	resp := ReqResp{Data: out}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/text/tokenize", TokenizeText)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
