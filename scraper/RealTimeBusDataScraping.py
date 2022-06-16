# THIS SCRIPT WILL SCRAPE THE GTFS-R DATA AND INSERT THE DATA IN THE COLLECTION
# DATA WILL BE SCRAPED EVERY MINUTE 


import urllib.request
import time
import certifi
import ssl
import json
import os
import configparser
from pymongo import MongoClient

# error checking when connecting to MongoClient
try:
    from pymongo import MongoClient
except ImportError:
    raise ImportError('PyMongo is not installed')

# load config file
print('reading configurations')
config = configparser.ConfigParser()
config.read('config/scrapercfg.ini')
connectionsconfig = config['scraper']

def gtfs_r():
    # connecting to MongoDB & gtfs_data
    uri = connectionsconfig['uri']
    url = connectionsconfig['url']
    hdr = connectionsconfig['hdr']

    cluster = MongoClient(uri)
    db = cluster["BusData"]  # use a database called "BusData"
    collection = db["realTimeData"]  # and inside that DB, a collection called "real-timeData"

    try:
        print("making the request")
        req = urllib.request.Request(url, headers=hdr)

        req.get_method = lambda: 'GET'
        response = urllib.request.urlopen(req, context=ssl.create_default_context(cafile=certifi.where()))

        print("reading the data")
        # reading the API response
        gtfs_data = response.read()

        if isinstance(gtfs_data, dict):
    # loading the response into a json file
            json_response = json.loads(gtfs_data.decode('utf-8'))

# inserting the data in mongodb collection
            collection.insert_one(json_response)

        # print("loading the data to json file")
        # # loading the response into a json file
        # json_response = json.loads(gtfs_data.decode('utf-8'))

            # dropping the collection to have only most recent data
            collection.drop()

            # inserting the data in mongodb collection
            collection.insert_one(json_response)

    except Exception as e:
        print(e)
    else:
        print("Data inserted successfully.")

    # close the connection
    finally:
        cluster.close()

    # real-time data will be scraped every 10 minutes
    time.sleep(1 * 60)


while True:
    gtfs_r()
