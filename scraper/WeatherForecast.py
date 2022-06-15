# SCRIPT FOR SCRAPING FORECAST WEATHER DATA
# DATA IS SCRAPED EVERY 2 HRS AND INSERTED IN MONGODDB - FORECAST COLLECTION
# TODO:
#  (1) REMOVE DUPLICATES BEFORE THEY ARE INSERTED IN THE COLLECTION


# importing libraries
import requests
import json
import time
from pymongo import MongoClient
import connectionsconfig


# function to insert the forecast weather to the mongodb collection
def weather_forecast_main(self):
    response = requests.get(connectionsconfig.urlForecast)
    data = response.text

    # testing to ensure the data was scraped
    if response.status_code != 200:
        print('Failed to get data:', response.status_code)
    else:
        print('Data is: ', data)

    # parsing response text to json format
    print('[*] Parsing response text')
    data = json.loads(response.text)

    # connecting to mongodb collection
    print('[*] Pushing data to MongoDB ')
    cluster = MongoClient(connectionsconfig.uri)
    db = cluster["Weather"]
    collection = db["forecast"]

    # inserting data in mongodb
    try:
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
    weather_forecast_main(connectionsconfig.url)
