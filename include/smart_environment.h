#ifndef _SMART_ENVIRONMENT_H_
#define _SMART_ENVIRONMENT_H_

#define CLIENT_PORT 2018
#define SERVER_CONN_PORT 2017
#define SERVER_POLL_PORT 2018

#define CLIENT_ID_LEN sizeof(client_id)
#define SERVER_ID_LEN sizeof(server_id)

#define CLIENT_INIT_MSG_LEN CLIENT_ID_LEN + IPV6_ADDR_MAX_STR_LEN + 1
#define SERVER_RESP_MSG_LEN SERVER_ID_LEN + IPV6_ADDR_MAX_STR_LEN + 1

// statische globals, nach außen hin nicht sichtbar
const char server_id[] = "riot-swp-2017-se-app-server";
const char client_id[] = "riot-swp-2017-se-app-client";

#endif
