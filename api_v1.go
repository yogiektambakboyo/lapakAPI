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

type dataTime struct {
	TimeMinute string `json:"timeminute"`
	TimeHour   string `json:"timehour"`
	TimeDay    string `json:"timeday"`
}

type summaryMTD struct {
	Periodid    string `json:"periodid"`
	SPVName     string `json:"spvname"`
	TargetSales string `json:"targetsales"`
	ActualSales string `json:"actsales"`
	TargetCall  string `json:"targetcall"`
	ActualCall  string `json:"actcall"`
	ActualPCall string `json:"actpcall"`
	TargetGDP   string `json:"targetgdp"`
	ActualGDP   string `json:"actgdp"`
	TargetIFF   string `json:"targetiff"`
	ActualIFF   string `json:"actiff"`
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

type storeMasterFull struct {
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

type dataVolume struct {
	Name   string `json:"name"`
	Target string `json:"target"`
	Actual string `json:"actual"`
}

type dataVolumeDetail struct {
	Name       string `json:"name"`
	Target     string `json:"target"`
	Actual     string `json:"actual"`
	Procent    string `json:"procent"`
	Sellercode string `json:"sellercode"`
}

type dataStartCall struct {
	Sellername string `json:"sellername"`
	Timein     string `json:"timein"`
	Timeclock  string `json:"timeclock"`
	Sellercode string `json:"sellercode"`
}

type dataVolumeDetailIFF struct {
	Name       string `json:"name"`
	Target     string `json:"target"`
	Actual     string `json:"actual"`
	Procent    string `json:"procent"`
	Sellercode string `json:"sellercode"`
	Sellername string `json:"sellername"`
}

type dataIFF struct {
	Sellercode string `json:"sellercode"`
	Sellername string `json:"sellername"`
	Iff0       string `json:"iff0"`
	Iff1       string `json:"iff1"`
	Iff2       string `json:"iff2"`
	Iff3       string `json:"iff3"`
	Iff4       string `json:"iff4"`
	Iff5       string `json:"iff5"`
	Iff6       string `json:"iff6"`
}

type dataComplaint struct {
	Sellercode    string `json:"sellercode"`
	Sellername    string `json:"sellername"`
	Storename     string `json:"storename"`
	Storecode     string `json:"storecode"`
	Longitude     string `json:"longitude"`
	Latitude      string `json:"latitude"`
	Longitudecall string `json:"longitudecall"`
	Latitudecall  string `json:"latitudecall"`
}

type dataSales struct {
	Sellercode string `json:"sellercode"`
	Sellername string `json:"sellername"`
}


type dataChannel struct {
	ChannelCode string `json:"channelcode"`
	ChannelDesc string `json:"channeldesc"`
}

type dataWorkDay struct {
	Currworkday string `json:"currworkday"`
	Workday     string `json:"workday"`
}

type colWorkDay struct {
	Message     string        `json:"message"`
	DataWorkDay []dataWorkDay `json:"data"`
	Status      string        `json:"status"`
}

type dataIFFDetail struct {
	Sellername string `json:"sellername"`
	Sellercode string `json:"sellercode"`
	Iffname    string `json:"iffname"`
	Target     string `json:"target"`
	Actual     string `json:"actual"`
}

type colIFF struct {
	Message       string          `json:"message"`
	DataIFFDetail []dataIFFDetail `json:"data"`
	Status        string          `json:"status"`
}


type colStoreMasterFull struct {
	Message       string          `json:"message"`
	DataStoreMasterFull []storeMasterFull `json:"data"`
	Status        string          `json:"status"`
}

type colStoreMasterFullv3 struct {
	Message       string          `json:"message"`
	DataStoreMasterFull []storeMasterFullv3 `json:"data"`
	Status        string          `json:"status"`
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

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	r.GET("/master/:name", func(c *gin.Context) {
		c.File("/home/yogi/go/src/github.com/pstp-api/app-release.apk")
	})

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))
	
	r.POST("/getSalesAll", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")

		dbname = sellerDivision(xusername)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select * from (select '%' as sellercode,'%' as sellername from dist_branch where branchcode='01' union select s.sellercode,s.sellername from dist_seller s join dist_spv v on v.spvcode=s.spvcode where v.spvcode='" + xusername + "' and s.active='1') d order by sellername "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var results []dataSales

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataSales{
				Sellercode: sellercode,
				Sellername: sellername,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})
	
	
	r.POST("/getChannel", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")

		dbname = sellerDivision(xusername)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sc.channelcode,sc.channeldesc  from store_channel sc join sv_validation_channel svc on svc.channelcode = sc.channelcode and svc.active = '1' where sc.channelcode not like '%/NOO' and sc.channelcode not like '%/CB%' and left(sc.channelcode,2)=right(left('"+xusername+"',5),2)  order by sc.channelcode  "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var channelcode string
		var channeldesc string
		var results []dataChannel

		for rows.Next() {
			err = rows.Scan(&channelcode, &channeldesc)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataChannel{
				ChannelCode: channelcode,
				ChannelDesc: channeldesc,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/dataTime", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")

		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select " +
			" sum(((secondtime-(((secondtime::int)%60)))/60)::int) as timeminute, " +
			" sum(((secondtime-(((secondtime::int)%60)))/3600)::int) as timehour, " +
			" sum((((secondtime-(((secondtime::int)%60)))/3600)/24)::int) as timeday " +
			" from (select sellercode,sellername, " +
			" 	  sum(EXTRACT(EPOCH FROM (timeout::timestamp - timein::timestamp))) as secondtime " +
			"   from ( " +
			" 	  select distinct r.sellercode as sellercode,r.sellername,h.timein,h.timeout   " +
			"   from store_master m " +
			"   join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'   " +
			"   join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month', current_date) as date)               " +
			"   WHERE r.active='1' and r.spvcode='" + xusername + "'  " +
			" ) a group by sellercode,sellername  " +
			" ) s "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var timehour string
		var timeminute string
		var timeday string
		var results []dataTime

		for rows.Next() {
			err = rows.Scan(&timeminute, &timehour, &timeday)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataTime{
				TimeMinute: timeminute,
				TimeHour:   timehour,
				TimeDay:    timeday,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/dataTimeASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")

		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select " +
			" sum(((secondtime-(((secondtime::int)%60)))/60)::int) as timeminute, " +
			" sum(((secondtime-(((secondtime::int)%60)))/3600)::int) as timehour, " +
			" sum((((secondtime-(((secondtime::int)%60)))/3600)/24)::int) as timeday " +
			" from (select sellercode,sellername, " +
			" 	  sum(EXTRACT(EPOCH FROM (timeout::timestamp - timein::timestamp))) as secondtime " +
			"   from ( " +
			" 	  select distinct r.sellercode as sellercode,r.sellername,h.timein,h.timeout   " +
			"   from store_master m " +
			"   join dist_seller r on r.sellercode=m.sellercode  " +
			"   join dist_seller r2 on r2.sellercode=r.spvcode  and r2.sellercode like '" + xsalescode + "'   " +
			"   join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month', current_date) as date)               " +
			"   WHERE r.active='1' and r2.spvcode='" + xusername + "'  " +
			" ) a group by sellercode,sellername  " +
			" ) s "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var timehour string
		var timeminute string
		var timeday string
		var results []dataTime

		for rows.Next() {
			err = rows.Scan(&timeminute, &timehour, &timeday)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataTime{
				TimeMinute: timeminute,
				TimeHour:   timehour,
				TimeDay:    timeday,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getMTDSummary", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsellercode := c.PostForm("sellercode")

		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select a.periodid,a.spvname,sum(coalesce(t.salestarget,0)) as targetsales,a.actsales,a.targetcall,a.actcall,a.actpcall,a.targetgdp,a.actgdp,a.targetiff,a.actiff from ( " +
			" select o.periodid,v.spvname,v.spvcode,sum(o.salestarget) targetsales,sum(o.salesvol/1.1) as actsales, " +
			" sum(o.targetcall) as targetcall,sum(o.callmtd) as actcall, " +
			" sum(o.prodcallmtd) as actpcall,sum(o.sbdtarget) as targetgdp,sum(o.goldenpoint) as actgdp, " +
			" sum(o.ifasttarget) as targetiff,sum(o.ifastpoint) as actiff " +
			" from store_objective o  " +
			" join store_master m on m.storecode=o.storecode and m.channelcode not like '%/NOO' " +
			" join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsellercode + "' " +
			" join dist_spv v on v.spvcode=r.spvcode and v.spvcode='" + xusername + "' " +
			" where o.periodid::character varying=to_char(now()::date,'YYYYMM') " +
			" group by o.periodid,v.spvname,v.spvcode	 ) a " +
			" join dist_seller l on l.spvcode=a.spvcode " +
			" left join seller_target t on t.sellercode=l.sellercode and t.periodid=a.periodid " +
			" group by a.periodid,a.spvname,a.actsales,a.targetcall,a.actcall,a.actpcall,a.targetgdp,a.actgdp,a.targetiff,a.actiff "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var periodid string
		var spvname string
		var targetsales string
		var targetiff string
		var actiff string
		var actsales string
		var targetcall string
		var actcall string
		var actpcall string
		var targetgdp string
		var actgdp string
		var results []summaryMTD

		for rows.Next() {
			err = rows.Scan(&periodid, &spvname, &targetsales, &actsales, &targetcall, &actcall, &actpcall, &targetgdp, &actgdp, &targetiff, &actiff)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := summaryMTD{
				Periodid:    periodid,
				SPVName:     spvname,
				TargetSales: targetsales,
				TargetIFF:   targetiff,
				ActualIFF:   actiff,
				ActualSales: actsales,
				TargetCall:  targetcall,
				ActualCall:  actcall,
				ActualPCall: actpcall,
				TargetGDP:   targetgdp,
				ActualGDP:   actgdp,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getMTDSummaryASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsellercode := c.PostForm("sellercode")

		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select a.periodid,a.spvname,sum(coalesce(t.salestarget,0)) as targetsales,sum(a.actsales) actsales,sum(a.targetcall) as targetcall,sum(a.actcall) actcall,sum(a.actpcall) as actpcall,sum(a.targetgdp) targetgdp,sum(a.actgdp) actgdp,sum(a.targetiff) targetiff,sum(a.actiff) actiff from ( " +
			"select o.periodid,v.spvname,r2.sellercode as spvcode,r.sellercode,sum(o.salestarget) targetsales,sum(coalesce(o.salesvol,0)/1.1) as actsales, " +
			" sum(o.targetcall) as targetcall,sum(o.callmtd) as actcall, " +
			" sum(o.prodcallmtd) as actpcall,sum(o.sbdtarget) as targetgdp,sum(o.goldenpoint) as actgdp, " +
			" sum(o.ifasttarget) as targetiff,sum(o.ifastpoint) as actiff " +
			" from store_objective o  " +
			" join store_master m on m.storecode=o.storecode and m.channelcode not like '%/NOO' " +
			" join dist_seller r on r.sellercode=m.sellercode " +
			" join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsellercode + "' " +
			" join dist_spv v on v.spvcode=r2.spvcode and v.spvcode='" + xusername + "' " +
			" where o.periodid::character varying=to_char(now()::date,'YYYYMM') " +
			" group by o.periodid,v.spvname,r2.sellercode,r.sellercode	) a " +
			" join dist_seller l on l.spvcode=a.spvcode and l.sellercode=a.sellercode  " +
			" left join seller_target t on t.sellercode=l.sellercode and t.periodid=a.periodid " +
			" group by a.periodid,a.spvname "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var periodid string
		var spvname string
		var targetsales string
		var targetiff string
		var actiff string
		var actsales string
		var targetcall string
		var actcall string
		var actpcall string
		var targetgdp string
		var actgdp string
		var results []summaryMTD

		for rows.Next() {
			err = rows.Scan(&periodid, &spvname, &targetsales, &actsales, &targetcall, &actcall, &actpcall, &targetgdp, &actgdp, &targetiff, &actiff)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := summaryMTD{
				Periodid:    periodid,
				SPVName:     spvname,
				TargetSales: targetsales,
				TargetIFF:   targetiff,
				ActualIFF:   actiff,
				ActualSales: actsales,
				TargetCall:  targetcall,
				ActualCall:  actcall,
				ActualPCall: actpcall,
				TargetGDP:   targetgdp,
				ActualGDP:   actgdp,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/checkVersion", func(c *gin.Context) {
		xminor := c.PostForm("minor")

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		if xminor == "1" {
			sqlstring = " select versionno as version,link,description,readmelink from pstp_version where active='1' limit 1 "
		} else {
			sqlstring = " select versionno as version,link,description,readmelink from pstp_version where active='1' limit 1 "
		}

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var version string
		var link string
		var description string
		var readmelink string
		var results []checkVersion

		for rows.Next() {
			err = rows.Scan(&version, &link, &description, &readmelink)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := checkVersion{
				Version:     version,
				Link:        link,
				Description: description,
				Readmelink:  readmelink,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getStartCall", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sellercode,sellername,to_char(timein,'DD-MM-YYYY HH24:MI') as timein,case when to_char(timein,'HH24MI')::int<0700 or to_char(timein,'HH24MI')::int<0800 then '07.00' when to_char(timein,'HH24MI')::int>0800 and to_char(timein,'HH24MI')::int<0900 then '08.00' when to_char(timein,'HH24MI')::int>0900 and to_char(timein,'HH24MI')::int<1000 then '09.00' else '10.00' end as timeclock   from (select distinct r.sellercode as sellercode,r.sellername,min(h.timein) over (partition by r.sellercode) as timein from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on 1=1 join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "' join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r.spvcode='" + xusername + "') a "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var timein string
		var timeclock string
		var results []dataStartCall

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername, &timein, &timeclock)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataStartCall{
				Sellercode: sellercode,
				Sellername: sellername,
				Timein:     timein,
				Timeclock:  timeclock,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getStartCallSummary", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select coalesce(to_char(avg(timein),'HH24:MI'),'00:00') as avgtime from ( " +
			" select sellercode,sellername,avg(times) as timein from ( " +
			" 	select sellercode,sellername,to_timestamp(timein, 'YYYY-MM-DD HH24:MI')::time as times from ( " +
			" 		select date,sellercode,sellername,to_char(timein,'YYYY-MM-DD HH24:MI') as timein    " +
			" 		from ( " +
			" 			select distinct date,r.sellercode as sellercode,r.sellername,min(h.timein) over (partition by h.date,r.sellercode) as timein  " +
			" 			from store_master m   " +
			" 			join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'  " +
			" 			join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month',current_date) as date) " +
			" 			WHERE r.active='1' and r.spvcode='" + xusername + "' " +
			" 		) a  " +
			" 	) g order by sellername " +
			" ) v group by sellercode,sellername " +
			" ) g		 "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var avgtime string
		avgtime = "00:00"

		for rows.Next() {
			err = rows.Scan(&avgtime)
			if err != nil {
				// handle this error
				panic(err)
			}
		}

		defer db.Close()

		c.String(http.StatusOK, avgtime)

	})

	r.POST("/getStartCallSummaryASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select coalesce(to_char(avg(timein),'HH24:MI'),'00:00') as avgtime from ( " +
			" select sellercode,sellername,avg(times) as timein from ( " +
			" 	select sellercode,sellername,to_timestamp(timein, 'YYYY-MM-DD HH24:MI')::time as times from ( " +
			" 		select date,sellercode,sellername,to_char(timein,'YYYY-MM-DD HH24:MI') as timein    " +
			" 		from ( " +
			" 			select distinct date,r.sellercode as sellercode,r.sellername,min(h.timein) over (partition by h.date,r.sellercode) as timein  " +
			" 			from store_master m   " +
			" 			join dist_seller r on r.sellercode=m.sellercode " +
			" 			join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  " +
			" 			join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month',current_date) as date) " +
			" 			WHERE r.active='1' and r2.spvcode='" + xusername + "' " +
			" 		) a  " +
			" 	) g order by sellername " +
			" ) v group by sellercode,sellername " +
			" ) g		 "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var avgtime string
		avgtime = "00:00"

		for rows.Next() {
			err = rows.Scan(&avgtime)
			if err != nil {
				// handle this error
				panic(err)
			}
		}

		defer db.Close()

		c.String(http.StatusOK, avgtime)

	})

	r.POST("/getDurationCallSummary", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select coalesce((((avg(avgtime)-(((avg(avgtime)::int)%60)))/3600)),0)::numeric(5,2) as avgtime from (select sellercode,sellername,avg(secondtime)::int as avgtime " +
			" from (select sellercode,sellername,dated,         " +
			" 	  sum(EXTRACT(EPOCH FROM (timeout::timestamp - timein::timestamp))) as secondtime     " +
			" 	  from (        " +
			" 		  select distinct r.sellercode as sellercode,h.date as dated,r.sellername,h.timein,h.timeout       " +
			" 		  from store_master m     " +
			" 		  join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'       " +
			" 		  join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month', current_date) as date)                   " +
			" 		  WHERE r.active='1' and r.spvcode='" + xusername + "'    " +
			" 	  ) a group by sellercode,sellername,dated   ) s  where secondtime>0  group by sellercode,sellername ) j "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var avgtime string
		avgtime = "00:00"

		for rows.Next() {
			err = rows.Scan(&avgtime)
			if err != nil {
				// handle this error
				panic(err)
			}
		}

		defer db.Close()

		c.String(http.StatusOK, avgtime)

	})

	r.POST("/getDurationCallSummaryASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select coalesce((((avg(avgtime)-(((avg(avgtime)::int)%60)))/3600)),0)::numeric(5,2) as avgtime from (select sellercode,sellername,avg(secondtime)::int as avgtime " +
			" from (select sellercode,sellername,dated,         " +
			" 	  sum(EXTRACT(EPOCH FROM (timeout::timestamp - timein::timestamp))) as secondtime     " +
			" 	  from (        " +
			" 		  select distinct r.sellercode as sellercode,h.date as dated,r.sellername,h.timein,h.timeout       " +
			" 		  from store_master m     " +
			" 		  join dist_seller r on r.sellercode=m.sellercode " +
			" 		  join dist_seller r2 on r2.sellercode=r.spvcode  and r2.sellercode like '" + xsalescode + "'       " +
			" 		  join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month', current_date) as date)                   " +
			" 		  WHERE r.active='1' and r2.spvcode='" + xusername + "'    " +
			" 	  ) a group by sellercode,sellername,dated   ) s  where secondtime>0  group by sellercode,sellername ) j "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var avgtime string
		avgtime = "00:00"

		for rows.Next() {
			err = rows.Scan(&avgtime)
			if err != nil {
				// handle this error
				panic(err)
			}
		}

		defer db.Close()

		c.String(http.StatusOK, avgtime)

	})

	r.POST("/getStopCallSummary", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select coalesce(to_char(avg(timein),'HH24:MI'),'23:59') as avgtime from ( " +
			" select sellercode,sellername,avg(times) as timein from ( " +
			" 	select sellercode,sellername,to_timestamp(timein, 'YYYY-MM-DD HH24:MI')::time as times from ( " +
			" 		select date,sellercode,sellername,to_char(timein,'YYYY-MM-DD HH24:MI') as timein    " +
			" 		from ( " +
			" 			select distinct date,r.sellercode as sellercode,r.sellername,max(h.timein) over (partition by h.date,r.sellercode) as timein  " +
			" 			from store_master m   " +
			" 			join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'  " +
			" 			join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month',current_date) as date) " +
			" 			WHERE r.active='1' and r.spvcode='" + xusername + "' " +
			" 		) a  " +
			" 	) g order by sellername " +
			" ) v group by sellercode,sellername " +
			" ) g		 "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var avgtime string
		avgtime = "23:59"

		for rows.Next() {
			err = rows.Scan(&avgtime)
			if err != nil {
				// handle this error
				panic(err)
			}
		}

		defer db.Close()

		c.String(http.StatusOK, avgtime)

	})

	r.POST("/getStopCallSummaryASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select coalesce(to_char(avg(timein),'HH24:MI'),'23:59') as avgtime from ( " +
			" select sellercode,sellername,avg(times) as timein from ( " +
			" 	select sellercode,sellername,to_timestamp(timein, 'YYYY-MM-DD HH24:MI')::time as times from ( " +
			" 		select date,sellercode,sellername,to_char(timein,'YYYY-MM-DD HH24:MI') as timein    " +
			" 		from ( " +
			" 			select distinct date,r.sellercode as sellercode,r.sellername,max(h.timein) over (partition by h.date,r.sellercode) as timein  " +
			" 			from store_master m   " +
			" 			join dist_seller r on r.sellercode=m.sellercode " +
			" 			join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  " +
			" 			join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date>=cast(date_trunc('month',current_date) as date) " +
			" 			WHERE r.active='1' and r2.spvcode='" + xusername + "' " +
			" 		) a  " +
			" 	) g order by sellername " +
			" ) v group by sellercode,sellername " +
			" ) g		 "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var avgtime string
		avgtime = "23:59"

		for rows.Next() {
			err = rows.Scan(&avgtime)
			if err != nil {
				// handle this error
				panic(err)
			}
		}

		defer db.Close()

		c.String(http.StatusOK, avgtime)

	})

	r.POST("/getStopCall", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sellercode,sellername,to_char(timein,'DD-MM-YYYY HH24:MI') as timein,case when to_char(timein,'HH24MI')::int<1600 then '15.00' when to_char(timein,'HH24MI')::int>16.00 and to_char(timein,'HH24MI')::int<1700 then '16.00' when to_char(timein,'HH24MI')::int>1700 and to_char(timein,'HH24MI')::int<1800 then '17.00' when to_char(timein,'HH24MI')::int>=1800 then '18.00' else '14.00' end as timeclock   from (select distinct r.sellercode as sellercode,r.sellername,max(h.timein) over (partition by r.sellercode) as timein	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on 1=1  join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "' join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r.spvcode='" + xusername + "') a "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var timein string
		var timeclock string
		var results []dataStartCall

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername, &timein, &timeclock)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataStartCall{
				Sellercode: sellercode,
				Sellername: sellername,
				Timein:     timein,
				Timeclock:  timeclock,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getDurationCall", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sellercode,sellername,((secondtime-(((secondtime::int)%60)))/60)::int as timein,((secondtime-(((secondtime::int)%60)))/3600)::int as timeclock from (select sellercode,sellername,sum(EXTRACT(EPOCH FROM (timeout::timestamp - timein::timestamp))) as secondtime from (select distinct r.sellercode as sellercode,r.sellername,h.timein,h.timeout	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on 1=1  join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date 		WHERE r.active='1' and r.spvcode='" + xusername + "' ) a group by sellercode,sellername ) s "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var timein string
		var timeclock string
		var results []dataStartCall

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername, &timein, &timeclock)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataStartCall{
				Sellercode: sellercode,
				Sellername: sellername,
				Timein:     timein,
				Timeclock:  timeclock,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolume", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " SELECT name,sum(coalesce(target,0)) as target,sum(coalesce(actual,0)/1.1) as actual from  " +
			"( " +
			"SELECT cast('-' as character varying) as name, " +
			"CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol,0)))/COALESCE(targetcall_r,0),2) END as target, " +
			"coalesce(total,0) as actual  " +
			"FROM dist_seller r  " +
			"JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'   " +
			"JOIN store_visit v on v.storecode=s.storecode  " +
			"JOIN ( (select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1  " +
			"LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM')  " +
			"LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date  " +
			"WHERE r.active='1' and r.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate    " + weekstr + " " +
			"union all " +
			"(SELECT cast('-' as character varying) as name, " +
			"coalesce(0,0) as target, " +
			"coalesce(total,0) as actual  " +
			"FROM dist_seller r  " +
			"JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'   " +
			"JOIN store_deviated v on v.storecode=s.storecode and v.date=now()::date " +
			"JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and i.date=now()::date  " +
			"WHERE r.active='1' and r.spvcode='" + xusername + "' ) " +
			") d group by name "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolumeSegmented", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual from (SELECT cast(c.channeldesc as character varying) as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual FROM dist_seller r JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'  JOIN store_visit v on v.storecode=s.storecode join store_channel c on c.channelcode=s.channelcode JOIN ( (select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date WHERE r.active='1' and r.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate  " + weekstr + " ) d group by name  "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getCall", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  		select '-' as name,coalesce(sum(target),0) as target,coalesce(sum(actual),0) as actual " +
			" from ( " +
			" 	select r.sellercode as sellercode,r.sellername as name,case when m.channelcode like '%/NOO' then count(null) else count(v.storecode) end as target,count(h.storecode) as actual,case when coalesce(count(v.storecode),0)=0 then count(h.storecode) else ((coalesce(count(h.storecode),0)*100)/coalesce(count(v.storecode),0))::int end as procent  " +
			" 	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate  " +
			" 	join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + "  " +
			" 	join store_master m on m.storecode=v.storecode  " +
			" 	join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  " +
			" 	left join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r.spvcode='" + xusername + "'  " +
			" 	group by r.sellercode,r.sellername,m.channelcode  " +
			" ) f  "

		//
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getPCall", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  	 select '-' as name,coalesce(sum(target),0) as target,coalesce(sum(actual),0) as actual " +
			" from ( " +
			" select r.sellercode as sellercode,r.sellername as name,case when m.channelcode like '%/NOO' then count(null) else count(distinct v.storecode) end  as target,count(distinct h.storecode) as actual   " +
			" from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate  " +
			" join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate   " + weekstr + " " +
			" join store_master m on m.storecode=v.storecode  " +
			" join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  " +
			" left join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date  " +
			" WHERE r.active='1' and r.spvcode='" + xusername + "'  " +
			" group by r.sellercode,r.sellername,m.channelcode " +
			" union all " +
			"  select r.sellercode as sellercode,r.sellername as name,count(null) as target,count(distinct h.storecode) as actual  " +
			" from store_deviated v  " +
			" join store_master m on m.storecode=v.storecode  " +
			" join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  " +
			" join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date  " +
			"  WHERE r.active='1' and r.spvcode='" + xusername + "' and v.date=now()::date " +
			" group by r.sellercode,r.sellername  " +
			" ) f		"

		//
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolumeBreakdown", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select * from (SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual,cast(case when sum(target)<=0 then sum(actual/1.1)*100/1  else sum(actual/1.1)*100/sum(target) end as numeric(16,0)) as procent,sellercode from (	SELECT r.sellercode,r.sellername as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual FROM dist_seller r JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'  JOIN store_visit v on v.storecode=s.storecode JOIN ( ( select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date WHERE r.active='1' and r.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + " ) d group by name,sellercode ) f order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFF", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string
		var strsales string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		strsales = ""

		if xsalescode == "%" {
			strsales = ""
		} else {
			strsales = xsalescode
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select s.name,s.target,coalesce(t.storecode,0) as actual from ( " +
			" 	select i.initdesc as name,count(distinct m.storecode) as target from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate " +
			" 	join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate  " + weekstr +
			" 	join store_master m on m.storecode=v.storecode " +
			" 	join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "' " +
			" 	join initiative_master i on i.periodid::character varying=to_char(now()::date,'YYYYMM') and i.channelcode=m.channelcode	" +
			" 	WHERE r.active='1' and r.spvcode='" + xusername + "' " +
			" 	group by i.initdesc ) s " +
			" left join (	" +
			" 	select m.initdesc as name,count(strcode) as storecode from fninitiativestatustoday2(now()::date,'" + xusername + "','" + strsales + "') i " +
			"	join initiative_master m on m.id=i.intid group by m.initdesc " +
			" ) t on t.name=s.name "
		rows, err := db.Query(sqlstring)

		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFFMTD", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string
		var strsales string

		strsales = ""

		if xsalescode == "%" {
			strsales = ""
		} else {
			strsales = xsalescode
		}

		sqlstring = " select s.name,s.target,coalesce(t.storecode,0) as actual from ( " +
			" 	select i.initdesc as name,count(distinct m.storecode) as target from store_master m  " +
			" 	join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "' " +
			" 	join initiative_master i on i.periodid::character varying=to_char(now()::date,'YYYYMM') and i.channelcode=m.channelcode	" +
			" 	WHERE m.channelcode not like '%/NOO' and m.storeflag!='D' and r.active='1' and r.spvcode='" + xusername + "' " +
			" 	group by i.initdesc ) s " +
			" left join (	" +
			" 	select m.initdesc as name,count(strcode) as storecode from fninitiativestatus2(to_char(now()::date,'YYYYMM')::int,'" + xusername + "','" + strsales + "') i " +
			"	join initiative_master m on m.id=i.intid group by m.initdesc " +
			" ) t on t.name=s.name "
		rows, err := db.Query(sqlstring)

		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFFMTDDetailSPV", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		xiffname := c.PostForm("iffname")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string
		var strsales string

		strsales = ""

		if xsalescode == "%" {
			strsales = ""
		} else {
			strsales = xsalescode
		}

		sqlstring = " SELECT " +
			" sdt.*  " +
			" FROM " +
			" ( " +
			" SELECT " +
			" salesname, " +
			" salescode, " +
			" s.NAME, " +
			" SUM ( s.target ) AS total_target, " +
			" COALESCE ( T.storecode, 0 ) AS actual  " +
			" FROM " +
			" ( " +
			" SELECT " +
			" r.sellercode AS salescode, " +
			" r.sellername AS salesname, " +
			" i.initdesc AS NAME, " +
			" COUNT ( DISTINCT M.storecode ) AS target  " +
			" FROM " +
			" store_master " +
			" M JOIN dist_seller r ON r.sellercode = M.sellercode  " +
			" AND r.sellercode LIKE '" + xsalescode + "' " +
			" JOIN initiative_master i ON i.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' )  " +
			" AND i.channelcode = M.channelcode  " +
			" WHERE " +
			" M.channelcode NOT LIKE'%/NOO'  " +
			" AND M.storeflag != 'D'  " +
			" AND r.active = '1'  " +
			" AND r.spvcode = $1  " +
			" GROUP BY " +
			" i.initdesc, " +
			" r.sellercode, " +
			" r.sellername, " +
			" i.ID  " +
			" ) s " +
			" LEFT JOIN ( " +
			" SELECT M.initdesc AS NAME, " +
			" sm.sellercode, " +
			" COUNT ( strcode ) AS storecode  " +
			" FROM " +
			" fninitiativestatus2 ( to_char( now( ) :: DATE, 'YYYYMM' ) :: INT, '" + xusername + "', '" + strsales + "' ) i " +
			" JOIN initiative_master M ON M.ID = i.intid " +
			" JOIN store_master sm ON sm.storecode = i.strcode  " +
			" GROUP BY " +
			" M.initdesc, " +
			" sm.sellercode  " +
			" ) T ON T.NAME = s.NAME  " +
			" AND T.sellercode = s.salescode  " +
			" WHERE " +
			" s.NAME = $2  " +
			" GROUP BY " +
			" s.salesname, " +
			" s.salescode, " +
			" T.storecode, " +
			" s.NAME  " +
			" ) sdt  " +
			" ORDER BY " +
			" sdt.actual DESC "
		rows, err := db.Query(sqlstring, xusername, xiffname)

		//log.Fatal(xiffname)

		var sellername string
		var sellercode string
		var iffname string
		var target string
		var actual string
		var results []dataIFFDetail

		if err != nil {
			panic(err)

			result := dataIFFDetail{
				Sellercode: "",
				Sellername: "",
				Iffname:    "",
				Target:     "",
				Actual:     "",
			}
			results = append(results, result)

			colIFFInit := colIFF{
				Message:       "OK",
				DataIFFDetail: results,
				Status:        "1",
			}

			c.JSON(http.StatusOK, colIFFInit)
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&sellername, &sellercode, &iffname, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataIFFDetail{
				Sellercode: sellercode,
				Sellername: sellername,
				Iffname:    iffname,
				Target:     target,
				Actual:     actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		colIFFInit := colIFF{
			Message:       "OK",
			DataIFFDetail: results,
			Status:        "1",
		}

		c.JSON(http.StatusOK, colIFFInit)

	})

	r.POST("/getGDP", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select  '-' as name,coalesce(sum(g.sbdtarget),0) as target,coalesce(sum(g.gp),0) as actual from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate  " + weekstr + "  join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'  join fngoldenpointtoday(now()::date,'" + xusername + "','') g on g.strcode=m.storecode WHERE r.active='1' and r.spvcode='" + xusername + "'  "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getDeviation", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select  '-' as name,coalesce(count(v.storecode),0) as target,coalesce(count(g.storecode),0) as actual from store_master m join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'  join getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate on 1=1 left join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate  " + weekstr + "  and v.storecode=m.storecode left join store_deviated g on g.storecode=m.storecode and g.date=now()::date and g.sellercode=r.sellercode WHERE r.active='1' and r.spvcode='" + xusername + "' "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getComplaint", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	 select m.storecode,m.name as storename,r.sellercode,r.sellername,case when coalesce(m.longitude,'0')='' then '0' else coalesce(m.longitude,'0') end as longitude, " +
			" case when coalesce(m.latitude,'0')='' then '0' else coalesce(m.latitude,'0') end latitude,coalesce(g.longitude,0) as longitudecall,coalesce(g.latitude,0) as latitudecall  " +
			" from store_master m  " +
			" join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  " +
			" join getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate on 1=1 " +
			" join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate   " + weekstr + "  and v.storecode=m.storecode  " +
			" left join store_callhist g on g.storecode=m.storecode and g.date=now()::date and g.sellercode=r.sellercode  " +
			" WHERE r.active='1' and r.spvcode='" + xusername + "' and m.channelcode not like '%/NOO' "
		/**" union all " +
		" select m.storecode,m.name as storename,r.sellercode,r.sellername,case when coalesce(m.longitude,'0')='' then '0' else coalesce(m.longitude,'0') end as longitude, " +
		" case when coalesce(m.latitude,'0')='' then '0' else coalesce(m.latitude,'0') end latitude,coalesce(g.longitude,0) as longitudecall,coalesce(g.latitude,0) as latitudecall  " +
		" from store_master m  " +
		" join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  " +
		" join store_deviated v on  v.storecode=m.storecode and date=now()::date " +
		" join store_callhist g on g.storecode=m.storecode and g.date=now()::date and g.sellercode=r.sellercode  " +
		" WHERE r.active='1' and r.spvcode='" + xusername + "' "**/

		//

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var storecode string
		var storename string
		var longitude string
		var longitudecall string
		var latitude string
		var latitudecall string
		var results []dataComplaint

		for rows.Next() {
			err = rows.Scan(&storecode, &storename, &sellercode, &sellername, &longitude, &latitude, &longitudecall, &latitudecall)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataComplaint{
				Sellercode:    sellercode,
				Sellername:    sellername,
				Storecode:     storecode,
				Storename:     storename,
				Latitude:      latitude,
				Latitudecall:  latitudecall,
				Longitude:     longitude,
				Longitudecall: longitudecall,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	//ASM Menu

	r.POST("/getStartCallASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sellercode,sellername,to_char(timein,'DD-MM-YYYY HH24:MI') as timein,case when to_char(timein,'HH24MI')::int<0700 or to_char(timein,'HH24MI')::int<0800 then '07.00' when to_char(timein,'HH24MI')::int>0800 and to_char(timein,'HH24MI')::int<0900 then '08.00' when to_char(timein,'HH24MI')::int>0900 and to_char(timein,'HH24MI')::int<1000 then '09.00' else '10.00' end as timeclock   from (select distinct r2.sellercode as sellercode,r2.sellername,min(h.timein) over (partition by r2.sellercode) as timein from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on 1=1 join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode JOIN dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "' join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "') a "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var timein string
		var timeclock string
		var results []dataStartCall

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername, &timein, &timeclock)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataStartCall{
				Sellercode: sellercode,
				Sellername: sellername,
				Timein:     timein,
				Timeclock:  timeclock,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getStopCallASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sellercode,sellername,to_char(timein,'DD-MM-YYYY HH24:MI') as timein,case when to_char(timein,'HH24MI')::int<1600 then '15.00' when to_char(timein,'HH24MI')::int>16.00 and to_char(timein,'HH24MI')::int<1700 then '16.00' when to_char(timein,'HH24MI')::int>1700 and to_char(timein,'HH24MI')::int<1800 then '17.00' when to_char(timein,'HH24MI')::int>=1800 then '18.00' else '14.00' end as timeclock   from (select distinct r2.sellercode as sellercode,r2.sellername,max(h.timein) over (partition by r2.sellercode) as timein	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on 1=1  join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode JOIN dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "' join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "') a "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var timein string
		var timeclock string
		var results []dataStartCall

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername, &timein, &timeclock)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataStartCall{
				Sellercode: sellercode,
				Sellername: sellername,
				Timein:     timein,
				Timeclock:  timeclock,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getDurationCallASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sellercode,sellername,((secondtime-(((secondtime::int)%60)))/60)::int as timein,((secondtime-(((secondtime::int)%60)))/3600)::int as timeclock from (select sellercode,sellername,sum(EXTRACT(EPOCH FROM (timeout::timestamp - timein::timestamp))) as secondtime from (select distinct r2.sellercode as sellercode,r2.sellername,h.timein,h.timeout	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on 1=1  join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode  JOIN dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date 		WHERE r.active='1' and r2.spvcode='" + xusername + "' ) a group by sellercode,sellername ) s "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var timein string
		var timeclock string
		var results []dataStartCall

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername, &timein, &timeclock)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataStartCall{
				Sellercode: sellercode,
				Sellername: sellername,
				Timein:     timein,
				Timeclock:  timeclock,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolumeASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual from (SELECT cast('-' as character varying) as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual FROM dist_seller r join dist_seller r2 on r2.sellercode=r.spvcode JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'  JOIN store_visit v on v.storecode=s.storecode JOIN ( (select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + " union all SELECT cast('-' as character varying) as name,coalesce(0,0) as target,coalesce(total,0) as actual FROM dist_seller r join dist_seller r2 on r2.sellercode=r.spvcode JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '"+ xsalescode +"' JOIN store_deviated v on v.storecode=s.storecode and v.date=now()::date LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and i.date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' ) d group by name "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolumeSegmentedASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual from (SELECT cast(c.channeldesc as character varying) as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual FROM dist_seller r  JOIN dist_seller r2 on r2.sellercode=r.spvcode   JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'  JOIN store_visit v on v.storecode=s.storecode join store_channel c on c.channelcode=s.channelcode JOIN ( (select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + " ) d group by name  "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getCallASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  		select '-' as name,coalesce(sum(target),0) as target,coalesce(sum(actual),0) as actual from (select r2.sellercode as sellercode,r2.sellername as name,count(v.storecode) as target,count(h.storecode) as actual,case when coalesce(count(v.storecode),0)=0 then count(h.storecode) else ((coalesce(count(h.storecode),0)*100)/coalesce(count(v.storecode),0))::int end as procent from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode JOIN dist_seller r2 on r2.sellercode=r.spvcode   and r2.sellercode like '" + xsalescode + "'  left join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' group by r2.sellercode,r2.sellername ) f  "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getPCallASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  	select '-' as name,coalesce(sum(target),0) as target,coalesce(sum(actual),0) as actual from (select r2.sellercode as sellercode,r2.sellername as name,count(v.storecode) as target,count(h.storecode) as actual,case when coalesce(count(v.storecode),0)=0 then count(h.storecode) else ((coalesce(count(h.storecode),0)*100)/coalesce(count(v.storecode),0))::int end as procent from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode  JOIN dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "' left join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' group by r2.sellercode,r2.sellername ) f  "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolumeBreakdownASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select * from (SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual,cast(case when sum(target)<=0 then sum(actual/1.1)*100/1  else sum(actual/1.1)*100/sum(target) end as numeric(16,0)) as procent,sellercode from (	SELECT r2.sellercode,r2.sellername as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual FROM dist_seller r JOIN dist_seller r2 on r2.sellercode=r.spvcode JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'  JOIN store_visit v on v.storecode=s.storecode JOIN ( ( select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + " ) d group by name,sellercode ) f order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFFASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select sellercode from dist_seller where active='1' and spvcode='" + xusername + "' and sellercode like '" + xsalescode + "'  "
		rowsseller, errseller := db.Query(sqlstring)
		if errseller != nil {
			panic(errseller)
		}

		var sellercode string
		var querycompleted string
		var cunion int
		cunion = 0

		for rowsseller.Next() {
			errseller = rowsseller.Scan(&sellercode)
			if errseller != nil {
				// handle this error
				panic(errseller)
			}

			if cunion > 0 {
				querycompleted = querycompleted + " union select * from fninitiativestatustoday2(now()::date,'" + sellercode + "','') "
			} else {
				querycompleted = querycompleted + " select * from fninitiativestatustoday2(now()::date,'" + sellercode + "','') "
			}
			cunion++
		}

		defer rowsseller.Close()

		sqlstring = " select s.name,s.target,coalesce(t.storecode,0) as actual from ( " +
			" select i.initdesc as name,count(distinct m.storecode) as target from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate " +
			"join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr +
			"join store_master m on m.storecode=v.storecode " +
			"join dist_seller r on r.sellercode=m.sellercode	" +
			"join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "' " +
			"join initiative_master i on i.periodid::character varying=to_char(now()::date,'YYYYMM') and i.channelcode=m.channelcode	" +
			"WHERE r.active='1' and r2.spvcode='" + xusername + "'	" +
			"group by i.initdesc " +
			") s " +
			"left join (" +
			"	select m.initdesc as name,count(strcode) as storecode from (" +
			"		" + querycompleted + "" +
			"	) i join initiative_master m on m.id=i.intid group by m.initdesc " +
			") t on t.name=s.name		 "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFFASMMTD", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " 	select sellercode from dist_seller where active='1' and spvcode='" + xusername + "' and sellercode like '" + xsalescode + "'  "
		rowsseller, errseller := db.Query(sqlstring)
		if errseller != nil {
			panic(errseller)
		}

		var sellercode string
		var querycompleted string
		var cunion int
		cunion = 0

		for rowsseller.Next() {
			errseller = rowsseller.Scan(&sellercode)
			if errseller != nil {
				// handle this error
				panic(errseller)
			}

			if cunion > 0 {
				querycompleted = querycompleted + " union select * from fninitiativestatus2(to_char(now()::date,'YYYYMM')::int,'" + sellercode + "','') "
			} else {
				querycompleted = querycompleted + " select * from fninitiativestatus2(to_char(now()::date,'YYYYMM')::int,'" + sellercode + "','') "
			}
			cunion++
		}

		defer rowsseller.Close()

		sqlstring = " select s.name,s.target,coalesce(t.storecode,0) as actual from ( " +
			" select i.initdesc as name,count(distinct m.storecode) as target from store_master m " +
			"join dist_seller r on r.sellercode=m.sellercode	" +
			"join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "' " +
			"join initiative_master i on i.periodid::character varying=to_char(now()::date,'YYYYMM') and i.channelcode=m.channelcode	" +
			"WHERE m.channelcode not like '%/NOO' and m.storeflag!='D' and r.active='1' and r2.spvcode='" + xusername + "'	" +
			"group by i.initdesc " +
			") s " +
			"left join (" +
			"	select m.initdesc as name,count(distinct strcode) as storecode from (" +
			"		" + querycompleted + "" +
			"	) i join initiative_master m on m.id=i.intid group by m.initdesc " +
			") t on t.name=s.name		 "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFFDetailASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select sellercode from dist_seller where active='1' and spvcode='" + xusername + "'  "
		rowsseller, errseller := db.Query(sqlstring)
		if errseller != nil {
			panic(errseller)
		}

		var sellercodex string
		var querycompleted string
		var cunion int
		cunion = 0

		for rowsseller.Next() {
			errseller = rowsseller.Scan(&sellercodex)
			if errseller != nil {
				// handle this error
				panic(errseller)
			}

			if cunion > 0 {
				querycompleted = querycompleted + " union select * from fninitiativestatustoday2(now()::date,'" + sellercodex + "','') "
			} else {
				querycompleted = querycompleted + " select * from fninitiativestatustoday2(now()::date,'" + sellercodex + "','') "
			}
			cunion++
		}

		defer rowsseller.Close()

		sqlstring = " select s.name,s.target,coalesce(t.storecode,0) as actual,case when coalesce(t.storecode,0)=0 then 0 else (coalesce(t.storecode,1)*100)/coalesce(s.target,0) end as procent,s.sellercode,s.sellername from (select r2.sellercode,r2.sellername,i.initdesc as name,count(distinct m.storecode) as target " +
			" from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate " +
			" join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr +
			" join store_master m on m.storecode=v.storecode " +
			" join dist_seller r on r.sellercode=m.sellercode join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "' 	" +
			" join initiative_master i on i.periodid::character varying=to_char(now()::date,'YYYYMM') and i.channelcode=m.channelcode	WHERE r.active='1' and r2.spvcode='" + xusername + "' " +
			" group by i.initdesc,r2.sellercode,r2.sellername ) s " +
			" left join (	" +
			"	select r2.sellercode,m.initdesc as name,count(strcode) as storecode from " +
			" ( " + querycompleted +
			") i " +
			"join initiative_master m on m.id=i.intid  join store_master s on s.storecode=i.strcode  " +
			"join dist_seller l on l.sellercode=s.sellercode join dist_seller r2 on r2.sellercode=l.spvcode  " +
			"group by r2.sellercode,m.initdesc " +
			") t on t.name=s.name and t.sellercode=s.sellercode " +
			"order by s.sellercode,s.name "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var sellername string
		var results []dataVolumeDetailIFF

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode, &sellername)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetailIFF{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
				Sellername: sellername,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getGDPASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select sellercode from dist_seller where active='1' and spvcode='" + xusername + "'  "
		rowsseller, errseller := db.Query(sqlstring)
		if errseller != nil {
			panic(errseller)
		}

		var sellercode string
		var querycompleted string
		var cunion int
		cunion = 0

		for rowsseller.Next() {
			errseller = rowsseller.Scan(&sellercode)
			if errseller != nil {
				// handle this error
				panic(errseller)
			}

			if cunion > 0 {
				querycompleted = querycompleted + " union select * from fngoldenpointtoday(now()::date,'" + sellercode + "','') "
			} else {
				querycompleted = querycompleted + " select * from fngoldenpointtoday(now()::date,'" + sellercode + "','') "
			}
			cunion++
		}

		defer rowsseller.Close()

		sqlstring = " 	select  '-' as name,coalesce(sum(g.sbdtarget),0) as target,coalesce(sum(g.gp),0) as actual from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " 	join store_master m on m.storecode=v.storecode 	join dist_seller r on r.sellercode=m.sellercode join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  join ( " + querycompleted + " ) g on g.strcode=m.storecode WHERE r.active='1' and r2.spvcode='" + xusername + "'  "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getDeviationASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select  '-' as name,coalesce(count(v.storecode),0) as target,coalesce(count(g.storecode),0) as actual from store_master m join dist_seller r on r.sellercode=m.sellercode join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  join getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate on 1=1 left join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " and v.storecode=m.storecode left join store_deviated g on g.storecode=m.storecode and g.date=now()::date and g.sellercode=r.sellercode WHERE r.active='1' and r2.spvcode='" + xusername + "' "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var results []dataVolume

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolume{
				Name:   name,
				Target: target,
				Actual: actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getComplaintASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select m.storecode,m.name as storename,r2.sellercode,r2.sellername,case when coalesce(m.longitude,'0')='' then '0' else coalesce(m.longitude,'0') end as longitude,case when coalesce(m.latitude,'0')='' then '0' else coalesce(m.latitude,'0') end latitude,coalesce(g.longitude,0) as longitudecall,coalesce(g.latitude,0) as latitudecall from store_master m join dist_seller r on r.sellercode=m.sellercode  join dist_seller r2 on r2.sellercode=r.spvcode   and r.sellercode like '" + xsalescode + "' join getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate on 1=1 join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " and v.storecode=m.storecode left join store_callhist g on g.storecode=m.storecode and g.date=now()::date and g.sellercode=r.sellercode WHERE r.active='1' and r2.spvcode='" + xusername + "' "
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var storecode string
		var storename string
		var longitude string
		var longitudecall string
		var latitude string
		var latitudecall string
		var results []dataComplaint

		for rows.Next() {
			err = rows.Scan(&storecode, &storename, &sellercode, &sellername, &longitude, &latitude, &longitudecall, &latitudecall)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataComplaint{
				Sellercode:    sellercode,
				Sellername:    sellername,
				Storecode:     storecode,
				Storename:     storename,
				Latitude:      latitude,
				Latitudecall:  latitudecall,
				Longitude:     longitude,
				Longitudecall: longitudecall,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolumeDetailASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		xsession := c.PostForm("session")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		if xsession != "" {
			rowsinsses, errsses := db.Query("INSERT INTO public.login_session(session, sellercode, description) VALUES ($1, $2, 'VOLUME DETAIL'); ", xsession, xusername)
			if errsses != nil {
				// handle this error better than this
				fmt.Print(errsses)
			}

			defer rowsinsses.Close()
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select * from (SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual,cast(case when sum(target)<=0 then sum(actual/1.1)*100/1  else sum(actual/1.1)*100/sum(target) end as numeric(16,0)) as procent,sellercode from (	SELECT r2.sellercode,r2.sellername as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual FROM dist_seller r  JOIN dist_seller r2 on r2.sellercode=r.spvcode  JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'  JOIN store_visit v on v.storecode=s.storecode JOIN ( ( select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + "  union all SELECT r2.sellercode,r2.sellername as name,coalesce(0,0) as  target,coalesce(total,0) as actual FROM dist_seller r JOIN dist_seller r2 on r2.sellercode=r.spvcode JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '"+ xsalescode +"' JOIN store_deviated v on v.storecode=s.storecode and v.date=now()::date JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and i.date=now()::date WHERE r.active='1' and r2.spvcode='"+ xusername +"'   ) d group by name,sellercode ) f order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	// End ASM Menu

	r.POST("/getVolumeDetail", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		xsession := c.PostForm("session")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		if xsession != "" {
			rowsinsses, errsses := db.Query("INSERT INTO public.login_session(session, sellercode, description) VALUES ($1, $2, 'VOLUME DETAIL'); ", xsession, xusername)
			if errsses != nil {
				// handle this error better than this
				fmt.Print(errsses)
			}

			defer rowsinsses.Close()
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  select * from ( " +
			" SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual, " +
			" cast(case when sum(target)<=0 then sum(actual/1.1)*100/1  else sum(actual/1.1)*100/sum(target) end as numeric(16,0)) as procent,sellercode  " +
			" from (  " +
			" 	SELECT sellercode,name,target,sum(actual) as actual from (SELECT r.sellercode,r.sellername as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual " +
			" 	FROM dist_seller r  " +
			" 	JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'   " +
			" 	JOIN store_visit v on v.storecode=s.storecode  " +
			" 	JOIN ( ( select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1  " +
			" 	LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM')  " +
			" 	LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date  " +
			" 	WHERE r.active='1' and r.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate  " + weekstr + "  ) g group by sellercode,name,target " +
			" 	union all " +
			" 	 SELECT r.sellercode,r.sellername as name,coalesce(0,0) as target,coalesce(total,0) as actual  " +
			" 	FROM dist_seller r  " +
			" 	JOIN store_master s on s.sellercode=r.sellercode  and s.sellercode like '" + xsalescode + "'  " +
			" 	JOIN store_deviated v on v.storecode=s.storecode and v.date=now()::date " +
			" 	JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and i.date=now()::date  " +
			" 	WHERE r.active='1' and r.spvcode='" + xusername + "' " +
			" ) d group by name,sellercode  " +
			") f order by procent asc "

		//
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFFDetail", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select s.name,s.target,coalesce(t.storecode,0) as actual,case when coalesce(t.storecode,0)=0 then 0 else (coalesce(t.storecode,1)*100)/coalesce(s.target,0) end as procent,s.sellercode,s.sellername from (select r.sellercode,r.sellername,i.initdesc as name,count(distinct m.storecode) as target from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate  " + weekstr + " join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "' 	join initiative_master i on i.periodid::character varying=to_char(now()::date,'YYYYMM') and i.channelcode=m.channelcode	WHERE r.active='1' and r.spvcode='" + xusername + "' group by i.initdesc,r.sellercode,r.sellername ) s left join (	select l.sellercode,m.initdesc as name,count(strcode) as storecode from fninitiativestatustoday2(now()::date,'" + xusername + "','') i join initiative_master m on m.id=i.intid  join store_master s on s.storecode=i.strcode  join dist_seller l on l.sellercode=s.sellercode  group by l.sellercode,m.initdesc ) t on t.name=s.name and t.sellercode=s.sellercode order by s.sellercode,s.name "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var sellername string
		var results []dataVolumeDetailIFF

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode, &sellername)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetailIFF{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
				Sellername: sellername,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getGDPDetail", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select name,target,actual,procent,sellercode from (select  r.sellercode,r.sellername as name,sum(g.sbdtarget) as target,sum(g.gp) as actual,case when coalesce(sum(g.sbdtarget),0)=0 then sum(g.gp) else ((coalesce(sum(g.gp),0)*100)/coalesce(sum(g.sbdtarget),0))::int end as procent from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "' join fngoldenpointtoday(now()::date,'" + xusername + "','') g on g.strcode=m.storecode WHERE r.active='1' and r.spvcode='" + xusername + "' group by r.sellercode,r.sellername ) f order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getGDPDetailASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select sellercode as sellercodex from dist_seller where active='1' and spvcode='" + xusername + "'  "
		rowsseller, errseller := db.Query(sqlstring)
		if errseller != nil {
			panic(errseller)
		}

		var sellercodex string
		var querycompleted string
		var cunion int
		cunion = 0

		for rowsseller.Next() {
			errseller = rowsseller.Scan(&sellercodex)
			if errseller != nil {
				// handle this error
				panic(errseller)
			}

			if cunion > 0 {
				querycompleted = querycompleted + " union select * from fngoldenpointtoday(now()::date,'" + sellercodex + "','') "
			} else {
				querycompleted = querycompleted + " select * from fngoldenpointtoday(now()::date,'" + sellercodex + "','') "
			}
			cunion++
		}

		defer rowsseller.Close()

		sqlstring = " select name,target,actual,procent,sellercode from ( select  r2.sellercode,r2.sellername as name,case when coalesce(sum(g.sbdtarget),0)=0 then coalesce(sum(g.gp),0) else ((coalesce(sum(g.gp),0)*100)/coalesce(sum(g.sbdtarget),0))::int end as procent, 	coalesce(sum(g.sbdtarget),0) as target,coalesce(sum(g.gp),0) as actual 	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate  " + weekstr + " join store_master m on m.storecode=v.storecode 	join dist_seller r on r.sellercode=m.sellercode join dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  	join ( " + querycompleted + ") g on g.strcode=m.storecode WHERE r.active='1' and r2.spvcode='" + xusername + "'  group by r2.sellercode,r2.sellername) x order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getDeviationDetail", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select name,target,actual,procent,sellercode from (select  r.sellername as name,coalesce(count(v.storecode),0) as target,coalesce(count(g.storecode),0) as actual,case when (coalesce(count(v.storecode),0)+coalesce(count(g.storecode),0))=0 then 0 else ((coalesce(count(g.storecode),0)*100)/((coalesce(count(v.storecode),0)+coalesce(count(g.storecode),0))))::int   end as procent ,r.sellercode as sellercode from store_master m join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  join getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate on 1=1 left join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " and v.storecode=m.storecode left join store_deviated g on g.storecode=m.storecode and g.date=now()::date and g.sellercode=r.sellercode WHERE r.active='1' and r.spvcode='" + xusername + "'  group by r.sellercode,r.sellername  ) f order by procent desc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getDeviationDetailASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select name,target,actual,procent,sellercode from (select  r2.sellername as name,coalesce(count(v.storecode),0) as target,coalesce(count(g.storecode),0) as actual,case when (coalesce(count(v.storecode),0)+coalesce(count(g.storecode),0))=0 then 0 else ((coalesce(count(g.storecode),0)*100)/((coalesce(count(v.storecode),0)+coalesce(count(g.storecode),0))))::int   end as procent ,r2.sellercode as sellercode from store_master m join dist_seller r on r.sellercode=m.sellercode  JOIN dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  join getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate on 1=1 left join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + "  and v.storecode=m.storecode left join store_deviated g on g.storecode=m.storecode and g.date=now()::date and g.sellercode=r.sellercode WHERE r.active='1' and r2.spvcode='" + xusername + "'  group by r2.sellercode,r2.sellername  ) f order by procent desc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getGDPDetailStore", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 	select name,target,actual,procent,sellercode from (select  m.storecode as sellercode,m.name as name,sum(g.sbdtarget) as target,sum(g.gp) as actual,case when coalesce(sum(g.sbdtarget),0)=0 then sum(g.gp) else ((coalesce(sum(g.gp),0)*100)/coalesce(sum(g.sbdtarget),0))::int end as procent 	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate  " + weekstr + "  join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'  join fngoldenpointtoday(now()::date,'" + xusername + "','') g on g.strcode=m.storecode 	WHERE r.active='1' and r.spvcode='" + xusername + "' and r.sellercode='" + xsalescode + "' group by m.storecode,m.name ) f order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getCallDetailStore", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  select  name,target,actual,procent,sellercode " +
			" from (    " +
			" 	select  m.storecode as sellercode,m.name as name,case when m.channelcode like '%/NOO' then count(null) else count(v.storecode) end as target,count(h.storecode) as actual,  " +
			" 	case when count(h.storecode)>0 then 100 else 0 end  as procent " +
			" 	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate      " +
			" 	join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate   " + weekstr + " " +
			" 	join store_master m on m.storecode=v.storecode           " +
			" 	join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'     " +
			" 	left join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date     " +
			" 	WHERE r.active='1' and r.spvcode='" + xusername + "'    " +
			" 	group by m.storecode,m.name     " +
			" 	union all " +
			" 	select  m.storecode as sellercode,m.name as name,0 as target,count(h.storecode) as actual,case when count(h.storecode)>0 then 100 else 0 end  as procent   " +
			" 	from store_deviated v            " +
			" 	join store_master m on m.storecode=v.storecode     " +
			" 	join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'      " +
			" 	join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date    " +
			" 	WHERE r.active='1' and r.spvcode='" + xusername + "' and v.date=now()::date      " +
			" 	group by m.storecode,m.name  " +
			" ) f order by procent asc "

		//
		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getCallDetail", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  select name,target,actual,case when coalesce(target,0)=0 then actual else ((coalesce(actual,0)*100)/coalesce(target,0))::int end as procent,sellercode " +
			" from (      " +
			" 	select sellercode,name,sum(target) target,sum(actual) as actual       " +
			" 	from " +
			" 	( " +
			" 		select r.sellercode as sellercode,r.sellername as name,     " +
			" 		case when m.channelcode like '%/NOO' then count(null) else count(v.storecode) end as target,  " +
			" 		count(h.storecode) as actual       " +
			" 		from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate " +
			"		join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate   " + weekstr + "    " +
			" 		join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'      " +
			" 		left join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date         " +
			" 		WHERE r.active='1' and r.spvcode='" + xusername + "'    " +
			" 		group by r.sellercode,r.sellername,m.channelcode    " +
			" 	) f group by sellercode,name      " +
			" ) g order by procent asc     "

		//

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getCallDetailASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  	select name,target,actual,procent,sellercode from (select r2.sellercode as sellercode,r2.sellername as name,count(v.storecode) as target,count(h.storecode) as actual,case when coalesce(count(v.storecode),0)=0 then count(h.storecode) else ((coalesce(count(h.storecode),0)*100)/coalesce(count(v.storecode),0))::int end as procent 	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate 	join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode JOIN dist_seller r2 on r2.sellercode=r.spvcode and r2.sellercode like '" + xsalescode + "'  left join store_callhist h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date	WHERE r.active='1' and r2.spvcode='" + xusername + "' group by r2.sellercode,r2.sellername ) f order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getPCallDetail", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "   select * from ( select name,sum(target) as target,sum(actual) as actual,case when coalesce(sum(target),0)=0 then sum(actual)::int else ((sum(actual)*100)/sum(target))::int   end  as procent,sellercode from ( " +
			" select r.sellercode as sellercode,r.sellername as name, " +
			" case when m.channelcode like '%/NOO' then count(null) else count(distinct v.storecode) end as target, " +
			" count(distinct h.storecode) as actual,0 as procent " +
			" from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate  " +
			" join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate   " + weekstr + " " +
			" join store_master m on m.storecode=v.storecode  " +
			" join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'   " +
			" left join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date  " +
			" WHERE r.active='1' and r.spvcode='" + xusername + "'  " +
			" group by r.sellercode,r.sellername,m.channelcode " +
			" union all " +
			" select r.sellercode as sellercode,r.sellername as name, " +
			" 0 as target, " +
			" count(distinct h.storecode) as actual,0 as procent " +
			" from store_deviated v   " +
			"  join store_master m on m.storecode=v.storecode  " +
			" join dist_seller r on r.sellercode=m.sellercode and r.sellercode like '" + xsalescode + "'   " +
			" join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date  " +
			" WHERE r.active='1' and r.spvcode='" + xusername + "' and v.date=now()::date  " +
			" group by r.sellercode,r.sellername,m.channelcode " +
			" ) f group by name,sellercode) g order by procent asc "

		//

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getPCallDetailASM", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  select name,target,actual,procent,sellercode from (select r2.sellercode as sellercode,r2.sellername as name,count(v.storecode) as target,count(h.storecode) as actual,case when coalesce(count(v.storecode),0)=0 then count(h.storecode) else ((coalesce(count(h.storecode),0)*100)/coalesce(count(v.storecode),0))::int end as procent from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate " + weekstr + " join store_master m on m.storecode=v.storecode join dist_seller r on r.sellercode=m.sellercode JOIN dist_seller r2 on r2.sellercode=r.spvcode  and r2.sellercode like '" + xsalescode + "'  left join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date WHERE r.active='1' and r2.spvcode='" + xusername + "' group by r2.sellercode,r2.sellername ) f order by procent asc "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getPCallDetailStore", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " 		 select name,target,actual,procent,sellercode from " +
			" ( " +
			" 	select  m.storecode as sellercode,m.name as name,case when m.channelcode like '%/NOO' then count(null) else count(v.storecode) end as target,count(h.storecode) as actual,case when count(h.storecode)>0 then 100 else 0 end as procent  " +
			" 	from getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate  " +
			" 	join store_visit v on (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=selldate   " + weekstr +
			" 	join store_master m on m.storecode=v.storecode  " +
			" 	join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  " +
			" 	left join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date        " +
			" 	WHERE r.active='1' and r.spvcode='" + xusername + "'  " +
			" 	group by m.storecode,m.name  " +
			" 	union all " +
			" 	select  m.storecode as sellercode,m.name as name,count(v.storecode) as target,count(h.storecode) as actual,case when count(h.storecode)>0 then 100 else 0 end as procent  " +
			" 	from store_deviated v  " +
			" 	join store_master m on m.storecode=v.storecode  " +
			" 	join dist_seller r on r.sellercode=m.sellercode  and r.sellercode like '" + xsalescode + "'  " +
			" 	left join order_master h on h.storecode=m.storecode and h.sellercode=m.sellercode and h.date=now()::date        " +
			" 	WHERE r.active='1' and r.spvcode='" + xusername + "' and v.date=now()::date " +
			" 	group by m.storecode,m.name  " +
			" ) f order by procent asc		"

		//

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var sellercode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &sellercode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: sellercode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getIFFDetailSQL", func(c *gin.Context) {
		//xusername := c.PostForm("spvcode")
		//xsalescode := c.PostForm("sellercode")
		xsql := c.PostForm("sql")

		xuseragent := c.Request.Header.Get("User-Agent")

		if xuseragent == "5uPErV1sIon_8CP_m0biL3" {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			var sqlstring string

			if strings.Count(strings.ToLower(xsql), "delete") <= 0 && strings.Count(strings.ToLower(xsql), "update") <= 0 && strings.Count(strings.ToLower(xsql), "insert") <= 0 {
				sqlstring = xsql
			}

			rows, err := db.Query(sqlstring)
			if err != nil {
				panic(err)
			}

			defer rows.Close()

			var sellercode string
			var sellername string
			var iff0 string
			var iff1 string
			var iff2 string
			var iff3 string
			var iff4 string
			var iff5 string
			var iff6 string
			var results []dataIFF

			for rows.Next() {
				err = rows.Scan(&sellercode, &sellername, &iff0, &iff1, &iff2, &iff3, &iff4, &iff5, &iff6)
				if err != nil {
					// handle this error
					panic(err)
				}
				result := dataIFF{
					Sellercode: sellercode,
					Sellername: sellername,
					Iff0:       iff0,
					Iff1:       iff1,
					Iff2:       iff2,
					Iff3:       iff3,
					Iff4:       iff4,
					Iff5:       iff5,
					Iff6:       iff6,
				}
				results = append(results, result)

			}

			defer db.Close()

			c.JSON(http.StatusOK, results)
		}
	})

	r.POST("/getVolumeDetailStore", func(c *gin.Context) {
		xsalescode := c.PostForm("sellercode")
		dbname = sellerDivision(xsalescode)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  select * from " +
			" (SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual,cast(case when sum(target)<=0 then sum(actual/1.1)*100/1  else sum(actual/1.1)*100/sum(target) end as numeric(16,0)) as procent,storecode " +
			" from (	 " +
			" SELECT s.storecode,s.name as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual 	 " +
			" FROM dist_seller r  " +
			" JOIN store_master s on s.sellercode=r.sellercode  " +
			" JOIN store_visit v on v.storecode=s.storecode  " +
			" JOIN ( ( select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1  " +
			" LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM')  " +
			" LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date  " +
			" WHERE r.active='1' and r.sellercode='" + xsalescode + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + " " +
			" union all " +
			" SELECT s.storecode,s.name as name,0 as target,coalesce(total,0) as actual 	 " +
			" FROM dist_seller r  " +
			" JOIN store_master s on s.sellercode=r.sellercode  " +
			" JOIN store_deviated v on v.storecode=s.storecode and v.date=now()::date " +
			" JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and i.date=now()::date  " +
			" WHERE r.active='1' and r.sellercode='" + xsalescode + "'  " +
			" ) d group by name,storecode ) f order by procent asc "

		//

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var storecode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &storecode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: storecode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getVolumeDetailStoreSegmented", func(c *gin.Context) {
		xsalescode := c.PostForm("sellercode")
		xusername := c.PostForm("spvcode")
		xsegment := c.PostForm("segment")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = "  select * from (	SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual,cast(case when sum(target)<=0 then sum(actual/1.1)*100/1  else sum(actual/1.1)*100/sum(target) end as numeric(16,0)) as procent,storecode 	from (			SELECT s.storecode,s.name as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual 			FROM dist_seller r 		JOIN store_master s on s.sellercode=r.sellercode 		JOIN store_visit v on v.storecode=s.storecode 		join store_channel c on c.channelcode=s.channelcode and c.channeldesc='" + xsegment + "'		JOIN ( ( select getFJP(to_char(activeselldate_noclose('" + xsalescode + "'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 		LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') 		LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date 		WHERE r.active='1' and spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + " ) d group by name,storecode ) f order by procent asc"

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var storecode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &storecode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: storecode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	r.POST("/getWeekNo", func(c *gin.Context) {

		var weekstr string
		weekstr = " 1=1 "
		xusername := c.PostForm("spvcode")
		if len(xusername) > 5 {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			var sqlstring string

			sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

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
				weekstr = " and wk1='1' "
			} else if week == "2" {
				weekstr = " and wk2='1' "
			} else if week == "3" {
				weekstr = " and wk3='1' "
			} else if week == "4" {
				weekstr = " and wk4='1' "
			}

			defer db.Close()
		}

		c.String(http.StatusOK, weekstr)
	})

	r.POST("/getVolumeDetailSegmentStore", func(c *gin.Context) {
		xchannelname := c.PostForm("channelname")
		xusername := c.PostForm("spvcode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = "  	SELECT week FROM public.calendar where date=now()::date		"

		rowsweek, errweek := db.Query(sqlstring)
		if errweek != nil {
			panic(errweek)
		}

		defer rowsweek.Close()

		var week string
		var weekstr string

		for rowsweek.Next() {
			errweek = rowsweek.Scan(&week)
			if errweek != nil {
				// handle this error
				panic(errweek)
			}
		}

		weekstr = " "
		if week == "1" {
			weekstr = " and wk1='1' "
		} else if week == "2" {
			weekstr = " and wk2='1' "
		} else if week == "3" {
			weekstr = " and wk3='1' "
		} else if week == "4" {
			weekstr = " and wk4='1' "
		}

		sqlstring = " select * from (SELECT name,sum(target) as target,sum(coalesce(actual,0)/1.1) as actual,cast(case when sum(target)<=0 then sum(actual/1.1)*100/1  else sum(actual/1.1)*100/sum(target) end as numeric(16,0)) as procent,storecode from (SELECT cast(s.name as character varying) as name,CASE WHEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2)<0 THEN 0 WHEN COALESCE(targetcall_r,0)=0 THEN round((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)),2) ELSE round(((COALESCE(salestarget,0))-(COALESCE(salesvol/1.1,0)))/COALESCE(targetcall_r,0),2) END as target,coalesce(total,0) as actual,s.storecode FROM dist_seller r JOIN store_master s on s.sellercode=r.sellercode JOIN store_visit v on v.storecode=s.storecode join store_channel c on c.channelcode=s.channelcode and c.channeldesc='" + xchannelname + "' JOIN ( (select getFJP(to_char(activeselldate_noclose('01/01/0848'),'YYYY-MM-DD')) as selldate ) ) a on 1=1 LEFT JOIN store_objective o on o.storecode=s.storecode and cast(o.periodid as character varying)=to_char(now()::date,'YYYYMM') LEFT JOIN order_master i on i.storecode=s.storecode and i.sellercode=r.sellercode and date=now()::date WHERE r.active='1' and r.spvcode='" + xusername + "' and (v.sun||CAST(v.mon AS character varying)||v.tue||v.wed||v.thu||v.fri||v.sat)=a.selldate " + weekstr + " ) d group by name,storecode ) a order by procent  "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var name string
		var target string
		var actual string
		var procent string
		var storecode string
		var results []dataVolumeDetail

		for rows.Next() {
			err = rows.Scan(&name, &target, &actual, &procent, &storecode)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataVolumeDetail{
				Name:       name,
				Target:     target,
				Actual:     actual,
				Procent:    procent,
				Sellercode: storecode,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

	})

	// End generate

	r.POST("/getWorkDay", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		if len(xusername) > 5 {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			var sqlstring string

			sqlstring = "  	select coalesce(sum(case when date<=current_date then workday else 0 end),0) as currworkday,coalesce(sum(workday),0) as workday from calendar where date between date_trunc('month', current_date)::date and (date_trunc('MONTH', (current_date)::date) + INTERVAL '1 MONTH - 1 day')::DATE;	"

			rowsworkday, errworkday := db.Query(sqlstring)
			if errworkday != nil {
				panic(errworkday)
			}

			var curworkday string
			var workday string
			var dataWorkDayInits []dataWorkDay

			for rowsworkday.Next() {
				errworkday = rowsworkday.Scan(&curworkday, &workday)
				if errworkday != nil {
					// handle this error
					panic(errworkday)
				}
			}

			dataWorkDayInit := dataWorkDay{
				Currworkday: curworkday,
				Workday:     workday,
			}

			dataWorkDayInits = append(dataWorkDayInits, dataWorkDayInit)

			colWorkDayInit := colWorkDay{
				Message:     "OK",
				DataWorkDay: dataWorkDayInits,
				Status:      "1",
			}

			defer rowsworkday.Close()
			defer db.Close()

			c.JSON(http.StatusOK, colWorkDayInit)
		} else {
			var dataWorkDayInits []dataWorkDay

			dataWorkDayInit := dataWorkDay{
				Currworkday: "0",
				Workday:     "1",
			}

			dataWorkDayInits = append(dataWorkDayInits, dataWorkDayInit)

			colWorkDayInit := colWorkDay{
				Message:     "Your Request Invalid",
				DataWorkDay: dataWorkDayInits,
				Status:      "0",
			}
			c.JSON(http.StatusOK, colWorkDayInit)
		}
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

			rows, err := db.Query(" select s.id as spvcode,s.name  as spvname,s.username as loginid,s.password,s.branch_id ,b.remark as branchname,'1' as versionupdate,'1' as forceupdate,'1' as weekno,2 as bearer from sales s join branch b on b.id=s.branch_id where s.username = $1 and s.password=$2 and s.active='1' ", xusername, xpassword)
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
	
	// Begin Logout
	r.POST("/logout", func(c *gin.Context) {

		xusername := c.PostForm("sellercode")
		xsession := c.PostForm("session")
		dbname = sellerDivision(xusername)

		xuseragent := c.Request.Header.Get("User-Agent")

		if xuseragent == "5uPErV1sIon_8CP_m0biL3" {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			if xsession != "" {
				rowsinsses, errsses := db.Query("INSERT INTO public.login_session(session, sellercode, description) VALUES ($1, $2, 'LOGOUT'); ", xsession, xusername)
				if errsses != nil {
					// handle this error better than this
					fmt.Print(errsses)
				}

				defer rowsinsses.Close()
			}
			
			defer db.Close()
		}

		c.String(http.StatusOK, "1")
	})
	// End Logout

	// Begin Tracing
	r.POST("/tracing", func(c *gin.Context) {

		xusername := c.PostForm("sellercode")
		xsession := c.PostForm("session")
		xdesc := c.PostForm("description")
		dbname = sellerDivision(xusername)

		xuseragent := c.Request.Header.Get("User-Agent")

		if xuseragent == "5uPErV1sIon_8CP_m0biL3" {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			if xsession != "" {
				rowsinsses, errsses := db.Query("INSERT INTO public.login_session(session, sellercode, description) VALUES ($1, $2, $3); ", xsession, xusername, xdesc)
				if errsses != nil {
					// handle this error better than this
					fmt.Print(errsses)
				}

				defer rowsinsses.Close()

				rowsxses, errsxses := db.Query("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname!='postgres' AND pid <> pg_backend_pid() AND state = 'idle' AND state_change < current_timestamp - INTERVAL '240' MINUTE; ")
				if errsxses != nil {
					// handle this error better than this
					fmt.Print(errsxses)
				}

				defer rowsxses.Close()
			}
			defer db.Close()
		}

		c.String(http.StatusOK, "1")
	})
	// End Tracing

	// Begin MTD Detail

	type dataMTDDetail struct {
		Periodid            string `json:"periodid"`
		Sellercode          string `json:"sellercode"`
		Sellername          string `json:"sellername"`
		Targetsales         string `json:"targetsales"`
		Actsales            string `json:"actsales"`
		Porcent_rankspv_mtd string `json:"porcent_rankspv_mtd"`
	}

	type colMTDDetail struct {
		Message       string          `json:"message"`
		DataMTDDetail []dataMTDDetail `json:"data"`
		Status        string          `json:"status"`
	}

	r.POST("/getMTDDetailSPV", func(c *gin.Context) {
		xusername := c.PostForm("username")
		xsellercode := c.PostForm("sellercode")
		xtype := c.PostForm("cmdtype")
		//xuseragent := "5uPErV1sIon_8CP_m0biL3"
		xuseragent := c.Request.Header.Get("User-Agent")
		dbname = sellerDivision(xusername)

		var result []dataMTDDetail
		var sqlquery string

		if xuseragent == "5uPErV1sIon_8CP_m0biL3" {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			switch xtype {
			case "call":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales,case when targetsales<=0 then ( actsales * 100 ) else (( actsales * 100 ) / targetsales) end AS porcent_rankspv_mtd FROM ( SELECT  o.periodid,  r.sellercode,  r.sellername,  SUM ( o.targetcall ) targetsales,  SUM ( o.callmtd ) AS actsales    FROM  store_objective o  JOIN store_master M ON M.storecode = o.storecode   AND M.channelcode NOT LIKE '%/NOO'  JOIN dist_seller r ON r.sellercode = M.sellercode   AND r.sellercode LIKE '" + xsellercode + "'  JOIN dist_spv v ON v.spvcode = r.spvcode   AND v.spvcode = $1    WHERE  o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' )    GROUP BY  o.periodid,   r.sellercode  ) b ORDER BY Porcent_rankspv_mtd DESC "
			case "pcall":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales,case when targetsales<=0 then ( actsales * 100 ) else (( actsales * 100 ) / targetsales) end AS porcent_rankspv_mtd FROM ( SELECT	o.periodid,r.sellercode,r.sellername,SUM ( o.targetcall ) targetsales,SUM ( o.prodcallmtd ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode AND M.channelcode NOT LIKE '%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode AND r.sellercode LIKE '" + xsellercode + "'  JOIN dist_spv v ON v.spvcode = r.spvcode AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid,r.sellercode ) b ORDER BY Porcent_rankspv_mtd DESC "
			case "volume":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales,case when targetsales<=0 then ( actsales * 100 ) else (( actsales * 100 ) / targetsales) end AS porcent_rankspv_mtd FROM ( SELECT	o.periodid,r.sellercode,r.sellername,SUM ( o.salestarget ) targetsales,SUM( o.salesvol/1.1 ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode AND M.channelcode NOT LIKE '%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode AND r.sellercode LIKE '" + xsellercode + "'  JOIN dist_spv v ON v.spvcode = r.spvcode AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid,r.sellercode ) b ORDER BY Porcent_rankspv_mtd DESC "
			case "gdp":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales,case when targetsales<=0 then ( actsales * 100 ) else (( actsales * 100 ) / targetsales) end AS porcent_rankspv_mtd FROM ( SELECT	o.periodid,r.sellercode,r.sellername,SUM ( o.sbdtarget ) targetsales,SUM ( o.goldenpoint ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode AND M.channelcode NOT LIKE '%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode AND r.sellercode LIKE '" + xsellercode + "'  JOIN dist_spv v ON v.spvcode = r.spvcode AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid,r.sellercode ) b ORDER BY Porcent_rankspv_mtd DESC "
			default:
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales,case when targetsales<=0 then ( actsales * 100 ) else (( actsales * 100 ) / targetsales) end AS porcent_rankspv_mtd FROM ( SELECT	o.periodid,r.sellercode,r.sellername,SUM ( o.salestarget ) targetsales,SUM ( o.salesvol ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode AND M.channelcode NOT LIKE '%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode AND r.sellercode LIKE '" + xsellercode + "'  JOIN dist_spv v ON v.spvcode = r.spvcode AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid,r.sellercode ) b ORDER BY Porcent_rankspv_mtd DESC "
			}

			rows, err := db.Query(sqlquery, xusername)
			if err != nil {
				// handle this error better than this
				panic(err)
			}

			defer rows.Close()
			var periodid string
			var sellercode string
			var sellername string
			var targetsales string
			var actsales string
			var porcent_rankspv_mtd string

			for rows.Next() {
				err = rows.Scan(&periodid, &sellercode, &sellername, &targetsales, &actsales, &porcent_rankspv_mtd)
				if err != nil {
					// handle this error
					panic(err)

					dataMTDDetail := dataMTDDetail{
						Periodid:            "",
						Sellercode:          "",
						Sellername:          "",
						Targetsales:         "",
						Actsales:            "",
						Porcent_rankspv_mtd: "",
					}

					// Parse JSON
					result = append(result, dataMTDDetail)

					colMTDDetailInit := colMTDDetail{
						Message:       "Failed Server Fetch Data",
						DataMTDDetail: result,
						Status:        "0",
					}

					c.JSON(http.StatusOK, colMTDDetailInit)

				}

				dataMTDDetail := dataMTDDetail{
					Periodid:            periodid,
					Sellercode:          sellercode,
					Sellername:          sellername,
					Targetsales:         targetsales,
					Actsales:            actsales,
					Porcent_rankspv_mtd: porcent_rankspv_mtd,
				}

				// Parse JSON
				result = append(result, dataMTDDetail)
			}

			colMTDDetailInit := colMTDDetail{
				Message:       "OK",
				DataMTDDetail: result,
				Status:        "1",
			}

			c.JSON(http.StatusOK, colMTDDetailInit)
			defer db.Close()
		} else {

			dataMTDDetail := dataMTDDetail{
				Periodid:            "",
				Sellercode:          "",
				Sellername:          "",
				Targetsales:         "",
				Actsales:            "",
				Porcent_rankspv_mtd: "",
			}

			// Parse JSON
			result = append(result, dataMTDDetail)

			colMTDDetailInit := colMTDDetail{
				Message:       "Failed Request Parameter",
				DataMTDDetail: result,
				Status:        "0",
			}

			c.JSON(http.StatusOK, colMTDDetailInit)

		}
	})

	r.POST("/getMTDDetailASM", func(c *gin.Context) {
		xusername := c.PostForm("username")
		xsellercode := c.PostForm("sellercode")
		xtype := c.PostForm("cmdtype")
		//xuseragent := "5uPErV1sIon_8CP_m0biL3"
		xuseragent := c.Request.Header.Get("User-Agent")
		dbname = sellerDivision(xusername)

		var result []dataMTDDetail
		var sqlquery string

		if xuseragent == "5uPErV1sIon_8CP_m0biL3" {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			switch xtype {
			case "call":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales, ( actsales * 100 ) / targetsales AS porcent_rankspv_mtd FROM ( SELECT  o.periodid,  r2.sellercode,  r2.sellername,  SUM ( o.targetcall ) targetsales,  SUM ( o.callmtd ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode  AND M.channelcode NOT LIKE'%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode JOIN dist_seller r2 ON r2.sellercode = r.spvcode  AND r2.sellercode LIKE '" + xsellercode + "' JOIN dist_spv v ON v.spvcode = r2.spvcode  AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid, r2.sellercode  ) b ORDER BY porcent_rankspv_mtd DESC "
			case "pcall":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales, ( actsales * 100 ) / targetsales AS porcent_rankspv_mtd FROM ( SELECT  o.periodid,  r2.sellercode,  r2.sellername,  SUM ( o.targetcall ) targetsales,  SUM ( o.prodcallmtd ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode  AND M.channelcode NOT LIKE'%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode JOIN dist_seller r2 ON r2.sellercode = r.spvcode  AND r2.sellercode LIKE '" + xsellercode + "' JOIN dist_spv v ON v.spvcode = r2.spvcode  AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid, r2.sellercode  ) b ORDER BY porcent_rankspv_mtd DESC "
			case "volume":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales, ( actsales * 100 ) / targetsales AS porcent_rankspv_mtd FROM ( SELECT  o.periodid,  r2.sellercode,  r2.sellername,  SUM ( o.salestarget ) targetsales,  SUM ( o.salesvol/1.1 ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode  AND M.channelcode NOT LIKE'%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode JOIN dist_seller r2 ON r2.sellercode = r.spvcode  AND r2.sellercode LIKE '" + xsellercode + "' JOIN dist_spv v ON v.spvcode = r2.spvcode  AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid, r2.sellercode  ) b ORDER BY porcent_rankspv_mtd DESC "
			case "gdp":
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales, ( actsales * 100 ) / targetsales AS porcent_rankspv_mtd FROM ( SELECT  o.periodid,  r2.sellercode,  r2.sellername,  SUM ( o.sbdtarget ) targetsales,  SUM ( o.goldenpoint ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode  AND M.channelcode NOT LIKE'%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode JOIN dist_seller r2 ON r2.sellercode = r.spvcode  AND r2.sellercode LIKE '" + xsellercode + "' JOIN dist_spv v ON v.spvcode = r2.spvcode  AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid, r2.sellercode  ) b ORDER BY porcent_rankspv_mtd DESC "
			default:
				sqlquery = " SELECT periodid,sellercode,sellername,targetsales,actsales, ( actsales * 100 ) / targetsales AS porcent_rankspv_mtd FROM ( SELECT  o.periodid,  r2.sellercode,  r2.sellername,  SUM ( o.salestarget ) targetsales,  SUM ( o.salesvol ) AS actsales FROM store_objective o JOIN store_master M ON M.storecode = o.storecode  AND M.channelcode NOT LIKE'%/NOO' JOIN dist_seller r ON r.sellercode = M.sellercode JOIN dist_seller r2 ON r2.sellercode = r.spvcode  AND r2.sellercode LIKE '" + xsellercode + "' JOIN dist_spv v ON v.spvcode = r2.spvcode  AND v.spvcode = $1 WHERE o.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) GROUP BY o.periodid, r2.sellercode  ) b ORDER BY porcent_rankspv_mtd DESC "
			}

			rows, err := db.Query(sqlquery, xusername)
			if err != nil {
				// handle this error better than this
				panic(err)
			}

			defer rows.Close()
			var periodid string
			var sellercode string
			var sellername string
			var targetsales string
			var actsales string
			var porcent_rankspv_mtd string

			for rows.Next() {
				err = rows.Scan(&periodid, &sellercode, &sellername, &targetsales, &actsales, &porcent_rankspv_mtd)
				if err != nil {
					// handle this error
					panic(err)

					dataMTDDetail := dataMTDDetail{
						Periodid:            "",
						Sellercode:          "",
						Sellername:          "",
						Targetsales:         "",
						Actsales:            "",
						Porcent_rankspv_mtd: "",
					}

					// Parse JSON
					result = append(result, dataMTDDetail)

					colMTDDetailInit := colMTDDetail{
						Message:       "Failed Server Fetch Data",
						DataMTDDetail: result,
						Status:        "0",
					}

					c.JSON(http.StatusOK, colMTDDetailInit)

				}

				dataMTDDetail := dataMTDDetail{
					Periodid:            periodid,
					Sellercode:          sellercode,
					Sellername:          sellername,
					Targetsales:         targetsales,
					Actsales:            actsales,
					Porcent_rankspv_mtd: porcent_rankspv_mtd,
				}

				// Parse JSON
				result = append(result, dataMTDDetail)
			}

			colMTDDetailInit := colMTDDetail{
				Message:       "OK",
				DataMTDDetail: result,
				Status:        "1",
			}

			c.JSON(http.StatusOK, colMTDDetailInit)
			defer db.Close()
		} else {

			dataMTDDetail := dataMTDDetail{
				Periodid:            "",
				Sellercode:          "",
				Sellername:          "",
				Targetsales:         "",
				Actsales:            "",
				Porcent_rankspv_mtd: "",
			}

			// Parse JSON
			result = append(result, dataMTDDetail)

			colMTDDetailInit := colMTDDetail{
				Message:       "Failed Request Parameter",
				DataMTDDetail: result,
				Status:        "0",
			}

			c.JSON(http.StatusOK, colMTDDetailInit)

		}
	})

	r.POST("/getIFFASMMTDDetail", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xsalescode := c.PostForm("sellercode")
		xiffname := c.PostForm("iffname")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " 	select sellercode from dist_seller where active='1' and spvcode='" + xusername + "' and sellercode like '" + xsalescode + "'  "
		rowsseller, errseller := db.Query(sqlstring)
		if errseller != nil {
			panic(errseller)
		}

		var sellercodex string
		var querycompleted string
		var cunion int
		cunion = 0

		for rowsseller.Next() {
			errseller = rowsseller.Scan(&sellercodex)
			if errseller != nil {
				// handle this error
				panic(errseller)
			}

			if cunion > 0 {
				querycompleted = querycompleted + " union select * from fninitiativestatus2(to_char(now()::date,'YYYYMM')::int,'" + sellercodex + "','') "
			} else {
				querycompleted = querycompleted + " select * from fninitiativestatus2(to_char(now()::date,'YYYYMM')::int,'" + sellercodex + "','') "
			}
			cunion++
		}

		defer rowsseller.Close()

		sqlstring = "  SELECT salesname as sellername,salescode as sellercode,s.NAME as iffname,s.target,COALESCE ( T.storecode, 0 ) AS actual " +
			" FROM " +
			" 	( " +
			" 	SELECT DISTINCT abc.salesname,abc.salescode,( abc.NAME ) AS NAME,SUM ( abc.target ) AS target  " +
			" 	FROM " +
			" 	( " +
			" 	SELECT DISTINCT " +
			" 		( r2.sellercode ) AS salescode,r2.sellername AS salesname,i.initdesc AS NAME,COUNT ( DISTINCT M.storecode ) AS target FROM " +
			" 	store_master " +
			" 		M JOIN dist_seller r ON r.sellercode = M.sellercode " +
			" 	JOIN dist_seller r2 ON r2.sellercode = r.spvcode " +
			" 	JOIN dist_spv r3 ON r3.spvcode = r2.spvcode AND r3.spvcode LIKE '%' " +
			" 	JOIN initiative_master i ON i.periodid :: CHARACTER VARYING = to_char( now( ) :: DATE, 'YYYYMM' ) AND i.channelcode = M.channelcode  " +
			" WHERE M.channelcode NOT LIKE '%/NOO' AND M.storeflag != 'D' AND r.active = '1' AND r3.spvcode = '" + xusername + "'  " +
			" GROUP BY i.initdesc,r2.sellercode,r2.sellername,i.ID ORDER BY i.initdesc DESC  " +
			" ) abc  " +
			" GROUP BY abc.NAME,abc.salesname,abc.salescode " +
			" ) s " +
			" LEFT JOIN ( " +
			" SELECT M.initdesc AS NAME,COUNT ( DISTINCT strcode ) AS storecode FROM " +
			" ( " +
			" " + querycompleted + " " +
			" ) i JOIN initiative_master M ON M.ID = i.intid GROUP BY M.initdesc  " +
			" ) T ON T.NAME = s.NAME WHERE s.NAME = '" + xiffname + "' and s.salescode like '" + xsalescode + "' ORDER BY s.target DESC  "

		rows, err := db.Query(sqlstring)

		var sellername string
		var sellercode string
		var iffname string
		var target string
		var actual string
		var results []dataIFFDetail

		if err != nil {
			panic(err)

			result := dataIFFDetail{
				Sellercode: "",
				Sellername: "",
				Iffname:    "",
				Target:     "",
				Actual:     "",
			}
			results = append(results, result)

			colIFFInit := colIFF{
				Message:       "OK",
				DataIFFDetail: results,
				Status:        "1",
			}

			c.JSON(http.StatusOK, colIFFInit)
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&sellername, &sellercode, &iffname, &target, &actual)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataIFFDetail{
				Sellercode: sellercode,
				Sellername: sellername,
				Iffname:    iffname,
				Target:     target,
				Actual:     actual,
			}
			results = append(results, result)
		}

		defer db.Close()

		colIFFInit := colIFF{
			Message:       "OK",
			DataIFFDetail: results,
			Status:        "1",
		}

		c.JSON(http.StatusOK, colIFFInit)

	})
	
	
	r.POST("/getStoreMasterFull3", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string


		sqlstring = " select distinct  sc.channeldesc||' ('||sc.channelcode||')' as channelcode,sm.storecode,sm.name as storename,sm.address, coalesce(sm.city,'-') as city,sm.sellercode,coalesce(sm.phoneno,'') as phoneno,coalesce(sm.whatsappno,'') as whatsappno,sm.longitude,sm.latitude,coalesce(sm.ktp,'-') as netizenid,getfjpname(sv.mon::character::smallint,sv.tue::character::smallint,sv.wed::character::smallint,sv.thu::character::smallint,sv.fri::character::smallint,sv.sat::character::smallint,sv.sun::character::smallint,sv.wk1::character::smallint,sv.wk2::character::smallint,sv.wk3::character::smallint,sv.wk4::character::smallint) as fjp, case when svs.storecode is null then '0' else '1' end as isvalidated,st.target_validationspv from store_master sm join setting st on st.mainbranch = '01' join store_channel sc on sc.channelcode = sm.channelcode  join store_visit sv on sv.storecode = sm.storecode join dist_seller t on t.sellercode=sm.sellercode left join sv_validation_store svs on svs.storecode =sm.storecode and svs.spvcode = t.spvcode and svs.dated between date_trunc('month', current_date)::date and  now()::date where sm.channelcode not like '%NOO%' and sm.storeflag!='D' and (t.spvcode,'/','') = '"+ xusername +"' order by sm.name "
		rows, err := db.Query(sqlstring)

		//log.Fatal(xiffname)

		var storecode string
		var storename string
		var address string
		var city string
		var sellercode string
		var phoneno string
		var whatsappno string
		var longitude string
		var latitude string
		var netizenid string
		var fjp string
		var channelcode string
		var isvalidated string
		var targetvalidation string
		var results []storeMasterFull

		if err != nil {
			panic(err)

			result := storeMasterFull{
				Sellercode : "",
				Storecode   : "",
				Storename   : "",
				Address     : "",
				Channeldesc : "",
				Longitude   : "",
				Latitude   : "",
				City    	: "",
				Phoneno    	: "",
				Whatsappno  : "",
				Netizenid   : "",
				Fjp 		: "",
				Isvalidated : "",
				TargetValidation : "",			}
			results = append(results, result)

			colStoreMasterFull := colStoreMasterFull{
				Message:       "OK",
				DataStoreMasterFull: results,
				Status:        "1",
			}

			c.JSON(http.StatusOK, colStoreMasterFull)
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&channelcode,&storecode,&storename,&address, &city,&sellercode,&phoneno,&whatsappno,&longitude,&latitude,&netizenid,&fjp,&isvalidated,&targetvalidation)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := storeMasterFull{
				Sellercode : sellercode,
				Storecode   : storecode,
				Storename   : storename,
				Address     : address,
				Channeldesc : channelcode,
				Longitude   : longitude,
				Latitude   : latitude,
				City    	: city,
				Phoneno    	: phoneno,
				Whatsappno  : whatsappno,
				Netizenid   : netizenid,
				Fjp 		: fjp,
				Isvalidated : isvalidated,
				TargetValidation : targetvalidation,
			}
			results = append(results, result)
		}

		defer db.Close()

		colStoreMasterFull := colStoreMasterFull{
			Message:       "OK",
			DataStoreMasterFull: results,
			Status:        "1",
		}

		c.JSON(http.StatusOK, colStoreMasterFull)

	})
	
	r.POST("/getStoreMasterFull_v2", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string


		sqlstring = " select distinct  sc.channeldesc||' ('||sc.channelcode||')' as channelcode,sm.storecode,sm.name as storename,sm.address, coalesce(sm.city,'-') as city,sm.sellercode,coalesce(sm.phoneno,'') as phoneno,coalesce(sm.whatsappno,'') as whatsappno,sm.longitude,sm.latitude,coalesce(sm.ktp,'-') as netizenid,getfjpname(sv.mon::character::smallint,sv.tue::character::smallint,sv.wed::character::smallint,sv.thu::character::smallint,sv.fri::character::smallint,sv.sat::character::smallint,sv.sun::character::smallint,sv.wk1::character::smallint,sv.wk2::character::smallint,sv.wk3::character::smallint,sv.wk4::character::smallint) as fjp, case when coalesce(sm.isverified::int,0)=1::int then '2' when sm.storecode='01/01/000HFL' then '2' when svs.storecode is null then '0' else '1' end as isvalidated,st.target_validationspv from store_master sm join setting st on st.mainbranch = '01' join store_channel sc on sc.channelcode = sm.channelcode  join store_visit sv on sv.storecode = sm.storecode join dist_seller t on t.sellercode=sm.sellercode left join sv_validation_store svs on svs.storecode =sm.storecode and svs.dated between now()-interval'180 days' and  now()::date where  sm.channelcode not like '%USB%' and sm.channelcode not like '%PAR%' and sm.channelcode not like '%NOO%' and sm.storeflag!='D' and replace(t.spvcode,'/','') = '"+ xusername +"' order by sm.name "
		rows, err := db.Query(sqlstring)

		//log.Fatal(xiffname)

		var storecode string
		var storename string
		var address string
		var city string
		var sellercode string
		var phoneno string
		var whatsappno string
		var longitude string
		var latitude string
		var netizenid string
		var fjp string
		var channelcode string
		var isvalidated string
		var targetvalidation string
		var results []storeMasterFull

		if err != nil {
			panic(err)

			result := storeMasterFull{
				Sellercode : "",
				Storecode   : "",
				Storename   : "",
				Address     : "",
				Channeldesc : "",
				Longitude   : "",
				Latitude   : "",
				City    	: "",
				Phoneno    	: "",
				Whatsappno  : "",
				Netizenid   : "",
				Fjp 		: "",
				Isvalidated : "",
				TargetValidation : "",			}
			results = append(results, result)

			colStoreMasterFull := colStoreMasterFull{
				Message:       "OK",
				DataStoreMasterFull: results,
				Status:        "1",
			}

			c.JSON(http.StatusOK, colStoreMasterFull)
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&channelcode,&storecode,&storename,&address, &city,&sellercode,&phoneno,&whatsappno,&longitude,&latitude,&netizenid,&fjp,&isvalidated,&targetvalidation)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := storeMasterFull{
				Sellercode : sellercode,
				Storecode   : storecode,
				Storename   : storename,
				Address     : address,
				Channeldesc : channelcode,
				Longitude   : longitude,
				Latitude   : latitude,
				City    	: city,
				Phoneno    	: phoneno,
				Whatsappno  : whatsappno,
				Netizenid   : netizenid,
				Fjp 		: fjp,
				Isvalidated : isvalidated,
				TargetValidation : targetvalidation,
			}
			results = append(results, result)
		}

		defer db.Close()

		colStoreMasterFull := colStoreMasterFull{
			Message:       "OK",
			DataStoreMasterFull: results,
			Status:        "1",
		}

		c.JSON(http.StatusOK, colStoreMasterFull)

	})
	
	r.POST("/getStoreMasterFull_v3", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		dbname = sellerDivision(xusername)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string


		sqlstring = " select distinct  sc.channeldesc||' ('||sc.channelcode||')' as channelcode,sm.storecode,sm.name as storename,sm.address, coalesce(sm.city,'-') as city,sm.sellercode,coalesce(sm.phoneno,'') as phoneno,coalesce(sm.whatsappno,'') as whatsappno,sm.longitude,sm.latitude,coalesce(sm.ktp,'-') as netizenid,getfjpname(sv.mon::character::smallint,sv.tue::character::smallint,sv.wed::character::smallint,sv.thu::character::smallint,sv.fri::character::smallint,sv.sat::character::smallint,sv.sun::character::smallint,sv.wk1::character::smallint,sv.wk2::character::smallint,sv.wk3::character::smallint,sv.wk4::character::smallint) as fjp, case when coalesce(sm.isverified::int,0)=1::int then '2' when sm.storecode='01/01/000HFL' then '2' when svs.storecode is null then '0' else '1' end as isvalidated,st.target_validationspv,case when to_char(coalesce(svs.dated,'2022-01-01'::date),'YYYYMM')=to_char(now()::date,'YYYYMM') then 1 else 0 end as ismtd  from store_master sm join setting st on st.mainbranch = '01' join store_channel sc on sc.channelcode = sm.channelcode  join store_visit sv on sv.storecode = sm.storecode join dist_seller t on t.sellercode=sm.sellercode left join sv_validation_store svs on svs.storecode =sm.storecode and svs.dated between now()-interval'180 days' and  now()::date where  sm.channelcode not like '%USB%' and sm.channelcode not like '%PAR%' and sm.channelcode not like '%NOO%' and sm.storeflag!='D' and replace(t.spvcode,'/','') = '"+ xusername +"' order by sm.name "
		rows, err := db.Query(sqlstring)

		//log.Fatal(xiffname)

		var storecode string
		var storename string
		var address string
		var city string
		var sellercode string
		var phoneno string
		var whatsappno string
		var longitude string
		var latitude string
		var netizenid string
		var fjp string
		var channelcode string
		var isvalidated string
		var targetvalidation string
		var ismtd string
		var results []storeMasterFullv3

		if err != nil {
			panic(err)

			result := storeMasterFullv3{
				Sellercode : "",
				Storecode   : "",
				Storename   : "",
				Address     : "",
				Channeldesc : "",
				Longitude   : "",
				Latitude   : "",
				City    	: "",
				Phoneno    	: "",
				Whatsappno  : "",
				Netizenid   : "",
				Fjp 		: "",
				Isvalidated : "",
				TargetValidation : "",
				IsMtd :"",
			}
			results = append(results, result)

			colStoreMasterFullv3 := colStoreMasterFullv3{
				Message:       "OK",
				DataStoreMasterFull: results,
				Status:        "1",
			}

			c.JSON(http.StatusOK, colStoreMasterFullv3)
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&channelcode,&storecode,&storename,&address, &city,&sellercode,&phoneno,&whatsappno,&longitude,&latitude,&netizenid,&fjp,&isvalidated,&targetvalidation,&ismtd)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := storeMasterFullv3{
				Sellercode : sellercode,
				Storecode   : storecode,
				Storename   : storename,
				Address     : address,
				Channeldesc : channelcode,
				Longitude   : longitude,
				Latitude   : latitude,
				City    	: city,
				Phoneno    	: phoneno,
				Whatsappno  : whatsappno,
				Netizenid   : netizenid,
				Fjp 		: fjp,
				Isvalidated : isvalidated,
				TargetValidation : targetvalidation,
				IsMtd	: ismtd,
			}
			results = append(results, result)
		}

		defer db.Close()

		colStoreMasterFullv3 := colStoreMasterFullv3{
			Message:       "OK",
			DataStoreMasterFull: results,
			Status:        "1",
		}

		c.JSON(http.StatusOK, colStoreMasterFullv3)

	})



	// Begin Insert WhatsApp
	r.POST("/insertwhatsapp", func(c *gin.Context) {

		xid := c.PostForm("id")
		xphone := c.PostForm("phone")
		xmessage := c.PostForm("message")
		xpushName := c.PostForm("pushName")
		//xthumbProfile := c.PostForm("thumbProfile")
		//xgroupId := c.PostForm("groupId")
		//xgroupSubject := c.PostForm("groupSubject")
		xtimestamp := c.PostForm("timestamp")
		xcategory := c.PostForm("category")
		xreceiver := c.PostForm("receiver")

		if xtimestamp != "" {
			psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

			db, err := sql.Open("postgres", psqlInfo)
			if err != nil {
				log.Fatal(err)
			}

			rowsinsses, errsses := db.Query("INSERT INTO public.bot_log(phone, message, pushname, message_created, category, receiver, message_id) VALUES ($1, $2, $3, $4, $5, $6, $7); ", xphone, xmessage, xpushName, xtimestamp, xcategory, xreceiver, xid)
			if errsses != nil {
				// handle this error better than this
				fmt.Print(errsses)
			}

			defer rowsinsses.Close()
			defer db.Close()
			c.String(http.StatusOK, "1")
		}else{
			c.String(http.StatusOK, "3")
		}
	})
	// End Logout


	r.POST("/insertValidationFalse", func(c *gin.Context) {

		xspvcode	:= c.PostForm("spvcode")
		xtime_in 	:= c.PostForm("time_in")
		xtime_out 	:= c.PostForm("time_out")
		xstorecode 	:= c.PostForm("storecode")
		xstorecode_2 	:= c.PostForm("storecode")
		xdated 	:= c.PostForm("dated")
		xisvalid_exist := c.PostForm("isvalid_exist")
		xfilephoto 	:= c.PostForm("filephoto")
		xnotes 	:= c.PostForm("notes")
		xlatitude 	:= c.PostForm("latitude")
		xlongitude 	:= c.PostForm("longitude")
		dbname = sellerDivision(xspvcode)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		rowsinsses, errsses := db.Query("INSERT INTO public.sv_validation_store (storecode, spvcode, dated, isvalid_exist, filephoto, createddate, notes, time_in, time_out,longitude, latitude,hash_data) VALUES($1, $2, $3, $4, $5, now(), $6, $7, $8, $9, $10,md5(now()||$11)); ", xstorecode, xspvcode, xdated, xisvalid_exist, xfilephoto, xnotes, xtime_in,xtime_out, xlongitude, xlatitude,xstorecode_2)
		if errsses != nil {
			// handle this error better than this
			fmt.Print(errsses)
		}

		defer rowsinsses.Close()
		defer db.Close()

		c.String(http.StatusOK, "1")
	})
	
	r.POST("/insertValidationFalse_v2", func(c *gin.Context) {

		xspvcode	:= c.PostForm("spvcode")
		xtime_in 	:= c.PostForm("time_in")
		xtime_out 	:= c.PostForm("time_out")
		xstorecode 	:= c.PostForm("storecode")
		xstorecode_2 	:= c.PostForm("storecode")
		xdated 	:= c.PostForm("dated")
		xisvalid_exist := c.PostForm("isvalid_exist")
		xfilephoto 	:= c.PostForm("filephoto")
		xnotes 	:= c.PostForm("notes")
		xlatitude 	:= c.PostForm("latitude")
		xlongitude 	:= c.PostForm("longitude")
		dbname = sellerDivision(xspvcode)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		rowsinsses, errsses := db.Query("INSERT INTO public.sv_validation_store (storecode, spvcode, dated, isvalid_exist, filephoto, createddate, notes, time_in, time_out,longitude, latitude,hash_data) VALUES($1, $2, $3, $4, $5, now(), $6, $7, $8, $9, $10,md5(now()||$11)); ", xstorecode, xspvcode, xdated, xisvalid_exist, xfilephoto, xnotes, xtime_in,xtime_out, xlongitude, xlatitude,xstorecode_2)
		if errsses != nil {
			// handle this error better than this
			fmt.Print(errsses)
		}

		defer rowsinsses.Close()
		defer db.Close()

		c.String(http.StatusOK, "1")
	})
	
	r.POST("/insertValidation", func(c *gin.Context) {

		xspvcode	:= c.PostForm("spvcode")
		xtime_in 	:= c.PostForm("time_in")
		xtime_out 	:= c.PostForm("time_out")
		xstorecode 	:= c.PostForm("storecode")
		xstorecode_2 	:= c.PostForm("storecode")
		xdated 	:= c.PostForm("dated")
		xisvalid_exist := c.PostForm("isvalid_exist")
		xisvalid_name := c.PostForm("isvalid_name")
		xisvalid_address := c.PostForm("isvalid_address")
		xisvalid_geolocation := c.PostForm("isvalid_geolocation")
		xisvalid_channel := c.PostForm("isvalid_channel")
		xfilephoto 	:= c.PostForm("filephoto")
		xnotes 	:= c.PostForm("notes")
		xlatitude 	:= c.PostForm("latitude")
		xlongitude 	:= c.PostForm("longitude")
		xlatitude_mst 	:= c.PostForm("latitude_mst")
		xlongitude_mst := c.PostForm("longitude_mst")
		xrequest_name := c.PostForm("request_name")
		xrequest_address := c.PostForm("request_address")
		xrequest_channel := c.PostForm("request_channel")
		dbname = sellerDivision(xspvcode)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		rowsinsses, errsses := db.Query("INSERT INTO public.sv_validation_store (storecode, spvcode, dated, isvalid_exist, filephoto, createddate, notes, time_in, time_out,longitude, latitude,isvalid_name,isvalid_address,isvalid_geolocation, isvalid_channel, longitude_mst,latitude_mst,request_name,request_address,request_channel,hash_data) VALUES($1, $2, $3, $4, $5, now(), $6, $7, $8, $9, $10, $11, $12,$13,$14,$15,$16,$17,$18,$19,md5(now()||$20)); ", xstorecode, xspvcode, xdated, xisvalid_exist, xfilephoto, xnotes, xtime_in,xtime_out, xlongitude, xlatitude,xisvalid_name,xisvalid_address,xisvalid_geolocation, xisvalid_channel, xlongitude_mst,xlatitude_mst,xrequest_name,xrequest_address,xrequest_channel,xstorecode_2)
		if errsses != nil {
			// handle this error better than this
			fmt.Print(errsses)
		}

		defer rowsinsses.Close()
		defer db.Close()

		c.String(http.StatusOK, "1")
	})
	
	r.POST("/insertValidation_v2", func(c *gin.Context) {
		
		xspvcode	:= c.PostForm("spvcode")
		xtime_in 	:= c.PostForm("time_in")
		xtime_out 	:= c.PostForm("time_out")
		xstorecode 	:= c.PostForm("storecode")
		xstorecode_2 	:= c.PostForm("storecode")
		xdated 	:= c.PostForm("dated")
		xisvalid_exist := c.PostForm("isvalid_exist")
		xisvalid_name := c.PostForm("isvalid_name")
		xisvalid_address := c.PostForm("isvalid_address")
		xisvalid_geolocation := c.PostForm("isvalid_geolocation")
		xisvalid_channel := c.PostForm("isvalid_channel")
		xfilephoto 	:= c.PostForm("filephoto")
		xnotes 	:= c.PostForm("notes")
		xlatitude 	:= c.PostForm("latitude")
		xlongitude 	:= c.PostForm("longitude")
		xlatitude_mst 	:= c.PostForm("latitude_mst")
		xlongitude_mst := c.PostForm("longitude_mst")
		xrequest_name := c.PostForm("request_name")
		xrequest_address := c.PostForm("request_address")
		xrequest_channel := c.PostForm("request_channel")
		
		xisvalid_handphone := c.PostForm("isvalid_handphone")
		xrequest_handphone := c.PostForm("request_handphone")
		xrequest_whatsapp := c.PostForm("request_whatsapp")
		xrequest_address_province := c.PostForm("request_address_province")
		xrequest_address_city := c.PostForm("request_address_city")
		xrequest_address_district := c.PostForm("request_address_district")
		xrequest_address_village := c.PostForm("request_address_village")
		
		dbname = sellerDivision(xspvcode)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		rowsinsses, errsses := db.Query("INSERT INTO public.sv_validation_store (storecode, spvcode, dated, isvalid_exist, filephoto, createddate, notes, time_in, time_out,longitude, latitude,isvalid_name,isvalid_address,isvalid_geolocation, isvalid_channel, longitude_mst,latitude_mst,request_name,request_address,request_channel,isvalid_handphone,request_handphone,request_whatsapp,request_address_province,request_address_city,request_address_district,request_address_village,hash_data) VALUES($1, $2, $3, $4, $5, now(), $6, $7, $8, $9, $10, $11, $12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,md5(now()||$27)); ", xstorecode, xspvcode, xdated, xisvalid_exist, xfilephoto, xnotes, xtime_in,xtime_out, xlongitude, xlatitude,xisvalid_name,xisvalid_address,xisvalid_geolocation, xisvalid_channel, xlongitude_mst,xlatitude_mst,xrequest_name,xrequest_address,xrequest_channel,xisvalid_handphone,xrequest_handphone,xrequest_whatsapp,xrequest_address_province,xrequest_address_city,xrequest_address_district,xrequest_address_village,xstorecode_2)
		if errsses != nil {
			// handle this error better than this
			fmt.Print(errsses)
		}

		defer rowsinsses.Close()
		defer db.Close()

		c.String(http.StatusOK, "1")
	})
	
	r.POST("/insertValidationKTP", func(c *gin.Context) {

		xspvcode	:= c.PostForm("spvcode")
		xspvcode       = strings.ReplaceAll(xspvcode, "/", "")
		xnetizenid 	:= c.PostForm("netizenid")
		xnetizenname	:= c.PostForm("netizenname")
		xstorecode	:= c.PostForm("storecode")
		xnetizenaddress 	:= c.PostForm("netizenaddress")
		xnetizenphoto 	:= c.PostForm("photo")
		xreason 	:= c.PostForm("reason")
		
		dbname = sellerDivision(xspvcode)

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		rowsinsses, errsses := db.Query("insert into store_revise(storecode,sellercode,tglentry,netizenid,photofile,netizenname,netizenaddress,reason,ismailed,status,remark) values ($1,$2,now(),$3,$4,$5,$6,$7,'1',0,'KTP_via_SPV_Apps') ", xstorecode, xspvcode, xnetizenid, xnetizenphoto, xnetizenname, xnetizenaddress, xreason)
		if errsses != nil {
			// handle this error better than this
			fmt.Print(errsses)
		}

		defer rowsinsses.Close()
		defer db.Close()

		c.String(http.StatusOK, "1")
		
	})
	
	r.POST("/getHistoryEdit", func(c *gin.Context) {
		xusername := c.PostForm("spvcode")
		xusername = strings.ReplaceAll(xusername, "/", "")

		dbname = sellerDivision(xusername)
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Fatal(err)
		}

		var sqlstring string

		sqlstring = " select sm.name as sellercode,to_char(sr.tglentry,'YYYY-MM-DD HH24:MI') as sellername  from store_revise sr join store_master sm on sm.storecode = sr.storecode where sr.remark = 'KTP_via_SPV_Apps' and replace(sr.sellercode,'/','')='"+xusername+"' and tglentry >= now()-interval'7 days' order by sr.tglentry "

		rows, err := db.Query(sqlstring)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		var sellercode string
		var sellername string
		var results []dataSales

		for rows.Next() {
			err = rows.Scan(&sellercode, &sellername)
			if err != nil {
				// handle this error
				panic(err)
			}
			result := dataSales{
				Sellercode: sellercode,
				Sellername: sellername,
			}
			results = append(results, result)
		}

		defer db.Close()

		c.JSON(http.StatusOK, results)

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
