package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"encoding/csv"
	"encoding/json"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"code.google.com/p/go.net/publicsuffix"
)

const (
	PREVIEW int = 0
	REFRESH int = 1
	EDIT    int = 2
	ERROR   int = 3
)

const VALUE_SEPARATOR = ","
const KEY_VALUE_SEPARATOR = "="
const PARAM_SEPARATOR = ";"
const PREVIEW_MAX_ROWS int = 300

var mode int

//var size int
var params string

type HttpRequest struct {
	uri *url.URL
	reqType string
	reqHeader string
	username string
	password string
	reqBody string
}

var debugFlag bool = false

func sendDSInfoBlock(httpRequest HttpRequest) {
	var URIInfo string
	var reqTypeInfo string
	var reqHeaderInfo string
	var usernameInfo string
	var passwordInfo string
	var reqBodyInfo string

	URIInfo = "URI;" + httpRequest.uri.String() + ";true;"
	reqTypeInfo = "TYPE;" + httpRequest.reqType + ";true;"
	reqHeaderInfo = "HEADER;" + httpRequest.reqHeader + ";true;"
	usernameInfo = "USERNAME;" + httpRequest.username + ";true"
	passwordInfo = "PASSWORD;" + httpRequest.password + ";true"
	reqBodyInfo = "BODY;" + httpRequest.reqBody + ";true"

	fmt.Println("beginDSInfo")
	fmt.Println(URIInfo)
	fmt.Println(reqTypeInfo)
	fmt.Println(reqHeaderInfo)
	fmt.Println(usernameInfo)
	fmt.Println(passwordInfo)
	fmt.Println(reqBodyInfo)
	//fmt.Println("dummy;" + "dummy param value appears" + ";false")
	fmt.Println("csv_first_row_has_column_names;true;true")
	fmt.Println("endDSInfo")
}

func sendDataBlock(httpRequest HttpRequest) {
	fmt.Println("beginData")
	if debugFlag {
		readParams()
	} else {
		readData(httpRequest)
	}
	fmt.Println("endData")
}

func readParams() {
	csvparamout := csv.NewWriter(os.Stdout)

	var csvheader []string
	csvheader = append(csvheader, "params")
	csvheader = append(csvheader, "mode")
	csvparamout.Write(csvheader)

	var paramsRow []string
	paramsRow = append(paramsRow, params)
	paramsRow = append(paramsRow, strconv.Itoa(mode))
	csvparamout.Write(paramsRow)

	csvparamout.Flush()
}

func readData(httpRequest HttpRequest) {

	uri := httpRequest.uri
	reqType := charUnescape(httpRequest.reqType)
	reqHeader := charUnescape(httpRequest.reqHeader)
	username := charUnescape(httpRequest.username)
	password := charUnescape(httpRequest.password)
	reqBody := charUnescape(httpRequest.reqBody)

	cookieOptions := cookiejar.Options {
		PublicSuffixList: publicsuffix.List,
	}
	
	jar, err := cookiejar.New(&cookieOptions)
	check(err)

	client := http.Client{
		Jar: jar,
	}

	//split headers from input text
	var headers []string
	var tokens []string
	headers = strings.Split(reqHeader, PARAM_SEPARATOR)


	req, err := http.NewRequest(reqType, uri.String(), strings.NewReader(reqBody))
	// req.Header.Set("Connection", "Keep-Alive")
	//req.Header.Add()
	for i := 0; i < len(headers); i++ {
		tokens = strings.Split(headers[i], `:`)
		req.Header.Add(tokens[0], tokens[1])
	}
	req.SetBasicAuth(username, password)

	for _, value := range jar.Cookies(uri) {
		req.AddCookie(value)
		//fmt.Println("Value of cookie: ", value)
	}
	//fmt.Println("Final Request is: ", req)
	
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()


	//parsing response body to csv
	var jsonOut []map[string]interface{}

	err = json.Unmarshal(body, &jsonOut)
	check(err)
	//fmt.Println(jsonOut)

	csvout := csv.NewWriter(os.Stdout)
	keys := []string{}

	if true {
		var header []string

		for k, _ := range jsonOut[0] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, l := range keys {
			header = append(header, l)
		}
		csvout.Write(header)
	}

	//TODO assuming subsequent JSON Objects have the same keys as the firsts
	for index, _ := range jsonOut {
		if true {
			var record []string

			for _, value := range keys {
				record = append(record, jsonOut[index][value].(string))
			}

			csvout.Write(record)
		}
	}
	csvout.Flush()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseArguments(args []string) {
	mode = ERROR

	for i := 0; i < len(args); i++ {
		if strings.ToLower(args[i]) == "-mode" && i+1 < len(args) {
			//fmt.Println("mode")

			if strings.ToLower(args[i+1]) == "preview" {
				mode = PREVIEW

			} else if strings.ToLower(args[i+1]) == "edit" {
				mode = EDIT

			} else if strings.ToLower(args[i+1]) == "refresh" {
				mode = REFRESH
			}

			// else if strings.ToLower(args[i]) == "-size" && i+1 < len(args) {
			//size, _ := strconv.ParseInt(args[i+1], 10, 64)
			//fmt.Println(size)

		} else if strings.ToLower(args[i]) == "-params" && i+1 < len(args) {
			params = args[i+1]
			//fmt.Println(params)
		}
	}
	//fmt.Println(args)
}

func charEscape(escapeStr string) string {
	r := strings.NewReplacer("\n", "%0A", "\r", "%0D", "\"", "%22", ";", "%3B")
	return r.Replace(escapeStr)
}

func charUnescape(unescapeStr string) string {
	r := strings.NewReplacer("%0A", "\n", "%0D", "\r", "%22", "\"", "%3B", ";")
	return r.Replace(unescapeStr)
}

func main() {
	args := os.Args[1:]
	parseArguments(args)

	var httpRequest HttpRequest

	var urlValue *walk.TextEdit
	var reqTypeValue *walk.TextEdit
	var usernameValue *walk.TextEdit
	var passwordValue *walk.TextEdit
	var reqHeaderValue *walk.TextEdit
	var reqBodyValue *walk.TextEdit

	if mode == REFRESH || mode == EDIT {
		//fmt.Println("REFRESH OR EDIT")
		//fmt.Println(params)
		var lines []string
		var tokens []string
		lines = strings.Split(params, PARAM_SEPARATOR)

		for i := 0; i < len(lines); i++ {
			tokens = strings.Split(lines[i], KEY_VALUE_SEPARATOR)

			if strings.ToLower(tokens[0]) == "uri" {
				var err error
				httpRequest.uri, err = url.Parse(tokens[1])
				check(err)
			}
			if strings.ToLower(tokens[0]) == "type" {
				httpRequest.reqType = tokens[1]
			}
			if strings.ToLower(tokens[0]) == "header" {
				httpRequest.reqHeader = tokens[1]
			}
			if strings.ToLower(tokens[0]) == "username" {
				httpRequest.username = tokens[1]
			}
			if strings.ToLower(tokens[0]) == "password" {
				httpRequest.password = tokens[1]
			}
			if strings.ToLower(tokens[0]) == "body" {
				httpRequest.reqBody = tokens[1]
			}
			//fmt.Println(tokens)
		}
		sendDataBlock(httpRequest)

		//fmt.Println(lines)

	} else if mode == PREVIEW {
		//fmt.Println("PREVIEW")

		MainWindow{
			Title:   "Enter the URL:",
			MinSize: Size{600, 400},
			Layout:  VBox{},
			Children: []Widget{
				Label{
					Text: `URL`,
					},
				TextEdit{ 
					AssignTo: &urlValue,
					Text: `http://localhost:3000/books.json`,
					},
				Label{
					Text: `Type`,
					},
				TextEdit{
					AssignTo: &reqTypeValue,
					Text: `GET`,
					},
				Label{
					Text: `Headers`,
					},
				TextEdit{
					AssignTo: &reqHeaderValue,
					Text: `Content-Type:application/json;Accept:*/*`,
					},
				Label{
					Text: `Username`,
					},
				TextEdit{
					AssignTo: &usernameValue,
					Text: ``,
					},
				Label{
					Text: `Password`,
					},
				TextEdit{
					AssignTo: &passwordValue,
					Text: ``,
					},
				Label{
					Text: `Request Body`,
					},
				TextEdit{
					AssignTo: &reqBodyValue,
					Text: ``,
					},
				PushButton{
					Text: "Enter",
					OnClicked: func() {
						var err error
						httpRequest.uri, err = url.Parse(strings.ToLower(urlValue.Text()))
						check(err)
						httpRequest.reqType = charEscape(reqTypeValue.Text())
						httpRequest.reqHeader = charEscape(reqHeaderValue.Text())
						httpRequest.username = charEscape(usernameValue.Text())
						httpRequest.password = charEscape(passwordValue.Text())
						httpRequest.reqBody = charEscape(reqBodyValue.Text())
					},
				},
			},
		}.Run()

		sendDSInfoBlock(httpRequest)
		sendDataBlock(httpRequest)
	}
}
