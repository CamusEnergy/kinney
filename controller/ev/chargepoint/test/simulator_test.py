# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

import time

from classes.simulator import CPInstance
from classes.vehicle import Vehicle

cpi = CPInstance("12345", 5, "Somewhere over the rainbow", 75, "mall")
v = Vehicle("Goofy")
v.setSpecs("Prius", 100, 6, 1, 80)
v.setCharge(50)
free_port = cpi.get_free_port()
if (free_port is not None):
    cpi.plugin(free_port.ID, v)
    time.sleep(10)
    print("Slept 10 seconds")
    print(cpi.get_load())
    time.sleep(120)
    print("Slept 120 seconds")
    print(cpi.get_load())