# SCRIPT FOR SCRAPING FORECAST WEATHER DATA
# DATA IS SCRAPED EVERY 2 HRS AND INSERTED IN MONGODDB - FORECAST COLLECTION


# importing libraries
import requests
import json
import time
from pymongo import MongoClient
import os
import configparser


# load config file
print('reading configurations')
config = configparser.ConfigParser()
config.read('config/scrapercfg.ini')
connectionsconfig = config['scraper']


# function to insert the forecast weather to the mongodb collection
def weather_forecast_main():
    urlForecast = connectionsconfig['urlForecast']
    urlForecast = urlForecast + "?lat=%s&lon=%s&appid=%s&units=metric"
    urlForecast = urlForecast % (
            connectionsconfig['lat'],
            connectionsconfig['lon'],
            connectionsconfig['api_key_forecast'])

    response = requests.get(urlForecast)
    data = response.text

    # testing to ensure the data was scraped
    if response.status_code != 200:
        print('Failed to get data:', response.status_code)
    #else:
    #    print('Data is: ', data)

    # parsing response text to json format
    print('[*] Parsing response text')
    data = json.loads(response.text)

    # connecting to mongodb collection
    print('[*] Pushing data to MongoDB ')
    cluster = MongoClient(connectionsconfig['uri'])
    db = cluster["Weather"]
    collection = db["Forecast"]

    # inserting data in mongodb
    try:
        # dropping the forecast collection
        collection.drop()
        # creating a new collection and inserting new data
        collection.insert_one(data)
    except Exception as ex:
        print(ex)
    else:
        print("Data inserted successfully")

    # close the connection
    finally:
        cluster.close()

    # forecast will be scraped every 2 hours
    time.sleep(120 * 60)


while True:
    weather_forecast_main()
