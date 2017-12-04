package rank

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	URLARG_USERNAME = "userid"
	URLARG_PASSWORD = "password"
)

func ResponseRank(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	defer func() {
		if p := recover(); p != nil {
			result, _ := CreateResponce("444", fmt.Sprintf("%v", p))
			fmt.Fprintf(w, "%s", result)
			log.Printf("main : ResponseRank : error=[%v]", p)
		}
	}()
	kv, err := spliteRankUrlArguments(r.URL.Query())
	if err != nil {
		result, _ := CreateResponce("404", err.Error())
		fmt.Fprintf(w, "%s", result)
		return
	}
	path := r.URL.Path
	fmt.Printf("path=%s,%s=%s\n", path, URLARG_USERNAME, kv[URLARG_USERNAME])
	var result string
	if path == "/rank/new_scores/" {
		result, err = StuScores(kv[URLARG_USERNAME], kv[URLARG_PASSWORD], 0)
	}
	if path == "/rank/all_scores/" {
		result, err = StuScores(kv[URLARG_USERNAME], kv[URLARG_PASSWORD], 1)
	}
	if err != nil {
		result, _ = CreateResponce("444", fmt.Sprintf("%v", err))
		fmt.Fprintf(w, "%s", result)
		return
	}
	fmt.Fprintf(w, "%s", result)
}

func spliteRankUrlArguments(value url.Values) (map[string]string, error) {
	kv := make(map[string]string)
	usernameValue := value.Get(URLARG_USERNAME)
	if usernameValue == "" {
		return nil, fmt.Errorf("%s", URLARG_USERNAME+" argument error.")
	}
	kv[URLARG_USERNAME] = usernameValue
	passwordValue := value.Get(URLARG_PASSWORD)
	if passwordValue == "" {
		return nil, fmt.Errorf("%s", URLARG_PASSWORD+" argument error.")
	}
	kv[URLARG_PASSWORD] = passwordValue
	return kv, nil
}
