from enum import Enum
import os
import json
import base64
import subprocess
import threading


def runJob(payload_dumps):
    # Call the payload using faas-cli invoke and send to provider
    res = subprocess.check_output(f"echo {payload_dumps} | faas-cli invoke figlet", shell=True)

    # Take the JSON string of result and transform it into a JSON object and index the response body
    response = json.loads(json.loads(res)["ResponseBody"])

    # Check for error in response object
    try:
        err = response["error"] 
        if err:
            print("ERROR HERE \n")
    except KeyError:
        # Isolate the result from res to find the output of faas-cli invoke
        result = response["result"]
        print("RESULT IS: " + result + "\n")
    return

# Worker class
class Worker:
    def __init__(self, id, ip, status):
        self.id = id
        self.ip = ip
        self.status = status

# Enumerated statuses for each worker
class Status(Enum):
    READY = 1
    RUNNING = 2
    POWEROFF = 3

if __name__ == "__main__":
    os.system("faas-cli store deploy figlet")

    # Worker objects for each worker
    all_workers = [
        Worker(1, "192.168.1.20", Status.READY),
        Worker(2, "192.168.1.21", Status.READY),
        Worker(3, "192.168.1.22", Status.READY),

    ]

    # Open python script with function to run on a worker and make it into a JSON tring
    with open("print.py", "r") as f:
        src_code = f.read()
        src_code = json.dumps(src_code)

    threads = []
    for i in range(3):
        print(i)
        # Package an encrypted JSON string into a payload JSON 
        payload = {
            "\"fid\"": "\"test\"",
            "\"src\"": "\"" + base64.b64encode(src_code.encode('ascii')).decode('ascii') + "\"",
            "\"params\"": f"\"HELLLOOOO {i}\"",
            "\"lang\"": "\"micropython\"",
            "\"worker\"": f"\"http://{all_workers[i].ip}:8080/run\""
        }
        payload_dumps = (json.dumps(payload))
        print(payload_dumps + "\n")

        t = threading.Thread(target=runJob, args=(payload_dumps,))
        threads.append(t)
        t.start()


