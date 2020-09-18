FROM golang:latest 
 
WORKDIR /app 

COPY ./ /app

RUN go mod download && go build -o main .

EXPOSE 8181

ENTRYPOINT go run main.go