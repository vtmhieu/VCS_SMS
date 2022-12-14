package controllers

import (
	"net"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/vtmhieu/VCS_SMS/models"
	"gorm.io/gorm"
)

type Server_controller struct {
	DB *gorm.DB
}

func New_server_controller(DB *gorm.DB) Server_controller {
	return Server_controller{DB}
}

// create a new server
func (sc *Server_controller) CreateServer(ctx *gin.Context) {
	var payload *models.Create_server
	currentUser := ctx.MustGet("currentUser").(models.User)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	var new_server models.Server
	newServer := models.Server{
		Server_id:    payload.Server_id,
		Server_name:  payload.Server_name,
		Status:       payload.Status,
		User_id:      currentUser.User_id,
		Created_time: now,
		Last_updated: now,
		Ipv4:         payload.Ipv4,
	}

	new_server = newServer
	results := sc.DB.Create(&newServer)
	if results.Error != nil {
		if strings.Contains(results.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": results.Error.Error()})
			return
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "data": new_server})
}

// Create many servers at one time
func (sc *Server_controller) CreatemanyServer(ctx *gin.Context) {
	var payload *models.Create_many_server

	currentUser := ctx.MustGet("currentUser").(models.User)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	succesfull := 0
	failed := 0
	now := time.Now()
	var resp models.Response_API
	var Responses_success []models.Response_API
	var Responses_failed []models.Response_API
	for x, y := range payload.Create_server {

		newServer := models.Server{
			Server_id:    y.Server_id,
			Server_name:  y.Server_name,
			Status:       y.Status,
			User_id:      currentUser.User_id,
			Created_time: now,
			Last_updated: now,
			Ipv4:         y.Ipv4,
		}

		results := sc.DB.Create(&newServer)
		if results.Error != nil {
			if strings.Contains(results.Error.Error(), "duplicate key") {
				failed++
				resp.Status = "500"
				resp.Number = x + 1
				resp.Server_ID = newServer.Server_id
				Responses_failed = append(Responses_failed, resp)
				// ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "number": x + 1, "message": results.Error.Error()})
				continue
			}
			continue
		} else {
			succesfull++
			resp.Status = "200"
			resp.Number = x + 1
			resp.Server_ID = newServer.Server_id
			Responses_success = append(Responses_success, resp)
		}

		// ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "number": x + 1, "data": new_server})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": gin.H{"successful": succesfull, "failed": failed, "List of successful": Responses_success, "List of failed": Responses_failed}})

}

//update server

func (sc *Server_controller) UpdateServer(ctx *gin.Context) {
	server_id := ctx.Param("server_id")
	var payload *models.Update_server

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}
	var updatedServer models.Server

	result := sc.DB.First(&updatedServer, "server_id = ?", server_id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no server found"})
		return
	}

	now := time.Now()

	server_to_updated := models.Server{
		Server_name:  payload.Server_name,
		Status:       payload.Status,
		User_id:      updatedServer.User_id,
		Created_time: updatedServer.Created_time,
		Last_updated: now,
		Ipv4:         payload.Ipv4,
	}

	sc.DB.Model(&updatedServer).Updates(server_to_updated)

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": updatedServer})

}

//get specific server

func (sc *Server_controller) GetServer(ctx *gin.Context) {
	server_id := ctx.Param("server_id")

	var foundServer models.Server
	result := sc.DB.First(&foundServer, "server_id=?", server_id)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "failed to find server"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "server found", "data": foundServer})
}

//get all server

func (sc *Server_controller) GetAllServer(ctx *gin.Context) {
	var from = ctx.DefaultQuery("from", "1")
	var to = ctx.DefaultQuery("to", "10000")
	var sort = ctx.DefaultQuery("sort", "")
	var type_sort = ctx.DefaultQuery("type", "")

	int_from, _ := strconv.Atoi(from)
	int_to, _ := strconv.Atoi(to)
	var Servers []models.Server
	if sort != "" {
		if type_sort == "desc" {
			results := sc.DB.Offset(int_from - 1).Order(sort + " DESC").Limit(int_to - int_from + 1).Find(&Servers)
			if results.Error != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no connection"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "number of servers": len(Servers), "data": Servers})
		} else if type_sort == "asc" {
			results := sc.DB.Offset(int_from - 1).Order(sort + " ASC").Limit(int_to - int_from + 1).Find(&Servers)
			if results.Error != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no connection"})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "number of servers": len(Servers), "data": Servers})
		}

	} else {
		results := sc.DB.Offset(int_from - 1).Limit(int_to - int_from + 1).Find(&Servers)

		if results.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no connection"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "number of servers": len(Servers), "data": Servers})

	}

}

//delete server

func (sc *Server_controller) DeleteServer(ctx *gin.Context) {
	server_id := ctx.Param("server_id")

	var Server_to_delete models.Server

	result := sc.DB.Offset(0).Delete(&Server_to_delete, "server_id=?", server_id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "no server id found."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"server_id": server_id, "status": http.StatusOK, "message": "Deleted"})
}

//delete all servers

func (sc *Server_controller) Delete_all_servers(ctx *gin.Context) {
	var Servers []models.Server
	sc.DB.Offset(0).Find(&Servers)

	results := sc.DB.Offset(0).Delete(&Servers)
	if results.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": results.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "results": "all data has been deleted successfully"})
}

func (sc *Server_controller) Check(ctx *gin.Context) {
	server_ipv4 := ctx.Param("ipv4")

	var servers []models.Server
	//find all servers with ipv4 = ...
	sc.DB.Find(&servers, "ipv4=?", server_ipv4)
	online := 0
	offline := 0
	for _, server := range servers {
		if strings.ToLower(server.Status) == "online" {
			online++
		} else if strings.ToLower(server.Status) == "offline" {
			offline++
		} else {
			continue
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"server_ipv4": server_ipv4, "result": gin.H{"online": online, "offline": offline}})

}

//create excel

func (sc *Server_controller) Export_Excel(ctx *gin.Context) {
	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	//sort and filter
	var from = ctx.DefaultQuery("from", "1")
	var to = ctx.DefaultQuery("to", "10000")
	var sort = ctx.DefaultQuery("sort", "")
	var type_sort = ctx.DefaultQuery("type", "")

	int_from, _ := strconv.Atoi(from)
	int_to, _ := strconv.Atoi(to)
	var Servers []models.Server
	if sort != "" {
		if type_sort == "desc" {
			results := sc.DB.Offset(int_from - 1).Order(sort + " DESC").Limit(int_to - int_from + 1).Find(&Servers)
			if results.Error != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no connection"})
				return
			}
		} else if type_sort == "asc" {
			results := sc.DB.Offset(int_from - 1).Order(sort + " ASC").Limit(int_to - int_from + 1).Find(&Servers)
			if results.Error != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no connection"})
				return
			}
		}
	} else {
		results := sc.DB.Offset(int_from - 1).Limit(int_to - int_from + 1).Find(&Servers)

		if results.Error != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "no connection"})
			return
		}
	}
	f.SetCellValue("Sheet1", "A1", "Server ID")
	f.SetCellValue("Sheet1", "B1", "Server Name")
	f.SetCellValue("Sheet1", "C1", "Status")
	f.SetCellValue("Sheet1", "D1", "Created Time")
	f.SetCellValue("Sheet1", "E1", "Updated Time")
	f.SetCellValue("Sheet1", "F1", "Ipv4")

	for i, c := range Servers {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), c.Server_id)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), c.Server_name)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), c.Status)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), c.Created_time)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), c.Last_updated)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), c.Ipv4)
		// Set active sheet of the workbook.
	}
	f.SetActiveSheet(index)
	// Save xlsx file by the given path.
	if err := f.SaveAs("New Server.xlsx"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "file has been created successfully"})

	// file, _ := ctx.FormFile("New Server.xlsx")
	// log.Println(file.Filename)

	// // Upload the file to specific dst.
	// ctx.SaveUploadedFile(file, "./assets/upload/"+uuid.New().String()+filepath.Ext(file.Filename))

	// ctx.JSON(http.StatusOK, gin.H("'%s' uploaded!", file.Filename))
}

func (sc *Server_controller) Post_by_excel(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}
	// log.Println(file.Filename)
	var server models.Server

	server.User_id = currentUser.User_id
	f, err := excelize.OpenFile(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}
	row, err := f.Rows("Sheet1")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}
	failed := 0
	succesfull := 0
	x := 0
	now := time.Now()
	var resp models.Response_API
	var Responses_success []models.Response_API
	var Responses_failed []models.Response_API
	for row.Next() {
		for i, value := range row.Columns() {
			if i == 0 {
				server.Server_id = value
			} else if i == 1 {
				server.Server_name = value
			} else if i == 2 {
				server.Status = value
			} else if i == 4 {
				server.Created_time = now
			} else if i == 5 {
				server.Last_updated = now
			} else if i == 3 {
				server.Ipv4 = value
			} else {
				continue
			}
		}
		x++
		results := sc.DB.Create(&server)
		if results.Error != nil {
			if strings.Contains(results.Error.Error(), "duplicate key") {
				failed++
				resp.Status = "500"
				resp.Number = x
				resp.Server_ID = server.Server_id
				Responses_failed = append(Responses_failed, resp)
				// ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "number": x + 1, "message": results.Error.Error()})
				continue
			}
			continue
		} else {
			succesfull++
			resp.Status = "200"
			resp.Number = x
			resp.Server_ID = server.Server_id
			Responses_success = append(Responses_success, resp)
		}

		// ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "number": x + 1, "data": new_server})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": gin.H{"successful": succesfull, "failed": failed, "List of successful": Responses_success, "List of failed": Responses_failed}})
}

func (sc *Server_controller) Check_on_off(ctx *gin.Context) {
	// host, _ := os.Hostname()
	currentUser := ctx.MustGet("currentUser").(models.User)
	// host := "127.0.0.1"
	var servers []models.Server
	sc.DB.Offset(0).Find(&servers)
	for _, server := range servers {
		timeout := time.Second * 30
		conn, err := net.DialTimeout("ip4:1", server.Ipv4, timeout)
		if err != nil {

			now := time.Now()
			server_to_update := models.Server{
				Server_id:    server.Server_id,
				Server_name:  server.Server_name,
				Status:       "Offline",
				User_id:      currentUser.User_id,
				Created_time: server.Created_time,
				Last_updated: now,
				Ipv4:         server.Ipv4,
			}
			sc.DB.Model(&server).Updates(server_to_update)
			ctx.JSON(http.StatusInternalServerError, gin.H{"IpV4": server.Ipv4, "message": "Offline", "updated": server_to_update})

		}
		if conn != nil {
			defer conn.Close()
			now := time.Now()
			server_to_update := models.Server{
				Server_id:    server.Server_id,
				Server_name:  server.Server_name,
				Status:       "Online",
				User_id:      currentUser.User_id,
				Created_time: server.Created_time,
				Last_updated: now,
				Ipv4:         server.Ipv4,
			}
			sc.DB.Model(&server).Updates(server_to_update)
			ctx.JSON(http.StatusOK, gin.H{"IpV4": server.Ipv4, "message": "Online", "updated": server_to_update})

		}
	}
}

func (sc *Server_controller) Daily_return(ctx *gin.Context) {

	var payload *models.Daily_API
	//check JSON input
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	start := payload.Start
	end := payload.End
	time1, err := time.ParseInLocation("2006-01-02", start, time.Local)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Could not parse time"})
		return
	}

	time2, err := time.ParseInLocation("2006-01-02", end, time.Local)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "Could not parse time"})
		return
	}
	now := time.Now()
	// duration2 := time2.Sub(now)
	duration1 := time1.Sub(now)
	// duration3 := time2.Sub(time1)

	if !time1.Before(time2) {
		ctx.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "The end day must be after the start day"})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "The request has been sent successfully"})

		go sc.sendemail(ctx, duration1, time2, payload.List_Email)
	}
}

func (sc *Server_controller) sendemail(ctx *gin.Context, duration1 time.Duration, time2 time.Time, list_email []string) {
	time.Sleep(1 * duration1)
	now := time.Now()
	for duration := true; duration; duration = (now.Before(time2)) {
		var servers []models.Server
		sc.DB.Offset(0).Find(&servers)

		from := "vtmhieu111@gmail.com"
		password := "sducehbiurfbsszu"

		toEmailAddress := list_email
		for _, email := range toEmailAddress {
			to := []string{email}

			host := "smtp.gmail.com"
			port := "587"
			address := host + ":" + port

			subject := "Subject: Daily update status\n"
			body := "This is the status of server today:\n"
			var mess string
			online := 0
			offline := 0
			for _, server := range servers {
				// mess = "server id: " + server.Server_id + "\n" + "server status: " + server.Status + "\n\n"
				if strings.ToLower(server.Status) == "online" {
					online++
				} else if strings.ToLower(server.Status) == "offline" {
					offline++
				} else {
					continue
				}
			}

			mess = "\nThe number of servers is: " + strconv.Itoa(len(servers)) + "\n" + "The number of Online servers is: " + strconv.Itoa(online) + "\n" + "The number of Offline servers is: " + strconv.Itoa(offline)
			body += mess
			message := []byte(subject + body)

			auth := smtp.PlainAuth("", from, password, host)

			err := smtp.SendMail(address, auth, from, to, message)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "something went wrong"})
				return
			}
		}
		time.Sleep(30 * time.Minute)
	}
}

//dial tcp 127.0.0.1:6500:

//how to check ipv4 in network

// func main() {
// 	host, _ := os.Hostname()
// 	fmt.Println("Host", host)
// 	addrs, _ := net.LookupIP(host)
// 	for _, addr := range addrs {
// 		if ipv4 := addr.To4(); ipv4 != nil {
// 			fmt.Println("IPv4: ", ipv4)
// 		}
// 	}
//get ipv4 from list here

// 	conn, err := net.Dial("ip4:1", "192.168.1.204")
// 	if err != nil {
// 		fmt.Println("Error in net.Dial", err)
// 		return
// 	}

//check connection
// 	conn.Close()
// 	fmt.Println("Successful")
// }
