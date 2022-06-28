# THIS SCRIPT WILL SCRAPE THE GTFS-R DATA AND INSERT THE DATA IN THE COLLECTION EVERY 4HRS 

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

    # # connecting to MongoDB
    cluster = MongoClient(uri)
    db = cluster["BusData"]  # use a database called "BusData"
    collection = db["storeGtfrs"]  # and inside that DB, a collection called "bus"

    try:
        print("making the request & getting data")
        time.sleep(1*60)
        response = requests.get(url, headers=http_header)
        data = response.text

        print("loading the response into a json file")
        json_response = json.loads(data)
        

        # inserting the data in mongodb collection
        print("inserting data")
        collection.insert_one(json_response)

                                # Aggregation
        cursor = collection.aggregate([{"$project" : {"_id":0}},
                                      {"$unwind": "$Entity"},
                                       {"$out": "storeGtfrs"}
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

    # real-time data will be scraped every 3hrs
    time.sleep(10800 * 60)


while True:
    gtfs_r()
