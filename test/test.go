package test

import (
	"github.com/XMatrixStudio/BlogReaper/service"
	"github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/globalsign/mgo/bson"
	"os"
)

var loginUser1 service.TestLoginParameters
var loginUser2 service.TestLoginParameters
var loginUser3 service.TestLoginParameters

func init() {
	err := os.Chdir("..")
	if err != nil {
		panic(err)
	}
	loginUser1 = service.TestLoginParameters{
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
	loginUser2 = service.TestLoginParameters{
		TokenRes: violetSdk.TokenRes{
			UserID: loginUser1.UserID,
			Token:  "00000000000",
		},
		UserInfoRes: violetSdk.UserInfoRes{
			Email: "b@xmatrix.studio",
			Name:  "XMatrix2",
			Info:  violetSdk.UserInfo{},
		},
	}
	loginUser3 = service.TestLoginParameters{
		TokenRes: violetSdk.TokenRes{
			UserID: bson.NewObjectId().Hex(),
			Token:  "3216549870",
		},
		UserInfoRes: violetSdk.UserInfoRes{
			Email: "admin@xmatrix.studio",
			Name:  "Admin",
			Info:  violetSdk.UserInfo{},
		},
	}
}
