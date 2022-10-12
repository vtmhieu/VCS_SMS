package routes

import "github.com/vtmhieu/VCS_SMS/controllers"

type Excel_route_controller struct {
	Excelcontroller controllers.Excel_controller
}

func New_route_excel_controller(excelcontroller controllers.Excel_controller) Excel_route_controller {
	return Excel_route_controller{excelcontroller}
}

func (c *Excel_route_controller) Excel_Route() {}
