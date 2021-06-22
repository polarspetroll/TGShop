package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/polarspetroll/gocolor"
)

var baseurl = "https://api.telegram.org/bot" + os.Getenv("TOKEN") + "/"

func SetWebhook(uri string) {
	uri = url.QueryEscape(uri)
	requrl := fmt.Sprintf(`%vsetWebhook?url=%v&allowed_updates=["message"]`, baseurl, uri)
	res, err := http.Get(requrl)
	defer res.Body.Close()
	CheckErr(err)
	body, err := ioutil.ReadAll(res.Body)
	CheckErr(err)
	var resout Webhookres
	err = json.Unmarshal(body, &resout)
	CheckErr(err)
	if !resout.Ok {
		log.Fatal(gocolor.ColorString("error setting webhook", "red", "bold"))
	}
	fmt.Println(gocolor.ColorString(resout.Description, "cyan", "italic"))
}

func SendPhoto(filename, caption string, chatid int64) (err error) {
	var client http.Client
	caption = url.QueryEscape(caption)
	uri := fmt.Sprintf("%vsendPhoto?chat_id=%v&caption=%v", baseurl, chatid, caption)
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	label, err := w.CreateFormField("label")
	if err != nil {
				return err
	}
	label.Write([]byte("photo"))
	summary, err := w.CreateFormField("summary")
	if err != nil {
				return err
	}
	summary.Write([]byte("file"))
	fw, err := w.CreateFormFile("photo", filename)
	if err != nil {
				return err
	}
	fd, err := os.Open(filename)
	if err != nil {
				return err
	}
	defer fd.Close()
	_, err = io.Copy(fw, fd)
	if err != nil {
				return err
	}
	w.Close()
	req, err := http.NewRequest("POST", uri, buf)
	if err != nil {
				return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return err
}

func SendMessage(chatid int64, message string) {
	uri := fmt.Sprintf("%vsendMessage?chat_id=%v&text=%v", baseurl, chatid, url.QueryEscape(message))
	http.Get(uri)
}

func CheckErr(err error) {
	if err != nil {
		er := fmt.Sprintf("%v", err)
		log.Fatal(gocolor.ColorString(er, "red", "bold"))
	}
}
