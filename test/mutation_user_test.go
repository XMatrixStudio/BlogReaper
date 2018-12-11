package test

import (
	"encoding/json"
	"github.com/XMatrixStudio/BlogReaper/app"
	"github.com/XMatrixStudio/BlogReaper/resolver"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type Response struct {
	Data interface{} `json:"data"`
}

func TestMutationResolver_CreateLoginURL(t *testing.T) {
	handler := app.TestApp()
	req, err := http.NewRequest("POST", "/graphql", strings.NewReader(`{"query":"mutation{createLoginUrl(backUrl:\"https://blog.xmatrix.studio/\")}"}`))
	if err != nil {
		t.Error(err.Error())
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := Response{}
	json.Unmarshal(w.Body.Bytes(), &res)
	parse, err := url.Parse(res.Data.(map[string]interface{})["createLoginUrl"].(string))
	if err != nil {
		t.Error(err.Error())
	}
	if parse.Host+parse.Path != "oauth.xmatrix.studio/Verify/Authorize/" {
		t.Errorf("Login url should be '%s', actually '%s'", "oauth.xmatrix.studio/Verify/Authorize/", parse.Host+parse.Path)
	}
	query := parse.Query()
	if len(query["responseType"]) != 1 || query["responseType"][0] != "code" {
		t.Errorf("Query responseType error: %v", query["responseType"])
	}
	if len(query["clientId"]) != 1 {
		t.Errorf("Query clientId error: %v", query["clientId"])
	}
	if len(query["state"]) != 1 || query["state"][0] != resolver.DefaultResolver().Session.GetString("state") {
		t.Errorf("State error: %v", query["state"])
	}
	if len(query["redirectUrl"]) != 1 || query["redirectUrl"][0] != "https://blog.xmatrix.studio/" {
		t.Errorf("RedirectUrl error: %v", query["redirectUrl"])
	}
}

func TestMutationResolver_Login(t *testing.T) {
	// TODO
}

func TestMutationResolver_Logout(t *testing.T) {
	// TODO
}
