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

#include "mpl3115a2.h"

int read_press(void *dev, phydat_t *res);

const saul_driver_t mpl3115a2_saul_driver;
