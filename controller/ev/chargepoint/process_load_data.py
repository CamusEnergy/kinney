#!/usr/bin/python

import getopt
import json
import os
import sys 

from classes.charge_session import ChargeSession
from classes.full_port import FullPort


charge_sessions = {}
#DATA_KEY = '   data'
DATA_KEY = 'data'

def process_file(full_filename, csv_filename):
    num_lines = 0
    with open(csv_filename, "a") as csv_out:
        with open(full_filename) as json_file:
            while (True):
                line = json_file.readline()
                if not line:
                    break
                num_lines = num_lines + 1
                print(str(line))
                if (line != "\n"):
                    process_line(line, csv_out)
                
    print("Num lines read = " + str(num_lines))

def process_line(line, csv_file):
    global DATA_KEY, charge_sessions
    if (line != "\n"):
        data = json.loads(line)
        timestamp = data['ts']
        sgID = data[DATA_KEY]['sgID']
        for sd in data[DATA_KEY]['stationData']:
            stationID = sd['stationID']
            for p in sd['Port']:
                port_load = p['portLoad']
                if (port_load > 0.0):
                    portID = p['portNumber']
                    credentialID = p['credentialID']
                    full_PortID = FullPort.buildID(sgID, stationID, portID)
                    data_str = "{:10.6f}, {:>25s},   {:2.3f}, {:>25s} \n". format(timestamp, credentialID, port_load ,full_PortID)
                    csv_file.write(data_str)



def process_day(dayDir):
    print("enter process_day")
    print(dayDir)
    csv_filename = os.path.join(dayDir, "charge_data.csv")
    print("CSV filename: " + str(csv_filename))
    for hh in os.listdir(dayDir):
        print("hour = " + str(hh))
        if ((not hh.endswith(".csv")) and
            (not hh.endswith(".xls"))):
            hourDir = os.path.join(dayDir,hh)
            for dd in os.listdir(hourDir):
                if (dd.startswith("load2") and dd.endswith(".json")):
                    filename = os.path.join(hourDir, dd)
                    print("Processing file:" + str(filename))
                    process_file(filename, csv_filename)
                else:
                    continue

def main(argv):
   inputDayDir= ''
   outputfile = ''
   try:
      opts, args = getopt.getopt(argv,"hi:o:",["iDir=","ofile="])
   except getopt.GetoptError:
      print("process_load_data.py -i <inputDayDir> -o <outputfile>")
      sys.exit(2)
   for opt, arg in opts:
      if opt == '-h':
         print("process_load_data.py -i <inputDayDir> -o <outputfile>")
         sys.exit()
      elif opt in ("-i", "--iDir"):
         inputDayDir = arg
      elif opt in ("-o", "--ofile"):
         outputfile = arg
   print("Input Day Directory is " + inputDayDir)
   ## need to check it ends with csv if in the same directory etc.
   print("Ignoring Output file: " + outputfile)
   process_day(inputDayDir)

if __name__ == "__main__":
   main(sys.argv[1:])

