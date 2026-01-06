const HOST = "0.0.0.0"; // default 
const PORT = 8080;

const server = Bun.serve({
    hostname: HOST,
    port: PORT,
    routes: {
        "/": () => new Response("Hello, Bun Server!\n"),
    }
});
