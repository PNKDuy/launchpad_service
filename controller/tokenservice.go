package controller

import (
	"context"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/jmoiron/jsonq"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"io/ioutil"
	"launchpad_service/model/response"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var priceList []response.Response
var klinesList [][]interface{}
var urls = []string{
		"https://api.binance.com/api/v3/ticker/price",
		"https://api.binance.com/api/v3/ticker/24hr",
}
var binanceUrls = []string{
	"https://api.binance.com/api/v3/ticker/price?symbol=",
	"https://api.binance.com/api/v3/ticker/24hr?symbol=",
}
var mux = &sync.Mutex{}


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
 	token = strings.ReplaceAll(token, "%2C", ",")

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

// GetPriceAdvance
// @Summary get price by Token via Binance API
// @Tags token
// @Param token path string false "token"
// @Param currency path string true "currency"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {HTTPError} HTTPError
// @Router /token/price-by-advance/{token}/{currency} [GET]
func GetPriceAdvance(c echo.Context) error{
	token := strings.ToUpper(c.Param("token"))
	currency := strings.ToLower(c.Param("currency"))
	token = strings.ReplaceAll(token, "%2C", ",")
	if strings.EqualFold(token, "undefined") || strings.EqualFold(token, "{token}") {
		token = strings.ToUpper(response.DefaultList)
	}

	coinGeckoUrl := "https://api.coingecko.com/api/v3/simple/price?ids=tether&vs_currencies="
	res, _ := http.Get(coinGeckoUrl + currency)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var currencyPrice map[string]interface{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	_ = dec.Decode(&currencyPrice)
	jq := jsonq.NewQuery(currencyPrice)
	pairUsdtToken, _ := jq.Float("tether", currency)

	var curList []response.Currency
	tokenList := strings.Split(token,  ",")
	for i := range tokenList {
		for j := range priceList {
			if strings.EqualFold(tokenList[i] + "USDT", priceList[j].Symbol){
				tokenUsdtPrice, _ := strconv.ParseFloat(priceList[i].Price, 64)
				priceChange, _ := strconv.ParseFloat(priceList[i].PriceChange, 64)
				volume, _ := strconv.ParseFloat(priceList[i].Volume, 64)
				highPrice, _ := strconv.ParseFloat(priceList[i].HighPrice, 64)
				lowPrice, _ := strconv.ParseFloat(priceList[i].LowPrice, 64)

				price := pairUsdtToken * tokenUsdtPrice
				var cur response.Currency
				cur.Price = price
				cur.Symbol = tokenList[i]
				cur.PriceChange= priceChange * pairUsdtToken
				cur.PriceChangePercent = priceList[i].PriceChangePercent
				cur.Volume = volume * price
				cur.HighPrice = highPrice * pairUsdtToken
				cur.LowPrice = lowPrice * pairUsdtToken
				curList = append(curList, cur)
			}
		}

	}
	return c.JSON(http.StatusOK, curList)
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
		ticker := time.Tick(d)
		for range ticker {
			err := f()
			if err != nil {
				break
			}
		}
}

func GetPriceAndUpdateList() error {
	for _, url := range urls {
		err := getAPI(url)
		if err != nil {
			break
		}
	}
	return nil
}

func getAPI(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("154", err)
		return err
	}
	client := &http.Client{ }

	resp, err := client.Do(req)
	if err != nil {
		log.Println("162", err)
		return err
	}


	if strings.EqualFold(resp.Status, "429 Too Many Requests") {
		resp.Body.Close()
		return errors.New("429 Too Many Requests")
	}

	if err := json.NewDecoder(resp.Body).Decode(&priceList); err != nil {
		log.Println("176", err)
		resp.Body.Close()
		return err
	}
	resp.Body.Close()
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
	if err = json.Unmarshal(body, &klinesList); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, klinesList)
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


