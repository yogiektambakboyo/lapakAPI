package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
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
	var seller, dbselected string
	dbselected = "ex_template"
	seller = strings.ReplaceAll(sellercode, "/", "")
	switch seller[2:4] {
	case "01":
		dbselected = "ex_template"
		break
	case "04":
		dbselected = "ex_template"
		break
	}
	fmt.Println(dbselected)
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
