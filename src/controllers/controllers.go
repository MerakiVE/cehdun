package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	
	"github.com/merakiVE/CEHDUN/common"
	"github.com/merakiVE/CEHDUN/core/types"
	"os"
    "text/template"
)

func ConnectDB(ctx context.Context) {

	var dataDB types.DataBase

	err := ctx.ReadJSON(&dataDB)
	
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(types.ResponseAPI{
			Message: "Invalid data",
			Data:    nil,
			Errors:  nil,
		})
		return
	}
	
	response, err := common.GetTables(dataDB)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(types.ResponseAPI{
			Message: "Invalid data",
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	} 

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(types.ResponseAPI{
		Message: "Connection success",
		Data:    response,
		Errors:  nil,
	})
}

func BuildAPI(ctx context.Context) {

	var api types.Api

	err := ctx.ReadJSON(&api)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(types.ResponseAPI{
			Message: "Invalid data",
			Data:    nil,
			Errors:  nil,
		})
		return
	}

	const TEMPLATE_API = `package design

{{$api := .}}

import (
    . "github.com/goadesign/goa/design"
    . "github.com/goadesign/goa/design/apidsl"
    . "github.com/goadesign/oauth2/design"
)

var _ = API('{{$api.Name}}', func() {                     
    Title('{{$api.Title}}')           
    Description('{{$api.Description}}')
    Contact(func() {
        Name('{{$api.Contact.Name}}')
        Email('{{$api.Contact.Email}}')
    })

    Host('{{$api.Host}}':'{{$api.Port}}')
    Scheme('http')                             
	BasePath('{{$api.BasePath}}')

    Consumes('application/json')
    Produces('application/json')

    Origin("*", func() {
            Methods("GET", "POST", "PUT", "PATCH", "DELETE")
            Headers("content-type")
            MaxAge(600)
            Credentials()
    })
})
`
	const TEMPLATE_RESOURCES = ` package design

{{$resources := .}}

import (
    . "github.com/goadesign/goa/design"
    . "github.com/goadesign/goa/design/apidsl"
    . "github.com/goadesign/oauth2/design"
)

{{range $resource := $resources}}

var _ = Resource("{{$resource.Namegroup}}", func() {
        BasePath("{{$resource.Basepath}}")

        {{range $action := $resource.Actions}}

        Action("{{$action.Name}}", func() {
                Description("{{$action.Description}}")
                Routing({{$action.Method}}("{{$action.Route}}"))
                
                {{if $action.Params}}
                Params(func(){
                    {{range $param := $action.Params}}
                    Param("{{$param.Name}}}", {{$param.Type}}, "{{$param.Description}}")
                    {{end}}
                })
                {{end}}

                {{if $action.Payload.Name }}
                Payload({{$action.Payload.Name}})
                {{end}}
                
                {{range $response := $action.Responses}}
                Response({{$response.Name}})
                {{end}}
        })
        {{end}}

})

`

    t, err := template.New("Api Design").Parse(TEMPLATE_API)

	    if err != nil {
	        panic(err)
	    }

    w, err := os.Create("api.go")

        if err != nil {
            panic(err)
        }

        defer w.Close()

    err = t.Execute(w, api)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(types.ResponseAPI{
			Message: "Invalid data",
			Data:    nil,
			Errors:  []string{err.Error()},
		})
		return
	} 

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(types.ResponseAPI{
		Message: "API success",
		Data:    nil,
		Errors:  nil,
	})
}