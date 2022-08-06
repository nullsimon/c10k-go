package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var userCache *User
var productCache *Product

const (
	Quantity = 10000 * 10000 // 库存数量
)

func init() {
	gormDB, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db = gormDB
	if err != nil {
		panic("failed to connect database")
	}
	// 迁移 schema
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Order{})

	// get product
	var product Product
	where := map[string]interface{}{
		"code": "Sticker",
	}
	db.First(&product, where)
	if product.ID == 0 {
		// Create
		fmt.Printf("product not found\n")
		db.Create(&Product{Code: "Sticker", Price: 100, Quantity: Quantity})
	}
	// Read
	// var product Product
	db.Last(&product, "code = ?", "Sticker")
	fmt.Printf(`product: %v`, product)

	// init users
	for i := 0; i < 10000; i++ {
		var user User
		where = map[string]interface{}{
			"name": fmt.Sprintf("demo-%d", i),
		}
		db.First(&user, where)
		if user.ID == 0 {
			// Create
			fmt.Printf("user not found\n")
			db.Create(&User{Name: fmt.Sprintf("demo-%d", i)})
		}
	}

	// init user cache
	var user User
	db.First(&user, "name = ?", "demo-0")
	userCache = &user
	// init product cache
	productCache = &product
	// init quantity to redis
	//redis.InitQuantity(context.Background(), redis.NewClient(), Quantity)
}

// request handler in fasthttp style, i.e. just plain function.
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
}

func createOrderHandler(ctx *fasthttp.RequestCtx) {
	// mock get user, 10 percent not found
	user := userCache

	// get product, only one product, maybe cache it
	product := productCache
	// create order
	err := creatOrder(db, *user, *product)
	if err != nil {
		fmt.Fprintf(ctx, `error: %v`, err)
		return
	}
	ctx.Response.Header.SetStatusCode(200)
	fmt.Fprintf(ctx, `order created`)
}

func getProductHandler(ctx *fasthttp.RequestCtx) {
	var product Product
	where := map[string]interface{}{
		"code": string(ctx.QueryArgs().Peek("code")),
	}
	db.First(&product, where)
	if product.ID == 0 {
		fmt.Printf("product not found\n")
		return
	}
	ctx.Response.Header.SetStatusCode(200)
	fmt.Fprintf(ctx, `product: %v`, product)
}

func main() {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/order":
			if ctx.IsGet() {
				createOrderHandler(ctx)
			} else {
				ctx.Error("Method not allowed", 405)
			}
		case "/product":
			if ctx.IsGet() {
				getProductHandler(ctx)
			} else {
				ctx.Error("Method not allowed", 405)
			}
		default:
			fastHTTPHandler(ctx)
		}
	}
	fasthttp.ListenAndServe(":8087", m)
}
