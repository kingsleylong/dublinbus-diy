#!/bin/bash

python3 WeatherCurrent.py &
python3 WeatherForecast.py &
python3 RealTimeBusDataScraping.py &

printf "This is running on 19/06 @14:00"
