import os
import json
import base64
import subprocess

os.system("faas-cli store deploy figlet")

with open("print.py", "r") as f:
    src_code = f.read()
    src_code = json.dumps(src_code)

payload = {
    "\"fid\"": "\"test\"",
    "\"src\"": "\"" + base64.b64encode(src_code.encode('ascii')).decode('ascii') + "\"",
    "\"params\"": "\"HELLLOOOO\"",
    "\"lang\"": "\"micropython\""
}


payload_dumps = (json.dumps(payload))

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


