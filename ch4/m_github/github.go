//package github provides a Go API for the Github issue tracker
//See https://developer.github.com/v3/search/#search-issues
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
	"time"
)

//API url
const IssuesURL = "https://api.github.com/search/issues"

//structs for all outputs from api response
//names of all the struct fields must be capitalized because they
//are being exported
type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}
type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    //in Markdown format
}
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

//create template variable
// const templ = `{{.TotalCount}} issues:
// {{range .Items}}-------------------------------------------------------------------
// Number: {{.Number}}
// User: {{.User.Login}}
// Title: {{.Title | printf "%.64s"}}
// Age: {{.CreatedAt | daysAgo}} days
// {{end}}`

var issueList = template.Must(template.New("issuelist").Parse(`<h1>{{.TotalCount}} issues</h1>
<table>
	<tr style = 'text-align: left'>
		<th>#</th>
		<th>State</th>
		<th>User</th>
		<th>Title</th>
	</tr>
	{{range .Items}}
	<tr>
		<td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
		<td>{{.State}}</td>
		<td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
		<td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
	</tr>
	{{end}}
	</table>
`))

// var report = template.Must(template.New("issuelist").
// 	Funcs(template.FuncMap{"daysAgo": daysAgo}). //adds daysAgo to set of functions accesible
// 	Parse(templ))

func main() {
	//command line args sent to SearchIssues function
	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := issueList.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
	// //prints total number of issues
	// fmt.Printf("%d issues:\n", result.TotalCount)
	// //ranges through array to print each issue
	// for _, item := range result.Items {
	// 	fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	// }
}

//search using parameters indicated by command line arguments
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	//QueryEscape in order to place inside url string
	q := url.QueryEscape(strings.Join(terms, " "))
	//Get request concatenating command line args with base URL
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}

	//We must close resp.Body on all execution paths
	//(Chapter 5 presents 'defer', which makes this simpler)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}
