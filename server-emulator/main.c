#include <stdio.h>
#include <stddef.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include "periph/gpio.h"
#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "net/gcoap.h"
#include "shell.h"
#include "board.h"
#include "thread.h"
#include "xtimer.h"

#include "smart_environment.h"

#define MAX_BOARD_NUM 10
#define TIMEOUT (uint32_t)2000000

typedef struct SensorBoard {
    bool init;
	sock_udp_ep_t board_ep;
} SensorBoard;

SensorBoard boards[MAX_BOARD_NUM];
char own_addr[IPV6_ADDR_MAX_STR_LEN];
char connect_stack[THREAD_STACKSIZE_MAIN];
uint8_t conn_buf[APP_PING_MSG_LEN];
uint8_t req_buf[128];

/*
 * description: checks if string prefix is a prefix of string str
 * requires: stdbool.h, string.h
 * param[char* str]: string that should start with something
 * param[char* prefix]: the prefix to check for
 * return[bool]: true wenn prefix is a prefix of str, false else
 */
static inline bool startsWith(const char* str, const char* prefix) {
    char* substr = strstr(str, prefix);
    return substr == str;
}

/*
 * description: compares two IPv6 addresses
 * param[ipv6_addr_t* addr1]: first address for compare
 * param[ipv6_addr_t* addr2]: second address for compare
 * return[bool]: true if addresses are equal, false if not
 */
static inline bool compareIPv6Addr(ipv6_addr_t* addr1, ipv6_addr_t* addr2) {
	for(size_t i=0; i<16; i++) {
		if(addr1->u8[i] != addr2->u8[i]) {
			return false;
		}
	}
	
	return true;
}

/*
 * description: strncpy guaranteeing null termination
 * param[char* dst]: buffer to copy data into
 * param[const char* src]: string to be copied into dst
 * param[size_t num]: max number of chars to be written into dst
 * return[size_t]: number of chars written
 */
static inline size_t my_strncpy(char* dst, const char* src, size_t num) {
    size_t i;
    for(i=0; src[i] && i<(num-1); i++) {
        dst[i] = src[i];
    }
    
    dst[i] = 0;
    return i;
}


/**********************************COAP STUFF**********************************/

coap_pkt_t pdu;
uint8_t coap_buff[GCOAP_PDU_BUF_SIZE];

static void sensors_resp_handler(unsigned req_state, coap_pkt_t* pdu) {
    if (req_state == GCOAP_MEMO_TIMEOUT) {
        printf("timeout for msg ID %02u\n", coap_get_id(pdu));
        return;
    } else if (req_state == GCOAP_MEMO_ERR) {
        printf("error in response\n");
        return;
    }
    
    puts("response!");
    printf("%s\n", (char*)pdu->payload);
}

/**********************************COAP STUFF**********************************/

static void button_handler(void* args) {
    puts("Requesting data from boards...");
    
    for(uint8_t i=0; i<MAX_BOARD_NUM; i++) {        
        if(!boards[i].init) {
            printf("Slot %u not used.\n", i);
            continue;
        }
        boards[i].board_ep.port = GCOAP_PORT;
        
        size_t req_len = (size_t)gcoap_request(
            &pdu,
            coap_buff,
            GCOAP_PDU_BUF_SIZE,
            COAP_GET,
            "/.well-known/core" /* replace this with whatever resource you wish to request */
        );
        
        size_t bytes_sent = gcoap_req_send2(
            coap_buff,
            req_len,
            &boards[i].board_ep,
            &sensors_resp_handler
        );
        
        char addr_str[IPV6_ADDR_MAX_STR_LEN];
        ipv6_addr_to_str(
         	addr_str,
            (ipv6_addr_t*)&boards[i].board_ep.addr.ipv6,
            IPV6_ADDR_MAX_STR_LEN
        );
        printf(
            "Send request %u to addr %s. Size: %u\n",
            coap_get_id(&pdu),
            addr_str,
            bytes_sent
        );
    }
}


/*
 * Daemon running in the background and registering new boards
 */
static void* connect_thread_handler(void* args) {
    (void)args;
    
    sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
    local.port = SERVER_CONN_PORT;
    ipv6_addr_from_str((ipv6_addr_t *)&local.addr.ipv6, "ff02::1");
    sock_udp_t sock;
    
    bool sock_ready = false;
    if(sock_udp_create(&sock, &local, NULL, 0) < 0) {
        fputs("Error creating UDP socket!\n", stderr);
        sock_ready = false;
    } else {
        puts("Successfully created UDP socket.\nWaiting for boards...");
        sock_ready = true;
    }
    
    while(sock_ready) {
    	sock_udp_ep_t remote;
        ssize_t res = sock_udp_recv(
            &sock,
            conn_buf,
            APP_PING_MSG_LEN,
            SOCK_NO_TIMEOUT,
            &remote
        );
        
        if(res < 0) {
            fprintf(
                stderr,
                "Error during \"sock_udp_recv\": %d!\nAboarting.",
                res
            );
            sock_ready = false;
        } else if((res) >= 0) {
            if(!startsWith((const char*)conn_buf, app_id)) {
                continue;
            }
                
            size_t board_index;
            bool in_use = false;
            for(board_index=0;
                board_index<MAX_BOARD_NUM && boards[board_index].init;
                board_index++
            ) {
            	if(compareIPv6Addr(
            		(ipv6_addr_t*)&boards[board_index].board_ep.addr.ipv6,
            		(ipv6_addr_t*)&remote.addr.ipv6
            	)) {
            		printf("Ping from board: %u\n", board_index);
            		in_use = true;
            	}
            }
                
            if(board_index >= MAX_BOARD_NUM) {
                fputs("Max nr. of boards registered!\n", stderr);
                break; /* thread terminates! */
            }
            
            if(!in_use) {
            	boards[board_index].init = true;
            	memcpy(
            		(void*)&boards[board_index].board_ep,
            		(void*)&remote,
            		sizeof(sock_udp_ep_t)
            	);
            
            	char addr_str[IPV6_ADDR_MAX_STR_LEN];
            	ipv6_addr_to_str(
            		addr_str,
            		(ipv6_addr_t*)&boards[board_index].board_ep.addr.ipv6,
            		IPV6_ADDR_MAX_STR_LEN
            	);
            
           		printf("board found at: %s\n", addr_str);
       	        PingMsg* msg = (PingMsg*)conn_buf;
       	        
       	        fputs("with: ", stdout);
       	        if(msg->sensors & IR_TEMP_SENSOR) {
       	        	fputs("IR-Thermopile sensor, ", stdout);
       	        }
       	        
       	        if(msg->sensors & HUMID_SENSOR) {
       	        	fputs("humidity sensor, ", stdout);
       	        }
       	        
       	        if(msg->sensors & MAG_SENSOR) {
       	        	fputs("magnetometer, ", stdout);
       	        }
       	        
       	        if(msg->sensors & RGB_LIGHT_SENDSOR) {
       	        	fputs("color light sensor, ", stdout);
       	        }
       	        
       	        if(msg->sensors & PRESS_SENSOR) {
       	        	fputs("pressure sensor, ", stdout);
       	        }
       	        
       	        if(msg->sensors & ACC_SENSOR) {
       	        	fputs("accelerometer, ", stdout);
       	        }
           		puts("");
           		
            }
        }
    }

    return NULL;
}

int main(void) {
    printf("This is the server-emulator on \"%s\".\n", RIOT_BOARD);
        
    msg_t msg_queueq[8];
	msg_init_queue(msg_queueq, 8);
     
    for(size_t i=0; i<MAX_BOARD_NUM; i++) {
        boards[i].init = false;
    }
    
    ipv6_addr_t addr;
    ipv6_addr_from_str(&addr, "fe80::");

    ipv6_addr_t* ll_addr = NULL;
    gnrc_ipv6_netif_find_by_prefix(&ll_addr, &addr);
    
    ipv6_addr_to_str(own_addr, ll_addr, IPV6_ADDR_MAX_STR_LEN);
    printf("Link-local address: \"%s\".\n", own_addr);
    
    gpio_init_int(BUTTON_GPIO, GPIO_IN_PU, GPIO_FALLING, button_handler, NULL);
    puts("Button activated.");
    
    thread_create(
        connect_stack,
        sizeof(connect_stack),
        THREAD_PRIORITY_MAIN-1,
        THREAD_CREATE_STACKTEST,
        connect_thread_handler,
        NULL,
        "connect thread"
    );

    char line_buf[SHELL_DEFAULT_BUFSIZE];
    shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
    
    return EXIT_SUCCESS;
}
