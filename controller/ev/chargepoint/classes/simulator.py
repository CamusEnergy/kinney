# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

import math
import time

from classes.charge_session import ChargeSession
from classes.ev_exceptions import EVException
from classes.full_port import FullPort

# Types of locations
HOME_LOCATION = "home"
OFFICE_LOCATION = "office"
STORE_LOCATION = "store"

#USA
REGION_PREFIX = "1:"

# track all charging sessions
_charge_sessions = []


class CPInstance():
    '''A representation of a charge port
    attributes:
        locationType: a string, Mall, Office ..
        sgID: station group ID, typically a logical grouping
        num_ports: 1 or more
        ports: a list of FullPort objects constructed based on num_ports
        address: string 
        capacity: total charge capacity
        shed_state: True or False
        allowed_load: curtailed capacity
        shed_percent: amount of curtailment as a percentage of total capacity
        last_update: last simulation time .. updated every time an action is taken against the instance
    '''
    locationType = None
    sgID = None
    group_name = "ChargePoint Simulator"
    num_ports = 0
    address = ""
    ports = []
    capacity = None
    shed_state = False
    allowed_load = None
    shed_percent = 0
    last_update = None

    def __init__(self, sgID, num_ports, address, capacity, locationType=STORE_LOCATION):
        self.locationType = "Office"
        self.sgID = sgID
        self.address = address
        self.capacity = capacity
        self.num_ports = num_ports
        self.last_update = time.time()
        full_load = (capacity / num_ports) # is this valid?
        # 2 ports per station
        for i in range(0, num_ports):
            stationID = REGION_PREFIX + str(math.floor((i+1)/2))
            portID = ((i+1) % 2)
            port = FullPort(sgID, stationID, portID)
            port.set_specs(address, full_load)
            self.ports.append(port)

    def fast_forward(self):
        now = time.time()
        elapsed_time_mins = (now - self.last_update)/60
        print("Simulator fast forward interval: " + str(elapsed_time_mins))
        # granularity minutes
        if (elapsed_time_mins >= 1.0):
            for p in self.ports:
                if p.occupied():
                    p.update_charge(now) 
            self.last_update = now

    # Matching output to ChargePoint Instance wrapper
    def get_total_load(self):
        total_load = 0.0
        for p in self.ports():
            if p.occupied():
                total_load = total_load + p.get_load()
        return total_load

    # absolute amount, percentage amount
    # TODO could curtail at the level of an individual station and port
    def shed(self, percent, amount):
        if ((percent is None) and
            (amount is None)):
            raise EVException("One of amount or percent must be specified", 303)
        self.shed_state = True
        if (amount is not None):
            self.allowed_load = amount
            self.shed_percent = None
        else: 
            self.allowed_load = self.capacity * percent * 0.01
            self.shed_percent = percent

    def clear(self):
        self.fast_forward()
        self.shed_state = False
        self.shed_percent = 0
        self.shed_amount = self.capacity

    def get_free_port(self):
        for p in self.ports:
            if (p.available()):
                return p
        return None

    def plugin(self, vehicle):
        self.fast_forward()
        free_port = self.get_free_port()
        if (free_port is not None):
            cs = ChargeSession(vehicle, free_port.ID)
            free_port.charge_session = cs 
            return cs

    def unplug(self, vehicle):
        self.fast_forward()
        # find the port where vehicle was located, free it
        # updated sessions with the vehicle charging data
        for p in self.ports:
            if (p.get_vehicleID() == vehicle.ID):
                cs = p.unplug()
                print("Vehicle: " + vehicle.ID +
                      " unplugged. Port: " + p.ID + " free")
                _charge_sessions.append(cs)

    def get_num_charge_sessions(self):
        count = 0
        for p in self.ports:
            if p.occupied():
                count = count + 1
        return count

    def get_load(self):
        return self.fmt_station_group()
            
    def fmt_stations(self):
        curr_station = None
        port_list = []
        station_fmt_list = []
        for port in self.ports:
            if (port.stationID != curr_station):
                if (len(port_list) > 0):
                    station_fmt_list.append(fmt_station(port_list))
                port_list = [port]
                curr_station = port.stationID
            else:
                port_list.append(port)
        return station_fmt_list

    def fmt_station_group(self):
        station_data_list = self.fmt_stations()
        return {
            'responseCode': '100',
            'responseText': 'API input request executed successfully.',
            'sgID': self.sgID,
            'numStations': (self.num_ports/2),
            'groupName': "Sim" + str(self.sgID),
            'sgLoad':  "Decimal(" + str(self.get_total_load()) + ")",
            'stationData': station_data_list
            }


def fmt_port(aPort):
    return {'portNumber': aPort.ID,
            'userID': 'Anonymous',
            'credentialID': aPort.get_vehicleID(),
            'shedState': aPort.get_shed_state(),
            'portLoad': "Decimal(" + str(aPort.get_load()) + ")",
            'allowedLoad': "Decimal(" + str(aPort.get_allowed_load()) + ")",
            'percentShed': aPort.get_percent_shed()
            }

def fmt_station(port_list):
    port0 = port_list[0]
    stationID = port0.stationID
    station_load = 0.0
    port_fmt_list = []
    for port in port_list:
        station_load = station_load + port.get_load()
        port_fmt_list.append(fmt_port(port))
    return {'stationID': stationID,
            'stationName':  str(stationID) + "Sim",
            'Address':  str(stationID) + " Somewhere",
            'stationLoad': "Decimal(" + str(station_load) + ")",
            'ports': port_fmt_list
            }
