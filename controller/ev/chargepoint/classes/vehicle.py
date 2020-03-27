# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0


class Vehicle():
    '''A representation of an electric vehicle
    attributes:
        ID: unique identifier, constructor only requires this field
        model: vehicle model, which would characterize its charge rate, capacity etc
        capacity: its full charge capacity
        current_charge: current charge
        charge_rate
        trickle_charge_rate
        trickle_point: often times this is around 80% of its full charge capacity
    '''
    ID: None
    model: "Unknown"
    capacity: None
    current_charge: 0.0
    # vehicle dependant
    charge_rate: None
    trickle_charge_rate = None
    trickle_charge_point: None

    def __init__(self, ID):
        self.ID = ID

    def set_specs(self, model, capacity,
                 charge_rate, trickle_charge_rate,
                 trickle_percent):
        self.model = model
        self.capacity = capacity
        self.charge_rate = charge_rate
        self.trickle_charge_rate = trickle_charge_rate
        self.tickle_charge_point = trickle_percent * capacity

    def set_charge(self, charge):
        self.current_charge = charge

    def get_load(self):
        if (self.current_charge == self.capacity):
            return 0.0
        else:
            if (self.current_charge < self.trickle_charge_point):
                return self.charge_rate
            else:
                return self.trickle_charge_rate
