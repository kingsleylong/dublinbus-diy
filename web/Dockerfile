# syntax=docker/dockerfile:1
FROM python:3.9.13-slim-bullseye
WORKDIR /web
COPY web/* /web/
EXPOSE 80/tcp
VOLUME ["/web"]
CMD ["python3", "-m", "http.server", "80", "--directory", "/web/"]
