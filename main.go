package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type customError struct {
	errorMessage string
}

func (ce *customError) Error() string {
	return ce.errorMessage
}

func checkErr(err error, message string) {
	if err != nil {
		fmt.Printf("%s\n", message)
		os.Exit(1)
	}
}

func getArguments() (float64, string, string) {
	args := os.Args
	var (
		from, to string
		amount   float64
	)

	if len(args) < 4 {
		argsError := customError{errorMessage: "Please specify a pair of currencies."}
		checkErr(&argsError, argsError.Error())
	}

	if len(args) < 4 {
		from = args[1]
		to = args[2]
		amount = 1
	} else {
		var err error
		amount, err = strconv.ParseFloat(args[1], 64)
		checkErr(err, "Could not understand the command.")

		from = args[2]
		to = args[3]
	}

	return amount, from, to
}

func getExchangeRate(requestUrl string) float64 {
	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	checkErr(err, "Could not make HTTP request. Check you internet connection, API key, or try again later.")

	res, err := http.DefaultClient.Do(req)
	checkErr(err, "Could not make HTTP request. Check you internet connection, API key, or try again later.")

	resBody, err := io.ReadAll(res.Body)
	checkErr(err, "Could not make HTTP request. Check you internet connection, API key, or try again later.")

	var jsonData map[string]interface{}
	err = json.Unmarshal(resBody, &jsonData)
	checkErr(err, "Could not retrieve data. Try again later.")

	rawValue, ok := jsonData["Realtime Currency Exchange Rate"]
	if !ok {
		fmt.Printf("Could not retrieve data. Try again later.\n")
		os.Exit(1)
	}

	exchangeRate, err := strconv.ParseFloat(rawValue.(map[string]interface{})["5. Exchange Rate"].(string), 64)
	checkErr(err, "Could not retrieve data. Try again later.")

	return exchangeRate
}

func main() {
	if len(os.Args) == 1 {
		argsError := customError{errorMessage: "Please specify arguments."}
		checkErr(&argsError, argsError.Error())
	}

	if os.Args[1] == "key" {
		// Add key
		if len(os.Args) < 3 {
			fmt.Println("You need to specify a key")
			os.Exit(1)
		}

		key := fmt.Sprintf("%s\n", os.Args[2])

		homeDirectory, err := os.UserHomeDir()
		checkErr(err, "Could not save the key. Try again later.")

		f, err := os.Create(fmt.Sprintf("%s/.forex", homeDirectory))
		defer f.Close()
		checkErr(err, "Could not save the key. Try again later.")

		_, err = f.Write([]byte(key))
		checkErr(err, "Could not save the key. Try again later.")
		_ = f.Sync()

		return
	}

	// Check for the key
	homeDirectory, err := os.UserHomeDir()
	checkErr(err, "Could not retrieve the home directory. Please check your installation.")
	apiKeyByes, err := os.ReadFile(fmt.Sprintf("%s/.forex", homeDirectory))
	checkErr(err, "Could not retrieve the API key. Please make sure you add it with the `key` argument.")
	apiKey := fmt.Sprintf("%q", strings.TrimRight(string(apiKeyByes), "\n"))

	// Make request
	amount, from, to := getArguments()
	requestUrl := fmt.Sprintf("http://alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=%s&apikey=%s", from, to, apiKey)
	exchangeRate := getExchangeRate(requestUrl)

	fmt.Printf("%.2f %s = %.2f %s\n", amount, strings.ToUpper(from), amount*exchangeRate, strings.ToUpper(to))
}
