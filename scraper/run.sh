#!/bin/bash

#python3 WeatherCurrent.py 
#python3 WeatherForecast.py 
#python3 RealTimeBusDataScraping.py 
printf "Creating an array with scripts"

scripts=( "WeatherCurrent.py" "WeatherForecast" "RealTimeBusDataScraping" )

printf "Creating for loop to run all scripts"
for x in scripts
    python3 ${scripts[@]}
done

printf "This is running on 19/06 @19:30"
