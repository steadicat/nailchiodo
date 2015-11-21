package web

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/", staticFileHandler)
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, strings.TrimSuffix(r.URL.Path, "/"), http.StatusMovedPermanently)
		return
	}

	c := appengine.NewContext(r)
	lang := getLanguage(w, r)

	var filename string
	if r.URL.Path == "/" {
		filename = fmt.Sprintf("static/index.%s.html", lang)
	} else {
		filename = fmt.Sprintf("static%s.%s.html", r.URL.Path, lang)
	}

	log.Infof(c, "[http] %s?lang=%s -> %s", r.URL.Path, lang, filename)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Warningf(c, "[file error] %v", err)
		notFoundHandler(c, w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(file)
}

func notFoundHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	file, err := ioutil.ReadFile("static/notfound.html")
	_, err = w.Write(file)
	if err != nil {
		log.Warningf(c, "[file error] %v", err)
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
