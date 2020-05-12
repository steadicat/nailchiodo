package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

func main() {
	http.HandleFunc("/", staticFileHandler)
	appengine.Main()
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.Header.Get("Host"), ".nailchiodo.com") {
		http.Redirect(w, r, "https://nailchiodo.com"+r.URL.RequestURI(), http.StatusMovedPermanently)
		return
	}

	c := appengine.NewContext(r)
	lang := getLanguage(w, r)
	trimmedPath := strings.Trim(r.URL.Path, "/")

	redirects, err := getRedirects(c)
	if err != nil {
		log.Printf("[error getting redirects] %v", err)
		errorHandler(c, w, r, lang)
		return
	}
	if redirects[trimmedPath] != "" {
		http.Redirect(w, r, redirects[trimmedPath], http.StatusMovedPermanently)
		return
	}

	branchName := fmt.Sprintf("static/%s/index.%s.html", trimmedPath, lang)
	leafName := fmt.Sprintf("static/%s.%s.html", trimmedPath, lang)
	var fileName string

	if strings.HasSuffix(r.URL.Path, "/") {
		_, err := os.Stat(leafName)
		if err == nil {
			log.Printf("[http] %s exists, redirecting", leafName)
			http.Redirect(w, r, "/"+trimmedPath, http.StatusMovedPermanently)
			return
		}
		fileName = branchName
	} else {
		_, err := os.Stat(branchName)
		if err == nil {
			log.Printf("[http] %s exists, redirecting", branchName)
			http.Redirect(w, r, "/"+trimmedPath+"/", http.StatusMovedPermanently)
			return
		}
		fileName = leafName
	}

	log.Printf("[http] %s?lang=%s -> %s", r.URL.Path, lang, fileName)

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("[not found] %v", err)
		notFoundHandler(c, w, r, lang)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[error responding] %v", err)
		errorHandler(c, w, r, lang)
	}
}

func notFoundHandler(c context.Context, w http.ResponseWriter, r *http.Request, lang string) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, err := ioutil.ReadFile(fmt.Sprintf("static/notfound.%s.html", lang))
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[file error] %v", err)
		errorHandler(c, w, r, lang)
	}
}

func errorHandler(c context.Context, w http.ResponseWriter, r *http.Request, lang string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, err := ioutil.ReadFile(fmt.Sprintf("static/error.%s.html", lang))
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[file error on error!] %v", err)
		fmt.Fprintf(w, "Error!")
	}
}

func getLanguage(w http.ResponseWriter, r *http.Request) string {
	var preferred []language.Tag
	var err error
	save := false

	r.ParseForm()
	lang := r.Form.Get("lang")
	if lang != "" {
		preferred, _, err = language.ParseAcceptLanguage(lang)
		if err != nil {
			// log err
		} else {
			save = true
		}
	}

	if preferred == nil {
		cookie, err := r.Cookie("lang")
		if err == nil {
			preferred, _, err = language.ParseAcceptLanguage(cookie.Value)
			if err != nil {
				// log err
			}
		}
	}

	if preferred == nil {
		preferred, _, err = language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		if err != nil {
			// log err
		}
	}

	matcher := language.NewMatcher([]language.Tag{
		language.English,
		language.Italian,
		language.French,
	})
	code, _, _ := matcher.Match(preferred...)
	base, _ := code.Base()

	if save {
		now := time.Now()
		expires := now.AddDate(1, 0, 0)
		http.SetCookie(w, &http.Cookie{
			Name:    "lang",
			Value:   base.String(),
			Path:    "/",
			Expires: expires,
		})
	}

	return base.String()
}

func getRedirects(c context.Context) (map[string]string, error) {
	var redirects map[string]string

	_, err := memcache.JSON.Get(c, "redirects", &redirects)
	if err == nil {
		log.Print("[redirects] Cache hit")
		return redirects, nil
	}

	if err != memcache.ErrCacheMiss {
		log.Print("[redirects] Error reading from cache", err.Error())
	} else {
		log.Print("[redirects] Cache miss")
	}

	redirectsFile, err := os.Open("redirects.json")
	if err != nil {
		log.Print("[redirects] Error opening redirects file", err.Error())
		return nil, err
	}

	jsonParser := json.NewDecoder(redirectsFile)
	err = jsonParser.Decode(&redirects)
	if err != nil {
		log.Print("[redirects] Error parsing redirects file", err.Error())
		return nil, err
	}

	err = memcache.JSON.Set(c, &memcache.Item{
		Key:    "redirects",
		Object: &redirects,
	})
	if err != nil {
		log.Print("[redirects] Error storing in memcache", err.Error())
	}

	return redirects, nil
}
