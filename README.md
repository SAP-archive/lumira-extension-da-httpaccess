HTTPAccess - SAP Lumira Data Access Extension
=================================================

* Connect Lumira to Web APIs and create documents.
* Use this connector to send a HTTP request and parse the JSON response into a Lumira document. 
* Parameters like request URL, request type, multiple headers, basic auth and body can be modified before sending a request.
* `master` branch: JSON parser reads all the different keys in the objects and adds int, float, string and boolean values to the document. Ignores arrays and nested structures by default. 
* `master2` branch: Example to manually specify the structure of the response and specific keys to parse instead of reading all of them.
* Some comments are marked with a CHANGE tag where you could modify code to suit your specific requirements. You could customize the GUI, create a series of requests instead of one or parse a complex response into a csv table.

Usage
-------
* Replace the `SAPLumira.ini` in your installation folder to the one present in the extras folder.
* Create a `daextensions` folder in `<lumira-install-root>/Desktop`
* Copy the `httpaccess-dae-lumira.exe` to the daextensions folder
* Open `SAP Lumira > New Document > External Data Source > Next > Select httpaccess-dae-lumira`
* Change the default parameters if you'd like
* Press `Enter` and close the GUI window
* Press `Create` to create a document
* Save the document and refresh in future to update values

GUI Parameters
--------
* `URL`: Enter the complete URL
* `Type`: GET or POST.
* `Headers`: Header type and value are separated by `:`. Multiple headers are separated by `;`
* `Username` and `Password`: base64 encodes both values to append a Basic Auth header to the request
* `Request Body`: Any valid JSON. The newline characters are encoded before being saved as a DSInfo parameter. 

Environment setup
-------------------
* Setup Golang environment
* Move this repo to the src directory
* go build to generate the executable

Debug flag
-----------
Change `var debugFlag bool = false` to `true` to read in the stored DSInfo parameters into a Lumira document in place of the data retrieved from the data source.

Extras
-------
####apimocker/apimocker.exe
* Utility to serve GET requests locally for testing the extension with a sample JSON file
* Place a JSON file in the `apimocker/public` folder and run this executable
* Send a GET request to `http://localhost:3000/<json-filename>` to get a response

####addExt.bat
* `addExt.bat httpaccess` to copy the generated executable into the `daextensions` folder

####master2 branch
* Parse specific keys/nested keys you'd like to use from the response
* Example code only works with response from `http://jsonplaceholder.typicode.com/users`
* Please use [JsonUtils](https://github.com/bashtian/jsonutils) to update the `struct` for your response

Resources
-----------
* Developer Guide - [SAP Lumira 1.20 Data Access Extension SDK Developer guide](http://help.sap.com/businessobject/product_guides/vi01/en/lum_120_dae_dev_en.pdf)
* Blog post - [Using the Data Access Extension SDK](http://scn.sap.com/community/lumira/blog/2014/10/14/using-the-data-access-extension-sdk--sap-lumira)
* Webinar - [Working with the Data Access Extension SDK and Demo](https://www.youtube.com/watch?v=oaUdztW5lKc)


License
---------

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

 [1]: https://github.com/SAP/httpaccess-dae-lumira
