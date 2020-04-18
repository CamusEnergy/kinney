# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

# Facility type helps to determine importance with respect energy curtailment

from enum import Enum

class Facility(Enum):
    UNKNOWN = 1000
    HOME = 2000
    HOME_KITCHEN = 2001
    SCHOOL = 3000
    DATACENTER = 4000
    TELECOMMUNICATIONS = 4100
    OFFICE = 4200
    FACTORY = 4300
    RETAIL = 4400
    PARKING_LOT = 4500
    ROAD_SIDE = 4600
    MOBILE = 5000
    HOSPITAL = 6000
    TRANSPORT = 7000
    
