import requests, sys

from models import EndpointTest

def test_register(t: EndpointTest, endpoint: str, username: str, password: str):
    print(f"Testing '{endpoint}' endpoint...")
    print(f"Payload - Username = {username}, Password = {password}")
    payload = {"username": username, "password": password}
    try: # Connect to endpoint
        req = requests.post(t.baseurl + endpoint, json=payload)
    except Exception as e: # Failed to connect
        print(f"Error connecting to {t.baseurl}{endpoint}:\n", e)
        print(e)
        sys.exit(1)
    if req.status_code > 299: # Failed to register
        print(f"Error while registering to endpoint '{endpoint}'")
        print(f"Status code: {req.status_code}")
        try:
            resp = req.json()
        except Exception as e:
            print(f"Error decoding json response")
            print(e)
            sys.exit(1)
        print(f"Server response: {resp}")
        sys.exit(1)
    # Succesful register
    cont = req.headers.get("Content-Type")
    if cont == None:
        print(f"WARNING: no content-type header")
    try:
        resp = req.json()
    except Exception as e: # Failed to get server response/malformed json
            print(f"Error decoding json response")
            print(e)
            sys.exit(1)
    print(f"{endpoint} test passed - server response:\n    {resp}")
    print("")
