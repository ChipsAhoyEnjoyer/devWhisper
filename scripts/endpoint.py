import requests
import uuid
import sys

class EndpointTest:
    def __init__(
        self, 
        baseurl: str,
        header: dict[str, str] | None = None
    ):
        self.baseurl = baseurl
        self.header = header or {}

    def test_register(self, endpoint: str, username: str, password: str):
        print(f"Testing '{endpoint}' endpoint...")
        print(f"Payload - Username = {username}, Password = {password}")
        payload = {"username": username, "password": password}
        try: # Connect to endpoint
            req = requests.post(self.baseurl + endpoint, json=payload)
        except Exception as e: # Failed to connect
            print(f"Error connecting to {self.baseurl}{endpoint}:\n", e)
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

    def test_login(self, endpoint: str, username: str, password:str) -> tuple[dict, uuid.UUID]:
        print(f"Testing '{endpoint}' endpoint...")
        print(f"Paylod - Username = {username}, Password = {password}")
        payload = {"username": username, "password": password}

        try: # Connect to endpoint
            req = requests.post(self.baseurl + endpoint, json=payload)
        except Exception as e: # Failed to connect
            print(f"Error connecting to {self.baseurl}{endpoint}:\n", e)
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

    def test_delete_user(self, endpoint: str, id: uuid.UUID, headers: dict[str, str]):
        print(f"Testing '{endpoint}' endpoint...")
        print(f"Paylod - jwt = {headers.get('Authorization')}, bearer = {headers.get('Refresh-Token')}")
        payload = {"id": str(id)}

        try: # Connect to endpoint
            req = requests.delete(self.baseurl + endpoint, headers=headers, json=payload)
        except Exception as e: # Failed to connect
            print(f"Error connecting to {self.baseurl}{endpoint}:\n", e)
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


def main():
    print("")
    print("-- Endpoint Testing Starting --")
    user = {"username": "AdminTest", "password": "AdminPass"}
    tester = EndpointTest("http://localhost:7777")
    # Tests
    tester.test_register("/register", user["username"], user["password"])
    headers, id = tester.test_login("/login", user["username"], user["password"])
    tester.test_delete_user("/deleteUser", id, headers)

if __name__ == "__main__":
    main()
