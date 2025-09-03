from register import test_register
from login import test_login
from delete import test_delete_user
from models import EndpointTest


def main():
    print("")
    print("-- Endpoint Testing Starting --")
    user = {"username": "AdminTest", "password": "AdminPass"}
    tester = EndpointTest("http://localhost:7777")
    # Tests
    test_register(tester, "/register", user["username"], user["password"])
    tester.header, id = test_login(tester, "/login", user["username"], user["password"])
    test_delete_user(tester, "/deleteUser", id)

if __name__ == "__main__":
    main()
