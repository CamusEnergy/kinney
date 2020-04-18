# Copyright 2020 program was created VMware, Inc.
# SPDX-License-Identifier: Apache-2.0

import abc

from classes.category import Category
from classes.curtail import Curtail
from classes.facility import Facility
from classes.function import Function

class Context():
    temperature = 70
    pandemic = False
    sports_event = False
    festive_season = False

class LoadInterface(metaclass=abc.ABCMeta):
    @classmethod
    def __subclasshook__(cls, subclass):
        return (hasattr(subclass, 'get_Load') and 
                callable(subclass.get_load) and 
                hasattr(subclass, 'get_load_need') and 
                callable(subclass.get_load_need) or
                hasattr(subclass, 'set_curtail_load') and 
                callable(subclass.set_curtail_load) or
                hasattr(subclass, 'get_curtail_load') and 
                callable(subclass.get_curtail_load) or
                hasattr(subclass, 'get_criticality') and 
                callable(subclass.get_criticality) or
                NotImplemented)

class BaseLoad():
    '''Base class representing an electrical load
    attributes:
        ID: unique identifier, constructor only requires this field
        name: Name of the load
        description: human readable, kitchen light or swimming pool pump
        address: address
        GPS_coords: GPS co-ordinates
        facility: one of school, hospital, retail, office, factory ..
        category: purpose the load serves: comfort, entertainment, communication
        function: stove, air conditioner, computer, ev_charger
        deferrable: boolean, can we delay running this function without any loss
        active: is this load active
        max_load: 0.0  this might be a cap from supply, pre-arranged limit ..
        curtail_load: 0.0
        load: current power usage
        safety_critical: is it important/essential for safety reasons? example light in a stairwell
        labels: any additional fixed string descriptors
        children: list of one or more LoadEndPoints. Parent load equals sum of children load   
    '''
    ID = None
    name = None
    address = None
    GPS_coords = None
    facility = Facility.UNKNOWN
    category = Category.UNKNOWN
    function = Function.UNKNOWN
    deferrable = False
    active = False
    max_load = 0.0
    curtail_load = 0.0
    load = 0.0
    safety_critical = False
    descriptors = list()
    sub_loads = list()

    def __init__(self,
                 id,
                 max_load,
                 safety_critical=False,
                 facility=Facility.UNKNOWN,
                 category=Category.UNKNOWN,
                 function=Function.UNKNOWN,
                 name="NameMe",
                 description=""):
        self.ID = id
        self.name = name
        self.description = description
        self.max_load = max_load
        self.safety_critical = safety_critical
        self.function = function
        self.category = category
        self.facility = facility
        self.descriptors = list()
        self.sub_loads = list()

    def set_curtail_load(self, absolute_amount = 0.0, percentage_amount = 0.0):
        if (absolute_amount == 0.0) and (percentage_amount == 0.0): 
            self.curtail_load = 0.0 # equivalent to clearing curtailment
        elif ((absolute_amount > 0.0) and (absolute_amount < self.max_load)):
            self.curtail_load = absolute_amount
        else:
            self.curtail_load = percentage_amount * 0.01 * self.max_load
    
    def get_curtail_load(self):
        return self.curtail_load
   
   # actual dynamic consumption, current load
    def get_load(self):
        if (len(self.sub_loads) == 0):
            return self.load
        else:
            load_sum = 0
            for sb in self.sub_loads:
                load_sum = load_sum + sb.get_load()
            return load_sum

    def add_sub_load(self, aLoadObject):
        self.sub_loads.append(aLoadObject)
        
    def get_load_need(self, level=Curtail.GREEN, context=None):
        # TODO consider context such as temperature, pandemic etc
        if (len(self.sub_loads) == 0):
            #leaf load device
            load_need = self.load
            if ((level == Curtail.GREEN) or
                self.safety_critical or
                (level == Curtail.BROWNOUT) and (not self.deferrable)):
                return self.load
            else:
                return 0.0    
        else:  
            load_need = 0.0
            for sb in self.sub_loads:
                load_need = load_need + sb.get_load_need(level=level, context=context)
            return load_need
