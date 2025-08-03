# syntax=docker/dockerfile:1
FROM golang:1.24

WORKDIR /app

COPY go.mod ./

COPY server.go ./
COPY ./cmd/ ./cmd

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /gensecret ./cmd/gensecret.go 
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./server.go 

# run the script
RUN chmod a+x /gensecret
RUN --mount=type=secret,id=actorName,env=ACTOR_NAME\
    --mount=type=secret,id=ethnicity,env=ETHNICITY\
    /gensecret -actorName="$ACTOR_NAME" -ethnicity="$ETHNICITY"

EXPOSE 8080

CMD ["/server"]
