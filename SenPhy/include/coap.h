/**
 * Copyright (C) 2017 Freie Universität Berlin
 * Authors: Adrian Herrmann <adrian.herrmann@fu-berlin.de>
 *          Thomas Tegethoff <cone@zedat.fu-berlin.de>
 *
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 * 
 * Careful: This header should not be included from coap.c because of multiple definition of the variable coap_listener
 */

#ifndef __SE_SWP_2017_COAP_H___
#define __SE_SWP_2017_COAP_H___

/**
 * @brief       The main function for the ping thread which sends regular announcement messages
 * 
 * @params[in]  Not used
 * @return      Not used
 */
void* ping_handler(void* args);

/**
 * @brief       The resource listener for gcoap
 */
extern gcoap_listener_t coap_listener;

#endif