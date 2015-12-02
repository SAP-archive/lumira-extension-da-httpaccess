/*
Copyright 2015, SAP SE

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/    
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
	//"sort"
	"github.com/hashicorp/go.net/publicsuffix"
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

//set this to true if you want to display stored DSInfo parameters in the Lumira document
var debugFlag bool = false

//send DSInfo parameters in the right format to Lumira in Preview Mode only
func sendDSInfoBlock(httpRequest HttpRequest) {
	var URIInfo string
	var reqTypeInfo string
	var reqHeaderInfo string
	var usernameInfo string
	var passwordInfo string
	var reqBodyInfo string

	URIInfo = "URI;" + httpRequest.uri.String() + ";true;"
	reqTypeInfo = "TYPE;" + httpRequest.reqType + ";true;"
	
	//multiple headers stored in one line as headername:headervalue; 
	//these header separators are encoded before being sent here 
	reqHeaderInfo = "HEADER;" + httpRequest.reqHeader + ";true;"
	
	//basic auth header will be created after base64 encoding
	usernameInfo = "USERNAME;" + httpRequest.username + ";true"
	passwordInfo = "PASSWORD;" + httpRequest.password + ";true"

	//new line characters in the body are encoded before storing here
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

//used in all three modes to send data to lumira
func sendDataBlock(httpRequest HttpRequest) {
	fmt.Println("beginData")

	//read debug flag to switch between displaying DSInfo parameters or 
	//data in the lumira document
	if debugFlag {
		readParams()
	} else {
		readData(httpRequest)
	}
	fmt.Println("endData")
}

//retrieve and send parameters to the document instead of data
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

//get data from source, convert it to csv and send to lumira
func readData(httpRequest HttpRequest) {

	uri := httpRequest.uri
	reqType := charUnescape(httpRequest.reqType)
	reqHeader := charUnescape(httpRequest.reqHeader)
	username := charUnescape(httpRequest.username)
	password := charUnescape(httpRequest.password)
	reqBody := charUnescape(httpRequest.reqBody)

	//cookie jar, can be used if multiple calls are made to retreive data
	cookieOptions := cookiejar.Options {
		PublicSuffixList: publicsuffix.List,
	}
	
	jar, err := cookiejar.New(&cookieOptions)
	check(err)

	client := http.Client{
		Jar: jar,
	}

	//CHANGE: Read cookies from the jar and append to subsequent requests if you 
	//are making multiple requests

	//extract individual headers
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

	//compute the basic auth header
	req.SetBasicAuth(username, password)

	for _, value := range jar.Cookies(uri) {
		req.AddCookie(value)
		//fmt.Println("Value of cookie: ", value)
	}
	//fmt.Println("Final Request is: ", req)
	
	//send request
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	//read the body of the response
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//convert the JSON response to CSV

	//CHANGE: please create a struct with the fields you'd like to parse 
	//from the response

	// JsonUtils from https://github.com/bashtian/jsonutils could save some manual work

	//CHANGE: this would only work for http://jsonplaceholder.typicode.com/users
	// please update this to your response structure
	type JsonOut []struct {
		Address struct {
			City string `json:"city"`
			Geo  struct {
				Lat float64 `json:"lat,string"`
				Lng float64 `json:"lng,string"`
			} `json:"geo"`
			Street  string `json:"street"`
			Suite   string `json:"suite"`
			Zipcode string `json:"zipcode"`
		} `json:"address"`
		Company struct {
			Bs          string `json:"bs"`
			CatchPhrase string `json:"catchPhrase"`
			Name        string `json:"name"`
		} `json:"company"`
		Email    string `json:"email"`
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Username string `json:"username"`
		Website  string `json:"website"`
	}

	//this method parses all the first level keys in all objects
	//we ignore nested values and arrays

	//var jsonOut []map[string]interface{}
	var jsonOut JsonOut;
	//refer to http://blog.golang.org/json-and-go

	//unmarshal response body to the map
	err = json.Unmarshal(body, &jsonOut)
	check(err)
	//fmt.Println(jsonOut)

	csvout := csv.NewWriter(os.Stdout)
	
	//CHANGE: parsing only a few fields
	if true {
		header := []string{"username","name","id","email","phone","company_name","city","zipcode"}
		csvout.Write(header)
	}

	for _, value := range jsonOut {
		if true {
			var record []string

			//CHANGE: change the list of fields to be parsed
			//Nested fields can be accessed and inserted into the flat table
			for _, d:= range []interface{}{ value.Username, value.Name, value.Id, value.Email, value.Phone, value.Company.Name, value.Address.City, value.Address.Zipcode } {

  				switch v := d.(type) {
					case string:
						record = append(record, v)
					case int:
						record = append(record, strconv.Itoa(v))
					case int64:
						record = append(record, strconv.FormatInt(v, 10))
					case float64:
						record = append(record, strconv.FormatFloat(v, 'f', -1, 64))
					case bool:
						if v {
							record = append(record, "true")	
						} else {
							record = append(record, "false")
						}
					default:
						//applies blank value to null, whitespace, arrays and nested structures in a JSON object
						//CHANGE: please add custom code here if you want a specific parsing behavior
						record = append(record, "")
					}
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

// parse arguments passed to the executable in various modes
func parseArguments(args []string) {
	mode = ERROR

	for i := 0; i < len(args); i++ {

		//detect and save the current mode of operation requested
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

		// read the params if sent in EDIT and REFRESH modes
		} else if strings.ToLower(args[i]) == "-params" && i+1 < len(args) {
			params = args[i+1]
			//fmt.Println(params)
		}
	}
	//fmt.Println(args)
}

//encoding characters that interfere when the SDK parses parameter values
//Carriage Return, Line Feed, semi-colon, and double-quotes
func charEscape(escapeStr string) string {
	r := strings.NewReplacer("\n", "%0A", "\r", "%0D", "\"", "%22", ";", "%3B")
	return r.Replace(escapeStr)
}

//decoding these characters before presenting the values to lumira
func charUnescape(unescapeStr string) string {
	r := strings.NewReplacer("%0A", "\n", "%0D", "\r", "%22", "\"", "%3B", ";")
	return r.Replace(unescapeStr)
}

//remove duplicates in the list
func removeDuplicates(headers *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *headers {
		if !found[x] {
			found[x] = true
			(*headers)[j] = (*headers)[i]
			j++
		}
	}
	*headers = (*headers)[:j]
}

//main
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

	// handle refresh and edit modes
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

	//handle preview mode
	} else if mode == PREVIEW {
		//fmt.Println("PREVIEW")

		//implementing the GUI interface to capture parameters from user
		//some fields have defaults
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
					Text: `http://jsonplaceholder.typicode.com/users`,
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
