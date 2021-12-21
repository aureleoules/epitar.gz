FROM golang
COPY src/ /app
WORKDIR /app
RUN go build -o /app/app
WORKDIR /output
CMD ["/app/app"]
