FROM alpine:latest

WORKDIR /app

COPY release/v1.0.3/linux/amd64/mtauthsrv /app/

EXPOSE 8080

CMD ["./mtauthsrv"]