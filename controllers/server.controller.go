package controllers

import (
	"net/http"
	"os/exec"
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
			ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "message": results.Error.Error()})
			return
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": new_server})
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
	for x, y := range payload.Create_server {
		var new_server models.Server
		newServer := models.Server{
			Server_id:    y.Server_id,
			Server_name:  y.Server_name,
			Status:       y.Status,
			User_id:      currentUser.User_id,
			Created_time: now,
			Last_updated: now,
			Ipv4:         y.Ipv4,
		}
		new_server = newServer
		results := sc.DB.Create(&newServer)
		if results.Error != nil {
			if strings.Contains(results.Error.Error(), "duplicate key") {
				failed++
				ctx.JSON(http.StatusConflict, gin.H{"status": "failed", "number": x + 1, "message": results.Error.Error()})
				continue
			}
			continue
		}
		succesfull++
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "number": x + 1, "data": new_server})
	}

	ctx.JSON(http.StatusOK, gin.H{"result": gin.H{"successful": succesfull, "failed": failed}})

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
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "no server found"})
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
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "failed to find server"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "server found", "data": foundServer})
}

//get all server

func (sc *Server_controller) GetAllServer(ctx *gin.Context) {
	var from = ctx.DefaultQuery("from", "1")
	var to = ctx.DefaultQuery("to", "10000")

	int_from, _ := strconv.Atoi(from)
	int_to, _ := strconv.Atoi(to)
	var Servers []models.Server

	results := sc.DB.Offset(int_from - 1).Limit(int_to - int_from + 1).Find(&Servers)

	if results.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "bad request", "message": "no connection"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "number of servers": len(Servers), "data": Servers})
}

//delete server

func (sc *Server_controller) DeleteServer(ctx *gin.Context) {
	server_id := ctx.Param("server_id")

	var Server_to_delete models.Server

	result := sc.DB.Offset(0).Delete(&Server_to_delete, "server_id=?", server_id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "bad", "message": "no server id found."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"server_id": server_id, "status": "deleted"})
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
	ctx.JSON(http.StatusOK, gin.H{"results": "all data has been deleted successfully"})
}

func (sc *Server_controller) Post_by_excel(ctx *gin.Context) {

}

// check on/off + update server
func (sc *Server_controller) Check_on_off(ctx *gin.Context) {
	var servers []models.Server
	sc.DB.Offset(0).Find(&servers)
	for _, server := range servers {
		out, _ := exec.Command("ping", server.Ipv4, "-c 5", "-i 3", "-w 10").Output()
		if strings.Contains(string(out), "Destination Host Unreachable") {
			server.Status = "Offline"
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "can't reach to server"})

		} else {
			ctx.JSON(http.StatusOK, gin.H{"message": server.Status})
		}
	}
}

// func raw_connect(host string, ports []string) {
// 	for _, port := range ports {
// 		timeout := time.Second
// 		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
// 		if err != nil {
// 			fmt.Println("Connecting error:", err)
// 		}
// 		if conn != nil {
// 			defer conn.Close()
// 			fmt.Println("Opened", net.JoinHostPort(host, port))
// 		}
// 	}
// }

// check if ipv4 has how many online and offline

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

	// Set value of a cell.
	var Servers []models.Server
	//get servers from DB
	sc.DB.Offset(0).Find(&Servers)
	for i, c := range Servers {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), c.Server_id)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), c.Server_name)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), c.Status)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), c.User_id)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), c.Created_time)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), c.Last_updated)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), c.Ipv4)
		// Set active sheet of the workbook.
	}
	f.SetActiveSheet(index)
	// Save xlsx file by the given path.
	if err := f.SaveAs("Server.xlsx"); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "file has been created successfully"})
}
