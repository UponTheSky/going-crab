# going-crab
A book on how to write an API server in [Bun](https://bun.com/) with fun.

## Before you read this book
Install [Bun](https://bun.com/docs/installation) first (it is *blazingly* fast and simple). Once installed, you can create a brand new project using `bun init` command elsewhere, or simply run `bun install` inside the directory of this book here to simply follow the written code here. However, I would like to recommend creating your own project(i.e., by `bun init`) and developing your own server freely while following the content of the book. 

If you have no experience in TypeScript, it is fine. I won't assume any deep knowledge in TypeScript, so any simple tutorial will do. What I would suggest is to simply follow the book and reference any material once you get stuck on the language. But at least you would definitely be happy with the TypeScript language server (LSP)(VS Code has it by default, for your information!).  

## Why This Project?
This project is for gathering and polishing my knowledge of backend development. 
Although there are tons of books out there serving the same or similar purposes, I believe this project has its specialty as follows:

1. Organic structure between the chapters: Whereas lots of books introduce topics as entirely separate chapters such as HTTP requests, logs, infrastructures, and so on, this book tries to start off from the perspective of a newbie who is about to dive into the backend development world. Each chapter is in conjunction with others based on the needs of a developer. However, this project can also be used as a reference book, as the later stage of the book will show self-contained content for the topics. 

2. Avoiding "behind the magic" as much as possible: Many IT-related books, let alone the backend development books, rely heavily on web frameworks and libraries. These may be essential for production code, but many of them hide how server applications work under the hood using *magic* methods and techniques. This book tries to avoid such frameworks and libraries if possible. 

This is why we chose [Bun](https://bun.com/) for developing a server application. Bun enables developing necessary functionality of a server with minimal dependencies. Although you have to know TypeScript a bit, I assume you are already an experienced programmer who has touched a bit of frontend technologies. If not, please check out some JavaScript and TypeScript books and materials first.

3. Reading a book with fun: I know; this is the hardest part. But I have read so many books with tons of words without any rewarding feeling. You need to have some fun while reading and coding. I tried my best to give readers fun through the writing style and the structure of the book, but ultimately it is up to you, readers, to have fun. 

## Contribution
Any contribution is welcome! Please make an issue first on the issue page of this repo, and create a PR from the issue. I will check it on a regular basis.
