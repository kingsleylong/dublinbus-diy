# syntax=docker/dockerfile:1
FROM python:3.9.13-slim-bullseye

# make a folder in the container for storing the python env file
WORKDIR /app/env
COPY requirements2.txt requirements.txt
# install the python dependencies
RUN pip list
RUN pip install -r requirements.txt
RUN pip list

WORKDIR /app/scripts/config
# this volume will be mounted in docker-compose.yml to a host config folder
VOLUME ["/app/scripts/config"]

WORKDIR /app/scripts/data
# this volume will be mounted in docker-compose.yml to a host data folder
VOLUME ["/app/scripts/data"]

WORKDIR /app/scripts
# copy the scraper scritps in the current folder
COPY *.py /app/scripts/

COPY run.sh
RUN chmod a+x run.sh

CMD ["./run.sh"]

# start the scraper scripts one by one
#RUN ["python3", "-u", "WeatherCurrent.py"]
#RUN ["python3", "-u", "WeatherForecast.py"]
#RUN ["python3", "-u", "RealTimeBusDataScraping.py"]




