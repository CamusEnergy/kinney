# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0


class EVException(Exception):
    err_code = None
    def __init__(self, message, code):

        # Call the base class constructor with the parameters it needs
        super(EVException, self).__init__(message)

        # Now for your custom code...
        self.err_code = code
