package parser

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strings"
)

type MoexResponse struct {
	Secur struct {
		Collums []string
		Data    [][]interface{} `json:"data"`
	} `json:"securities"`
}

func Search(ctx echo.Context) error {
	link := "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities.json"

	res, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	var moexResponse MoexResponse
	err = json.Unmarshal(body, &moexResponse)
	if err != nil {
		fmt.Println(err)
	}
	company := make(map[string]string)
	for _, v := range moexResponse.Secur.Data {
		if len(v) > 0 {
			ticker, ok := v[0].(string)
			if !ok {
				return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "missing code"})
			}
			name, ok := v[9].(string)
			if !ok {
				return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "missing company name"})
			}
			if strings.Contains(name, "ао") {
				name = strings.Replace(name, " ао", " обыкновенные акции ", 1)
			}
			if strings.Contains(name, "ап") {
				name = strings.Replace(name, " ап", " привилегированные акции ", 1)

			}
			company[name] = ticker

		}
	}

	return ctx.JSON(http.StatusOK, company)

}
