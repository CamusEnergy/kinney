from classes.category import Category
from classes.curtail import Curtail
from classes.facility import Facility
from classes.function import Function
from classes.load import BaseLoad



aLight = BaseLoad(1, 1, safety_critical=True, facility = Facility.HOME, category = Category.COMFORT, function=Function.LIGHT)
aLight.load = 1
bLight = BaseLoad(2, 1, facility = Facility.HOME, category = Category.COMFORT, function=Function.LIGHT)
bLight.load = 1
tv = BaseLoad(3, 100, facility = Facility.HOME, category = Category.ENTERTAINMENT, function=Function.TELEVISION)
tv.load = 40

home = BaseLoad(4, 400, safety_critical=True, facility = Facility.HOME, category = Category.COMFORT, function=Function.SHELTER)
home.add_sub_load(aLight)
home.add_sub_load(bLight)
home.add_sub_load(tv)

print(home)
print("Green")
print(home.get_load_need())
print("============")
print("Blackout")
print(home.get_load_need(level=Curtail.BLACKOUT))