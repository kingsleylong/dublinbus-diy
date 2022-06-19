#!/bin/bash

#python3 WeatherCurrent.py 
#python3 WeatherForecast.py 
#python3 RealTimeBusDataScraping.py 
printf "Creating an array with scripts 19/06 @20"

sudo python WeatherCurrent.py & sudo python WeatherForecast.py & sudo python RealTimeBusDataScraping.py

#scripts=( "WeatherCurrent.py" "WeatherForecast.py" "RealTimeBusDataScraping.py" )

#printf "Creating for loop to run all scripts"
#for x in ${scripts}
#    python ${scripts[@]} 
#done

printf "This is running on 19/06 @19:30"
