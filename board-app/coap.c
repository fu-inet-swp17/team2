#include "saul_reg.h"
#include "net/gcoap.h"
#include "senml.h"
#include "xtimer.h"

#include "smart_environment.h"

#define PING_TIMEOUT	10


static ssize_t sensors_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
    gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
    
    ssize_t payload_len = snprintf(
    	(char*)pdu->payload,
    	GCOAP_PDU_BUF_SIZE,
    	"%u",
    	(unsigned int)IR_TEMP_SENSOR | HUMID_SENSOR | MAG_SENSOR | 
    	RGB_LIGHT_SENDSOR | PRESS_SENSOR | ACC_SENSOR
    );
    
    puts("List of sensors requested.");
    return gcoap_finish(pdu, (size_t)payload_len, COAP_FORMAT_TEXT);
}

static int64_t temp_sum = 0; // TODO size of temp sum

static ssize_t temp_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
    gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
    phydat_t res;
    saul_reg_t* dev = saul_reg_find_type(SAUL_SENSE_TEMP);

    char dev_name[50]; // TODO dev name length
    strncpy(dev_name, dev->name, 50);
    saul_reg_read(dev, &res);

    temp_sum += res.val[0];

    uint8_t senml_len = 50; // TODO senml message length
    char senml_json_output[senml_len];
    senml_base_info_t base_info = {
        .version = SENML_SUPPORTED_VERSION,
        .base_name = dev_name,
        .base_time = xtimer_now_usec() / 1e6,
        .base_unit = (char*) ((int) res.unit),
        .base_value = 0,
        .base_sum = 0
    };
    senml_record_t records = {
        .update_time = 1,
        .value_sum = (double) temp_sum,
        .value_type = SENML_TYPE_FLOAT,
        .value.f = (double) res.val[0]
    };
    senml_pack_t pack = {
        &base_info,
        &records,
        1
    };

    senml_encode_json_s(&pack, senml_json_output, senml_len);
    
    ssize_t payload_len = snprintf(
    	(char*)pdu->payload,
    	GCOAP_PDU_BUF_SIZE,
    	"%s",
    	senml_json_output
    );

    puts("Data from thermometer requested.");
    return gcoap_finish(pdu, (size_t)payload_len, COAP_FORMAT_NONE);
}

static ssize_t humid_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
    gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
    phydat_t res;
    saul_reg_read(saul_reg_find_type(SAUL_SENSE_HUM), &res);
    
    ssize_t payload_len = snprintf(
    	(char*)pdu->payload,
    	GCOAP_PDU_BUF_SIZE,
    	"%d",
    	res.val[0]
    );

    puts("Data from humidity sensor requested.");
    return gcoap_finish(pdu, (size_t)payload_len, COAP_FORMAT_NONE);
}

static ssize_t mag_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
    gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
    phydat_t res;
    saul_reg_read(saul_reg_find_type(SAUL_SENSE_MAG), &res);
    
    ssize_t payload_len = snprintf(
    	(char*)pdu->payload,
    	GCOAP_PDU_BUF_SIZE,
    	"%d",
    	res.val[0]
    );

    puts("Data from magnetometer requested.");
    return gcoap_finish(pdu, (size_t)payload_len, COAP_FORMAT_NONE);
}

void* ping_handler(void* args) {
    (void)args;
    sock_udp_ep_t remote = SOCK_IPV6_EP_ANY;
    remote.port = SERVER_CONN_PORT;
    // ff02::1 -> addr für link-local broadcast
    ipv6_addr_from_str((ipv6_addr_t *)&remote.addr.ipv6, "ff02::1");
    PingMsg intro_msg;
    snprintf(intro_msg.app_id, APP_ID_LEN, "%s", app_id);
    
    intro_msg.sensors = IR_TEMP_SENSOR | HUMID_SENSOR | MAG_SENSOR |
    RGB_LIGHT_SENDSOR | PRESS_SENSOR | ACC_SENSOR;
    
    puts("Ping thread running");
    
    while(true) {
        ssize_t res = sock_udp_send(
            NULL,
            &intro_msg,
            APP_PING_MSG_LEN,
            &remote
        );
        
        if(res == -EAFNOSUPPORT) {
            fputs("Ping: EAFNOSUPPORT", stderr);
        } else if(res == -EHOSTUNREACH) {
            fputs("Ping: EHOSTUNREACH", stderr);
        } else if(res == -EINVAL) {
            fputs("Ping: EINVAL", stderr);
        } else if(res == -ENOMEM) {
            fputs("Ping: ENOMEM", stderr);
        } else if(res == -ENOTCONN) {
            fputs("Ping: ENOTCONN", stderr);
        } else if(res < 0) {
            fprintf(stderr, "Ping error: %d\n", res);
        }
        
        if(res < 0) {
            break;
        } else {
            xtimer_sleep(PING_TIMEOUT);
        }
    }
    
    puts("ping thread terminating!");
    
    return NULL;
}

const coap_resource_t coap_resources[] = {
    // ressource-path, ressource-type, response-handler
    {"/se-app/sensors", COAP_GET, &sensors_handler},
    {"/se-app/temp", COAP_GET, &temp_handler},
    {"/se-app/humid", COAP_GET, &humid_handler},
    {"/se-app/mag", COAP_GET, &mag_handler}
};

gcoap_listener_t coap_listener = {
    (coap_resource_t *)&coap_resources[0],
    sizeof(coap_resources) / sizeof(coap_resources[0]),
    NULL
};
