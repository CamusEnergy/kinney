# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

# Facility type helps to determine importance with respect energy curtailment

from enum import Enum

class Curtail(Enum):
    BLACKOUT = 100
    BROWNOUT = 200
    GREEN = 300
