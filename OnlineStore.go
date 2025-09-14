package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

type Product struct {
	ID    int
	Name  string
	Price int
}

type CartItem struct {
	Product Product
	Qty     int
}

var products = []Product{
	{1, "T-Shirt", 120000},
	{2, "Mug", 70000},
	{3, "Sticker", 20000},
}
var cart = []CartItem{}
var mu sync.Mutex

var page = `
<!DOCTYPE html>
<html>
<head>
<title>Online Store</title>
<style>
body { font-family: Arial; margin: 30px;}
</style>
</head>
<body>
<h1>Online Store</h1>
<h2>Products</h2>
<ul>
{{range .Products}}
<li>{{.Name}} - Rp{{.Price}} 
<form method="post" action="/add" style="display:inline;">
<input type="hidden" name="id" value="{{.ID}}">
<input type="number" name="qty" value="1" min="1" style="width:40px;">
<button>Add</button>
</form>
</li>
{{end}}
</ul>
<h2>Cart</h2>
<ul>
{{range .Cart}}
<li>{{.Product.Name}} x {{.Qty}}</li>
{{end}}
</ul>
<form method="post" action="/checkout">
<button>Checkout</button>
</form>
</body>
</html>
`

func main() {
	http.HandleFunc("/", handleStore)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/checkout", handleCheckout)
	http.ListenAndServe(":8088", nil)
}

func handleStore(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	tmpl, _ := template.New("store").Parse(page)
	tmpl.Execute(w, map[string]interface{}{"Products": products, "Cart": cart})
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("id"))
	qty, _ := strconv.Atoi(r.FormValue("qty"))
	for _, p := range products {
		if p.ID == id {
			found := false
			for i, it := range cart {
				if it.Product.ID == id {
					cart[i].Qty += qty
					found = true
					break
				}
			}
			if !found {
				cart = append(cart, CartItem{Product: p, Qty: qty})
			}
			break
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleCheckout(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	cart = []CartItem{}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}