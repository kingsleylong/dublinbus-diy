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
    collection = db["GTFSRdata"]  # and inside that DB, a collection called "bus"

    try:
        print("making the request & getting data")
        response = requests.get(url, headers=http_header)
        data = response.text

        print("loading the response into a json file")
        json_response = json.loads(data)
        
        # inserting the data in mongodb collection
        print("inserting data")
        collection.insert_one(json_response)

    except Exception as e:
        print(e)
    else:
        print("Data inserted successfully.")

    # close the connection
    finally:
        cluster.close()

    # real-time data will be scraped every 4hrs
    time.sleep(14400 * 60)


while True:
    gtfs_r()
