# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

# Facility type helps to determine importance with respect energy curtailment
from enum import Enum

class Function(Enum):
    UNKNOWN = 100000
    SHELTER=100001
    LIGHT = 200000
    STREET_LIGHT = 200100
    ELEVATOR = 300000
    ESCALATOR = 300001
    CONVEYOR_BELT = 300003
    WASHER_CLOTHES = 400000
    DRYER_CLOTHES = 400001
    WASHER_DISHES = 400002
    # Draw depends on active plates and their temperature
    STOVE = 400003
    # Draw depends upon contents within
    KITCHEN_REFRIGERATOR = 400003
    # Draw depends on speed
    FAN = 500000
    # Draw depends upon ambient temperature
    AIR_CONDITIONER = 500001
    DATACENTER_COOLER = 500002
    HEATER = 500003
    TELEVISION = 600000
   

    INTERNET = 700000

    # generic term to cover servers, routers
    COMPUTER = 800000
    CELLPHONE_TOWER = 800100
    RADIO_TOWER = 800200
    TELEVISION_TOWER = 800300
    MICROWAVE_TOWER = 800400

    ELECTRIC_TRAIN = 900000
    TRAM_LINE = 900100
    ELECTRIC_VEHICLE_CHARGER = 900200

 






