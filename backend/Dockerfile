FROM alpine:3.12

ENV USER=appuser
ENV APP=opa-center
ENV UID=1000
ENV GID=1000

RUN apk add --update ca-certificates && rm -rf /var/cache/apk/* && \
    addgroup -g $GID $USER && \
    adduser -D -g "" -h "/$USER" -G "$USER" -H -u "$UID" "$USER"

WORKDIR /$USER

COPY $APP /$USER/$APP
COPY static /$USER/static

RUN chown -R $UID:$GID /$USER

USER $USER

ENTRYPOINT [ "/appuser/opa-center" ]
