# iximuz-labs

This repository contains a collection of labs and experiments.

## Labs

*   **[L3 to L2 Networking](./labs/l3-to-l2/notes.md):** This lab explores the creation of Layer 2 networks using Linux networking primitives like network namespaces and virtual Ethernet (veth) pairs. It covers two main scenarios:
    *   **Simple L2 Bridge:** Demonstrates how to create a basic bridge connecting multiple hosts in the same broadcast domain.
    *   **VLAN-Aware Bridge:** Shows how to implement network segmentation using a single bridge and VLANs, effectively creating multiple isolated L2 networks.

    **Pain Points:** The lab highlights the complexity and manual effort required to configure these networks using `iproute2` commands. Managing namespaces, veth pairs, bridge settings, and VLAN configurations can be error-prone and difficult to scale without significant automation.

*   **[Working with Containers](./labs/containers/notes.md):** This lab provides practical notes on managing container images using `ctr`. It covers essential commands for listing, importing, mounting, and tagging images. Additionally, it delves into troubleshooting common Kubernetes issues, such as investigating killed pods by inspecting `kubelet` logs and using `crictl` to examine container states.



### Courses: How to install and Configure containerd on a linux Server

- The containerd release archive includes only the essentials components:
- ctr the command-line client
- containerd the daemon itself
- containerd-shim-runc-v2 an OCI container
