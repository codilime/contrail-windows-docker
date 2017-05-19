import sys
from contrail_vrouter_api.vrouter_api import ContrailVRouterApi

PORT_TYPE = 'NovaVMPort'

def main():
    operation = sys.argv[0]

    # TODO: ContrailVRouterApi has semaphores, but we call it like a batch command.
    # Therefore, we need to make sure that proper synchronization is realized in driver, so that
    # add/delete operations don't overlap.
    api = ContrailVRouterApi()

    if operation == "add":
        try:
            (operation, vmUuid, vifUuid, ifName, mac, dockerID) = sys.argv
            api.add_port(vmUuid, vifUuid, ifName, mac, port_type=PORT_TYPE,
                display_name=args.dockerId)
        except Exception:
            print("Mock no failure, since agent is not implemented yet.")
    elif operation == "delete":
        try:
            (operation, vmUuid) = sys.argv
            api.delete_port(vifUuid)
        except Exception:
            print("Mock no failure, since agent is not implemented yet.")
    else:
        raise ValueError("Invalid operation. Must be add or delete")


if __name__ == "__main__":
    main()