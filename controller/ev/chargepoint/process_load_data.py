from ev import ChargeSession, FullPort
import json

charge_sessions = {}
DATA_KEY = '   data'


def process_file(full_filename):
    global data
    num_lines = 0
    with open(full_filename) as f:
        while (True):
            line = f.readline()
            if not line:
                break
        # data = json.load(f) 
            num_lines = num_lines + 1
            print(str(line))
            print(len(line))
            print("----")
            if (line != "\n"):
                process_line(line)
    print("Num lines read = " + str(num_lines))

foo = None
def process_line(line):
    global DATA_KEY, charge_sessions, foo
    if (line != "\n"):
        data = json.loads(line)
        timestamp = data['ts']
        sgID = data[DATA_KEY]['sgID']
        for sd in data[DATA_KEY]['stationData']:
            stationID = sd['stationID']
            stationAddress = sd['Address']
            for p in sd['Port']:
                portLoad = p['portLoad']
                portID = p['portNumber']
                if (portLoad > 0.0):
                    credentialID = p['credentialID']
                    if credentialID in charge_sessions:
                        cs = charge_sessions[credentialID]
                        cs.update_charge(timestamp, portLoad)
                        print("Continuing charge session vehicleID : "
                              + credentialID
                              + " load : " + str(portLoad))
                    else:
                        full_port = FullPort(sgID, stationID, portID,
                                             stationAddress, portLoad)
                        cp = ChargeSession(full_port.ID, credentialID)
                        cp.update_charge(timestamp, portLoad)
                        charge_sessions[credentialID] = cp
                        foo = cp
                        print("New Charge Session detected vehicleID : "
                              + credentialID
                              + " load : " + str(portLoad))
                else: 
                    print("No vehicle charging at stationID: "
                          + str(stationID) + " portID: "
                          + str(portID))



process_file(
    "/Users/mbhandaru/kinney/controller/ev/chargepoint/ss2.json"
    )
for cs in charge_sessions:
    print(cs)
    print("-----")

print ("**********")
print(foo)