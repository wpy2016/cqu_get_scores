package server

import (
	"log"
	"net/http"
	"cqu/rank"
	"cqu/score"
)

const (
	HOST_URL        = "localhost:9296"
	SCORES          = "/scores/"
	RANK_NEW_SCORES = "/rank/new_scores/"
	RANK_ALL_SCORES = "/rank/all_scores/"
)

func StartServer() {
	http.HandleFunc(SCORES, score.ResponseScores)
	http.HandleFunc(RANK_NEW_SCORES, rank.ResponseRank)
	http.HandleFunc(RANK_ALL_SCORES, rank.ResponseRank)
	log.Fatal(http.ListenAndServe(HOST_URL, nil))
}

func AddHeader(w http.ResponseWriter)  {
	w.Header().Add("Access-Control-Allow-Origin", "http://cdqn.ccj.cqu.edu.cn")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
}