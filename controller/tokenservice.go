package controller

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/jmoiron/jsonq"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"io/ioutil"
	"launchpad_service/model/response"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var priceList []*response.Response


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

// GetPriceByCurrency
// @Summary get price by Token via Binance API
// @Tags token
// @Param token path string false "token"
// @Param currency path string true "currency"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {HTTPError} HTTPError
// @Router /token/price-by-currency/{token}/{currency} [GET]
func GetPriceByCurrency(c echo.Context) error{
 	token := c.Param("token")
 	currency := strings.ToLower(c.Param("currency"))
 	var resList []response.Currency
 	token = strings.ReplaceAll(token, "2%C", ",")

	url := "https://api.coingecko.com/api/v3/simple/price?ids=tether&vs_currencies="
	result, _ := http.Get(url + currency)

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	var currencyPrice map[string]interface{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	_ = dec.Decode(&currencyPrice)

	jq := jsonq.NewQuery(currencyPrice)
	pairUsdtToken, _ := jq.Float("tether", currency)

	if strings.EqualFold(token, "undefined") || strings.EqualFold(token, "{token}") {
		token = response.DefaultList
	}

	var tokenUsdtPrice, priceChange, volume, highPrice, lowPrice float64
	var priceChangePercent string

	tokenList := strings.Split(token,  ",")
	for i := range tokenList {
		tokenUsdtPair := tokenList[i] + "USDT"

		for i := range priceList {
			if strings.EqualFold(tokenUsdtPair, priceList[i].Symbol) {
				tokenUsdtPrice, _ = strconv.ParseFloat(priceList[i].Price, 64)
				priceChange, _ = strconv.ParseFloat(priceList[i].PriceChange, 64)
				priceChangePercent = priceList[i].PriceChangePercent
				volume, _ = strconv.ParseFloat(priceList[i].Volume, 64)
				highPrice, _ = strconv.ParseFloat(priceList[i].HighPrice, 64)
				lowPrice, _ = strconv.ParseFloat(priceList[i].LowPrice, 64)
				break
			}
		}

		price := pairUsdtToken * tokenUsdtPrice
		var res response.Currency
		res.Price = price
		res.Symbol = tokenList[i]
		res.PriceChangePercent = priceChangePercent
		res.PriceChange= priceChange * pairUsdtToken
		res.Volume = volume * price
		res.HighPrice = highPrice * pairUsdtToken
		res.LowPrice = lowPrice * pairUsdtToken
		resList = append(resList, res)
	}
	
	return c.JSON(http.StatusOK, resList)
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
			result = *priceList[i]
			break
		}
	}

	return c.JSON(http.StatusOK, result)
}

func DoEvery(d time.Duration, f func()error) {
	ticker := time.NewTicker(d)
	for range ticker.C {
		f()
	}
	ticker.Stop()
}

func GetPriceAndUpdateList() error {
	urls := []string{
		"https://api.binance.com/api/v3/ticker/price",
		"https://api.binance.com/api/v3/ticker/24hr",
	}
	for _, url := range urls {
		result, err := http.Get(url)
		if err != nil {
			return err
		}

		body, err := ioutil.ReadAll(result.Body)
		if body == nil {
			break
		}
		if err != nil {
			return err
		}
		if err = json.Unmarshal(body, &priceList); err != nil {
			return err
		}
		result.Body.Close()
	}

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


