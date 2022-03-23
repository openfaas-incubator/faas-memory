import json
import base64

with open("print.py", "r") as f:
        src_code = f.read()
        src_code = json.dumps(src_code)

print(base64.b64encode(src_code.encode('ascii')).decode('ascii'))