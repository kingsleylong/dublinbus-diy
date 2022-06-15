# SCRIPT IS TO INSERT THE STATIC DATA INTO THE MONGODB COLLECTION

import json
import pandas as pd
from pymongo import MongoClient
import connectionsconfig


try:
    # reading the csv files and creating a pandas df
    df = pd.read_csv("routes.txt", sep=",", decimal=',')
    df.replace({',', ' '}, {'"', ' '}, regex=True, inplace=True)

    # Write to a separate JSON file
    array_json = df.to_json(orient='index')

    # creating a json file from the data
    with open('json_data.json', 'w') as outfile:
        outfile.write(array_json)

    # connecting to mongodb
    cluster = MongoClient(connectionsconfig.uri)
    db = cluster["BusData"]  # use a database called "BusData"
    collection = db["bus"]  # and inside that DB, a collection called "bus"

    # opening the json file created to insert it in the mongodb collection
    with open('json_data.json') as file:
        file_data = json.load(file)

    collection.insert_one(file_data)
except Exception as e:
    print(e)
else:
    print("Data inserted successfully.")

# close the connection
finally:
    cluster.close()
