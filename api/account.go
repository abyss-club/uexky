package api

import (
	"net/http"

	"github.com/CrowsT/uexky/model"
	"github.com/julienschmidt/httprouter"
)

func accountApis(r *httprouter.Router) {
	r.GET("/account/", userInfo)
	r.POST("/account/", newUser)
}

func newUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	account, _ := model.NewAccount()
	jsonRes(w, &account)
}

func userInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	account, err := model.GetAccount(token)
	if err != nil {
		errRes(w, err)
	}
	jsonRes(w, account)
}
