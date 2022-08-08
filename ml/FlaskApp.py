# import necessary modules
import json

import bson
from flask import Flask, jsonify
import pandas as pd
import pickle
from pymongo import MongoClient
from pymongo.errors import ConnectionFailure
import configparser

# load config file
print('reading configurations')
config = configparser.ConfigParser()
config.read('config/scrapercfg.ini')
connectionsconfig = config['scraper']

# connecting to MongoDB 
uri = connectionsconfig['uri']

# flask app into variable
app = Flask(__name__)

@app.route('/prediction/<line>/<dir_>/<day>/<hour>/<month>/<departure_time>/<date_txt>', methods=['GET', 'POST'])
def get_prediction(line, dir_, day,hour, month,departure_time,date_txt): #full_date_hour
# allow prediction model on analytics page to take user inputs as prediction model parameters

# print url parameters 
    print('line:', line, ', direction:',dir_, ', day:',day,', hour:', hour, ', month:', month, ',      departure_time:', departure_time, ', date_txt:', date_txt)

# change date_txt to ensure it matches the hour in mongo
    date_split = date_txt.split(" ")
    time = date_split[1]
    time_split = time.split(":")
    time_split[1] = "00"
    time_split[2] = "00"
    format_time = time_split[0]+":"+time_split[1]+":"+time_split[2]
    format_date = date_split[0]+" "+format_time

# open pickle file and load into variable clf
    with open("/usr/local/dublinbus/data/ml/Pickles/" + f'RF_{line}_Model_dir{dir_}.pkl', 'rb') as pickle_file:
        clf = pickle.load(pickle_file)

# Establishing connection
    try:
        connect = MongoClient(uri)
        print(connect)
        print("Connected to MongoDB")
    except ConnectionFailure as e:
        print("Could not connect to MongoDB")
        print(e)

# get required data from mongodb
# Connecting or switching to the database
    db = connect.Weather

    # Switching to Forecast
    collection = db.Forecast

    # Select the temp feature 
    cursor = collection.aggregate ([
        {"$unwind" : {
            "path" : "$list"
        }},
        {
            "$match" :{
                "list.dt_txt" : format_date
            }

        },
        {"$project" :
            {
                "temp": "$list.main.temp",
                "_id": 0
            }
        }
])
    # get the values for temp
    for doc in cursor:
        for value in doc:
            mongo_temp = (doc[value])
            # print(mongo_temp)


        # get form values
    direction = f'{dir_}'#from go api
    day = f'{day}'#from go api
    hour = f'{hour}'#from go api
    month = f'{month}'#from go api
    departure_time = f'{departure_time}'#from go api
    temp = mongo_temp

    X = pd.DataFrame([[direction, day, hour, month, departure_time, temp]], columns=[
                         "direction", "weekday", "hour", "month", "departure_time", "temp"]).values
    
    print((X)[0])

        # generate prediction
    pred_seconds = clf.predict(X)[0]

    # read the error range (MAE) from csv file to get a range for the prediction
    data = pd.read_csv("/usr/local/dublinbus/data/ml/Pickles/" + f"line_{line}_rf_metrics_dir{dir_}.csv", sep = ":",names=[' metrics ',  'values'])
    value = data.iloc[7]
    value[1]

    pos_error_pred = pred_seconds + value[1]
    neg_error_pred = pred_seconds - value[1]

# changing the prediction to minutes
    pred_minutes = pred_seconds / 60
    pos_pred_minutes = pos_error_pred / 60
    neg_pred_minutes = neg_error_pred / 60

    full_predictions = pred_minutes, pos_pred_minutes, neg_pred_minutes
# return a json object
    return jsonify(full_predictions)

# run flask app
if __name__ == '__main__':
    app.run(threaded=False, host='0.0.0.0')
