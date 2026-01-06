# Chapter 1. Hello, Server!

Ok, let's start our journey to writing a server in Bun, with lot of fun!
We'll begin with the most common and widely used type of server, a HTTP server.
Bun as a JavaScript/TypeScript runtime [fully supports writing a HTTP(s) server](https://bun.com/docs/runtime/http/server). 

## 1. Run a simple HTTP server on port 8080
We start off by building a real, minimal HTTP server; but blocking its thread listening on a specific port. 
If you have initialized an empty Bun project correctly, there should be `index.ts` in the root directory. 

Now let's write the following code inside `index.ts` and run `bun run index.ts`:

```ts
const HOST = "0.0.0.0"; // default 
const PORT = 8080;

const server = Bun.serve({ // run the server
    hostname: HOST,
    port: PORT,
    routes: {
        "/": () => new Response("Hello, Bun Server!"),
    }
});

console.log(`Server running at ${server.url}`);
```

You can find exactly the same content [here in the official docs](https://bun.com/docs/runtime/http/server#changing-the-port-and-hostname). However, this is our first step for an entire server so I believe it is worth starting here.

- You specify hostname and port inside `Bun.serve()`. For simplicity, let us understand that the hostname of a server is a domain name of it (such as `localhost`), and the port is a serial number of a process running on the server. These are address information that clients need to know to send (HTTP) requests. Here we set the hostname as `0.0.0.0`, meaning that any request to the server is allowed. The port is arbitrary and here we set it as `8080`, but you'd better avoid well-known port numbers such as `80`(HTTP), `443`(HTTPS), or `22`(SSH). 
- `routes` is where we will going to list all the API endpoints we will expose to the frontend or other machines. 
- Once you create a server instance it is run asynchronously, such that the last line `console.log` is also executed. And you will see that the terminal is blocked by the server, which is awaiting requests.  

### Register endpoints with a simple response
Like any other server-side programming, writing a server with Bun requires you to register endpoints that clients hit with their HTTP requests. 

```ts
const server = Bun.serve({
    // [...]
    routes: {
        "/": () => new Response("Hello, Bun Server!\n"),
    }
});
```

Registering an endpoint to a Bun server is very straightforward; For the endpoint `/`, you will get the response text `"Hello, Bun Server!"`. Note that like other famous web frameworks, you don't imperatively write and flush your response bytes but simply return a new response object. 

### Test your server
Check out your brand-new server with `curl` on another terminal: `curl http://localhost:8080`. Once you get `Hello, Bun Server!` on your terminal, well done! You just wrote your first server in Bun!

## Summary
Writing an HTTP server in Bun is very simple and straightforward. Although there are quite a few options available for advanced configuration of a server, we would like to cover them in the later chapters. Also, note that we could also write not only HTTP but lower-level servers such as [TCP servers](https://bun.com/docs/runtime/networking/tcp), but those are beyond the interest of this book.

## Exercise
Add `/health` endpoint for health-checking. For the string message, `"Healthy"` would be alright  
