# Configure Your Server
This chapter is mainly about how to configure your server. Here the verb *"configure"* means that we would like to inject a set of parameters when the server is up and running, such that its behavior is affected by the parameters.

todo: add more words here

## Using Environment Variables
One of the easiest and a basic way of configuring your server(or any kind of application) is using the environment variables defined on the process in which your server is running.  

```go
if secret, ok := os.LookupEnv("secret"); ok {
    fmt.Printf("secret: %v\n", secret)
}
```
Here we use [`os.LookupEnv`](https://pkg.go.dev/os#LookupEnv) rather than [`os.GetEnv`](https://pkg.go.dev/os#Getenv) because we want to check the variable in a more systematic way. 
