# Configure Your Server
This chapter is mainly about how to configure your server. Here the verb *"configure"* means that we would like to inject a set of parameters when the server bootstraps and starts running, such that its behavior is affected by the parameters. 

Why don't you hardcode those parameters inside your project? There are several reasons for not doing that, but to list a few:

1. You don't want some sensitive information of your resources to be exposed to the public. Suppose you accidently revealed your company's shared OpenAI API Key to the public, and then the credit of your company will be soon depleted(or your company will be charged unimaginable amount of money, depending on your company's pricing plan). 

2. You want to change the way your server is being deployed, without changing its codebase at all. 
If you change your codebase, you have to pass several processes for it to be finally released. Tests, code review, build, pre-release, deploying on cloud, QA... Yes, it is horrible. But by tweaking a few variables in your config file, your job is done immediately. Only deploying matters now and you don't really care about the fresh new release of your server.

There are many ways to do so, depending on how you deploy your server. Using Kubernetes or cloud vendor services like AWS EC2, you may do configuration using bunch of `.yaml` files or tools like [Terraform](https://developer.hashicorp.com/terraform). In this chapter, we will build our own way of injecting such config information into the server application.  

## Using Environment Variables
One of the easiest and a basic way of configuring your server(or any kind of application) is using the environment variables defined on the process in which your server is running.  

```go
if secret, ok := os.LookupEnv("secret"); ok {
    fmt.Printf("secret: %v\n", secret)
}
```
Here we use [`os.LookupEnv`](https://pkg.go.dev/os#LookupEnv) rather than [`os.GetEnv`](https://pkg.go.dev/os#Getenv) because we want to check the variable in a more systematic way. 

However, you need to make sure that your code does not reveal any of the sensitive information. For example, if you allow client requests to inject a script to be executed on your server, your important information could be hacked by malicious attackers, by letting the script to extract the environment variables that are being used in your server. 

## Using Config Files - Docker Build
Obviously there would be multiple secret parameters that are difficult to be managed all, if the scale of the project becomes huge. Therefore, what people usually do is to provision a file in a simple format such as `.json`, `.yaml`, or `.toml`, and store the parameters inside the file in a structured way. Then your server application reads the entire file and use the parameters from it. 

Of course, the files could be hard-coded in the first place, but it could be dynamically generated at the build phase of the product, such as [Docker's build](https://docs.docker.com/build/building/secrets/). Then you can add the config file inside your version control like Git and maintain them conveniently.

In this chapter, we will briefly look through how we generate a json file using `docker build` with `--secret` option, and make the server application to read and use that json file inside the generated image. In a production environment, you might use already built-in solutions such as AWS ECS's [`TaskDefinition.json`](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/example_task_definitions.html). The purpose of this chapter is to show the readers how injecting secret parameters works in a simple way.

### Overall Workflow
Here is the overall workflow of how we inject secret information at a build phase dynamically:

1. Before building a Docker image, we prepare a executable file that generates a Json file given a set of flags. In our example we will use Go's [flag](https://pkg.go.dev/flag) package.

2. Next, when we run `docker build` command, the docker daemon runtime will run the executable file prepared in the previous step, so that the json file with sensitive information will be stored in a designated path, inside the image.

3. When the server runs, it will read the json file and use the secured information whenever it needs.  

Let's start from reading a json file inside our Go code.

### Reading a Json File Inside the Application
First, let's start with a simple example, where we try to read a json file with path `/tmp/gcrab_secrets.json`, with the following schema:

```json
{
    "actor_name": "Ken Jeong",
    "ethnicity": "Korean"
}
```

According to the schema above, we define a struct `Secret` and read the data from the given file path, `/tmp/gcrab_secrets.json`:

```go
type Secrets struct {
	ActorName string `json:"actor_name"`
	Ethnicity string `json:"ethnicity"`
}

func readSecretsFromJson(path string) (Secrets, error) {
	secretBytes, err := os.ReadFile(path)

	if err != nil && !errors.Is(err, io.EOF) {
		return Secrets{}, err
	}

	secrets := Secrets{}
	json.Unmarshal(secretBytes, &secrets) // you would be familiar with this function from basics-response chapter

	return secrets, nil
}
```

Next, check whether the function works okay by running the `main` function in `server.go`:

```go
func main() {
	// from file
	secretsPath := "/tmp/gcrab_secrets.json"
	secrets, err := readSecretsFromJson(secretsPath)

	if err != nil {
		log.Fatalf("reading json data from %v", secrets)
	}

	fmt.Printf("The secrets here: actor name - %v, ethnicity - %v", secrets.ActorName, secrets.Ethnicity)
}
```

### Building an Execution File Generating Json
But how could we generate the Json file dynamically? We will build an execution file that generates the Json file where the secrets are given as arguments. This time, we will use [flag](https://pkg.go.dev/flag) package.

Under a new directory `cmd`, we will create another file `gensecret.go`:

```go
// note that you must declare the package of the file as main, otherwise the built exec file won't run 
package main 

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

// define the Secrets struct again here, for simplicity
type Secrets struct {
	ActorName string `json:"actor_name"`
	Ethnicity string `json:"ethnicity"`
}

func main() {
	// you can simply use the default flagset, but here we define a new flagset for clarity
	flagSet := flag.NewFlagSet("gensecrets", flag.ExitOnError)

	// define the flags
	actorName := flagSet.String("actorName", "", "-actorName \"Ken Jeong\"")
	ethnicity := flagSet.String("ethnicity", "", "-ethnicity \"Korean\"")

	// parse the input text 
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatalf("parsing error: -actorName and -ethnicity flags must be provided: %v", err)
	}

	// encode the information into a Json string
	secrets := Secrets{ActorName: *actorName, Ethnicity: *ethnicity}
	encoded, err := json.Marshal(secrets)

	if err != nil {
		log.Fatalf("unexpected Json encoding error: %v", err)
	}

	// write the Json string to the designated path
	if err := os.WriteFile("/tmp/gcrab_secrets.json", encoded, 0644); err != nil {
		log.Fatalf("unexpected Json file writing error: %v", err)
	}
}
```

Note that you must set the `package` as `main`, otherwise the built executable file won't run. The other code lines are straightforward and does not need much explanation. If you want to understand details, please check out the comments inside the code.

After writing the file, we build the executable file by running the build command: `go build -o ./gensecret cmd/gensecret`. Don't forget to change the mode of the file, so that users can use it. Here, we allow all the users to be able to run the file, by running `chmod a+x ./gensecret`. 

After changing the mode of the file, try the executable file with secrets that mean something fun:

```sh
./gensecret --actorName Jimmy Yang --ethnicity Chinese
```

And check out the output:

```sh
cat /tmp/gcrab_secrets.json
# {"actor_name":"Jimmy Yang","ethnicity":"Chinese"}
```

### Writing `Dockerfile`
Now we have our logic of writing a Json file dynamically given pieces of text information. we'll write a [`Dockerfile`](https://docs.docker.com/guides/golang/build-images/) and build an image, where the json file will be generated and injected dynamically. 


```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.24

WORKDIR /app

COPY go.mod ./

COPY server.go ./
COPY ./cmd/ ./cmd

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /gensecret ./cmd/gensecret.go 
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./server.go 

# run the secret generating exec
RUN chmod a+x /gensecret
RUN --mount=type=secret,id=actorName,env=ACTOR_NAME\
    --mount=type=secret,id=ethnicity,env=ETHNICITY\
    /gensecret -actorName="$ACTOR_NAME" -ethnicity="$ETHNICITY"

EXPOSE 8080

CMD ["/server"]
```

We won't go to much details about Docker and `Dockerfile` here, but if you read the `Dockerfile` code above, the flow is straightforward; you basically copy the source code of the current directory to a certain directory(`WORKDIR`) inside the container, and build the executable files. Then we run the secret generator we created earlier, using a couple of secret parameters that we will mount through the Docker command arguments.

### Building an Application Image with Secrets
We build an image using the above `Dockerfile`, with the following command:

```sh
ACTOR_NAME="Ken Jeong" ETHNICITY="Korean" docker build --tag server --secret type=env,id=actorName,env=ACTOR_NAME --secret type=env,id=ethnicity,env=ETHNICITY .
```

Here we used `docker build --secret`'s [environment variables option](https://docs.docker.com/reference/cli/docker/buildx/build/#typeenv-usage). But there are other options. Please check out the link for more details.

Now let us run our "server"(technically there is no server-related logic here, only printing the secret injected at the build stage).

```sh
docker run server
# result: The secrets here: actor name - Ken Jeong, ethnicity - Korean  
```

Try another build with different secrets(note that we use `--no-cache` flag here!):

```sh
ACTOR_NAME="Jimmy Yang" ETHNICITY="Chinese" docker build --no-cache --tag server --secret type=env,id=actorName,env=ACTOR_NAME --secret type=env,id=ethnicity,env=ETHNICITY .
```

And running the same command `docker run server` will give `The secrets here: actor name - Jimmy Yang, ethnicity - Chinese`.

## Conclusion
The flow we showed in this chapter is very rough, and in production we don't really separately generate a json file containing secret information. Very mature and sophisticated tools are out there to help you build an application with secure configurations. Nevertheless, it is important to understand how sensitive information is provided inside an application, without relying on vulnerable environment variables defined at runtime.

## Exercises
1. Instead of using a json file, try using toml or yaml. That is, you need to rewrite functions `readSecretsFromToml` or `readSecretsFromYaml`. 

2. Reference the [`docker build --secret` page](https://docs.docker.com/build/building/secrets/#secret-mounts) again, and in this time, try using files rather than environment variables.
