package controllers

import (
	"gorm.io/gorm"
)

type Excel_controller struct {
	DB *gorm.DB
}

// func (ec *Excel_controller) Export_Excel(ctx *gin.Context) {
// 	f := excelize.NewFile()
// 	// Create a new sheet.
// 	index := f.NewSheet("Sheet1")

// 	// Set value of a cell.
// 	var Servers []models.Server
// 	//get servers from DB
// 	ec.DB.Offset(0).Find(&Servers)
// 	for i, c := range Servers {
// 		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), c.Server_id)
// 		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), c.Server_name)
// 		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), c.Status)
// 		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), c.Created_time)
// 		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), c.Last_updated)
// 		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), c.Ipv4)
// 		// Set active sheet of the workbook.
// 	}
// 	f.SetActiveSheet(index)
// 	// Save xlsx file by the given path.
// 	if err := f.SaveAs("Server.xlsx"); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "file has been created successfully"})
// }
