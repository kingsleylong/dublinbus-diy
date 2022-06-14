# syntax=docker/dockerfile:1
FROM python:3.9.13-slim-bullseye
WORKDIR /web
COPY web/* /web/
VOLUME ["/web"]
CMD ["python3", "-m", "http.server", "8080", "--directory", "/web/"]
