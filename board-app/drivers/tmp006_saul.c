/*
 * Copyright (C) 2016 Freie Universität Berlin
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

/**
 * @ingroup     driver_tmp006
 * @{
 *
 * @file
 * @brief       TMP006 adaption to the RIOT actuator/sensor interface
 *
 * @author      Adrian Herrmann <adrian.herrmann@fu-berlin.de>
 *
 * @}
 */

#include <string.h>
#include <stdio.h>

#include "saul.h"
#include "tmp006.h"

static int read_temp(void *dev, phydat_t *res)
{
    int16_t* rawv = 0;
    int16_t* rawt = 0;
    uint8_t* drdy = 0;
    tmp006_read((tmp006_t *)dev, rawv, rawt, drdy);

    float* tamb = 0;
    float* tobj = 0;
    tmp006_convert(*rawv, *rawt, tamb, tobj);

    res->val[0] = *tamb;
    res->val[1] = *tobj;
    res->unit = UNIT_TEMP_C;
    res->scale = 0;

    return 2;
}

const saul_driver_t tmp006_saul_driver = {
    .read = read_temp,
    .write = saul_notsup,
    .type = SAUL_SENSE_TEMP,
};
