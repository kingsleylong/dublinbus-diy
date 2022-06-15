# syntax=docker/dockerfile:1
FROM python:3.9.13-slim-bullseye
WORKDIR /app/config
WORKDIR /app/scripts
COPY *.py /app/scripts/
VOLUME ["/config"]
RUN ["python3", "WeatherCurrent.py"]
