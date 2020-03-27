# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0


import time

from classes.charge_session import ChargeSession
from classes.full_port import FullPort
import constants
import cp
from util import get_ENV_val




sgID = get_ENV_val("sgID")
stationID = get_ENV_val("stationID")
portID = get_ENV_val("portID")

res0 = cp.get_load(sgID)

spID = FullPort(None, stationID, portID).ID
# Experimental phase, to not affect live users, cap at 6.6
res1 = cp.shed_load(spID, absolute_amount=6.6, time_interval=10)
print("Sleeping .. a  bit for shed_load to take effect")
time.sleep(15)
LOAD2 = cp.get_load(spID)

res2 = cp.clear_shed(spID)
print(res2)
print("Sleeping .. a  bit for clear_shed to take effect")
time.sleep(15)
res3 = cp.get_load(spID)
print("Sleeping .. a  bit for clear_shed for station_group")
time.sleep(10)
res4 = cp.clear_shed(sgID)
print(res4)
res5 = cp.get_load(spID)
print(res5)
