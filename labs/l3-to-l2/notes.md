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
