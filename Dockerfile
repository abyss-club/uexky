FROM debian:latest
WORKDIR /app
COPY ./dist/uexky /app/uexky
CMD ["/app/uexky", "--help"]
