#!/bin/bash

exec python3 WeatherCurrent.py &
exec python3 WeatherForecast.py & 
exec python3 RealTimeBusDataScraping.py &
