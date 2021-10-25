FROM golang:1.13-alpine3.10 AS builder

WORKDIR /go/src
RUN mkdir -p /go/src/github.com/dongwenjuan/gerrit_event
ADD ./ /go/src/github.com/dongwenjuan/gerrit_event
RUN go build -o /go/bin/gerrit_event /go/src/github.com/dongwenjuan/gerrit_event/.


FROM alpine:3.7
LABEL maintainer="dwj <dong.wenjuan@zte.com.cn>"

COPY --from=builder /go/bin/gerrit_event /bin/gerrit_event
COPY id_rsa   /tmp

CMD  [ "/bin/gerrit_event" ]