from enum import Enum
import os
import json
import base64
import subprocess
import threading
import sys

NUM_JOBS = 15

# Worker class
class Worker:
    def __init__(self, id, ip, status):
        self.id = id
        self.ip = ip
        self.status = status
        self.queue = []

# Enumerated statuses for each worker
class Status(Enum):
    READY = 1
    RUNNING = 2
    POWEROFF = 3

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

if __name__ == "__main__":
    os.system("faas-cli store deploy figlet")

    # Worker objects for each worker
    all_workers = [
        Worker(1, "192.168.1.20", Status.READY),
        Worker(2, "192.168.1.21", Status.READY),
        Worker(3, "192.168.1.22", Status.READY),

    ]

    # Open python script with function to run on a worker and make it into a JSON string
    with open("print.py", "r") as f:
        src_code = f.read()
        src_code = json.dumps(src_code)

    threads = []

    queued_count = 0
    completed_count = 0
    while True:
        are_queues_empty = False
        # Loop through all the worker objects
        for worker in all_workers:
            if (queued_count < NUM_JOBS):
                # Package an encrypted JSON string into a payload JSON 
                payload = {
                    "\"fid\"": "\"test\"",
                    "\"src\"": "\"" + base64.b64encode(src_code.encode('ascii')).decode('ascii') + "\"",
                    "\"params\"": f"\"HELLLOOOO {worker.id}\"",
                    "\"lang\"": "\"micropython\"",
                    "\"worker\"": f"\"http://{worker.ip}:8080/run\""
                }
                payload_dumps = (json.dumps(payload))
                print(payload_dumps)
                worker.queue.append(payload_dumps)
                print("added job to worker ", worker.id, "and it is now length", len(worker.queue))
                queued_count += 1
                print(queued_count)

            # If the worker is ready and the queue is not empty, run the first job on the queue
            if (worker.status == Status.READY and worker.queue != []):
                # Set the current worker's status to running (set back to running after reboot)
                # worker.status = Status.RUNNING
                
                # Add the job to the thread and run it
                t = threading.Thread(target=runJob, args=(worker.queue[0],))
                threads.append(t)
                t.start()
                print("job running on worker", worker.id)
                # Pop the first payload from the queue
                worker.queue.pop(0)
                completed_count += 1
                
                # Once all the jobs have been completed
                if completed_count >= NUM_JOBS:
                    sys.exit()
        



