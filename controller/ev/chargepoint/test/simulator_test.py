# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

import time

from classes.simulator import CPInstance
from classes.vehicle import Vehicle

cpi = CPInstance("12345", 5, "Somewhere over the rainbow", 100, "mall")

# length test
if (cpi.num_ports == len(cpi.ports)):
    print("Pass: Simulator num ports test")
else:
    print("Fail: Simulator num ports test")

if (cpi.get_load() == 0.0):
    print("Pass: Simulator no load test")
else:
    print("Fail: Simulator no load test")

if (cpi.get_num_charge_sessions() == 0):
    print("Pass: Simulator no charge sessions test")
else:
    print("Fail: Simulator no charge sessions test")

v = Vehicle("Goofy")
v.set_specs("Prius", 100, 6, 1, 80)
v.set_charge(50.0)
free_port = cpi.get_free_port()
if (free_port is not None):
    print("Pass: Simulator free port test")
else:
    print("Fail: Simulator free port test")

plugin_port = cpi.plugin(v)
if (plugin_port is not None):
    print("Pass: Simulator vehicle plugin test")
else:
    print("Fail: Simulator vehicle plugin test")

if (cpi.get_num_charge_sessions() == 1):
    print("Pass: Simulator 1 charge session test")
else:
    print("Fail: Simulator 1 charge session test")

time.sleep(10)
print("Slept 10 seconds, less than a minute")
if (cpi.get_load() == 0.0):
    print("Pass: Simulator inadequate charge time test")
else:
    print("Fail: Simulator inadequate charge time test")

time.sleep(120)
load = cpi.get_load()
if (load > 0.0):
    print("Pass: Simulator adequate charge time test")
else:
    print("Fail: Simulator adequate charge time test")
