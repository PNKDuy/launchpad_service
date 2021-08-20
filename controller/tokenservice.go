package controller

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"io/ioutil"
	"launchpad_service/model/response"
	"net/http"
	"strings"
	"time"
)

var priceList []response.Response

// GetAllPrice
// @Summary get price by Token via Binance API
// @Tags token
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {HTTPError} HTTPError
// @Router /token/price [GET]
func GetAllPrice(c echo.Context) error {
	return c.JSON(http.StatusOK, priceList)
}

// GetPrice
// @Summary get price by Token via Binance API
// @Tags token
// @Param token path string true "token"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {HTTPError} HTTPError
// @Router /token/price/{token} [GET]
func GetPrice(c echo.Context) error {
	token := c.Param("token")
	//token = strings.ToUpper(token)
	var result response.Response
	for i := range priceList {
		if strings.EqualFold(priceList[i].Symbol, token) {
			result = priceList[i]
			break
		}
	}

	return c.JSON(http.StatusOK, result)
}

func DoEvery(d time.Duration, f func()error) {
	for range time.Tick(d) {
		f()
	}
}

func GetPriceAndUpdateList() error {
	urls := []string{
		"https://api.binance.com/api/v3/ticker/price",
		"https://api.binance.com/api/v3/ticker/24hr",
	}
	for _, url := range urls {
		result, _ := http.Get(url)

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(body, &priceList); err != nil {
			return err
		}
	}


	//fmt.Println(priceList[0].Price)
	return nil
}


// GetKlines
// @Summary get klines(candlestick) by symbol
// @Tags token
// @Param token path string true "token"
// @Param interval path string true "interval"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {HTTPError} HTTPError
// @Router /token/klines/{token}/{interval} [GET]
func GetKlines(c echo.Context) error {
	interval := c.Param("interval")
	token := c.Param("token")
	token = strings.ToUpper(token)
	url := "https://api.binance.com/api/v3/klines?symbol="
	result, _ := http.Get(url + token + "&interval="+interval)
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, string(body))
}

// GetTransaction
// @Summary get transaction by its hash
// @Tags token
// @Param hash path string true "hash"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {HTTPError} HTTPError
// @Router /token/transaction/{hash} [GET]
func GetTransaction(c echo.Context) error {
	hash := c.Param("hash")

	url := "https://blockchain.info/rawtx/"
	result, _ := http.Get(url+hash)
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": err.Error(),
		})
	}

	msg := &messaging.Message{
		Topic: hash,
		Data: map[string]string{
			"data": string(body),
		},
	}

	opt := option.WithCredentialsFile("generated-private-key.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx,nil, opt)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": err.Error(),
		})
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": err.Error(),
		})
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	res, err := client.Send(ctx, msg)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status": err.Error(),
		})
	}

	// Response is a message ID string.
	return c.JSON(http.StatusOK, map[string]string{
		"status": res,
	})
	//return c.JSON(http.StatusOK, tx)
}


