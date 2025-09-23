# Structured Application
When more features are added to the application and the business logic gets complicated, you can't no longer maintain one single huge `main.go` and `func main()`. It is time to structuralize the application based on [separating the concerns](https://en.wikipedia.org/wiki/Separation_of_concerns). In other words, "compartment the entire application, giving each element only one single responsibility". Or even in another expression, "decouple the parts of your applications as much as you can"(I can do this all day, such as "refactoring", "TDD", etc. but IMHO the phrase *"separation of concerns"* summarizes the concept in the most concise way).

#### Disclaimer: There is no **BEST** practice
The problem is, there are countless discussions and so-called *"best practices"* but the world seeminly haven't converged to a single idea. Everyone has one's own idea of how to tidy up application code, and I am pretty sure the readers do the same. 

So instead of *explaining* how to compartment an application, I would like to *show* my own way of structring it. The structure of the example application in this chapter is based on several well known backend frameworks such as [Django](https://www.djangoproject.com/) or [Spring](https://spring.io/), but I will try to simplify it without bloating the application with too much interfaces. 

## Controller, Service, and Repository
Let me introduce the three main components of our application - Controller, Service, and Repository. There are many different names that referring to these concepts, but I find these three terms the most intuitive to understand. You could find out that there are so many similar concepts in different names in the web frameworks widely used in the world.

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

Defining the controller interface would be like as follows.

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

Then we define a controller as a real-life example(note that we only show the update API due to the limit on the page area).

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

First, define the [DTO](https://en.wikipedia.org/wiki/Data_transfer_object) under `service/dto/dto.go`. The file is under the `service` package since how DTO is determined will be upto how the service layer needs it.

```go
// service/dto/dto.go

// ActorUpsertDto contains necessary information for creating and updating an Actor object
type ActorUpsertDto struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
```

Next, we decode the request body inside the API handler.

```go
// controller/actor.go

{
    Path:   "/{id}",
    Method: http.MethodPatch,
    HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
        // data parsing
        actorId := r.PathValue("id")
        defer r.Body.Close()

        dto := dto.ActorUpdateDto{}
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
The service layer contains the core business logic of the application. The below is the "ActorService" that handles recording actor's data to the given database. 

Let's define its interface, `IActorService` in `service/interface.go`, after first defining the `Actor` schema(or model) under `service/schema/actor.go`.

```go
// service/schema/actor.go

type Actor struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
```

```go
// service/interface.go

type actor = schema.Actor // this is for convenience

type IActorService interface {
	ReadAll() ([]actor, error)

	Read(id string) (actor, error)

	Create(dto dto.ActorUpsertDto) (actor, error)

	Update(id string, dto dto.ActorUpsertDto) (actor, error)

	Delete(id string) (actor, error)
}
```

Once we define the interface, the next step is to define its implementation. Note that for convinience, as usual, we will only show `Update()` method. For the other methods, please check out the complete code example. 

```go
// service/actor.go

type IActorRepository any

type ActorServiceImpl struct {
	repository IActorRepository
}

func (s *ActorServiceImpl) Update(id string, dto dto.ActorUpsertDto) (actor, error) {
	return actor{}, nil
}

func NewActorService(repository IActorRepository) *ActorServiceImpl {
	return &ActorServiceImpl{repository: repository}
}
```

For now, we mocked the `IActorRepository` and `Update()` method, for connecting the layers one by one. Now, let's remove the linter errors in `controller/actor.go` by connecting the controller layer(`ActorController`) and the service layer(`IActorService`).

#### Dependency Injection
But how to connect the two different layers? There is a widely used technique in the web framework world called ["dependency injection"](https://en.wikipedia.org/wiki/Dependency_injection), which is just - at least I think - a bit fancy name for a layer having its dependencies provisioned by the framework rather than the object(class, struct, etc.) itself creating them on its own.

It could be a little bit confusing concept when you only read the words, but the code is fairly simple. See `ActorController` again with the injected dependency, `IActorService`.

```go
// controller/actor.go

type ActorController struct {
	path     string
	handlers []ControllerHandler
}

func NewActorController(actorService service.IActorService) *ActorController {
	return &ActorController{
		path:         "/actors",
		handlers: []ControllerHandler{
		// ...
		}
	}
}
```

See that `IActorService` is *injected* to the contructor function of the controller. This implies that someone has to actually call this constructor function and connect the entities in dependency relationships somewhere.

#### Connecting Entities 
Once you inject the dependency to the controller, you'll see that `main.go` has lint errors, because we are not providing the required parameter to `NewActorController()` function:

```go
actorController := controller.NewActorController() // need ActorService
```

Let's generate an instance of `ActorServiceImpl` right above the function, and inject to the controller contructor.

```go
actorService := service.NewActorService(struct{}{})
actorController := controller.NewActorController(actorService)
```

It doesn't have to be `main()` where we combine the entities, but for now we put the code inside `main()`.

#### Encoding and Returning the Application Data
Now fill in the unfinished implementation of the controllers. But before we can finally test the API endpoints at least returning mock responses, we have to encode the returning processed data from the service layer.

Let's revisit the update API.

```go
HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
	// data parsing
	actorId := r.PathValue("id")
	defer r.Body.Close()

	dto := dto.ActorUpsertDto{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "json body parsing error", http.StatusBadRequest)
		return
	}

	// here goes the service layer
	// here we hand over the dto to the service layer
	updatedActor, err := actorService.Update(actorId, dto)

	if err != nil {
		// log the error
		// todo: need to separately define errors for invalid inputs, or if the actor doesn't exist
		// at the moment, we assume that it is 401
		http.Error(w, "the request contains invalid data", http.StatusBadRequest)
		return
	}

	// write the status header first(specific to when using encode/json's Encoder::Encode())
	w.WriteHeader(http.StatusOK)

	// return the json as response
	if err := json.NewEncoder(w).Encode(&updatedActor); err != nil {
		// json encode error
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
},
```

It is not perfect right at the moment, but the controller now delegates the tasks related to business logic to the injected `IActorService`. Also, note that
- we try to handle the error cases as much as possible, 
- in case of returning json response, `encode/json` provides a very handy way to do so, and
- we don't forget to specify writing the response status to the header. 

Now, please check whether the API endpoints return mock responses. Once the checking has finished, we move onto implementing the repository layer in a very simple way. 

### Repository
The repository layer records business data into "repositories" such as relational databases. However, preparing a separate database is a bit cumbersome, as we need to provision it using a separate process(although we will cover how to do so using [Docker](https://hub.docker.com/_/postgres) in a later chatper). So for brevity of explanation, we will implement a very simple key-value database(i.e. of type `map[string]Actor`). 

But at the same time, we will also make comments on where to start *database transactions* inside the service layer functions. Not only because it is common to use relational databases, but also a transaction defines a [unit of work](https://en.wikipedia.org/wiki/Unit_of_work)(a series of tasks that should be correctly conducted together, otherwise cancel them all). 

Implementing a simple database of type `map[string]Actor` is very straigtforward. As usual, we would like to define the interface first. And for the implementation, we want to only show the update function as an example, just for convenience. 

```go
// repository/interface.go

type actor = schema.Actor

type IActorRepository interface {
	ReadAll() ([]actor, error)

	Read(id string) (actor, error)

	Create(dto dto.ActorUpsertDto) (actor, error)

	Update(id string, dto dto.ActorUpsertDto) (actor, error)

	Delete(id string) (actor, error)
}
```

Since our application is a very simple CRUD app, the repository interface is the same as the service interface. The implementation, for sure, will be different.

```go
// repository/actor.go

var MockDB = map[string]service.Actor{
	"ken": {
		Id:   "ken",
		Name: "Ken Jeong",
		Role: "Lesley Chow",
	},
	"alan": {
		Id:   "alan",
		Name: "Zach Galifianakis",
		Role: "Alan Garner",
	},
}

type SimpleDB = map[string]service.Actor // for convenience

type ActorRepositoryImpl struct {
	db SimpleDB
}


func (s *ActorRepositoryImpl) Update(id string, dto dto.ActorUpsertDto) (actor, error) {
	// validation
	if _, ok := s.db[id]; !ok {
		return actor{}, fmt.Errorf("actor with id %v not found", id)
	}

	if dto.Name == "" || dto.Role == "" {
		return actor{}, errors.New("the name and role must be provided")
	}

	updatedActor := actor{Id: id, Name: dto.Name, Role: dto.Role}
	s.db[id] = updatedActor

	return updatedActor, nil
}

func NewActorRepository(db SimpleDB) *ActorRepositoryImpl {
	return &ActorRepositoryImpl{db: db}
}
```

Note that if we use transactions, then each method needs to have a transaction object(`sql.Tx`) as one of its arguments, such as `Update(tx *sql.Tx, id string, dto service.ActoUpsertDto)`. The transaction usually starts inside the service layer, and while the transaction is open the service layer interacts with the repository layer.

But since we don't have any transaction for our application, we want to complete implementing `ActorServiceImpl` with commenting where we should open transactions. Here is the complete version of `Update()` method.

```go
func (s *ActorServiceImpl) Update(id string, dto dto.ActorUpsertDto) (actor, error) {
	// TODO: validate the dto

	// NOTE: starts transaction here if the DB supports transactional operations
	updatedActor, err := s.actorRepository.Update(id, dto)

	// NOTE: ends transaction here if the DB supports transactional operations

	return updatedActor, err
}
```

Once we finish implementing until the service layer, we now have a running application at least with a mocked database. Please test the application on your own!

### Dependency Relations Between the Components
Now let's wrap up the entire structure of the application. You might have wondered why we have separate interfaces and implementations of the service and repository layers. 

The answer is simple: we want the entire application to rely on the interfaces, not the implementations. The interfaces - not only those starting with "`I`" prefix but also the schemas(`Actor` struct) and DTOs(`ActorUpsertDto`) are the direct translations of the business logic. Hence all the changes start from these, not from the implementations.

Note that we injected one entity to the other via its interface, not its implementation. So all the layer `struct`s are depending on the interfaces. Concrete implementations can be relatively easy to change. Even changing to a different database system requires changing the repository layer and only a small prortion of the other application parts. 

## Configurations and Refactoring
I'm afraid we're not done yet. We want to separate the entire application from the server itself. For now, we initialize the application and run it inside `main()` all at once. 

### Separate Server Runtime From the Application
The HTTP server itself is not much related to the application logic. We first separate the application initialization logic inside a separate file `app.go`.

```go
// app.go


```

### 

## Conclusion

## Exercise
We haven't completed yet - there are lots of things to be added, but to name a few:
- logging errors 
- separately defining error types(like `404NotFoundError`)
- input DTO validation
- unit and integration tests

It is up to the readers to implement any of these features, but once you decide to do this exercise, please add test code!
