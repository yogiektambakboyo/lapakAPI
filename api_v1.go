package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
)

var dbname = "smd"

var db = make(map[string]string)

type activeTrip struct {
	Dated  string `json:"dated"`
	Time_start   string `json:"time_start"`
	Time_end   string `json:"time_end"`
	Georeverse     string `json:"georeverse"`
	Duration string `json:"duration"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	Id    string `json:"id"`
}

type storeNOO struct {
	Branch_name string `json:"branch_name"`
	Branch_id string `json:"branch_id"`
	Sales_id string `json:"sales_id"`
	Sales_name string `json:"sales_name"`
	Customer_id string `json:"customer_id"`
	Customer_name string `json:"customer_name"`
	Address string `json:"address"`
	Visit_week string `json:"visit_week"`
	City string `json:"city"`
	Credit_limit string `json:"credit_limit"`
	Longitude string `json:"longitude"`
	Latitude string `json:"latitude"`
	Email string `json:"email"`
	Handphone string `json:"handphone"`
	Whatsapp_no string `json:"whatsapp_no"`
	Citizen_id string `json:"citizen_id"`
	Tax_id string `json:"tax_id"`
	Contact_person string `json:"contact_person"`
	Type string `json:"type"`
	Clasification string `json:"clasification"`
	Contact_person_level string `json:"contact_person_level"`
	Visit_day string `json:"visit_day"` 
}

type productOrderCheckout struct {
	Product_name  string `json:"product_name"`
	Brand_name   string `json:"brand_name"`
	Uom   string `json:"uom"`
	Price     string `json:"price"`
	Qty string `json:"qty"`
	Id    string `json:"id"`
	Order_no    string `json:"order_no"`
	Sales_id    string `json:"sales_id"`
	Customers_id    string `json:"customers_id"`
	Total    string `json:"total"`
	Seq   string `json:"seq"`
}

type productOrder struct {
	Product_name  string `json:"product_name"`
	Brand_name   string `json:"brand_name"`
	Uom   string `json:"uom"`
	Price     string `json:"price"`
	Qty string `json:"qty"`
	Id    string `json:"id"`
	Order_no    string `json:"order_no"`
	Seq    string `json:"seq"`
}

type orderSales struct {
	Product_id  string `json:"product_id"`
	Total   string `json:"total"`
	Seq   string `json:"seq"`
	Price     string `json:"price"`
	Qty string `json:"qty"`
	Id    string `json:"id"`
}

type orderMaster struct {
	Customer_name  string `json:"customer_name"`
	Order_no   string `json:"order_no"`
	Total   string `json:"total"`
}

type orderDetail struct {
	Product_name  string `json:"product_name"`
	Qty  string `json:"qty"`
	Price  string `json:"price"`
	Product_total  string `json:"product_total"`
	Customer_name  string `json:"customer_name"`
	Order_no   string `json:"order_no"`
	Total   string `json:"total"`
}

type colOrderDetail struct {
	Message     string        `json:"message"`
	Data []orderDetail `json:"data"`
	Status      string        `json:"status"`
}

type colOrderMaster struct {
	Message     string        `json:"message"`
	Data []orderMaster `json:"data"`
	Status      string        `json:"status"`
}

type colStoreNOO struct {
	Message     string        `json:"message"`
	Data []storeNOO `json:"data"`
	Status      string        `json:"status"`
}

type colProductOrder struct {
	Message     string        `json:"message"`
	Data []productOrder `json:"data"`
	Status      string        `json:"status"`
}

type activeVisit struct {
	Dated  string `json:"dated"`
	Time_start   string `json:"time_start"`
	Time_end   string `json:"time_end"`
	Georeverse     string `json:"georeverse"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	Id    string `json:"id"`
	Customer_Id    string `json:"customer_id"`
	Is_checkout    string `json:"is_checkout"`
	Sales_vol    string `json:"sales_vol"`
}

type colActiveTrip struct {
	Message     string        `json:"message"`
	Data []activeTrip `json:"data"`
	Status      string        `json:"status"`
}

type colActiveVisit struct {
	Message     string        `json:"message"`
	Data []activeVisit `json:"data"`
	Status      string        `json:"status"`
}

type storeMaster struct {
	Sellercode  string `json:"Sellercode"`
	Storecode   string `json:"Storecode"`
	Storename   string `json:"Storename"`
	Address     string `json:"Address"`
	Channeldesc string `json:"Channeldesc"`
	Longitude   string `json:"Longitude"`
	Latitude    string `json:"Latitude"`
}


type colStoreVisit struct {
	Message     string        `json:"message"`
	Data []storeVisit `json:"data"`
	Status      string        `json:"status"`
}

type storeVisit struct {
	Branch_name  string `json:"branch_name"`
	Branch_id   string `json:"branch_id"`
	Sales_id   string `json:"sales_id"`
	Sales_name     string `json:"sales_name"`
	Customer_id string `json:"customer_id"`
	Customer_name   string `json:"customer_name"`
	Address    string `json:"address"`
	Visit_day    string `json:"visit_day"`
	Visit_week    string `json:"visit_week"`
	Isvisit    string `json:"isvisit"`
}


type storeReg struct {
	Id   string `json:"id"`
	Address     string `json:"address"`
	Name string `json:"name"`
	IsApproved   string `json:"is_approved"`
	CreatedAt    string `json:"created_at"`
}

type colStoreReg struct {
	Message     string        `json:"message"`
	Data []storeReg `json:"data"`
	Status      string        `json:"status"`
}

type storeMasterFullv3 struct {
	Sellercode  string `json:"sellercode"`
	Storecode   string `json:"storecode"`
	Storename   string `json:"storename"`
	Address     string `json:"address"`
	Channeldesc string `json:"channeldesc"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	City    	string `json:"city"`
	Phoneno    	string `json:"phoneno"`
	Whatsappno  string `json:"whatsappno"`
	Netizenid   string `json:"netizenid"`
	Fjp    		string `json:"fjp"`
	Isvalidated string `json:"isvalidated"`
	TargetValidation string `json:"targetvalidation"`
	IsMtd string `json:"ismtd"`
}

type resultJSON struct {
	Username      string `json:"Username"`
	Password      string `json:"Password"`
	Deviceid      string `json:"Deviceid"`
	Version       string `json:"Version"`
	Status        string `json:"Status"`
	Name          string `json:"Name"`
	DownloadDate  string `json:"DownloadDate"`
	BranchID      string `json:"BranchID"`
	BranchName    string `json:"BranchName"`
	Code          string `json:"Code"`
	VersionUpdate string `json:"VersionUpdate"`
	ForceUpdate   string `json:"ForceUpdate"`
	WeekNo        string `json:"WeekNo"`
	Bearer        string `json:"Bearer"`
}

type salesMaster struct {
	Name        string `json:"Name"`
	Code        string `json:"Code"`
	NumStore    string `json:"NumStore"`
	LastTrained string `json:"LastTrained"`
}

type checkVersion struct {
	Version     string `json:"version"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Readmelink  string `json:"readmelink"`
}

type masterData struct {
	SalesMaster []salesMaster
	StoreMaster []storeMaster
}


func sellerDivision(sellercode string) string {
	var dbselected string
	dbselected = "smd"
	return dbselected
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	// sudo kill `sudo lsof -t -i:9001` Kill Port

	// Logging to a file.
	f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pongs")
	})

	// Ping test
	r.GET("/getImage/:image", func(c *gin.Context) {
		imagefile := c.Params.ByName("image")
		path := "/home/yogiaditya0130/lapakAPI/images/" + imagefile
		//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		c.Header("Content-Type", "image/jpeg")
		c.File(path)
	})

	// Login
	r.POST("/login", func(c *gin.Context) {

		xusername := c.PostForm("username")
		xpassword := c.PostForm("password")
		xdeviceid := c.PostForm("deviceid")
		xversion := c.PostForm("version")
		xsession := c.PostForm("session")
		dbname = sellerDivision(xusername)

		xuseragent := c.Request.Header.Get("User-Agent")

		var result []resultJSON

		if xuseragent == "5uPErV1sIon_8CP_m0biL3" {
			fmt.Printf("username: %s; password: %s; deviceid: %s; version is: %s", xusername, xpassword, xdeviceid, xversion)

			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			if xsession != "" {
				rowsinsses, errsses := db.Query("INSERT INTO public.login_session(session, sellercode, description) VALUES ($1, $2, 'LOGIN'); ", xsession, xusername)
				if errsses != nil {
					// handle this error better than this
					fmt.Print(errsses)
				}

				defer rowsinsses.Close()
			}

			rows, err := db.Query(" select s.id as spvcode,s.name  as spvname,s.username as loginid,s.password,s.branch_id ,b.remark as branchname,'1' as versionupdate,'1' as forceupdate,'1' as weekno,2 as bearer from sales s join branch b on b.id=s.branch_id where upper(s.username) = $1 and s.password=$2 and s.active='1' ", xusername, xpassword)
			if err != nil {
				// handle this error better than this
				panic(err)
			}

			defer rows.Close()
			var spvcode string
			var spvname string
			var loginid string
			var password string
			var branchcode string
			var branchname string
			var versionupdate string
			var forceupdate string
			var weekno string
			var bearer string
			var counter string
			counter = "0"

			for rows.Next() {
				err = rows.Scan(&spvcode, &spvname, &loginid, &password, &branchcode, &branchname, &versionupdate, &forceupdate, &weekno, &bearer)
				if err != nil {
					// handle this error
					panic(err)
				}
				counter = "1"
			}

			// Parse JSON

			currentTime := time.Now()

			result = []resultJSON{
				resultJSON{
					Username:      loginid,
					Password:      password,
					Deviceid:      xdeviceid,
					Version:       xversion,
					Status:        counter,
					Name:          spvname,
					DownloadDate:  currentTime.Format("2006-01-02"),
					BranchID:      branchcode,
					BranchName:    branchname,
					Code:          spvcode,
					VersionUpdate: versionupdate,
					ForceUpdate:   forceupdate,
					WeekNo:        weekno,
					Bearer:        bearer,
				},
			}

			defer db.Close()
		} else {

		}

		c.JSON(http.StatusOK, result)
	})
	// End Login

	r.POST("/getActiveTrip", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " SELECT st.id,st.dated,st.time_start,st.time_end,sd.longitude,sd.latitude,sd.georeverse,lpad(EXTRACT(HOUR  FROM (st.time_end  - st.time_start))::text, 2, '0') ||':'||lpad(EXTRACT(MINUTE  FROM (st.time_end  - st.time_start))::text, 2, '0') as duration  from sales_trip st join sales_trip_detail sd on sd.trip_id = st.id where st.dated = now()::date and st.sales_id = $1 and active = 1 order by sd.id desc limit 1"

		rows, err := db.Query(sqlstring,xsales_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var id string
		var dated string
		var time_start string
		var time_end string
		var longitude string
		var latitude string
		var georeverse string
		var duration string
		var counter int

		var results []activeTrip

		counter = 0

		for rows.Next() {
			err = rows.Scan(&id,&dated,&time_start,&time_end,&longitude,&latitude,&georeverse,&duration)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := activeTrip{
				Id: id,
				Dated: dated,
				Time_start: time_start,
				Time_end: time_end,
				Longitude: longitude,
				Latitude: latitude,
				Georeverse: georeverse,
				Duration: duration,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colActiveTrip{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	
	r.POST("/getOrderToday", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xdated := c.PostForm("dated")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select c.name as customer_name,om.order_no,sum(od.total) as total  from order_master om join customers c on c.id = om.customers_id join order_detail od on od.order_no=om.order_no  where om.dated = $1 and c.sales_id=$2 group by  c.name,om.order_no order by c.name"

		rows, err := db.Query(sqlstring,xdated,xsales_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var customer_name string
		var order_no string
		var total string
		var counter int

		var results []orderMaster

		counter = 0

		for rows.Next() {
			err = rows.Scan(&customer_name,&order_no,&total)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := orderMaster{
				Customer_name : customer_name,
				Order_no : order_no,
				Total : total,

			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colOrderMaster{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colOrderMaster{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getOrderDetail", func(c *gin.Context) {
		xorder_no := c.PostForm("order_no")

		dbname = sellerDivision(xorder_no)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select c.name as customer_name,om.order_no,om.total,ps.remark as product_name,od.qty,od.price,od.total as product_total  from order_master om join order_detail od on od.order_no = om.order_no join customers c on c.id = om.customers_id join product_sku ps on ps.id = od.product_id where om.order_no = $1 order by c.name"

		rows, err := db.Query(sqlstring,xorder_no)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var customer_name string
		var order_no string
		var total string
		var product_total string
		var product_name string
		var qty string
		var price string
		var counter int

		var results []orderDetail

		counter = 0

		for rows.Next() {
			err = rows.Scan(&customer_name,&order_no,&total,&product_name,&qty,&price,&product_total)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := orderDetail{
				Customer_name : customer_name,
				Order_no : order_no,
				Total : total,
				Product_total : product_total,
				Qty : qty,
				Price : price,
				Product_name : product_name,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colOrderDetail{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colOrderDetail{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	
	r.POST("/insertOrder", func(c *gin.Context) {
		var datas []productOrderCheckout

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := c.BindJSON(&datas)
		if err != nil {
			log.Fatal(err)
		}

		var results []activeTrip

		dbname = sellerDivision("01")
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " DELETE FROM public.order_detail WHERE order_no in (select order_no from order_master where customers_id = $1 and is_checkout = 0);	"
		rowsdel, errdel := db.Query(sqlstring,datas[0].Customers_id)

		if errdel != nil {
			log.Fatal(errdel)
		}

		defer rowsdel.Close()

		sqlstring = " DELETE FROM order_master where customers_id = $1 and is_checkout = 0;	"
		rowsdelm, errdelm := db.Query(sqlstring,datas[0].Customers_id)

		if errdelm != nil {
			log.Fatal(errdelm)
		}

		defer rowsdelm.Close()

		sqlstring = " INSERT INTO public.order_master (order_no, dated, customers_id, total, sales_id) VALUES($1, now()::date, $2, $3, $4); "

		rows, err := db.Query(sqlstring,datas[0].Order_no,datas[0].Customers_id,datas[0].Total,datas[0].Sales_id)
		defer rows.Close()
		if err != nil {
			//defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert order master",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{

			// obj is a JsonObject
			for i := range datas {
				obj := datas[i]

				sqlstring = " INSERT INTO public.order_detail(order_no, product_id, qty, price, total, seq) VALUES($1, $2, $3, $4, $5, $6);	"
				rowsupd, errupd := db.Query(sqlstring,obj.Order_no,obj.Id,obj.Qty,obj.Price,obj.Total,obj.Seq)
	
				if errupd != nil {
					log.Fatal(errupd)
				}
	
				defer rowsupd.Close()
				i = i + 1
			}

			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getProductOrder", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xcustomer_id := c.PostForm("customer_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select ps.id,ps.abbr as product_name,pb.remark as brand_name,u.remark as uom,pp.price,coalesce(od.qty,0) as qty,coalesce(od.seq,'99999999') as seq  from product_sku ps join product_brand pb on pb.id = ps.brand_id join product_uom pu on pu.product_id = ps.id join uom u on u.id = pu.uom_id join product_distribution pd on pd.product_id = ps.id and pd.active = 1	join sales s on s.branch_id = pd.branch_id and s.id = $1	join customers c on c.sales_id = s.id and c.id = $2 join product_price pp on pp.product_id = ps.id and pp.branch_id = pd.branch_id left join order_master om on om.customers_id = c.id and om.dated = now()::date  and is_checkout=0 left join order_detail od on od.order_no = om.order_no and od.product_id = ps.id order by ps.remark"

		rows, err := db.Query(sqlstring,xsales_id,xcustomer_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var id string
		var product_name string
		var brand_name string
		var uom string
		var price string
		var qty string
		var seq string
		var counter int

		var results []productOrder

		counter = 0

		for rows.Next() {
			err = rows.Scan(&id,&product_name,&brand_name,&uom,&price,&qty,&seq)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := productOrder{
				Id: id,
				Product_name: product_name,
				Brand_name: brand_name,
				Uom: uom,
				Price: price,
				Qty: qty,
				Seq: seq,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colProductOrder{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colProductOrder{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getProductOrderCheckout", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xcustomer_id := c.PostForm("customer_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select om.order_no,ps.id,ps.abbr as product_name,pb.remark as brand_name,u.remark as uom,pp.price,coalesce(od.qty,0) as qty  from product_sku ps join product_brand pb on pb.id = ps.brand_id join product_uom pu on pu.product_id = ps.id join uom u on u.id = pu.uom_id join product_distribution pd on pd.product_id = ps.id and pd.active = 1	join sales s on s.branch_id = pd.branch_id and s.id = $1	join customers c on c.sales_id = s.id and c.id = $2 join product_price pp on pp.product_id = ps.id and pp.branch_id = pd.branch_id join order_master om on om.customers_id = c.id and om.dated = now()::date and is_checkout=0 join order_detail od on od.order_no = om.order_no and od.product_id=ps.id order by ps.remark"

		rows, err := db.Query(sqlstring,xsales_id,xcustomer_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var id string
		var product_name string
		var brand_name string
		var uom string
		var price string
		var qty string
		var order_no string
		var counter int
		var results []productOrder

		counter = 0

		for rows.Next() {
			err = rows.Scan(&order_no,&id,&product_name,&brand_name,&uom,&price,&qty)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := productOrder{
				Id: id,
				Product_name: product_name,
				Brand_name: brand_name,
				Uom: uom,
				Price: price,
				Qty: qty,
				Order_no : order_no,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colProductOrder{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colProductOrder{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getVisitActive", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xcustomer_id := c.PostForm("customer_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sum(coalesce(od.total,0)) as sales_vol,sv.id,sv.dated,customer_id,to_char(time_start,'dd-mm-YYYY HH24:mi') as time_start,to_char(time_end,'dd-mm-YYYY HH24:mi')  as time_end,georeverse,longitude,latitude,coalesce(sv.is_checkout,0) as is_checkout  from sales_visit sv left join order_master om on om.customers_id = sv.customer_id and om.dated = sv.dated left join order_detail od on od.order_no=om.order_no where sv.dated = now()::date and sv.customer_id = $1 and sv.sales_id =$2  group by sv.id,sv.dated,customer_id,to_char(time_start,'dd-mm-YYYY HH24:mi'),to_char(time_end,'dd-mm-YYYY HH24:mi'),georeverse,longitude,latitude,coalesce(sv.is_checkout,0)"

		rows, err := db.Query(sqlstring,xcustomer_id,xsales_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var id string
		var dated string
		var time_start string
		var time_end string
		var longitude string
		var latitude string
		var georeverse string
		var customer_id string
		var is_checkout string
		var sales_vol string
		var counter int

		var results []activeVisit

		counter = 0

		for rows.Next() {
			err = rows.Scan(&sales_vol,&id,&dated,&customer_id,&time_start,&time_end,&georeverse,&longitude,&latitude,&is_checkout)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := activeVisit{
				Id: id,
				Dated: dated,
				Time_start: time_start,
				Time_end: time_end,
				Longitude: longitude,
				Latitude: latitude,
				Georeverse: georeverse,
				Customer_Id: customer_id,
				Is_checkout: is_checkout,
				Sales_vol: sales_vol,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colActiveVisit{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colActiveVisit{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getActiveTripAll", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " SELECT st.id,to_char(st.dated,'dd-mm-YYYY') as dated,to_char(st.time_start,'dd-mm-YYYY HH24:mi') as time_start,to_char(st.time_end,'dd-mm-YYYY HH24:mi')  as time_end,lpad(EXTRACT(HOUR  FROM (st.time_end  - st.time_start))::text, 2, '0') ||':'||lpad(EXTRACT(MINUTE  FROM (st.time_end  - st.time_start))::text, 2, '0') as duration  from sales_trip st where st.dated = now()::date and st.sales_id = $1 order by id asc"

		rows, err := db.Query(sqlstring,xsales_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var id string
		var dated string
		var time_start string
		var time_end string
		var longitude string
		var latitude string
		var georeverse string
		var duration string
		var counter int

		var results []activeTrip

		counter = 0

		for rows.Next() {
			err = rows.Scan(&id,&dated,&time_start,&time_end,&duration)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := activeTrip{
				Id: id,
				Dated: dated,
				Time_start: time_start,
				Time_end: time_end,
				Longitude: longitude,
				Latitude: latitude,
				Georeverse: georeverse,
				Duration: duration,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colActiveTrip{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getStoreRegAll", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select cr.id,cr.name,cr.address,to_char(cr.created_at,'dd-mm-YYYY HH24:MI:ss') as created_at ,cr.is_approved  from customers_registration cr where cr.is_approved = 0 and sales_id=$1"

		rows, err := db.Query(sqlstring,xsales_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var id string
		var name string
		var created_at string
		var address string
		var is_approved string

		var results []storeReg

		counter := 0

		for rows.Next() {
			err = rows.Scan(&id,&name,&address,&created_at,&is_approved)
			if err != nil {
				// handle this error
				panic(err)
			}

			result := storeReg{
				Id: id,
				Address: address,
				CreatedAt : created_at,
				Name: name,
				IsApproved: is_approved,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colStoreReg{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colStoreReg{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getStoreVisitToday", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select b.remark as branch_name,s.branch_id,s.id as sales_id,s.name as sales_name,c.id as customer_id,c.name as customer_name,c.address,c.visit_day,c.visit_week,coalesce(sv.id,0) as isvisit  from customers c join sales s on s.id = c.sales_id join branch b on b.id = s.branch_id left join sales_visit sv on sv.dated = now()::date and sv.sales_id = s.id and sv.customer_id = c.id where c.sales_id = $1 and c.visit_day is not null and c.visit_week is not null "

		rows, err := db.Query(sqlstring,xsales_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var branch_name string
		var branch_id string
		var sales_id string
		var sales_name string
		var customer_id string
		var customer_name string
		var address string
		var visit_day string
		var visit_week string
		var isvisit string

		var counter int

		var results []storeVisit

		counter = 0

		for rows.Next() {
			err = rows.Scan(&branch_name,&branch_id,&sales_id,&sales_name,&customer_id,&customer_name,&address,&visit_day,&visit_week,&isvisit)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := storeVisit{
				Branch_name: branch_name,
				Branch_id: branch_id,
				Sales_id: sales_id,
				Sales_name: sales_name,
				Customer_id: customer_id,
				Customer_name: customer_name,
				Address: address,
				Visit_day: visit_day,
				Visit_week: visit_week,
				Isvisit: isvisit,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colStoreVisit{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colStoreVisit{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	
	r.POST("/getStoreNOO", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select b.remark as branch_name,s.branch_id,s.id as sales_id,s.name as sales_name,c.id as customer_id,c.name as customer_name,c.address,coalesce(c.visit_day,'') as visit_day,coalesce(c.visit_week,'') as visit_week,0 as isvisit  from customers c join sales s on s.id = c.sales_id join branch b on b.id = s.branch_id join customers_segment cs on cs.id = c.segment_id and cs.remark like '%NOO%'  where c.sales_id = $1 ; "

		rows, err := db.Query(sqlstring,xsales_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var branch_name string
		var branch_id string
		var sales_id string
		var sales_name string
		var customer_id string
		var customer_name string
		var address string
		var visit_day string
		var visit_week string
		var isvisit string

		var counter int

		var results []storeVisit

		counter = 0

		for rows.Next() {
			err = rows.Scan(&branch_name,&branch_id,&sales_id,&sales_name,&customer_id,&customer_name,&address,&visit_day,&visit_week,&isvisit)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := storeVisit{
				Branch_name: branch_name,
				Branch_id: branch_id,
				Sales_id: sales_id,
				Sales_name: sales_name,
				Customer_id: customer_id,
				Customer_name: customer_name,
				Address: address,
				Visit_day: visit_day,
				Visit_week: visit_week,
				Isvisit: isvisit,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colStoreVisit{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colStoreVisit{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getStoreNOODetail", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xcustomer_id := c.PostForm("customer_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select b.remark as branch_name,s.branch_id,s.id as sales_id,s.name as sales_name,c.id as customer_id,c.name as customer_name,coalesce(c.address,'') address,coalesce(c.visit_week,'') as visit_week,coalesce(c.city,'') city,coalesce(c.credit_limit,0) as credit_limit,coalesce(c.longitude,'') as longitude,coalesce(c.latitude,'') as latitude,coalesce(c.email,'') as email,coalesce(c.handphone,'') as handphone,coalesce(c.whatsapp_no,'') as whatsapp_no,coalesce(c.citizen_id,'') citizen_id,coalesce(c.tax_id,'') tax_id,coalesce(c.contact_person,'') as contact_person,coalesce(c.type,'') as type, coalesce(c.clasification,'') as clasification,coalesce(c.contact_person_level,'') as contact_person_level,coalesce(c.visit_day,'') visit_day from customers c join sales s on s.id = c.sales_id join branch b on b.id = s.branch_id join customers_segment cs on cs.id = c.segment_id  where c.sales_id = $1 and c.id = $2 and cs.remark like '%NOO%' "

		rows, err := db.Query(sqlstring,xsales_id,xcustomer_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var branch_name string
		var branch_id string
		var sales_id string
		var sales_name string
		var customer_id string
		var customer_name string
		var address string
		var visit_week string
		var city string
		var credit_limit string
		var longitude string
		var latitude string
		var email string
		var handphone string
		var whatsapp_no string
		var citizen_id string
		var tax_id string
		var contact_person string
		var types string
		var clasification string
		var contact_person_level string
		var visit_day string

		var counter int

		var results []storeNOO

		counter = 0
		for rows.Next() {
			err = rows.Scan(&branch_name,&branch_id,&sales_id,&sales_name,&customer_id,&customer_name,&address,&visit_week,&city,&credit_limit,&longitude,&latitude,&email,&handphone,&whatsapp_no,&citizen_id,&tax_id,&contact_person,	&types,&clasification,&contact_person_level,	&visit_day)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := storeNOO{
				Branch_name: branch_name,
				Branch_id: branch_id,
				Sales_id: sales_id,
				Sales_name: sales_name,
				Customer_id: customer_id,
				Customer_name: customer_name,
				Address: address,
				Visit_day: visit_day,
				Visit_week: visit_week,
				Credit_limit: credit_limit,
				Longitude: longitude,
				Latitude: latitude,
				Email: email,
				City: city,
				Handphone: handphone,
				Whatsapp_no: whatsapp_no,
				Citizen_id: citizen_id,
				Tax_id: tax_id,
				Contact_person: contact_person,
				Clasification: clasification,
				Contact_person_level: contact_person_level,
				Type: types,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colStoreNOO{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colStoreNOO{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/updateNOO", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xcustomer_id := c.PostForm("customer_id")
		xname := c.PostForm("name")
		xaddress := c.PostForm("address")
		xphone_no := c.PostForm("phone_no")
		xcity := c.PostForm("city")
		xcredit_limit := c.PostForm("credit_limit")
		xemail := c.PostForm("email")
		xhandphone := c.PostForm("handphone")
		xwhatsapp_no := c.PostForm("whatsapp_no")
		xcitizen_id := c.PostForm("citizen_id")
		xtax_id := c.PostForm("tax_id")
		xcontact_person := c.PostForm("contact_person")
		xtype := c.PostForm("type")
		xcontact_person_job_position := c.PostForm("contact_person_job_position")
		xclasification := c.PostForm("clasification")
		xlongitude := c.PostForm("longitude")
		xlatitude := c.PostForm("latitude")
		xcontact_person_level := c.PostForm("contact_person_level")
		var results []activeTrip

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " update public.customers set segment_id=1,name=$1, address=$2, phone_no=$3, city=$4, credit_limit=$5, longitude=$6, latitude=$7, email=$8, handphone=$9, whatsapp_no=$10, citizen_id=$11, tax_id=$12, contact_person=$13, type=$14, clasification=$15, contact_person_job_position=$16, contact_person_level=$17 where id=$18 and segment_id=2;"

		rows, err := db.Query(sqlstring,xname,xaddress,xphone_no,xcity,xcredit_limit,xlongitude,xlatitude,xemail,xhandphone,xwhatsapp_no,xcitizen_id,xtax_id,xcontact_person,xtype,xclasification,xcontact_person_job_position,xcontact_person_level,xcustomer_id)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
			defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert reg detail",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{
			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/getActiveTripDetail", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xtrip_id := c.PostForm("trip_id")

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " SELECT st.id,st.dated,st.time_start,st.time_end,sd.longitude,sd.latitude,sd.georeverse,to_char(sd.created_at,'dd-MM-YYYY HH24:MI:ss') duration  from sales_trip st join sales_trip_detail sd on sd.trip_id = st.id where st.dated = now()::date and st.sales_id = $1 and sd.trip_id=$2 order by time_start asc"

		rows, err := db.Query(sqlstring,xsales_id,xtrip_id)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var id string
		var dated string
		var time_start string
		var time_end string
		var longitude string
		var latitude string
		var georeverse string
		var duration string
		var counter int

		var results []activeTrip

		counter = 0

		for rows.Next() {
			err = rows.Scan(&id,&dated,&time_start,&time_end,&longitude,&latitude,&georeverse,&duration)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := activeTrip{
				Id: id,
				Dated: dated,
				Time_start: time_start,
				Time_end: time_end,
				Longitude: longitude,
				Latitude: latitude,
				Georeverse: georeverse,
				Duration: duration,
			}
			results = append(results, result)
			counter = counter + 1
		}

		defer db.Close()

		if(counter>0){
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}else{
			colInit := colActiveTrip{
				Message:     "Failed, Data not found",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/insertReg", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xname := c.PostForm("name")
		xaddress := c.PostForm("address")
		xphone_no := c.PostForm("phone_no")
		xcity := c.PostForm("city")
		xnotes := c.PostForm("notes")
		xcredit_limit := c.PostForm("credit_limit")
		xemail := c.PostForm("email")
		xhandphone := c.PostForm("handphone")
		xwhatsapp_no := c.PostForm("whatsapp_no")
		xcitizen_id := c.PostForm("citizen_id")
		xtax_id := c.PostForm("tax_id")
		xcontact_person := c.PostForm("contact_person")
		xtype := c.PostForm("type")
		xcontact_person_job_position := c.PostForm("contact_person_job_position")
		xclasification := c.PostForm("clasification")
		xlongitude := c.PostForm("longitude")
		xlatitude := c.PostForm("latitude")
		xcontact_person_level := c.PostForm("contact_person_level")
		xphoto := c.PostForm("photo")
		var results []activeTrip

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " INSERT INTO public.customers_registration(name, address, phone_no, membership_id, abbr, branch_id, created_at, sales_id, city, notes, credit_limit, longitude, latitude, email, handphone, whatsapp_no, citizen_id, tax_id, contact_person, type, clasification, contact_person_job_position, contact_person_level, is_approved,photo) VALUES($1, $2, $3, 1, '', 1, now(), $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, 0,$20);"

		rows, err := db.Query(sqlstring,xname,xaddress,xphone_no,xsales_id,xcity,xnotes,xcredit_limit,xlongitude,xlatitude,xemail,xhandphone,xwhatsapp_no,xcitizen_id,xtax_id,xcontact_person,xtype,xclasification,xcontact_person_job_position,xcontact_person_level,xphoto)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
			defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert reg detail",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{
			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/insertActiveTripDetail", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xlongitude := c.PostForm("longitude")
		xlatitude := c.PostForm("latitude")
		xgeoreverse := c.PostForm("georeverse")
		xtrip_id := c.PostForm("trip_id")
		var results []activeTrip

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " INSERT INTO public.sales_trip_detail (trip_id, longitude, latitude, georeverse,created_by, created_at) VALUES($1, $2, $3, $4, $5, now()); "

		rows, err := db.Query(sqlstring,xtrip_id,xlongitude,xlatitude,xgeoreverse,xsales_id)
		defer rows.Close()
		if err != nil {
			defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert trip detail",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{
			sqlstring = " UPDATE public.sales_trip set time_end=now(),updated_at=now(),updated_by=$1 WHERE id=$2; "
			rowsupd, errupd := db.Query(sqlstring,xsales_id,xtrip_id)

			if errupd != nil {
				log.Fatal(err)
			}

			defer rowsupd.Close()
			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})



	r.POST("/insertTrip", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xsales_idx := c.PostForm("sales_id")
		xlongitude := c.PostForm("longitude")
		xlatitude := c.PostForm("latitude")
		xgeoreverse := c.PostForm("georeverse")
		xphoto := c.PostForm("photo")
		xnotes := c.PostForm("notes")
		var results []activeTrip

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " update sales_trip set active = 0 where sales_id = $1; "

		rows1, err1 := db.Query(sqlstring,xsales_id)
		defer rows1.Close()
		if err1 != nil {
			log.Fatal(err1)
		}

		sqlstring = "INSERT INTO public.sales_trip(dated, sales_id, time_start, time_end, active, created_by, created_at, photo, notes) VALUES(now()::date, $1, now(), now(), '1', $2, now(), $3, $4);"

		rows2, err2 := db.Query(sqlstring,xsales_id,xsales_idx,xphoto,xnotes)
		defer rows2.Close()
		if err2 != nil {
			log.Fatal(err2)
		}

		sqlstring = "SELECT max(id) as id from public.sales_trip where sales_id=$1 and dated=now()::date;"

		rows3, err3 := db.Query(sqlstring,xsales_id)
		defer rows3.Close()
		if err3 != nil {
			log.Fatal(err3)
		}

		var lastid string;

		for rows3.Next() {
			errd := rows3.Scan(&lastid)
			if errd != nil {
				// handle this error
				panic(errd)
			}
		}

		sqlstring = " INSERT INTO public.sales_trip_detail (trip_id, longitude, latitude, georeverse,created_by, created_at) VALUES($1, $2, $3, $4, $5, now()); "

		rows, err := db.Query(sqlstring,lastid,xlongitude,xlatitude,xgeoreverse,xsales_id)
		defer rows.Close()
		if err != nil {
			defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert trip detail",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{
			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/insertStoreVisit", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xsales_idx := c.PostForm("sales_id")
		xcustomer_id := c.PostForm("customer_id")
		xlongitude := c.PostForm("longitude")
		xlatitude := c.PostForm("latitude")
		xgeoreverse := c.PostForm("georeverse")
		xphoto := c.PostForm("photo")
		var results []activeTrip

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

	

		sqlstring = "INSERT INTO public.sales_visit(dated, sales_id, customer_id, time_start, time_end, created_at, created_by, georeverse, longitude, latitude, photo) VALUES(now()::date, $1, $2, now(), now(), now(), $3, $4, $5, $6, $7);"

		rows, err := db.Query(sqlstring,xsales_id,xcustomer_id,xsales_idx,xgeoreverse,xlongitude,xlatitude,xphoto)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}



		defer rows.Close()
		if err != nil {
			defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert visit",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{
			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/updateVisitActive", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xcustomer_id := c.PostForm("customer_id")
		xis_checkout := c.PostForm("is_checkout")
		var results []activeTrip

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

	
		sqlstring = "update public.sales_visit set time_end=now(),is_checkout=$1 where customer_id=$2 and dated=now()::date and sales_id=$3 and coalesce(is_checkout,0)<>1;"

		rows, err := db.Query(sqlstring,xis_checkout,xcustomer_id,xsales_id)
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}



		defer rows.Close()
		if err != nil {
			defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert visit",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{
			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/insertStopActiveTrip", func(c *gin.Context) {
		xsales_id := c.PostForm("sales_id")
		xlongitude := c.PostForm("longitude")
		xlatitude := c.PostForm("latitude")
		xgeoreverse := c.PostForm("georeverse")
		xtrip_id := c.PostForm("trip_id")
		var results []activeTrip

		dbname = sellerDivision(xsales_id)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " INSERT INTO public.sales_trip_detail (trip_id, longitude, latitude, georeverse,created_by, created_at) VALUES($1, $2, $3, $4, $5, now()); "

		rows, err := db.Query(sqlstring,xtrip_id,xlongitude,xlatitude,xgeoreverse,xsales_id)
		defer rows.Close()


		if err != nil {
			defer db.Close()
			colInit := colActiveTrip{
				Message:  "Failed insert trip detail",
				Data: results,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colInit)
			
		}else{
			sqlstring = " UPDATE public.sales_trip set active=0,time_end=now(),updated_at=now(),updated_by=$1 WHERE id=$2; "
			rowsupd, errupd := db.Query(sqlstring,xsales_id,xtrip_id)

			if errupd != nil {
				log.Fatal(err)
			}

			defer rowsupd.Close()
			defer db.Close()
			colInit := colActiveTrip{
				Message:     "OK",
				Data: results,
				Status:      "1",
			}
			c.JSON(http.StatusOK, colInit)
		}
	})

	r.POST("/confirmOrder", func(c *gin.Context) {
		xorder_no := c.PostForm("order_no")
		xdelivery_date := c.PostForm("delivery_date")
		xnotes := c.PostForm("notes")
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	UPDATE order_master set is_checkout=1,remark=$1,delivery_date=$2 where order_no=$3;		"

		rowsweek, errweek := db.Query(sqlstring,&xnotes,&xdelivery_date,&xorder_no)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()
		defer db.Close()

		c.String(http.StatusOK, "1")
	})

	
	r.POST("/getWeekNo", func(c *gin.Context) {

		var weekstr string
		weekstr = " 1=1 "
		xusername := c.PostForm("spvcode")
		if len(xusername) > 3 {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			var sqlstring string

			sqlstring = "  	SELECT week FROM public.calendar where dated=now()::date		"

			rowsweek, errweek := db.Query(sqlstring)
			if errweek != nil {
				panic(errweek)
			}

			defer rowsweek.Close()

			var week string

			for rowsweek.Next() {
				errweek = rowsweek.Scan(&week)
				if errweek != nil {
					// handle this error
					panic(errweek)
				}
			}

			if week == "1" {
				weekstr = "wk1"
			} else if week == "2" {
				weekstr = "wk2"
			} else if week == "3" {
				weekstr = " wk3 "
			} else if week == "4" {
				weekstr = "wk4"
			}

			defer db.Close()
		}

		c.String(http.StatusOK, weekstr)
	})

	r.POST("/uploadPhoto", func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("uploaded_file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		path := "/var/www/html/lapakERP/public/images/smd-image/" + file.Filename
		err := c.SaveUploadedFile(file, path)
		if(err != nil){
			log.Fatal(err)
		}

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	
	
	return r
}

func main() {
	r := setupRouter()
	// Check Port
	port := "8081"

	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't listen on port %q: %s", port, err)
		os.Exit(1)
	}

	ln.Close()
	fmt.Printf("TCP Port %q is available", port)

	// Listen and Server in 0.0.0.0:8081
	r.Run(":8081")
}
