FROM golang:alpine as builder

RUN apk --no-cache  add tzdata

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ /etc/timezone

ARG service

WORKDIR /src
COPY . .

RUN CGO_ENABLE=0 go build -v -o /bundle ${service}

FROM alpine:latest
COPY --from=builder /bundle /bundle
COPY config.yaml /

CMD ["/bundle"]