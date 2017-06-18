#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <stdbool.h>

#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "shell.h"
#include "saul.h"
#include "saul_reg.h"
#include "net/gcoap.h"
#include "thread.h"
#include "xtimer.h"

#include "smart_environment.h"
#include "coap.h"
#include "sensors.h"


char own_addr[IPV6_ADDR_MAX_STR_LEN];

int main(void) {
    printf("Smart environment app on %s\n", RIOT_BOARD);
	printf("%u\n", APP_ID_LEN);
	
	msg_t msg_queueq[8];
	msg_init_queue(msg_queueq, 8);

    saul_reg_t* dev = saul_reg;

    if (dev == NULL) {
        puts("The SAUL registry does not contain any devices.");
    } else {
        printf("The SAUL registry contains the following devices: ");
        while(1) {
            printf("%s", dev->name);
            if (dev->next == NULL) {
                break;
            }
            printf(", ");
            dev = dev->next;
        };
        printf("\n");
    }

    sensors_test();

    ipv6_addr_t addr;
    ipv6_addr_from_str(&addr, "fe80::");

    ipv6_addr_to_str(own_addr, &addr, IPV6_ADDR_MAX_STR_LEN);

    ipv6_addr_t* out = NULL;
    gnrc_ipv6_netif_find_by_prefix(&out, &addr);
    
    ipv6_addr_to_str(own_addr, out, IPV6_ADDR_MAX_STR_LEN);
    printf("own ipv6 addr: %s\n", own_addr);

	gcoap_register_listener(&coap_listener);
	puts("coap ready");

    // thread für pings
    thread_create(
        ping_stack,
        THREAD_STACKSIZE_DEFAULT,
        THREAD_PRIORITY_MAIN-1,
        0,
        ping_handler,
        NULL,
        "ping"
    );
    
    // shell starten damit Application nicht terminiert
    char line_buf[SHELL_DEFAULT_BUFSIZE];
    shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
        
    return EXIT_SUCCESS;
}
