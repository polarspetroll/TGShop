package app

import (
	"DB"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"telegram"
	"time"
)

type Product struct {
	Name  string
	Price string
	Id    int64
	Stat  bool
}
var wg sync.WaitGroup
var tmps *template.Template = template.Must(template.ParseGlob("templates/*.html"))

func Login(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		if r.Method == "GET" {
			tmps.ExecuteTemplate(w, "login.html", nil)
			return nil
		} else if r.Method == "POST" {
			username := r.PostFormValue("username")
			password := r.PostFormValue("password")
			if strings.Replace(username, " ", "", -1) == "" || strings.Replace(password, " ", "", -1) == "" || len(username) > 30 {
				tmps.ExecuteTemplate(w, "login.html", "Invalid usernam or password")
				return nil
			}
			st, err := DB.QueryLogin(username, password)
			if err != nil {
				return err
			}
			if st {
				c, err := DB.SetCookie(username)
				if err != nil {
					return err
				}
				cookie := http.Cookie{Name: "SID", Value: c, Expires: time.Now().Add(10 * time.Hour)}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/index", 302)
				return err
			} else {
				tmps.ExecuteTemplate(w, "login.html", "Invalid username or password")
				return err
			}
		} else {
			http.Error(w, "Method Not Allowed", 405)
			return nil
		}
		return nil
	}()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SID")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	stat, username := DB.GetCookie(cookie.Value)
	if !stat {
		http.Redirect(w, r, "/login", 302)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	tmps.ExecuteTemplate(w, "index.html", username)
}

func List(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		cookie, err := r.Cookie("SID")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return err
		}
		stat, _ := DB.GetCookie(cookie.Value)
		if !stat {
			http.Redirect(w, r, "login", 302)
			return err
		}
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", 405)
			return err
		}
		products, err := DB.ListQuery()
		if err != nil {
			return err
		}
		prd := []Product{}
		var pstat bool
		for i := 0; i < len(products.Id); i++ {
			if products.Stat[i] == 1 {
				pstat = true
			} else {
				pstat = false
			}
			prd = append(prd, Product{Name: products.Name[i], Price: products.Price[i], Stat: pstat, Id: products.Id[i]})
		}
		tmps.ExecuteTemplate(w, "list.html", prd)
		return err
	}()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		cookie, err := r.Cookie("SID")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return err
		}
		stat, _ := DB.GetCookie(cookie.Value)
		if !stat {
			http.Redirect(w, r, "login", 302)
			return err
		}
		if r.Method == "GET" {
			param := r.URL.Query()
			id := param.Get("id")
			pid, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				http.Error(w, "Product Not Found", 404)
				return err
			}
			out, err := DB.QueryById(pid)
			if err != nil {
				return err
			}
			if len(out.Name) == 0 {
				http.Error(w, "Product Not Found", 404)
				return err
			}
			tmps.ExecuteTemplate(w, "product.html", out)
		} else if r.Method == "POST" {
			name := r.PostFormValue("name")
			price := r.PostFormValue("price")
			status := r.PostFormValue("status")
			productid := r.URL.Query()
			id := productid.Get("id")
			v, i, pid := editinputcheck(name, price, status, id)
			if !v {
				w.Write([]byte("<script>alert('invalid inputs');window.location='/product?id=" + id + "'</script>"))
				return err
			}
			row, err := DB.Update(name, price, i, pid)
			if err != nil {
				return err
			}
			if row <= 1 {
				http.Redirect(w, r, "/product?id="+id, 302)
				DB.Del(fmt.Sprintf("product:%v", pid))
			}
		} else {
			http.Error(w, "Method Now Allowed", 405)
		}
		return err
	}()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func GroupMessage(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		cookie, err := r.Cookie("SID")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return err
		}
		stat, _ := DB.GetCookie(cookie.Value)
		if !stat {
			http.Redirect(w, r, "login", 302)
			return err
		}

		if r.Method == "GET" {
			tmps.ExecuteTemplate(w, "message.html", nil)
			return err
		} else if r.Method == "POST" {
			message := r.PostFormValue("message")
			if len(message) > 4000 { //telegram message limit
				tmps.ExecuteTemplate(w, "message.html", "Error, Message too long")
				return err
			}
			tmps.ExecuteTemplate(w, "message.html", "Done! However it may take some time to send the message to all users.")
			chatids, err := DB.GetList("chatids")
			if err != nil {
				return err
			}
			chatidlen := len(chatids)
			if chatidlen != 1 ||  chatidlen != 2 {
				for _, i := range chatids {
					cid, _ := strconv.ParseInt(i.(string), 10, 64)
					telegram.SendMessage(cid, message)
				}
				return err
			} else if chatidlen > 2 {
					rng := chatidlen / 2
					chatsone := chatids[:rng]
					chatstwo := chatids[rng + 1:]
					wg.Add(2)
					go ConcurrenceMsg(chatsone, message)
					go ConcurrenceMsg(chatstwo, message)
					wg.Wait()
				}
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
		return err
	}()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func NewProduct(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		cookie, err := r.Cookie("SID")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return err
		}
		stat, _ := DB.GetCookie(cookie.Value)
		if !stat {
			http.Redirect(w, r, "login", 302)
			return err
		}

		if r.Method == "GET" {
			tmps.ExecuteTemplate(w, "newproduct.html", nil)
			return err
		} else if r.Method == "POST" {
			name := r.PostFormValue("name")
			price := r.PostFormValue("price")
			status := r.PostFormValue("status")
			validate, stat, _ := editinputcheck(name, price, status, "0")
			if !validate {
				tmps.ExecuteTemplate(w, "newproduct.html", "Invalid product")
				return err
			}
			r.ParseMultipartForm(5 << 20)
			file, f, err := r.FormFile("image")
			if err != nil {
				http.Error(w, "Bad Request", 400)
				return err
			}
			defer file.Close()
			ext := path.Ext(f.Filename)
			if !CheckFile(ext) {
				tmps.ExecuteTemplate(w, "newproduct.html", "Invalid file format")
				return err
			} else if f.Size > 5242880 {
				tmps.ExecuteTemplate(w, "newproduct.html", "Error. File too large")
				return err
			}
			tmpfile, err := ioutil.TempFile("images", "*"+ext)
			if err != nil {
				return err
			}
			defer tmpfile.Close()
			fbyte, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			tmpfile.Write(fbyte)
			c, err := DB.Insert(name, price, tmpfile.Name(), stat)
			if err != nil {
				return err
			}
			if c != 1 {
				w.Write([]byte("Internal Server Error"))
				return err
			}
			tmps.ExecuteTemplate(w, "newproduct.html", "Done!")
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
		return err
	}()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func Del(w http.ResponseWriter, r *http.Request) {
	err := func() error {
		cookie, err := r.Cookie("SID")
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return err
		}
		stat, _ := DB.GetCookie(cookie.Value)
		if !stat {
			http.Error(w, http.StatusText(401), 401)
			return err
		}
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", 405)
			return err
		}
		param := r.URL.Query()
		id := param.Get("id")
		pid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid Product ID", 400)
			return err
		}
		row, err := DB.Delete(pid)
		if err != nil {
			return err
		}
		if row == 0 {
			http.Error(w, "Product Not Found", 404)
			return err
		}
		http.Redirect(w, r, "/list", 302)
		DB.Del(id)
		return err
	}()
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SID")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	DB.Del(cookie.Value)
	http.Redirect(w, r, "/login", 302)
}

//////////////////////////////////////////////////////////////////////////////////////////
func ConcurrenceMsg(chats []interface{}, message string) {
	for _, i := range chats {
		cid, _ := strconv.ParseInt(i.(string), 10, 64)
		telegram.SendMessage(cid, message)
	}
	wg.Done()
}

func CheckFile(ext string) bool {
	legalfiles := [4]string{".png", ".jpeg", ".jpg", ".gif"}
	arr := reflect.ValueOf(legalfiles)
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == ext {
			return true
		}
	}
	return false
}

func editinputcheck(name, price, status, id string) (validate bool, stat int, pid int64) {
	validate = true
	if len(name) > 100 || len(price) > 10 {
		validate = false
		stat = 0
		pid = 0
		return validate, stat, pid
	}
	inputs := [2]string{name, price}
	for i := range inputs {
		if strings.Replace(inputs[i], " ", "", -1) == "" {
			validate = false
			pid = 0
			stat = 0
			break
			return validate, stat, pid
		}
	}
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		validate = false
		stat = 0
		pid = 0
		return validate, stat, pid
	}
	if status == "1" {
		stat = 1
		return validate, stat, pid
	} else {
		stat = 0
		return validate, stat, pid
	}
}
