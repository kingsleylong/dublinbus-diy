#!/bin/bash

nohup python3 -u RealTimeBusDataScraping.py &
nohup python3 -u WeatherCurrent.py &
python3 -u WeatherForecast.py
