FROM alpine:3.7
LABEL maintainer="dwj <dong.wenjuan@zte.com.cn>"

COPY gerrit_event /bin/gerrit_event

ENTRYPOINT  [ "/bin/gerrit_event" ]
