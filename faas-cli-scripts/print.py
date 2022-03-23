import sys
import time

def handler(x):
    print(x)
    time.sleep(4)
    print("done sleeping")
handler(sys.argv[1])
