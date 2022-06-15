# syntax=docker/dockerfile:1
FROM python:3.9.13-slim-bullseye

WORKDIR /app/env
COPY requirements2.txt requirements.txt
RUN pip list
RUN pip install -r requirements.txt
RUN pip list

WORKDIR /app/config
VOLUME ["/app/config"]

WORKDIR /app/scripts
COPY *.py /app/scripts/

CMD ["python3", "WeatherCurrent.py"]
