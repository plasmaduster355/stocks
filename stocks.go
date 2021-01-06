package stocks

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type News struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		URLToImage  string    `json:"urlToImage"`
		PublishedAt time.Time `json:"publishedAt"`
		Content     string    `json:"content"`
	} `json:"articles"`
}

//Sets api key
func SetStockKey(key string) (err error) {
	err = ioutil.WriteFile("temp/stockKey.txt", []byte(key), 0600)
	if err != nil {
		err = nil
		f, err := os.Create("temp/stockKey.txt")
		f.Write([]byte(key))
		f.Close()
		return err
	}

	return err
}

func SetNewsKey(key string) (err error) {
	err = ioutil.WriteFile("temp/newsKey.txt", []byte(key), 0600)
	if err != nil {
		err = nil
		f, err := os.Create("temp/newsKey.txt")
		f.Write([]byte(key))
		f.Close()
		return err
	}

	return err
}

//Gets api key
func GetStockKey() (key string, err error) {
	b, err := ioutil.ReadFile("temp/stockKey.txt")
	return string(b), err
}

func GetNewsKey() (key string, err error) {
	b, err := ioutil.ReadFile("temp/newsKey.txt")
	return string(b), err
}

//GetDailyStockData ;)
func GetDailyStockData(location string, stockCode string) (data [][]string, err error) {
	var r *http.Response
	var k string
	if location == "av" {
		k, err = GetStockKey()
		r, err = http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&outputsize=full&symbol=" + stockCode + "&apikey=" + k + "&datatype=csv")

	}
	return loadcsv(r, stockCode)
}

//Create csv file and return it
func loadcsv(r *http.Response, filename string) (data [][]string, err error) {

	b, _ := ioutil.ReadAll(r.Body)
	ioutil.WriteFile("temp/"+filename+".csv", b, 0600)
	f, err := os.Open("temp/" + filename + ".csv")
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}
	return lines, nil
}
func StockNews(country string) (news string, err error) {
	k, err := GetNewsKey()
	r, err := http.Get("https://newsapi.org/v2/top-headlines?country=" + country + "&category=business&apiKey=" + k)
	b, _ := ioutil.ReadAll(r.Body)
	var n News
	json.Unmarshal(b, &n)
	for _, ne := range n.Articles {
		news = news + "|" + ne.Title + "(" + ne.Source.Name + ")"
	}
	return news, err
}
