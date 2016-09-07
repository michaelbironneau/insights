package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	//"github.com/labstack/echo/middleware"
	"io/ioutil"
	"strconv"
)

var queue = NewMemoryQueue(100)

func main() {
	e := echo.New()
	//e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	//	Root:   "reports",
	//	Browse: true,
	//}))
	e.GET("/hosts/:hostid/work-items/head", getWorkItem)
	e.POST("/hosts/:hostid/work-items/:id/success", uploadReport)
	e.POST("/hosts/:hostid/work-items/:id/failure", reportFailure)
	e.POST("/hosts/:hostid/work-items", createWorkItem)
	//temporary stuff - we'll make these better later on
	e.Static("/reports", "reports")
	//e := echo.New()

	e.Run(fasthttp.New(":8089"))

}

func createWorkItem(c echo.Context) error {
	var item WorkItem
	if err := c.Bind(&item); err != nil {
		return err
	}
	if _, err := queue.Push(&item); err != nil {
		return err
	}
	return c.NoContent(200)
}

func getWorkItem(c echo.Context) error {
	item, err := queue.Peek()
	if item == nil && err == nil {
		c.NoContent(200)
		return nil
	} else if err != nil {
		return err
	}
	return c.JSON(200, item)
}

func uploadReport(c echo.Context) error {
	content, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		return err
	}
	return ioutil.WriteFile("reports/"+c.Param("id")+".html", content, 0644)
}

func reportFailure(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	if err = queue.Fail(id); err != nil {
		return err
	}
	return c.NoContent(200)
}
