# THIS SCRIPT WILL SCRAPE THE GTFS-R DATA AND INSERT THE DATA IN THE COLLECTION
# DATA WILL BE SCRAPED EVERY MINUTE 


import urllib.request
import time
import certifi
import ssl
import json
import requests
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
    http_header = {"x-api-key":hdr}

    cluster = MongoClient(uri)
    db = cluster["BusData"]  # use a database called "BusData"
    collection = db["realTimeData"]  # and inside that DB, a collection called "real-timeData"

    try:
        print("P5 making the request & getting data")
        response = requests.get(url, headers=http_header)
        data = response.text

        print("loading the response into a json file")
        json_response = json.loads(data)

        # dropping the collection to have only most recent data
        collection.drop()

        print("inserting data")
        collection.insert_one(json_response)

                        # Aggregation
        print("Creating aggregation")
        cursor = collection.aggregate([{"$project" : {"_id":0}},
                                      {"$unwind": "$Entity"},
                                       {"$out": "realTimeData"}
                                  ])
                                  
        #  inserting the data in mongodb collection
        for document in cursor:
            collection.insert_many(document)


    except Exception as e:
        print(e)
    else:
        print("Data inserted successfully.")

    # close the connection
    finally:
        cluster.close()

    # real-time data will be scraped every 10 minutes
    time.sleep(10 * 60)


while True:
    gtfs_r()
