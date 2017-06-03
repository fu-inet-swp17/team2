/*
 * Copyright (C) 2016 Freie Universität Berlin
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

/**
 * @ingroup     driver_mpl3115a2
 * @{
 *
 * @file
 * @brief       MPL3115A2 adaption to the RIOT actuator/sensor interface
 *
 * @author      Adrian Herrmann <adrian.herrmann@fu-berlin.de>
 *
 * @}
 */

#include <string.h>
#include <stdio.h>

#include "saul.h"
#include "mpl3115a2.h"

static int read_press(void *dev, phydat_t *res)
{
    uint32_t* pres = 0;
    uint8_t* status = 0;
    mpl3115a2_read_pressure((mpl3115a2_t) *dev, pres, status);

    res->val[0] = *pres;
    res->unit = UNIT_PA;
    res->scale = 0;

    return 1;
}

const saul_driver_t mpl3115a2_saul_driver = {
    .read = read_press,
    .write = saul_notsup,
    .type = SAUL_SENSE_PRESS,
};
