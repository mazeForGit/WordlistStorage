package routers

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"
	
	"github.com/gin-gonic/gin"
	data "github.com/mazeForGit/Wordlist/data"
)

func Test_WordListGET_emptyData(t *testing.T) {
	router := gin.Default()
	router.GET("/wordlist", WordListGET)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/wordlist", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("You received a %v error.", w.Code)
	}

	var wl data.WordList
	json.Unmarshal(w.Body.Bytes(), &wl)
	if wl.LastUsedId != 0 && wl.Count != 0 {
		t.Errorf("got wrong response .. wl=%v", wl)
	}
}
