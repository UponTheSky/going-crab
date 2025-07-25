# Hello, Database!
Now it is time to talk about database. I have been trying to put this topic off on purpose, because setting up a database requires a lot of knowledge, thus covering this topic early in this book may make readers feel bored. 

However, since we have covered basic(and fundamental) parts of writing a HTTP server, we now can talk deeply about using databases in this chapter. For sure, database on its own is a gigantic subject, so we won't look into pure database things, such as how to write a fast SQL query. In this chapter, we will mainly talk about connecting our server to a running database server such as [Postgres](https://www.postgresql.org/).

If you are not familiar with Postgres and you want to know more, I recommend to check out Neon's [Postgres tutorial](https://neon.com/postgresql/tutorial) to grasp the overall concepts. However, it is not necessary to understand fully about Postgres to read and understand this chapter at all. You just need to know what SQL is and how to write basic CRUD operations with SQL. 

## Running a Postgres Server

### Install and Run Docker Postgres Container
For running a Postgres server without a painful and long process of installation, we will use [Docker](https://docs.docker.com/desktop/setup/install/mac-install/). Yeah I know, it is old and heavy, but it is one of the easiest way to run a container image, which suits our need. 

1. First, install Docker following the instruction on the [official docs](https://docs.docker.com/desktop/setup/install/mac-install/).

2. After installation, start Docker application.

3. Run ["docker run"](https://docs.docker.com/reference/cli/docker/container/run/) command:
```sh
docker run --name gcrab-postgres -e POSTGRES_USER=gcrab -e POSTGRES_PASSWORD=gcrab -d -p 5432:5432 postgres:latest
```

Let's break down the command briefly:
- `docker run --name gcrab-postgres`: Run a new container of an image named `postgres:latest`(at the very end of the command), with the name of the container as `gcrab-postgres`. If there is the image already pulled on your machine, Docker won't pull a new image; otherwise, it will pull the image `postgres:latest`.

- `-e POSTGRES_USER=gcrab -e POSTGRES_PASSWORD=gcrab`: Configure the database to be run on the container with the environment variables `POSTGRES_USER`, etc.

- `-d`: Run the container on a background process. It is just for convenience, but if you don't specify it then you have to open a new terminal process.

- `-p {hostPort}:{containerPort}`: Forwards the port used inside the container to the one of the host machine. That is, if you make a socket connection to `hostPort` on your machine, it is directed to `containerPort` that the container uses at the moment. Conventionally, people usually assign the port `5432` to Postgres. 

- `postgres:latest`: The image of the container to be run.

After waiting for a short period of time, you can check whether the container is running, by running the command [`docker ps`](https://docs.docker.com/reference/cli/docker/container/ls/). 

```sh
CONTAINER ID   IMAGE             COMMAND                  CREATED          STATUS          PORTS      NAMES
76a4c823fcdc   postgres:latest   "docker-entrypoint.s…"   46 seconds ago   Up 45 seconds   5432/tcp   gcrab-postgres
```

If you read `STATUS` column, we can see that the container is running okay. 

Or you can also check through the desktop UI if you have installed Docker Desktop.

### Create Database and Table
Now we will create the following elements for storing and reading data:
1. database: What? Another database in a *database*? We already have Postgres but then why do we need another one? Just like the term *server*, many IT people don't really take notice on the ambiguoity of the word *database*. Our Postgres *database* is actually a database **system**, which includes not only data but also the server on which engines and other applications run. It is a massive ecosystem. Whereas *database* inside our Postgres database indicates a logical partition for the data tables stored in the database storage. 

2. schema: Under a database(a logical parition), a schema is a namespace that separates the tables inside the database. You can join tables from different schemas.

3. table: An entity representing a collection of rows and columns containing information. 

First, run the following ["docker exec"](https://docs.docker.com/reference/cli/docker/container/exec/) command to run `psql` inside the container. `psql` is a program that enables us to communicate with Postgres easily(not necessarily tied to Docker). 

```sh
docker exec -it gcrab-postgres psql -U gcrab -W
```

Breaking down this command would be as follows:
- `docker exec`: Runs a command in a container, just like you run a command on terminal. In this case, you run the command `psql -U gcrab -W` in the container named `gcrab-postgres`, which we have just started running.
- `-it`: Applies the options `-i` and `-t` at the same time. `-i` enables users to provide inputs to the Stdin of the container process. And `-t` puts a terminal to the container, such that once you execute `docker exec` on a terminal, you can use that terminal(if the application you run needs a terminal, such as `sh` or `psql`). In short, putting `-it` flag makes you to run a command-line application where you can put in data in real-time.

- `gcrab-postgres`: The name of the container in which the command will be executed. 
- `psql -U gcrab -W`: Runs `psql` with the role named `gcrab` and asks the program to wait for the password to be provided. When you run this command, the terminal holds on and asks you the password for the role like this: `Password: `.

Once you manage to run `psql` in the container, now it is time to create a database, schema, and a table.

First, let's create a database with name `"hangover"`: 

```sh
gcrab=# CREATE DATABASE hangover;
CREATE DATABASE # this is the result
```

Next, connect to the database `hangover`, with the password of the current role `gcrab`:

```sh
gcrab=# \c hangover;
Password: 
You are now connected to database "hangover" as user "gcrab".
```

Next, just like creating the database, create a schema:

```sh
hangover=# CREATE SCHEMA general;
CREATE SCHEMA

# check whether the schema has been created:
hangover=# \dn
       List of schemas
  Name   |       Owner       
---------+-------------------
 general | gcrab
 public  | pg_database_owner
(2 rows)

# move 
```

Finally, under the schema `general`, create a table `role` with fiels `id`, `name`, and `actor_name`:

```sh
hangover=# CREATE TABLE IF NOT EXISTS general.role (
    id INTEGER PRIMARY KEY, 
    name TEXT NOT NULL, 
    actor_name TEXT NOT NULL
); 
CREATE TABLE

# check whether the table has been created as expected:
hangover=# \d general.role
 id         | integer |           | not null | 
 name       | text    |           | not null | 
 actor_name | text    |           | not null | 
```

### A Word on Production Database
So far we have provisioned a database server using Docker. However, in a real production environment, you may have to interact with a database server either on your on-premise server or on a cloud vendor like AWS. Depending on how you provision and configure those servers, you may have to provide extra information or 3rd party SDK on the backend side. In this chapter we try to specify minimum essential information for simplicity, but please check out when you deploy your application for your own good! 

## Connecting Go Server to Postgres Server
Wow. It's been a long *prelude*. We haven't touched any Go code yet(This is why I feel infra-related jobs such as DevOps or cloud engineering is super hard, lol). From this section, we will dive into Go code that handles database connection and several basic operations. 

### Installing Go Postgres Driver
Go doesn't provide any built-in database driver, so we have to prepare our own. [pgx](https://github.com/jackc/pgx) database driver is one of the recommended drivers on the [official wiki page](https://go.dev/wiki/SQLDrivers), and it is still under active development and support. We will use this driver throughout the entire book.

Run the following command to install the driver: `go get github.com/jackc/pgx/v5`. You'll see a few lines are added to `go.mod` file, including `github.com/jackc/pgx/v5 v5.x.x`(here `x` is arbitrary, depending on which version of the driver you install).

### Making Connection to Database
Let's connect to the running Postgres database server, and handle the data within our Go code. First, you need to provide valid information about the database server.

```go
dbUrl := "postgres://gcrab:gcrab@localhost:5432/hangover"
config, err := pgx.ParseConfig(dbUrl)

if err != nil {
    log.Fatal(err)
}

// [...]
```

Here `dbUrl` is the database url containing the necessary information that needs to be provided to the database server for secure connections. About the format of the URL, [Introduction to PostgreSQL connection URIs](https://www.prisma.io/dataguide/postgresql/short-guides/connection-uris) by Prisma team gives a very good and detailed explanation. But in short, the format is `postgres://{user_name}:{password}@{host}:{port}/{database_name}`. 

If the format is not valid, then the `err` returned by `pgx.ParseConfig` is not `nil` and it looks like:

```sh
cannot parse `postgre://gcrab:xxxxxx@localhost:5432/hangover`: failed to parse as keyword/value (invalid keyword/value)
exit status 1
```

After providing a valid URL, we make a connection to the database server as follows:

```go
// [...]

connCtx := context.Background()
conn, err := pgx.ConnectConfig(connCtx, config)

if err != nil {
    log.Fatal(err)
}

// don't forget to close the connection after using it
defer func() {
    closeCtx := context.Background()

    if err := conn.Close(closeCtx); err != nil {
        log.Fatal(err)
    }
}()
```

Here `pgx.ConnectConfig` returns a non-nil error if there is a connection problem. For example, if the database server is not running, the error message would be like:

```sh
failed to connect to `user=gcrab database=hangover`:
[::1]:5432 (localhost): dial error: dial tcp [::1]:5432: connect: connection refused
127.0.0.1:5432 (localhost): dial error: dial tcp 127.0.0.1:5432: connect: connection refused
[::1]:5432 (localhost): dial error: dial tcp [::1]:5432: connect: connection refused
127.0.0.1:5432 (localhost): dial error: dial tcp 127.0.0.1:5432: connect: connection refused
exit status 1
```

Also, note that we pass the empty context here for simple explanation, but you can also use something like `context.WithTimeout`.

If you haven't met any error by now, it's time to execute SQL operations and queries.

### Inserting, Querying, Updating, and Deleting Rows
Now that our Go server has been connected to the database, it is time to play with data. Here we will cover four basic operations: inserting, querying, updating, and deleting rows.

#### Inserting

Let's create a few rows first:

```go
// [...]
rows := []Role{
    {1, "Leslie Chow", "Ken Jeong"},
    {2, "Philip Wenneck", "Bradley Cooper"},
    {3, "Stuart Price", "Ed Helms"},
    {4, "Zach Galifianakis", "Alan Garner"},
}

// make a transaction for inserting multiple rows
txCtx := context.Background()
tx, err := conn.Begin(txCtx)

if err != nil {
    log.Fatalf("opening transaction error: %v", err)
}

// don't forget to rollback in case something happens; it doesn't affect successful results
defer func() {
    rbCtx := context.Background()
    tx.Rollback(rbCtx)
}()

execCtx := context.Background()

for _, role := range rows {
    tag, err := conn.Exec(
        execCtx,
        "INSERT INTO general.role(id, name, actor_name) VALUES ($1, $2, $3)",
        role.Id, role.Name, role.ActorName, // we pass the parameters here for avoding SQL injection
    )

    if err != nil {
        log.Fatalf("executing insert operation error: %v", err)
    }

    if tag.RowsAffected() == 0 {
        log.Fatalf("inserting the current row has not been successful: %v", role)
    }
}

commitCtx := context.Background()

if err = tx.Commit(commitCtx); err != nil {
    log.Fatalf("commit error: %v", err)
}
// [...]
```
We first start a transaction, and make several executions using a loop. We could have used [CopyFrom](https://pkg.go.dev/github.com/jackc/pgx/v5#Conn.CopyFrom), but for consistency we stick to `conn.Exec()`. Note also that we don't pass any formatted string with parameters injected using `fmt.Sprintf`, in order to avoid SQL injection attacks. Instead, we pass the parameters directly to `conn.Exec()`. 

#### Querying
Once we have created a new row and inserted to the database, it's time to check whether the operation has been executed correctly.

```go
// define this struct corresponding to the table separately:
type Role struct {
	Id        int
	Name      string
	ActorName string
}


// [...]
execCtx = context.Background() // another context for running the query
firstRow, err := conn.Query(execCtx, "SELECT id, name, actor_name FROM general.role WHERE id = $1", rows[0].Id)

if err != nil {
    log.Fatal("sending the query and initializing the rows has not been successful", err)
}

// the function in the second argument converts the raw row data to a Go data type
row, err := pgx.CollectExactlyOneRow(firstRow, func(rawRow pgx.CollectableRow) (Role, error) {
    value, err := rawRow.Values()

    if err != nil || len(value) != 3 {
        return Role{}, err
    }

    id := value[0].(int32)
    name := value[1].(string)
    actorName := value[2].(string)

    return Role{int(id), name, actorName}, nil
})

if err != nil {
    log.Fatalf("either 0, or more than one rows exist: %v", err)
}
// [...]
```
Here we used `pgx.CollectExectlyOneRow` as a helper function for checking a single row conveniently. There are several ways to do the same thing using this `pgx` package, but here we only show a couple of examples here; Please check out the [documentation page](https://pkg.go.dev/github.com/jackc/pgx) for details.

#### Updating and Deleting
Now let's update a row and delete another row at the same time, under the same transaction:

```go
// [...]
// start transaction
txCtx = context.Background()
tx, err = conn.Begin(txCtx)

if err != nil {
    log.Fatalf("opening transaction error: %v", err)
}

defer func() {
    rbCtx := context.Background()
    tx.Rollback(rbCtx)
}()

// update
execCtx = context.Background()
tag, err := conn.Exec(execCtx, "UPDATE general.role SET actor_name = $1 WHERE id = $2", strings.ToUpper(rows[0].ActorName), rows[0].Id)

if err != nil {
    log.Fatalf("executing update operation error: %v", err)
}

if tag.RowsAffected() == 0 {
    log.Fatalf("row with id %v has not been updated", rows[0].Id)
}

// delete
execCtx = context.Background()
tag, err = conn.Exec(execCtx, "DELETE FROM general.role WHERE LOWER(name) LIKE $1", "%"+"teddy"+"%")

if err != nil {
    log.Fatalf("executing delete operation error: %v", err)
}

if tag.RowsAffected() == 0 {
    log.Fatalln("row with name 'Teddy' has not been deleted")
}

// commit
commitCtx = context.Background()

if err = tx.Commit(commitCtx); err != nil {
    log.Fatalf("commit error: %v", err)
}
// [...]
```

The code has almost the same structure as the one in [Inserting](#inserting). Simple(and yes, unfortunately, Teddy is not a member of the Wolfpack).

By far we have covered the very basic operations that are used in database management. Of course, there are way more features in the database itself(not only Postgres though!). Since this book is not a book dedicated to databases, you may have to search for further materials if you are interested. You must be able to handle databases if you are a backend developer, anyways.

## Conclusion
If you get used to the basic operations covered in this chapter, you would probably get by the daily backend engineering tasks, although you may have to look up the documentations of the package `sqlx`, Postgres, or the cloud services like [AWS RDS](https://aws.amazon.com/rds/). Database is a huge subject on its own, so I highly recommend taking some time to study the field deeply. 

Before we wrap up the chapter, I would like to mention a few topics from database management for production. 

### A Word on Database Migration
The word "migration" reads as a bit daunting, but it means simply that there is a series of logs of SQL operations that have been executed ever since the database has been created, such as which column has been added, or on which columns we have indexes. Everytime there is a change in the database schema, you use a migration tool to make a log, so that people can keep track and revert the changes if necessary.

Although we haven't cover the migration tools, there are many tools out there, and you don't have to be tied to tools written in Go! Please check out on your own!

### A Word on ORMs
ORM(Object-Relational Mapping) allows developers define tables and other database-specific features in programming languages. For example, you have [Hibernate](https://hibernate.org/) in Java, and [SQLAlchemy](https://www.sqlalchemy.org/) in Python. Of course, Go also has a few famous ORM tools such as [Gorm](https://gorm.io/index.html).

It is debatable whether using ORM is good or not. Some people would prefer for being able to leverage the rich features of a mature programming language(in our case it would be Go), and other reasons such as easier migration to other databases. Some people would not prefer it for its innate N+1 problem, and some additional mental overburden for learning a new tools rather than sticking to simple SQL statements. Each has its own pros and cons, and it is upto you to decide.

## Exercise
Thanks for reading this long chapter(I didn't mean it in the first place...)! Since you have spent a lot of time on reading through the whole chapter, this time I won't ask you extra problems to solve. 

Instead, I would like you to *refactor* the long and winding `server.go`'s `main()` function. There is no clear answer, but I expect it to make sense :) 
