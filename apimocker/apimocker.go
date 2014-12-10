/*
Copyright 2014, SAP SE

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

import "github.com/go-martini/martini"
//import "path/filepath"
//import "fmt"

func main() {
	//var booksFilename string

	//booksFilename, err := filepath.Abs(".\\books.json")
	//check(err)

	// URI = "http://localhost:3000/books.json"

	m := martini.Classic()
	m.Get("/", func() string {
		return "<a href=\"books.json\">Books</a>"
		})
	m.Use(martini.Static("public"))
	m.Run()
}

func check(e error){
	if e != nil {
		panic(e)
	}
}