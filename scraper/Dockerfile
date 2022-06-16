# syntax=docker/dockerfile:1
FROM python:3.9.13-slim-bullseye

WORKDIR /app/env
COPY requirements2.txt requirements.txt
RUN pip list
RUN pip install -r requirements.txt
RUN pip list
COPY rt_leavetimes_DB_2018.txt

WORKDIR /app/scripts/config
VOLUME ["/app/scripts/config"]

WORKDIR /app/scripts
COPY *.py /app/scripts/

CMD ["python3", "-u", "WeatherCurrent.py"]
CMD ["python3", "-u", "WeatherForecast.py"]
CMD ["python3", "-u", "RealTimeBusDataScraping.py"]
CMD ["python3", "-u", "HistBusData.py"]
