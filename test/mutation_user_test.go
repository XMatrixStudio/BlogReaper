package test

import (
	"encoding/json"
	"encoding/xml"
	"github.com/XMatrixStudio/BlogReaper/app"
	"github.com/XMatrixStudio/BlogReaper/resolver"
	"github.com/XMatrixStudio/BlogReaper/service"
	"github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type Response struct {
	Errors []Error     `json:"errors"`
	Data   interface{} `json:"data"`
}

type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
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

func createLoginUrl(t *testing.T, handler http.Handler) (state string) {
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
	query := parse.Query()
	if len(query["state"]) != 1 {
		t.Errorf("State error: %v", query["state"])
	}
	return query["state"][0]
}

func login(t *testing.T, handler http.Handler, userInfo service.TestLoginParameters, state string) Response {
	bytes, err := xml.Marshal(userInfo)
	if err != nil {
		t.Error(err.Error())
	}
	code := strings.Replace(string(bytes), `"`, `\\"`, -1)
	query := `
mutation {
	login(code: \"` + code + `\", state: \"` + state + `\") {
		email
		info {
			name
			avatar
			bio
			gender
		}
	}
}`
	query = strings.Replace(query, "\n", " ", -1)
	query = strings.Replace(query, "\t", " ", -1)
	req, err := http.NewRequest("POST", "/graphql", strings.NewReader(`{"query":"`+query+`"}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := Response{}
	json.Unmarshal(w.Body.Bytes(), &res)
	return res
}

func TestMutationResolver_Login_ErrorState(t *testing.T) {
	handler := app.TestApp()
	loginParams := service.TestLoginParameters{
		TokenRes: violetSdk.TokenRes{
			UserID: bson.NewObjectId().Hex(),
			Token:  "1234567890",
		},
		UserInfoRes: violetSdk.UserInfoRes{
			Email: "a@xmatrix.studio",
			Name:  "XMatrix",
			Info:  violetSdk.UserInfo{},
		},
	}
	res := login(t, handler, loginParams, "1234567890")
	if len(res.Errors) == 0 || res.Errors[0].Message != "error_state" {
		t.Errorf("Error expected %s, actual %s", "error_state", res.Errors[0].Message)
	}
}

func TestMutationResolver_Login_CreateUser(t *testing.T) {
	handler := app.TestApp()
	loginParams := service.TestLoginParameters{
		TokenRes: violetSdk.TokenRes{
			UserID: bson.NewObjectId().Hex(),
			Token:  "1234567890",
		},
		UserInfoRes: violetSdk.UserInfoRes{
			Email: "a@xmatrix.studio",
			Name:  "XMatrix",
			Info:  violetSdk.UserInfo{},
		},
	}
	res := login(t, handler, loginParams, "")
	if len(res.Errors) != 0 {
		t.Error(res.Errors)
	}
	data := res.Data.(map[string]interface{})["login"].(map[string]interface{})
	if data["email"] != loginParams.Email {
		t.Errorf("Email expected %s, actual %s", loginParams.Email, data["email"])
	} else if data["info"].(map[string]interface{})["name"] != loginParams.Name {
		t.Errorf("Name expected %s, actual %s", loginParams.Name, data["info"].(map[string]interface{})["name"])
	}
}

func TestMutationResolver_Logout(t *testing.T) {
	// TODO
}
