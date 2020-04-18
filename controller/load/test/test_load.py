from classes.category import Category
from classes.curtail import Curtail
from classes.facility import Facility
from classes.function import Function
from classes.load import BaseLoad



aLight = BaseLoad(1, 1, safety_critical=True, facility = Facility.HOME, category = Category.COMFORT, function=Function.LIGHT)
bLight = BaseLoad(2, 1, facility = Facility.HOME, category = Category.COMFORT, function=Function.LIGHT)
tv = BaseLoad(3, 100, facility = Facility.HOME, category = Category.ENTERTAINMENT, function=Function.TELEVISION)

home = BaseLoad(4, 102, safety_critical=True, facility = Facility.HOME, category = Category.COMFORT, function=Function.SHELTER)
home.children={aLight, bLight, tv}

print(home)
print(home.get_load_need())