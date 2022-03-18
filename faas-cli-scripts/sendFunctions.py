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


if __name__ == "__main__":
    os.system("faas-cli store deploy figlet")

    # IP Addresses for all workers
    all_workers = {
        1: "192.168.1.20",
        2: "192.168.1.21",
        3: "192.168.1.22"
    }

    # Open python script with function to run on a worker and make it into a JSON tring
    with open("print.py", "r") as f:
        src_code = f.read()
        src_code = json.dumps(src_code)

    threads = []
    for i in range(1,4):
        # Package an encrypted JSON string into a payload JSON 
        payload = {
            "\"fid\"": "\"test\"",
            "\"src\"": "\"" + base64.b64encode(src_code.encode('ascii')).decode('ascii') + "\"",
            "\"params\"": f"\"HELLLOOOO {i}\"",
            "\"lang\"": "\"micropython\"",
            "\"worker\"": f"\"http://{all_workers[i]}:8080/run\""
        }
        payload_dumps = (json.dumps(payload))
        print(payload_dumps + "\n")

        t = threading.Thread(target=runJob, args=(payload_dumps,))
        threads.append(t)
        t.start()


