# Structured Application
When more features are added to the application and the business logic gets complicated, you can't no longer maintain one single huge `main.go` and `func main()`. It is time to structuralize the application based on [separating the concerns](https://en.wikipedia.org/wiki/Separation_of_concerns). In other words, "compartment the entire application, giving each element only one single responsibility". Or even in another expression, "decouple the parts of your applications as much as you can"(I can do this all day, such as "refactoring", "TDD", etc. but IMHO the phrase *"separation of concerns"* summarizes the concept in the most concise way).

#### Disclaimer: There is no **BEST** practice
The problem is, there are countless discussions and so-called *"best practices"* but the world seeminly haven't converged to a single idea. Everyone has one's own idea of how to tidy up application code, and I am pretty sure the readers do the same. So instead of *explaining* how to compartment an application, I would like to *show* my own way of refactoring the code. So this chapter is a kind of disclaimer. However, I would try to follow some well-known practices and keep things simple. 

## Controller, Service, and Repository
The architecture pattern to be introduced in this chapter consists of three main components - Controller, Service, and Repository. There is a name called the "Controller-Service-Repository pattern" found on the Internet, but I decided not to use this name since I could not find much reliable resources. All I can say is that it is suspectable that the names are originated from [Spring](https://docs.spring.io/spring-framework/reference/web/webmvc.html) or [Angular](https://angular.dev/guide/di/dependency-injection). Nevertheless, you could see so many similar patterns in other web frameworks, so here the names don't matter much.

To briefly summarize what those components are:
- The **Controller** layer exposes the APIs to the clients, receiving (HTTP)requests, and return (HTTP)responses. It decodes the incoming request data, delegates the data handling tasks to the service layer, and encodes the processed data from the service layer for returning the (HTTP)response.
- The **Service** layer handles the decoded data from the controller layer, according to the business logic. It communicates with the repository layer, internal microservices, or other third party APIs if necessary. Then it returns the processed data back to the controller.
- The **Repository** layer stores the data from the service layer such that it records the current business state. It is usually composed with the separate dedicated database server(s) such as Postgres. 

Now, let's design a simple CRUD application(called "Actor") following the architecture from ground-up. 

### Controller
Let's start simple. A minimal controller requires 
- a path string for the API address(usually recognized as a prefix), 
- a HTTP method for that API, and
- a set of methods, where each of the methods responds to an individual API endpoints.

Defining the controller interface would be like as follows:

```go
// controllers/interface.go

type Controller interface {
	Path() string                  // prefix path
	Handlers() []ControllerHandler // individual API handlers
}

type ControllerHandler struct {
	Path        string           // individual subpath
	Method      string           // allowed method for the subpath
	HandlerFunc http.HandlerFunc // individual http handler
}
```

Then we define a controller as a real-life example(note that we only show the update API due to the limit on the page area):

```go
// controller/actor.go

type ActorController struct {
	path     string
	handlers []ControllerHandler
}

func (c *ActorController) Path() string {
	return c.path
}

func (c *ActorController) Handlers() []ControllerHandler {
	return c.handlers
}

// NewActorController generates a new NewActorController instance.
//
// Inside the `handler` field, we specify the individual API handlers.
func NewActorController() *ActorController {
	return &ActorController{
		path: "/actors",
		handlers: []ControllerHandler{
            // NOTE: we omit the other handlers for the limit on the page area
            // please check out the complete code example in the repository
			{
				Path:   "/{id}",
				Method: http.MethodPatch,
				HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
					// here goes the service layer
					fmt.Fprintln(w, "PATCH /actors/{id}")
				},
			},
		},
	}
}
```

For now, let's register the handlers defined in `controller/actor.go` inside `main.go`.

```go
// main.go

func registerController(mux *http.ServeMux, c controller.Controller) {
	for _, handler := range c.Handlers() {
		pattern := fmt.Sprintf("%v %v", handler.Method, filepath.Join(c.Path(), handler.Path))

		mux.Handle(pattern, handler.HandlerFunc)
	}
}

func main() {
	// mux
	mux := http.NewServeMux()

	actorController := controller.NewActorController()
	registerController(mux, actorController)

	// listener
	listener, err := net.Listen("tcp", ":8080") 

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	// server
	server := &http.Server{Handler: mux} 
	defer server.Close()

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
```

Please check yourself whether the APIs work using tools such as `curl`. After checking that everything works okay, let's talk about decoding the request data before moving onto the service layer.

#### Decoding the Request Data
We prefer handing over the decoded data to the service layer, such that the service layer can purely handle the business logic rather than starting from the raw `r.Body` data. 

Let's look at the example of updating an actor's data(`PATCH /actors/{id}` path).

First, define the [DTO](https://en.wikipedia.org/wiki/Data_transfer_object) under `service/dto.go`:

```go
// service/dto.go

// ActorUpsertDto contains necessary information for creating and updating an Actor object
type ActorUpsertDto struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
```

Next, we decode the request body inside the API handler:

```go
// controller/actor.go

{
    Path:   "/{id}",
    Method: http.MethodPatch,
    HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
        // data parsing
        actorId := r.PathValue("id")
        defer r.Body.Close()

        dto := service.ActorUpdateDto{}
        if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
            http.Error(w, "json body parsing error", http.StatusBadRequest)
        }

        // here goes the service layer
        // here we hand over the dto to the service layer
        fmt.Fprintln(w, "PATCH /actors/{id}")
    },
},
```

Since `actorId` is not being used, you won't be able to test the endpoint until the service layer is ready, which will be our next topic. 

### Service

#### Dependency Injection

#### Encoding the Application Data

### Repository

## Bootstrap and Configurations

## Conclusion

## Exercise



