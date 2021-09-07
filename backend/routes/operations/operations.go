package routes

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
)

/*
Operation

/api/v1/operation
POST: Create a new operation

/api/v1/operation/(id)
GET: Retrieve the given operation
*/

// Creating a REST layer as OperationResource to have all the operation

type Operation struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OperationResource struct {
	Operations map[string]Operation `json:"operations"`
}

// Creating a webservice for operation endpoints
func (o OperationResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1/operation").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{operation-id}").To(o.findOperation))
	ws.Route(ws.POST("").To(o.addOperation))
	container.Add(ws)
}

// GET info of operations depening upon the id
func (o OperationResource) findOperation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("operation-id")
	opr := o.Operations[id]
	if len(opr.Id) == 0 {
		response.AddHeader("Content-Type", "text/plain")
		err := response.WriteErrorString(http.StatusNotFound, "Operation not found!")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := response.WriteEntity(opr)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// POST to create an operation
func (o *OperationResource) addOperation(request *restful.Request, response *restful.Response) {
	opr := new(Operation)
	err := request.ReadEntity(&opr)
	if err == nil {
		o.Operations[opr.Id] = *opr
		err := response.WriteEntity(opr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		response.AddHeader("Content-Type", "text/plain")
		err := response.WriteErrorString(http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Add function to start up the server, running against dedicated port
// Usage of CurlyRouter is done because of the efficiency while using wildcards and expressions
// Not used in the main function, just here for the unit testing
func RunRestfulCurlyRouterServer() {
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	o := OperationResource{map[string]Operation{}}
	o.Register(wsContainer)

	log.Print("The server is up, and listening to port 8090 on your host.")
	server := &http.Server{Addr: ":8090", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}