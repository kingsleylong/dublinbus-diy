# SCRIPT FOR SCRAPING CURRENT WEATHER DATA

# DATA IS SCRAPED EVERY 10 MINS AND INSERTED IN MONGODDB - CURRENTWEATHER COLLECTION
# TODO:
#  (1) REMOVE DUPLICATES BEFORE THEY ARE INSERTED IN THE COLLECTION

# importing libraries
import requests
import json
import time
import connectionsconfig

# error checking when connecting to MongoClient
try:
    from pymongo import MongoClient
except ImportError:
    raise ImportError('PyMongo is not installed')


def weather_current_main():
    response = requests.get(connectionsconfig.urlWeather)
    data = response.text

    # testing to ensure the data was scraped
    if response.status_code != 200:
        print('Failed to get data:', response.status_code)
    else:
        print('Data is: ', data)

    # parsing response text to json format
    print('[*] Parsing response text')
    data = json.loads(response.text)

    # print("pushing data to mongodb with the functions")

    # inserting data to mongodb database
    print('[*] Pushing data to MongoDB ')
    cluster = MongoClient(connectionsconfig.uri)
    db = cluster["Weather"]
    collection = db["currentWeather"]

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

    # weather info will be scraped every 10 minutes
    time.sleep(10 * 60)


while True:
    weather_current_main()
