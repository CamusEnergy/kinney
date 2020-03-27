# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

import math

from classes.charge_session import ChargeSession
from classes.ev_exceptions import EVException
import constants

class FullPort():
    '''A representation of a charge port
    attributes:
        sgID: station group ID, typically a logical grouping
        stationID: a unique identifier
        portID: port ID
        ID: a fully qualified ID that is constructed from sgID, stationID, and portID
        address: its physical location
        full_load: how much charge the port can push out in unit time 
        shed_state: True or False, true when it is being curtailed
        allowed_load: less that full_load when it is being curtailed
        percent_shed: percentage shedding of full load
        load: the actual amount of charge port is supplying
        charge_session: a ChargeSession object when a vehicle is plugged in
    '''
    sgID = None
    stationID = None
    portID = None
    ID = None
    address = None
    full_load = None
    shed_state = False
    allowed_load = None
    percent_shed = 0
    load = 0.0
    charge_session = None

    def __init__(self, sgID, stationID, portID):
        self.ID = FullPort.buildID(sgID, stationID, portID)
        self.sgID = sgID
        self.stationID = stationID
        self.portID = portID

    def set_specs(self, address, full_load):
        self.address = address
        self.full_load = full_load
        self.allowed_load = full_load

    @staticmethod
    def fromID(IDstr):
        if ((IDstr is None) or (IDstr == "")):
            print("Invalid input provided to ChargePoint constructor")
            raise EVException("IDstr cannot be empty or None",
                              constants.ERR_INVALID_VALUE)
        IDlist = IDstr.split(constants.SEPARATOR)
        sgID = None
        stationID = None
        portID = None

        if constants.DEBUG:
            print(IDlist)
        if (len(IDlist) == 1):
            sgID = IDlist[0]
        else:  # len(IDlist) will be 3
            if (IDlist[0] != ""):
                sgID = IDlist[0]
            if (IDlist[1] != ""):
                stationID = IDlist[1]
            if (IDlist[2] != ""):
                portID = IDlist[2]
        if constants.DEBUG:
            print("id = " + IDstr + ", ")
            print("sgID = " + str(sgID) + ", ")
            print("stationID = " + str(stationID) + ", ")
            print("portID = " + str(portID) + "\n")
        return FullPort(sgID, stationID, portID)

    @staticmethod
    def buildID(sgID, stationID, portID):
        IDstr = ""
        sep = constants.SEPARATOR
        if ((sgID is None) and (stationID is None)):
            raise EVException("Specify at least Station Group or Station ID",
                              constants.ERR_INVALID_VALUE)
        if (sgID is not None):
            IDstr = str(sgID)
        if (stationID is not None):
            IDstr = IDstr + sep + str(stationID) + sep
        if (portID is not None):
            IDstr += str(portID)
        return IDstr

    def shed_amount(self, allowed_load):
        self.shed_state = True
        self.allowed_load = allowed_load
        self.shed_percent = None

    def shed(self, percent, amount):
        if ((percent is None) and
            (amount is None)):
            raise EVException("One of amount or percent must be specified", 303)
        self.shed_state = True
        if (amount is not None):
            self.allowed_load = amount
            self.shed_percent = None
        else: 
            self.allowed_load = self.full_load * percent * 0.01
            self.shed_percent = percent

    def get_shed_state(self):
        if self.shed:
            return 1
        else:
            return 0
        
    def get_allowed_load(self):
        return self.allowed_load

    def get_percent_shed(self):
        return self.percent_shed

    def clear(self):
        self.shed_state = False
        self.allowed_load = self.full_load
        self.shed_percent = 0

    def occupied(self):
        return (self.charge_session is not None)

    def available(self):
        return (self.charge_session is None)

    def get_charge_session(self):
        return self.charge_session

    def plug(self, vehicle):
        self.charge_session = ChargeSession(self.ID, vehicle)
        return self.charge_session

    def unplug(self):
        cs = self.charge_session
        self.charge_session = None
        self.load = 0.0
        return cs

    def set_load(self, load):
        self.load = load

    def get_load(self):
        return self.load

    def get_vehicleID(self):
        if (self.charge_session is not None):
            return self.charge_session.get_vehicle().ID
        else:
            return None

    def update_charge(self, now):
        if (self.charge_session is None):
            self.load = 0.0
        else:
            load = self.charge_session.update_charge(now, self.allowed_load)
            self.set_load(load)