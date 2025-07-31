#### Bash functions

create_bridge() {

  local nsname="$1"
  local ifname="$2"

  echo "Creating bridge ${nsname}/${ifname}"

  ip netns add ${nsname}
  ip netns exec ${nsname} ip link set lo up
  ip netns exec ${nsname} ip link add ${ifname} type bridge
  ip netns exec ${nsname} ip link set ${ifname} up
}

# The `create_bridge` function automates the creation of a network namespace 
# and a bridge interface within it. This setup is fundamental for creating isolated 
# network environments, where the bridge can connect multiple virtual interfaces.

# ASCII Diagram: Visualizing `create_bridge "ns-test" "br-test"`

# 1. Initial State:
# The system starts with the default network namespace, which contains the physical
# interfaces like `eth0`.

#    +----------------------------------------+
#    |         Default Network Namespace      |
#    |                                        |
#    |   +--------+      +--------+           |
#    |   |  eth0  |      |   lo   |           |
#    |   +--------+      +--------+           |
#    |                                        |
#    +----------------------------------------+

# 2. `ip netns add ns-test`:
# A new, isolated network namespace named "ns-test" is created. It starts with
# its own loopback interface `lo`, which is initially down.

#    +----------------------------------------+
#    |         Default Network Namespace      |
#    +----------------------------------------+
#
#    +----------------------------------------+
#    |          Network Namespace (ns-test)   |
#    |                                        |
#    |   +--------+                           |
#    |   |   lo   | (DOWN)                    |
#    |   +--------+                           |
#    |                                        |
#    +----------------------------------------+

# 3. `ip link add br-test type bridge` & `ip link set ... up`:
# Inside "ns-test", a bridge interface "br-test" is created, and both `lo` and
# `br-test` are brought up. The namespace is now ready to connect other interfaces.

#    +----------------------------------------+
#    |         Default Network Namespace      |
#    +----------------------------------------+
#
#    +------------------------------------------------+
#    |            Network Namespace (ns-test)         |
#    |                                                |
#    |   +------------------+      +--------+         |
#    |   |  Bridge (br-test) |      |   lo   | (UP)    |
#    |   |       (UP)       |      +--------+         |
#    |   +------------------+                         |
#    |           ^                                    |
#    |           |                                    |
#    |   (Ready to connect veth pairs)                |
#    |                                                |
#    +------------------------------------------------+

create_end_host() {
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

# The `create_end_host` function sets up a new network namespace, representing an
# end host, and connects it to an existing bridge in another namespace. This is
# achieved by creating a virtual Ethernet (veth) pair, which acts as a virtual
# patch cable between the two namespaces.

# ASCII Diagram: Visualizing `create_end_host "ns-h1" "veth-h1" "10.10.0.1/24" "ns-br" "br0"`

# 1. Initial State:
# We start with a namespace "ns-br" containing a bridge "br0". The default
# namespace is also present.

#    +----------------------------------------+
#    |         Default Network Namespace      |
#    +----------------------------------------+
#
#    +------------------------------------------------+
#    |            Network Namespace (ns-br)           |
#    |                                                |
#    |   +------------------+      +--------+         |
#    |   |   Bridge (br0)   |      |   lo   | (UP)    |
#    |   |       (UP)       |      +--------+         |
#    |   +------------------+                         |
#    |                                                |
#    +------------------------------------------------+

# 2. `ip netns add ns-h1`:
# A new namespace "ns-h1" is created for the end host.

#    +----------------------------------------+
#    |         Default Network Namespace      |
#    +----------------------------------------+
#
#    +----------------------------------------+     +----------------------------------------+
#    |          Network Namespace (ns-br)   |     |          Network Namespace (ns-h1)   |
#    |                                        |     |                                        |
#    |   +------------------+                 |     |   +--------+                           |
#    |   |   Bridge (br0)   |                 |     |   |   lo   | (UP)                      |
#    |   +------------------+                 |     |   +--------+                           |
#    +----------------------------------------+     +----------------------------------------+

# 3. `ip link add veth-h1 ... peer veth-h1b ...`:
# A veth pair is created. "veth-h1" is placed in "ns-h1" and "veth-h1b" in "ns-br".

#    +----------------------------------------+
#    |         Default Network Namespace      |
#    +----------------------------------------+
#
#    +----------------------------------------+     +-------------------------------------------+
#    |          Network Namespace (ns-br)   |     |          Network Namespace (ns-h1)      |
#    |                                        |     |                                           |
#    |   +------------------+  +-----------+  |     |   +-----------+      +--------+         |
#    |   |   Bridge (br0)   |  | veth-h1b  |  |     |   |  veth-h1  |      |   lo   | (UP)    |
#    |   +------------------+  +-----------+  |     |   +-----------+      +--------+         |
#    |                           |            |     |        |                                |
#    |                           +------------+-----+--------+                                |
#    |                                        |     |                                           |
#    +----------------------------------------+     +-------------------------------------------+

# 4. `ip link set ... master br0` & `ip addr add ...`:
# "veth-h1b" is attached to the bridge "br0". "veth-h1" gets an IP address.

#    +-------------------------------------------+     +------------------------------------------------------+
#    |           Network Namespace (ns-br)       |     |              Network Namespace (ns-h1)             |
#    |                                           |     |                                                      |
#    |   +------------------+                    |     |   +-------------------------+      +--------+        |
#    |   |  Bridge (br0)    | <--- attached ---- | ----+---| veth-h1 (10.10.0.1/24)  |      |   lo   | (UP)   |
#    |   | +-------------+  |                    |     |   +-------------------------+      +--------+        |
#    |   | |  veth-h1b   |  |                    |     |                                                      |
#    |   | +-------------+  |                    |     |                                                      |
#    |   +------------------+                    |     |                                                      |
#    +-------------------------------------------+     +------------------------------------------------------+
