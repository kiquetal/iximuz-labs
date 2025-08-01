# iximuz-labs

This repository contains a collection of labs and experiments.

## Labs

*   **L3 to L2 Networking:** This lab explores the creation of Layer 2 networks using Linux networking primitives like network namespaces and virtual Ethernet (veth) pairs. It covers two main scenarios:
    *   **Simple L2 Bridge:** Demonstrates how to create a basic bridge connecting multiple hosts in the same broadcast domain.
    *   **VLAN-Aware Bridge:** Shows how to implement network segmentation using a single bridge and VLANs, effectively creating multiple isolated L2 networks.

    **Pain Points:** The lab highlights the complexity and manual effort required to configure these networks using `iproute2` commands. Managing namespaces, veth pairs, bridge settings, and VLAN configurations can be error-prone and difficult to scale without significant automation.
