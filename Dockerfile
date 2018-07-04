FROM golang:1.9.0 as builder
WORKDIR /go/src/github.com/moooofly/harbor-go-client
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a

# Final image.
FROM scratch
LABEL maintainer "moooofly <centos.sf@gmail.com>"
COPY --from=builder /go/src/github.com/moooofly/harbor-go-client/harbor-go-client .
COPY conf /conf
ENTRYPOINT ["/harbor-go-client"]
CMD ["-h"]
