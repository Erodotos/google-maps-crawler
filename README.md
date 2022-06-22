# google-maps-crawler

### Usage
1. `go mod tidy`
2. `go build scraper.go`
3. `./scraper`

### Command-line arguments
- `-latitude` 
- `-longitude`
- `-numberOfPages` : the number of google maps pages of restaurants you want to crawl
- `-output` : the name of the json file where you want to output the list of the found restaurants.