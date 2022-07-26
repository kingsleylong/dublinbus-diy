#!/bin/bash

nohup python3 -u RealTimeBusDataScraping.py &
nohup python3 -u WeatherCurrent.py &
nohup python3 -u StoreGTFSR.py &
nohup python3 -u WeatherForecast.py &
python3 -u FlaskApp.py