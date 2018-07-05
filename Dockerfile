FROM scratch

ADD main ./

EXPOSE 5000

CMD ["/main", "-c", "config.json"]
