#!/bin/bash

nohup python3 -u RealTimeBusDataScraping.py &
nohup python3 -u WeatherCurrent.py &
nohup python3 -u StoreGTFSR.py &
python3 -u WeatherForecast.py
