#### Bash functions (Simple L2 Bridge)

create_bridge_simple() {

  local nsname="$1"
  local ifname="$2"

  echo "Creating bridge ${nsname}/${ifname}"

  ip netns add ${nsname}
  ip netns exec ${nsname} ip link set lo up
  ip netns exec ${nsname} ip link add ${ifname} type bridge
  ip netns exec ${nsname} ip link set ${ifname} up
}

create_end_host_simple() {
  local host_nsname="$1"
  local peer1_ifname="$2"
  local peer2_ifname="$2b"
  local peer1_ifaddr="$3"
  local bridge_nsname="$4"
  local bridge_ifname="$5"

  echo "Creating end host ${host_nsname} ${peer1_ifaddr} connected to ${bridge_nsname}/${bridge_ifname} bridge"

  # Create end host network namespace.
  ip netns add ${host_nsname}
  ip netns exec ${host_nsname} ip link set lo up

  # Create a veth pair connecting end host and bridge namespaces.
  ip link add ${peer1_ifname} netns ${host_nsname} type veth peer \
              ${peer2_ifname} netns ${bridge_nsname}
  ip netns exec ${host_nsname} ip link set ${peer1_ifname} up
  ip netns exec ${bridge_nsname} ip link set ${peer2_ifname} up

  # Setting host's IP address.
  ip netns exec ${host_nsname} ip addr add ${peer1_ifaddr} dev ${peer1_ifname}

  # Attach peer2 interface to the bridge.
  ip netns exec ${bridge_nsname} ip link set ${peer2_ifname} master ${bridge_ifname}
}


---

### L2 Network Diagram (Single Bridge)

This diagram illustrates a simple Layer 2 network. A single bridge, `br`, is created. Three hosts (`host10`, `host11`, `host12`), each in their own namespace, are then connected to this bridge. Because all hosts are connected to the same bridge and assigned IP addresses in the same subnet (`192.168.0.0/24`), they can all communicate directly with each other.

```ascii
                    +---------------------------------+
                    |      Bridge Namespace (br2)     |
                    |                                 |
                    | +--------------+  +----------+  |
                    | |  Bridge: br  |  |  lo (UP) |  |
                    | +--------------+  +----------+  |
                    +---------------------------------+
                               |      |      |
                 +-------------+      |      +-------------+
                 | (veth)             | (veth)             | (veth)
  +--------------+------------+  +--------------+------------+  +--------------+------------+
  |      Host (host10)      |  |      Host (host11)      |  |      Host (host12)      |
  |                         |  |                         |  |                         |
  | eth10: 192.168.0.10/24  |  | eth11: 192.168.0.11/24  |  | eth12: 192.168.0.12/24  |
  | lo (UP)                 |  | lo (UP)                 |  | lo (UP)                 |
  +-------------------------+  +-------------------------+  +-------------------------+
```

---

### L2 Network with VLANs

Here are the specific functions used to create a VLAN-aware bridge and connect hosts to it on specific VLANs.

#### Bash functions (VLAN-Aware L2 Bridge)

create_bridge_vlan() {
  local nsname="$1"
  local ifname="$2"

  echo "Creating bridge ${nsname}/${ifname}"

  ip netns add ${nsname}
  ip netns exec ${nsname} ip link set lo up
  ip netns exec ${nsname} ip link add ${ifname} type bridge
  ip netns exec ${nsname} ip link set ${ifname} up

  # Enable VLAN filtering on bridge.
  ip netns exec ${nsname} ip link set ${ifname} type bridge vlan_filtering 1
}

create_end_host_vlan() {
  local host_nsname="$1"
  local peer1_ifname="$2"
  local peer2_ifname="$2b"
  local vlan_vid="$3"
  local bridge_nsname="$4"
  local bridge_ifname="$5"

  echo "Creating end host ${host_nsname} connected to ${bridge_nsname}/${bridge_ifname} bridge (VLAN ${vlan_vid})"

  # Create end host network namespace.
  ip netns add ${host_nsname}
  ip netns exec ${host_nsname} ip link set lo up

  # Create a veth pair connecting end host and bridge namespaces.
  ip link add ${peer1_ifname} netns ${host_nsname} type veth peer \
              ${peer2_ifname} netns ${bridge_nsname}
  ip netns exec ${host_nsname} ip link set ${peer1_ifname} up
  ip netns exec ${bridge_nsname} ip link set ${peer2_ifname} up

  # Attach peer2 interface to the bridge.
  ip netns exec ${bridge_nsname} ip link set ${peer2_ifname} master ${bridge_ifname}

  # Put host into right VLAN
  ip netns exec ${bridge_nsname} bridge vlan del dev ${peer2_ifname} vid 1
  ip netns exec ${bridge_nsname} bridge vlan add dev ${peer2_ifname} vid ${vlan_vid} pvid ${vlan_vid}
}

This setup demonstrates how a single Linux bridge can be used to create multiple isolated Layer 2 networks. The `vlan_filtering 1` command turns the bridge into a VLAN-aware switch. The `create_end_host_vlan` function then assigns each host's connection to a specific VLAN using the `bridge vlan add` command. The `pvid` (Port VLAN ID) setting is key, as it tags all untagged traffic from the host with the correct VLAN ID, allowing the hosts to be unaware of the VLAN configuration.

This creates two separate broadcast domains on the same bridge. Hosts in VLAN 10 can only communicate with other hosts in VLAN 10, and hosts in VLAN 20 can only communicate with other hosts in VLAN 20.

#### Diagram: L2 Network with VLANs

```ascii
                    +---------------------------------------------------+
                    |              Bridge Namespace (bridge1)           |
                    |                                                   |
                    | +---------------------------------------------+   |
                    | |         VLAN-Aware Bridge: br1              |   |
                    | |                                             |   |
                    | |  +-----------+      +-----------+           |   |
                    | |  | VLAN 10   |      | VLAN 20   |           |   |
                    | |  +-----------+      +-----------+           |   |
                    | |      | | |              | | |               |   |
                    | +------|-|-|--------------|-|-|---------------+
                    +--------|-|-|--------------|-|-|-------------------+
                             | | |              | | |
         (veth) -------------+ | +------------  | | |
         (veth) --------------+--------------  | | |
         (veth) -----------------------------  | | |
                                               | | |
         (veth) -------------------------------+\ | +-------------
         (veth) ---------------------------------+---------------
         (veth) -------------------------------------------------

  +-------------------------+  +-------------------------+  +-------------------------+
  |      Host (host10)      |  |      Host (host11)      |  |      Host (host12)      |
  |                         |  |                         |  |                         |
  | eth10: 192.168.10.10/24 |  | eth11: 192.168.10.11/24 |  | eth12: 192.168.10.12/24 |
  +-------------------------+  +-------------------------+  +-------------------------+

  +-------------------------+  +-------------------------+  +-------------------------+
  |      Host (host20)      |  |      Host (host21)      |  |      Host (host22)      |
  |                         |  |                         |  |                         |
  | eth20: 192.168.20.20/24 |  | eth21: 192.168.20.21/24 |  | eth22: 192.168.20.22/24 |
  +-------------------------+  +-------------------------+  +-------------------------+
```
*(Note: The IP addresses are assumed based on a common convention where the third octet matches the VLAN ID.)*

### Command Deep Dive: `bridge vlan add`

Let's break down this command with an example. Assume we have the following:
- A bridge in a namespace called `br-namespace`
- A veth interface `veth-host10` that connects a host to the bridge.
- We want to assign this host to VLAN `100`.

The command would be:
```bash
ip netns exec br-namespace bridge vlan add dev veth-host10 vid 100 pvid 100
```

Here is what each part does:

1.  **`ip netns exec br-namespace`**: This tells the system to run the following command *inside* the network namespace named `br-namespace`. This is where our bridge device lives.

2.  **`bridge vlan add`**: This is the command to add a VLAN configuration to a port on a bridge.

3.  **`dev veth-host10`**: This specifies the network interface (`dev`) that we are configuring. In this case, it's the `veth-host10` interface, which acts as the bridge port for our host.

4.  **`vid 100`**: This sets the VLAN ID (`vid`). It configures the port to accept traffic that is tagged with VLAN ID `100`.

5.  **`pvid 100`**: This sets the Port VLAN ID (`pvid`). This is the crucial part for making the host "unaware" of the VLAN. Any traffic that arrives on the `veth-host10` port *without* a VLAN tag will be automatically tagged with VLAN ID `100`.

In short, this single command configures the `veth-host10` port on the bridge to be an **access port** for VLAN `100`. It ensures that all traffic from the connected host is correctly placed into VLAN 100, whether the host itself is tagging packets or not (it usually isn't).

#### Diagram: VLAN Access Port in Action

This diagram shows what happens to a single Ethernet frame as it travels from a host, across the `veth` pair, and into the VLAN-aware bridge.

```ascii
+---------------------------+                                +------------------------------------+
|   Host Namespace (host10) |                                |  Bridge Namespace (br-namespace)   |
|                           |                                |                                    |
| +-----------------------+ |                                | +--------------------------------+ |
| | Application           | |                                | |      VLAN-Aware Bridge (br1)   | |
| +-----------------------+ |                                | |                                | |
|           |               |                                | |              +---------------+ | |
|           v               |                                | |              | VLAN 100 DB   | | |
| +-----------------------+ |       Untagged Ethernet Frame  | |              +---------------+ | |
| |      eth0             | |      (No VLAN Header)          | |                     ^            | |
| +-----------------------+ |                                | |                     |            | |
|           |               | -----------------------------> | +---------------------+----------+ | |
+-----------|---------------+
| | Port: veth-host10   |          | |
            |                                                | | PVID=100            |          | |
            +------------------ veth pair -------------------+ | Tagging Happens Here! |          | |
                                                             | +---------------------+----------+ |
                                                             |           |                        | |
                                                             |           v                        | |
                                                             |   Forwarded to other VLAN 100 ports  |
                                                             +------------------------------------+
```

**Explanation of the Diagram:**

1.  **Host Sends Traffic:** An application inside the `host10` namespace sends data. The host's network stack creates a standard Ethernet frame. Since the host is not configured to be VLAN-aware, this frame has no VLAN tag.

2.  **Travels Over `veth`:** The untagged frame is sent out of the host's `eth0` interface and travels across the `veth` pair to the bridge namespace.

3.  **Arrives at the Bridge Port:** The frame arrives at the `veth-host10` interface, which is a port on the `br1` bridge.

4.  **The `PVID` Rule is Applied:** The bridge sees an untagged frame arrive on a port that has a `PVID` of `100`. The bridge immediately **adds a VLAN tag with ID 100** to the frame.

5.  **Internal Bridge Logic:** Now that the frame is tagged, the bridge knows it belongs to VLAN 100. It looks up its forwarding database and sends the frame out only to other ports that are also members of VLAN 100.

This is why the `pvid` setting is so critical for connecting VLAN-unaware devices to a VLAN-aware network. It automatically handles the VLAN tagging for them at the edge of the network (the bridge port).

### Scenario: Communication Between Two Hosts in the Same VLAN

Let's expand on the previous example to see how two hosts, `host10` and `host11`, communicate when they are in the same VLAN (e.g., VLAN 100).

```ascii
+---------------------------+                                +------------------------------------+---------------------------+
|   Host Namespace (host10) |                                |  Bridge Namespace (br-namespace)   |   Host Namespace (host11) |
|                           |                                |                                    |                           |
|      eth0: 10.1.1.10      |                                |      VLAN-Aware Bridge (br1)       |      eth0: 10.1.1.11      |
|           |               |       Untagged Frame           |                                    |           ^               |
|           v               | -----------------------------> | +---------------------+----------+ |           |               |
| (1) Send Ping             |                                | | Port: veth-host10   |          | |           |               |
+---------------------------+                                | | PVID=100            |          | |           |               |
                                                             | | (2) Tag Added (VID=100)          | |           |               |
                                                             | +---------------------+----------+ |           |               |
                                                             |           |                        | |           |               |
                                                             |           v (3) Forwarding         | |           |               |
                                                             |                                    | |           |               |
                                                             | +---------------------+----------+ |       Untagged Frame      |
                                                             | | Port: veth-host11   |          | | <------------------------ |
                                                             | | PVID=100            |          | | (5) Receive Ping          |
                                                             | | (4) Tag Stripped    |          | +---------------------------+
                                                             | +---------------------+----------+
                                                             +------------------------------------+
```

**Explanation of the Traffic Flow:**

1.  **`host10` Sends a Packet:** `host10` wants to ping `host11`. It creates a standard, untagged Ethernet frame with the destination MAC address of `host11`.

2.  **Frame Arrives at Bridge and is Tagged:** The untagged frame travels from `host10` to the bridge, arriving at the `veth-host10` port. Because this port is configured with `pvid 100`, the bridge adds a VLAN tag with ID `100` to the frame.

3.  **Bridge Forwards the Frame:** The bridge now sees a frame tagged for VLAN 100, destined for `host11`'s MAC address. It looks up its forwarding database and knows that this MAC address is reachable via the `veth-host11` port. Since `veth-host11` is also a member of VLAN 100, the bridge forwards the frame to that port.

4.  **Frame is Untagged on Egress:** The `veth-host11` port is also an access port. As the frame is sent out of this port towards `host11`, the bridge **strips the VLAN 100 tag**.

5.  **`host11` Receives the Packet:** `host11` receives a standard, untagged Ethernet frame, exactly as if it were on a simple, non-VLAN network with `host10`.

This entire process of tagging and untagging is completely transparent to the end hosts. They communicate as normal, while the VLAN-aware bridge enforces the network segmentation, ensuring the frame is only delivered to other members of the same VLAN.
