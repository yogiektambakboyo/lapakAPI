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

var dbname = "ex_template"

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

type colActiveTrip struct {
	Message     string        `json:"message"`
	Data []activeTrip `json:"data"`
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
	dbselected = "ex_template"
	fmt.Println(dbselected)
	fmt.Println(sellercode)
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

		sqlstring = " git ad"

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

		rows2, err2 := db.Query(sqlstring,xsales_id,xsales_id,xphoto,xnotes)
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
		path := "/home/yogiaditya0130/lapakAPI/images/" + file.Filename
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
