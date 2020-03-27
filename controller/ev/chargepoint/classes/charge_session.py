# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

import time

class ChargeSession():
    full_port_ID = None
    vehicle = None
    start = None
    end = None
    last_update = None
    total_charge = None

    def __init__(self, vehicle, full_port_ID):
        print("full_port_ID: " + full_port_ID)
        print("vehicle = ")
        print(vehicle)
        self.full_port_ID = full_port_ID
        self.vehicle = vehicle
        self.total_charge = 0.0
        now = time.time()
        self.start = now
        self.last_update = now

    def get_vehicle(self):
        return self.vehicle

    # TODO
    def update_charge(self, timestamp, load_KWhr):
        ret_load = 0.0
        interval_mins = (timestamp - self.last_update)/60
        print("elapsed_time = " + str(interval_mins))
        self.last_update = timestamp
        print("Enter compute_charge ts: " + str(timestamp)
              + " insta_load: " + str(load_KWhr))
        vehicle = self.get_vehicle()
        print("**********")
        print(vehicle)
        
        # approximate calculations here .. not taking care of trickle

        curr_charge = vehicle.current_charge
        add_charge = interval_mins * load_KWhr / 60  # hours/minutes
        if (vehicle.capacity < (curr_charge + add_charge)):
            add_charge = vehicle.capacity - curr_charge
            vehicle.set_charge(vehicle.capacity)
            ret_load = 0.0
        else:
            vehicle.set_charge(curr_charge + add_charge)
            ret_load = load_KWhr
        self.total_charge = self.total_charge + add_charge
        return ret_load
     



class ChargeSessions:
    sessions = None

    def __init__(self):
        self.sessions = {}

    def get_start(self, vehicleID, watt, now):
        start = now
        if (vehicleID in self.sessions.keys()):
            start = self.sessions[vehicleID]
            if (watt == 0):
                # charge session over
                del self.sessions[vehicleID]
        else:
            self.sessions[vehicleID] = start
        return start
       