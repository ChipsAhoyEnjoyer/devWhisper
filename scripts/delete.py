import requests, sys, uuid

from models import EndpointTest

def test_delete_user(t: EndpointTest, endpoint: str, id: uuid.UUID):
    print(f"Testing '{endpoint}' endpoint...")
    print(f"Paylod - jwt = {t.header.get('Authorization')}, bearer = {t.header.get('Refresh-Token')}")
    payload = {"id": str(id)}

    try: # Connect to endpoint
        req = requests.delete(t.baseurl + endpoint, headers=t.header, json=payload)
    except Exception as e: # Failed to connect
        print(f"Error connecting to {t.baseurl}{endpoint}:\n")
        print(e)
        sys.exit(1)
    if req.status_code > 299: # Failed to delete
        print(f"Error while delete at endpoint '{endpoint}'")
        print(f"Status code: {req.status_code}")
        try:
            resp = req.json()
        except Exception as e:
            print(f"Error decoding json response")
            print(e)
            sys.exit(1)
        print(f"Server response: {resp}")
    # Delete successful
    print(f"{endpoint} test passed - No content expected")
    print("")
