#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include "periph/gpio.h"
#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "shell.h"
#include "board.h"

typedef struct SensorBoard {
	ipv6_addr_t addr;
} SensorBoard;

void button_handler(void* args) {
	printf("hipp hipp hurra\n");
}

int main(void) {
	printf("This is the server-emulator on %s.\n", RIOT_BOARD);
	
	// funktioniert nur mit GPIO_IN_PU
	// mit GPIO_RISING wird der Handler 2 mal bei Boardstart aufgerufen,
	// Grund unbekannt 
	gpio_init_int(BUTTON_GPIO, GPIO_IN_PU, GPIO_FALLING, button_handler, NULL);
	printf("Button activated\n");
	
	/*
	sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
	
	sock_udp_t sock;
	local.port = 2017;
	if(sock_udp_create(&sock, &local, NULL, 0) < 0) {
		puts("Error creating UDP sock");
		return EXIT_FAILURE;
	} else {
		puts("Successfull createt UDP socket");
	}
	
	while(true) {
		ssize_t res;
		if((res = sock_udp_recv(&sock, buf, sizeof(buf), SOCK_NO_TIMEOUT, NULL)) >= 0) {
			printf("%.128s\n", buf);
		}
	}
	*/
	
	// TODO shell-commands bauen?
	char line_buf[SHELL_DEFAULT_BUFSIZE];
	shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
    
    return EXIT_SUCCESS;
}
