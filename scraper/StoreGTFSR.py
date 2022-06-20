# THIS SCRIPT WILL SCRAPE THE GTFS-R DATA AND INSERT THE DATA IN THE COLLECTION
# TODO:
#  (1) REMOVE DUPLICATES BEFORE THEY ARE INSERTED IN THE COLLECTION


import urllib.request
import time
import certifi
import ssl
import json
import connectionsconfig
from pymongo import MongoClient

# error checking when connecting to MongoClient
try:
    from pymongo import MongoClient
except ImportError:
    raise ImportError('PyMongo is not installed')


def gtfs_r():
    # # connecting to MongoDB
    cluster = MongoClient(connectionsconfig.uri)
    db = cluster["BusData"]  # use a database called "BusData"
    collection = db["GTFSRdata"]  # and inside that DB, a collection called "bus"

    try:
        req = urllib.request.Request(connectionsconfig.url, headers=connectionsconfig.hdr)

        req.get_method = lambda: 'GET'
        response = urllib.request.urlopen(req, context=ssl.create_default_context(cafile=certifi.where()))

        # reading the API response & loading the response into a json file
        json_response = json.loads(response.read())

        collection.create_index([('timestamp', 1)], unique=True)
        print("inserting data")

        # inserting the data in mongodb collection
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
