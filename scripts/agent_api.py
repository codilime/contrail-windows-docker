# For manual testing you can uncomment this and run a 
# python -m SimpleHTTPServer 9090 in another terminal
# the tests should pass.
# import socket
# print "Hello"
# s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# s.connect(("127.0.0.1", 9090))
# s.send("Test msg pls ignore")
# s.close()

import sys
from contrail_vrouter_api.vrouter_api import ContrailVRouterApi

PORT_TYPE = 'NovaVMPort'


def main():
    operation = sys.argv[0]

    # TODO: ContrailVRouterApi has semaphores, but we call it like a batch command.
    # Therefore, we need to make sure that proper synchronization is realized in driver, so that
    # add/delete operations don't overlap.
    api = ContrailVRouterApi()

    print("MOCK! Always returns success. TODO when Agent is ready.")

    if operation == "add":
        try:
            (operation, vmUuid, vifUuid, ifName, mac, dockerID) = sys.argv
            api.add_port(vmUuid, vifUuid, ifName, mac, port_type=PORT_TYPE,
                display_name=args.dockerId)
        except Exception:
            print("MOCK: not reporting FAIL, since agent is not implemented yet.")
    elif operation == "delete":
        try:
            (operation, vmUuid) = sys.argv
            api.delete_port(vifUuid)
        except Exception:
            print("MOCK not reporting FAIL, since agent is not implemented yet.")
    else:
        raise ValueError("Invalid operation. Must be add or delete")


if __name__ == "__main__":
    main()