#!/bin/bash

nohup python3 WeatherCurrent.py &
nohup python3 WeatherForecast.py &
nohup python3 RealTimeBusDataScraping.py &

printf "This is running on 19/06 @14:40"
