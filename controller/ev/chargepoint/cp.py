# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

from config import config
import constants
from datetime import datetime
from ev import ChargePoint, ChargeSessions, EVException, get_ENV_val
# default JSON serializer of zeep.helpers.serialize_object doesnt handle
# Decimal object but simplejson does.
import simplejson as json
import logging
import os
import time
from zeep import Client
import zeep.helpers
from zeep.wsse.username import UsernameToken

logging.Formatter.converter = time.gmtime
formatter = logging.basicConfig(
    format='%(asctime)s.%(msecs)03dZ %(levelname)s:%(name)s:%(message)s',
    datefmt='%Y-%m-%dT%H:%M:%S',
    level=logging.DEBUG)
logging.getLogger('zeep').setLevel(logging.ERROR)

# COMMON config
DEBUG = "True" == config["COMMON"]["DEBUG"]

# CHARGEPOINT Constants
WSDL = config["CHARGEPOINT"]["WSDL"]

# To be read from environment
_CP_USERNAME = "CP_USERNAME"
_CP_PASSWORD = "CP_PASSWORD"


# POLL config
_SAVE_TO_FILE = "True" == config["POLL"]["SAVE_TO_FILE"]
_INTERVAL = int(config["POLL"]["INTERVAL"])
_FILENAME = config["POLL"]["FILENAME"]
_SAVE_ROOT = config["POLL"]["SAVE_ROOT"]
_STREAM = "True" == config["POLL"]["STREAM"]

# Globals
_CHARGE_SESSIONS = ChargeSessions()
_CLIENT = None
_SAVE_PATH = None
_DATAFILE = None
_CURRENT_SAVE_PATH = None
_CURRENT_HOUR = None

# Environment Variables
_STATION_GROUP_ID = None
# list of stationIDs to monitor
_STATION_IDS = None

LOAD = None


# TODO throw error if not defined
def get_username_token():
    return UsernameToken(
        get_ENV_val(_CP_USERNAME),
        get_ENV_val(_CP_PASSWORD))


def get_file():
    """
    Takes a configured save root, and creates a directory
    structure that comprehends year, month, day, and hour
    and creates a file for appending in that directory with
    configured data file name
    """ 
    global _CURRENT_SAVE_PATH, _CURRENT_HOUR, _DATAFILE, _FILENAME

    now = datetime.now()
    if (_CURRENT_HOUR != now.hour):

        if (_DATAFILE is not None):
            if DEBUG:
                logging.debug("Closing file " + _FILENAME + "\n")
            # close old file
            _DATAFILE.close()

        # construct new file path and open new file
        save_path = os.path.join(_SAVE_ROOT, str(now.year))
        save_path = os.path.join(save_path, str(now.month))
        save_path = os.path.join(save_path, str(now.day))
        save_path = os.path.join(save_path, str(now.hour))
        _CURRENT_HOUR = now.hour
        # make directories as necessary
        try:
            if (not os.path.isdir(save_path)):
                os.makedirs(str(save_path), 0o777)
        except Exception as ee:
            logging.error(ee)
        filename = os.path.join(save_path, _FILENAME)
        if DEBUG:
            logging.debug("Writing load data to " + filename + "\n")
        _DATAFILE = open(filename, "a")
    return _DATAFILE


def save_to_file(load):
    global _DATAFILE
    sg_load = load['sgLoad']
    if DEBUG:
        print(
            "time: " + str(time.time()) +
            "  total_load @ sgID(" +
            str(load["sgID"]) + ") = " +
            sg_load + "\n"
        )
    if (sg_load == '0.000'):
        print("Skip writing to file, current load zero\n")
        return
    else:
        datafile = get_file()
        if DEBUG:
            print("Writing to file " + str(datafile))

        datafile.write(json.dumps({
            'ts': time.time(),
        '   data': zeep.helpers.serialize_object(load),
        }))
        datafile.flush()
        return


def init():
    global _CLIENT
    _CLIENT = Client(WSDL, wsse=get_username_token())
    if (_CLIENT is not None):
        logging.info("Created SOAP ChargePoint client")

    # save to file?
    if _SAVE_TO_FILE:
        _FILENAME = config["POLL"]["FILENAME"]
        _SAVE_ROOT = config["POLL"]["SAVE_ROOT"]


# TODO: handle errors
# Obtain load for all Hilltop Stations
# Note load response is of type <class 'zeep.objects.getLoadResponse'>
# and prints as a dictionary with a list of Station Data
# and is a genuine python object whose fields you can access

def _get_client():
    global _CLIENT
    if (_CLIENT is None):
        init()
    if DEBUG:
        print(_CLIENT)
    return _CLIENT


def _load_query(id):
    cp = ChargePoint.fromID(id)
    # constructing search query
    search_query = {}
    if (cp.sgID is not None):
        search_query["sgID"] = cp.sgID
    if (cp.stationID is not None):
        search_query["stationID"] = cp.stationID
    return search_query


def _get_load(client, search_query):
    return client.service.getLoad(search_query)


def process_load(load_res):
    res = []
    sgID = load_res['sgID']
    now = time.time()
    for station in load_res['stationData']:
        stationID = station['stationID']
        for port in station['Port']:
            cp = ChargePoint(sgID, stationID, port['portNumber'])
            vehicle = port['credentialID']
            watt = float(port['portLoad'])
            res.append({
                'point': cp.ID,
                'vehicle': vehicle,
                'watt': watt,
                'measured': now,
                'start': _CHARGE_SESSIONS.get_start(vehicle, watt, now)
            })
    return res


def get_load(id):
    try:
        load = _get_client().service.getLoad(_load_query(id))
        if DEBUG:
            print(load)
        print(str(process_load(load)))
        print("\n---------------\n")
        return load
    except Exception as e:
        errMsg = "Unable to getLoad data from " + WSDL + e
        logging.error(errMsg)
        return None


def shed_load(id, percent_amount=0, absolute_amount=0, time_interval=0):
    if DEBUG: 
        print("Enter shed_load")
    cp = ChargePoint.fromID(id)
    shed_query = {}
    # if station parameters present, construct the more restrictive query 
    if (cp.stationID is not None):
        station_query = {"stationID": cp.stationID}
        if (cp.portID is not None):
            port_query = {"portNumber": cp.portID}
            if (absolute_amount is not None):
                port_query["allowedLoadPerPort"] = absolute_amount
            else:
                # ns0:percentShedRange ??
                port_query["percentShedPerPort"] = percent_amount
            station_query["Ports"] = [{"Port": port_query}]
        shed_query["shedStation"] = station_query
    else:
        # construct a more general stationGroup query
        shed_query["shedGroup"] = {"sgID": cp.sgID}

    # shed amount
    if ((absolute_amount is None) and (percent_amount is None)):
        raise EVException("Both absolute_amount and Percent_amount cannot be None",
                            constants.ERR_MISSING_AMOUNT)

    if (absolute_amount is not None):
        shed_query["allowedLoadPerStation"] = absolute_amount
    else:
        # ns0:percentShedRange??
        shed_query["percentShedPerStation"] = percent_amount

    # shed interval, 0 means no limit
    shed_query["timeInterval"] = time_interval
    if DEBUG:
        print("shedQuery = " + str(shed_query))
    return _get_client().service.shedLoad(shed_query)


def clear_shed(id):
    if DEBUG:
        print("enter clear_shed " + id)
    shed_query = {}
    cp = ChargePoint.fromID(id)
    if (cp.stationID is not None):
        station_query = {"stationID": cp.stationID}
        if (cp.portID is not None):
            # technically can provide multiple ports
            station_query["Ports"] = [{"Port": {"portNumber": cp.portID}}]
        shed_query["shedStation"] = station_query
    else:
        shed_query["shedGroup"] = {"sgID": cp.sgID}

    if DEBUG:
        print("clear_shed query = " + str(shed_query))

    return _get_client().service.clearShedState(shed_query)


# Various Chargepoint API calls
def get_CPN_instances():
    return _get_client().service.getCPNInstances()


def streamData(url):
    print("Streaming not yet implemented")
    # Might want to anonymize data
    # Send up some averaged data


def poll_load(id):
    client = _get_client()
    search_query = _load_query(id)
    while True:
        try:
            load = _get_load(client, search_query)
            if _SAVE_TO_FILE:
                save_to_file(load)
            else:
                if _STREAM:
                    # TODO
                    print("Streaming not yet implemented")
        except Exception as ee:
            print("Error: poll_load encountered exception: " + str(ee))
        # seconds
        time.sleep(_INTERVAL * 60)
