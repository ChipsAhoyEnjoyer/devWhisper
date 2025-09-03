import uuid, requests, sys

from models import EndpointTest

def test_login(t: EndpointTest, endpoint: str, username: str, password:str) -> tuple[dict, uuid.UUID]:
    print(f"Testing '{endpoint}' endpoint...")
    print(f"Paylod - Username = {username}, Password = {password}")
    payload = {"username": username, "password": password}

    try: # Connect to endpoint
        req = requests.post(t.baseurl + endpoint, json=payload)
    except Exception as e: # Failed to connect
        print(f"Error connecting to {t.baseurl}{endpoint}:\n", e)
        print(e)
        sys.exit(1)
    if req.status_code > 299: # Failed to login
        print(f"Error while login to endpoint '{endpoint}'")
        print(f"Status code: {req.status_code}")
        try:
            resp = req.json()
        except Exception as e:
            print(f"Error decoding json response")
            print(e)
            sys.exit(1)
        print(f"Server response: {resp}")
        sys.exit(1)
    # Successful login
    jwt = req.headers.get("Authorization")
    bearer = req.headers.get("Refresh-Token")
    if jwt == None or bearer == None:
        print(f"Login unsuccessful: No JWT and/or Bearer returned")
        print(f"jwt: {jwt}")
        print(f"bearer: {bearer}")
        sys.exit(1)
    try:
        resp = req.json()
    except Exception as e: # Failed to get server response/malformed json
            print(f"Error decoding json response")
            print(e)
            sys.exit(1)
    id = resp.get("id")
    if id == None:
        print("Logged in but no user ID returned")
        print(f"{endpoint} test failed")
        sys.exit(1)
    try:
        uid = uuid.UUID(id)
    except Exception as e:
        print(f"Invalid uuid returned from server: {id}")
        sys.exit(1)
        
    print(f"/login test passed - server response")
    print(f"id: {uid}")
    print(f"bearer: {bearer}")
    print(f"jwt: {jwt}")
    print(f"{resp}")
    print("")

    return {"Authorization": jwt, "Refresh-Token": bearer}, uid
