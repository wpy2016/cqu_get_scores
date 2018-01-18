package score

import (
	"fmt"
	"strconv"
	"net/http"
	"log"
	"net/url"
	"zhihui/config"
	"cqu/server"
)

const (
	URLARG_ID       = "userid"
	URLARG_PASS     = "password"
	URLARG_SEMESTER = "semester"
	URLARG_YEAR     = "year"
)

func ResponseScores(w http.ResponseWriter, r *http.Request) {
	server.AddHeader(w)
	defer func() {
		if p := recover(); p != nil {
			result, _ := CreateErrorResponce("444", fmt.Sprintf("%v", p))
			fmt.Fprintf(w, "%s", result)
			log.Printf("main : ResponseScores : error=[%v]", p)
		}
	}()
	kv, err := spliteScoresUrlArguments(r.URL.Query())
	if err != nil {
		result, _ := CreateErrorResponce("404", err.Error())
		fmt.Fprintf(w, "%s", result)
		return
	}
	err = validateYearAndSemester(kv[URLARG_YEAR], kv[URLARG_SEMESTER])
	if err != nil {
		result, _ := CreateErrorResponce("404", err.Error())
		fmt.Fprintf(w, "%s", result)
		return
	}
	fmt.Println(URLARG_ID+"=", kv[URLARG_ID], URLARG_YEAR+"=", kv[URLARG_YEAR], URLARG_SEMESTER+"=", kv[URLARG_SEMESTER])
	result, _, err := StuScore(kv[URLARG_ID], kv[URLARG_PASS], kv[URLARG_YEAR], kv[URLARG_SEMESTER])
	if err != nil {
		result, _ = CreateErrorResponce("444", fmt.Sprintf("%v", err))
		fmt.Fprintf(w, "%s", result)
		return
	}
	fmt.Fprintf(w, "%s", result)
}

func spliteScoresUrlArguments(value url.Values) (map[string]string, error) {
	kv := make(map[string]string)
	idValue := value.Get(URLARG_ID)
	if idValue == "" {
		return nil, fmt.Errorf("%s", URLARG_ID+" argument error.")
	}
	kv[URLARG_ID] = idValue

	passValue := value.Get(URLARG_PASS)
	if passValue == "" {
		return nil, fmt.Errorf("%s", URLARG_PASS+" argument error.")
	}
	kv[URLARG_PASS] = passValue

	yearValue := value.Get(URLARG_YEAR)
	if yearValue == "" {
		return nil, fmt.Errorf("%s", URLARG_YEAR+" argument error.")
	}
	kv[URLARG_YEAR] = yearValue

	semesterValue := value.Get(URLARG_SEMESTER)
	if semesterValue == "" {
		return nil, fmt.Errorf("%s", URLARG_SEMESTER+" argument error.")
	}
	kv[URLARG_SEMESTER] = semesterValue
	return kv, nil
}

func validateYearAndSemester(year, semester string) error {
	if _, errY := strconv.Atoi(year); errY != nil {
		return fmt.Errorf("%s", URLARG_YEAR+" value error.")
	}
	if n, errN := strconv.Atoi(semester); errN == nil {
		if n < 0 || n > 1 {
			return fmt.Errorf("%s", URLARG_SEMESTER+" value error.out of 0 or 1")
		}
	} else {
		return fmt.Errorf("%s", URLARG_SEMESTER+" value error.out of 0 or 1")
	}
	return nil
}
