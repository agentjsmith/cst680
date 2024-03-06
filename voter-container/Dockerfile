FROM golang:latest as build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /voter-api-linux .

FROM scratch as runtime

WORKDIR /

COPY --from=build /voter-api-linux /

EXPOSE 1080

ENV VOTER_API_REDIS_ADDR=host.docker.internal:6379

ENTRYPOINT [ "/voter-api-linux" ]