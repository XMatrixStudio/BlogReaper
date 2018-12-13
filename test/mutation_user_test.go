package test

import (
	"encoding/json"
	"encoding/xml"
	"github.com/XMatrixStudio/BlogReaper/app"
	"github.com/XMatrixStudio/BlogReaper/service"
	"net/http"
	"net/http/httptest"
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

func login(t *testing.T, handler http.Handler, userInfo service.TestLoginParameters, state string) Response {
	bytes, err := xml.Marshal(userInfo)
	if err != nil {
		t.Fatal(err.Error())
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
	res := login(t, handler, loginUser1, "1234567890")
	if len(res.Errors) == 0 || res.Errors[0].Message != "error_state" {
		t.Fatalf("Error expected %s, actual %s", "error_state", res.Errors[0].Message)
	}
}

func TestMutationResolver_Login_AddUser(t *testing.T) {
	handler := app.TestApp()
	res := login(t, handler, loginUser3, "")
	if len(res.Errors) != 0 {
		t.Fatal(res.Errors)
	}
	data := res.Data.(map[string]interface{})["login"].(map[string]interface{})
	if data["email"] != loginUser3.Email {
		t.Fatalf("Email expected %s, actual %s", loginUser3.Email, data["email"])
	} else if data["info"].(map[string]interface{})["name"] != loginUser3.Name {
		t.Fatalf("Name expected %s, actual %s", loginUser3.Name, data["info"].(map[string]interface{})["name"])
	}
}

func TestMutationResolver_Login_SetToken(t *testing.T) {
	handler := app.TestApp()
	res := login(t, handler, loginUser1, "")
	if len(res.Errors) != 0 {
		t.Fatal(res.Errors)
	}
	data := res.Data.(map[string]interface{})["login"].(map[string]interface{})
	if data["email"] != loginUser1.Email {
		t.Fatalf("Email expected %s, actual %s", loginUser1.Email, data["email"])
	} else if data["info"].(map[string]interface{})["name"] != loginUser1.Name {
		t.Fatalf("Name expected %s, actual %s", loginUser1.Name, data["info"].(map[string]interface{})["name"])
	}
	res = login(t, handler, loginUser2, "")
	if len(res.Errors) != 0 {
		t.Fatal(res.Errors)
	}
	data = res.Data.(map[string]interface{})["login"].(map[string]interface{})
	if data["email"] != loginUser1.Email {
		t.Fatalf("Email expected %s, actual %s", loginUser1.Email, data["email"])
	} else if data["info"].(map[string]interface{})["name"] != loginUser1.Name {
		t.Fatalf("Name expected %s, actual %s", loginUser1.Name, data["info"].(map[string]interface{})["name"])
	}
}

func logout(t *testing.T, handler http.Handler) Response {
	query := `
mutation {
	logout
}`
	query = strings.Replace(query, "\n", " ", -1)
	query = strings.Replace(query, "\t", " ", -1)
	req, _ := http.NewRequest("POST", "/graphql", strings.NewReader(`{"query":"`+query+`"}`))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	res := Response{}
	json.Unmarshal(w.Body.Bytes(), &res)
	return res
}

func TestMutationResolver_Logout(t *testing.T) {
	handler := app.TestApp()
	res := login(t, handler, loginUser1, "")
	if len(res.Errors) != 0 {
		t.Fatal(res.Errors)
	}
	res = logout(t, handler)
	if len(res.Errors) != 0 {
		t.Fatal(res.Errors)
	}
	data := res.Data.(map[string]interface{})["logout"].(bool)
	if data == false {
		t.Fatal("Logout fail")
	}
}
