FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /chatbot .

FROM alpine:3.19
WORKDIR /app
COPY --from=build /chatbot .
COPY static ./static
EXPOSE 8080
CMD ["./chatbot"]
