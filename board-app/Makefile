# name of your application
APPLICATION = board-app

# If no BOARD is found in the environment, use this default:
BOARD ?= pba-d-01-kw2x #native

# This has to be the absolute path to the RIOT base directory:
RIOTBASE ?= $(CURDIR)/../RIOT

# Comment this out to disable code in RIOT that does safety checking
# which is not needed in a production environment but helps in the
# development process:
CFLAGS += -DDEVELHELP

# Change this to 0 show compiler invocation lines by default:
QUIET ?= 1

# Modules to include:
USEMODULE += shell
USEMODULE += shell_commands
USEMODULE += ps
USEMODULE += gnrc_netdev_default
USEMODULE += auto_init_gnrc_netif
USEMODULE += gnrc_ipv6_router_default
USEMODULE += gnrc_sock_udp
USEMODULE += ipv6_addr
USEMODULE += saul
USEMODULE += saul_reg

USEMODULE += mag3110
USEMODULE += mma8x5x
USEMODULE += tmp006
USEMODULE += mpl3115a2
USEMODULE += hdc1000
USEMODULE += tcs37727

INCLUDES += -I$(CURDIR)/drivers/include
INCLUDES += -I$(CURDIR)/include
INCLUDES += -I$(CURDIR)/../include

include $(RIOTBASE)/Makefile.include
