/**
 * Copyright (C) 2017 Freie Universität Berlin
 * Authors: Adrian Herrmann <adrian.herrmann@fu-berlin.de>
 *
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

#ifndef __SE_SWP_2017_FLOATINT_H__
#define __SE_SWP_2017_FLOATINT_H__

/**
 * @brief                  Creates a string with a float value without using floating point arithmetics
 * 
 * @params[out]    buf     Buffer where the string representation of the floating point number is to be written
 * @params[in]     val     Integer representation of the number to be converted
 * @params[in]     scale   Decimal scale of the number represented by val
 */
void int_to_floatstring(char* buf, int16_t val, int8_t scale);

/**
 * @brief                  Rounds a float represented by an integer value and a decimal scale to an integer
 * 
 * @params[in]     val     Integer representation of the number to be converted
 * @params[in]     scale   Decimal scale of the number represented by val
 * 
 * @return                 The rounded integer
 */
int32_t int_from_float(int16_t val, int8_t scale);

#endif