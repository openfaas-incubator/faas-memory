import os
import json
import base64

os.system("faas-cli store deploy figlet")

with open("shutdown.py", "r") as f:
    src_code = f.read()
    src_code = json.dumps(src_code)

print("SRCCODE", src_code)

payload = {
    "\"fid\"": "\"test\"",
    "\"src\"": "\"" + base64.b64encode(src_code.encode('ascii')).decode('ascii') + "\"",
    "\"params\"": "\"HELLLOOOO\"",
    "\"lang\"": "\"micropython\""
}

# print("HERE",payload)
payload_dumps = (json.dumps(payload))
print(payload_dumps + "\n")
os.system(f"echo {payload_dumps} | faas-cli invoke figlet")


# print(json.loads(dumpy)["src"])