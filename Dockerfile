FROM debian:latest

COPY ./dist/uexky /app/uexky
CMD ["/app/uexky", "--help"]
