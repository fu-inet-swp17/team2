/**
 * Copyright (C) 2017 Freie Universität Berlin
 * Author: Thomas Tegethoff cone@zedat.fu-berlin.de
 *
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

/**
 * @brief	driver for WS1812B LED stripes
 * Carfull! the NOP Slides in the 'ws2812b_show' function ar timed for the
 * PhyTEC WAVE KW22 (pba-d-01-kw2x), it might not work with other boards!
 * Also the output-pin (PORT_A 2) is hard-coded
 */

#ifndef __SE_2017_SWP_WS2812B_DRIVER_H__
#define __SE_2017_SWP_WS2812B_DRIVER_H__

#include "periph/gpio.h"
#include "color.h"

#ifdef __cplusplus
extern "C" {
#endif

#define MAX_LED_NUM	100

/**
 * @brief Descriptor struct for WS2812B LED stripe
 */
typedef struct {
    uint32_t led_count;  /**< count of LEDs to control */
    gpio_t pin;   /**< pin for data out */
    color_rgb_t leds[MAX_LED_NUM];
} ws2812b_stripe_t;

/**
 * @brief Initialize an LED stripe by assigning it a gpio pin
 *
 * @param[out]  dev         struct describing the WS2812B LED stripe
 * @param[in]   leds        LED count (just for further development, the LED
 * count is described by the macro 'MAX_LED_NUM')
 * @param[in]   data        gpio pin for data out
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
int ws2812b_init(ws2812b_stripe_t *dev, uint32_t led_count, gpio_t data_out);

/**
 * @brief sets the color of an single LED in an WS2812B-Stripe
 *
 * @param[out]  dev         struct describing the WS2812B LED stripe
 * @param[in]   index       position of the LED on the stripe
 * @param[in]   led         new color values for the LED
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
int ws2812b_set_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t led);

/**
 * @brief returns the color of an single LED in an WS2812B stripe
 *
 * @param[out]  dev         struct describing the WS2812B LED stripe
 * @param[in]   index       position of the LED on the stripe
 * @param[in]   led         container for the LED color value
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
int ws2812b_get_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t *led);

/**
 * @brief sets the given color values on the real LED stripe
 *
 * @param[out]  dev         struct describing the WS2812B LED stripe
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
void ws2812b_show(ws2812b_stripe_t *dev);

#ifdef __cplusplus
}
#endif

#endif // WS2812B_DRIVER

