# name of your application
APPLICATION = board-app

# If no BOARD is found in the environment, use this default:
BOARD ?= pba-d-01-kw2x #native
SERIAL ?= 02000203C3014E6A3EFDB392 #Thomas
# SERIAL ?= 02000203C36C4E143E90B3EC #Adrian

# This has to be the absolute path to the RIOT base directory:
RIOTBASE ?= $(CURDIR)/../RIOT

# Comment this out to disable code in RIOT that does safety checking
# which is not needed in a production environment but helps in the
# development process:
CFLAGS += -DDEVELHELP
CFLAGS += -DDEBUG_ASSERT_VERBOSE

# Change this to 0 show compiler invocation lines by default:
QUIET ?= 1

#Packeges to include:
USEPKG += nanocoap
USEPKG += jsmn

# Modules to include:
# We going to need some parts of RIOT
USEMODULE += fmt
USEMODULE += gcoap
USEMODULE += saul_default
USEMODULE += ipv6_addr
USEMODULE += ps
# Include GNRC
USEMODULE += gnrc_netdev_default
USEMODULE += auto_init_gnrc_netif
USEMODULE += gnrc_sock_udp
USEMODULE += gnrc_ipv6_default
USEMODULE += gnrc_icmpv6_echo
USEMODULE += gnrc_ipv6_router_default
# Also we want some statistics
USEMODULE += netstats_l2
USEMODULE += netstats_ipv6
USEMODULE += netstats_rpl

USEMODULE += shell
USEMODULE += shell_commands

INCLUDES += -I$(CURDIR)/drivers/include
INCLUDES += -I$(CURDIR)/include
INCLUDES += -I$(CURDIR)/../include

include $(RIOTBASE)/Makefile.include
