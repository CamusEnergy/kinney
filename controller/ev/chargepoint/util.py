# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

import constants
import os
import time


def get_ENV_val(str):
    if (str in os.environ):
        return os.environ[str]
    else:
        errMsg = "Missing environment variable: " + str
        raise EVException(errMsg, constants.ERR_ENV_VAR_MISSING)
    