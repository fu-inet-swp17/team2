#ifndef WS2812B_DRIVER
#define WS2812B_DRIVER

#include "periph/gpio.h"
#include "color.h"

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Descriptor struct for Neopixel LED stripe
 */
typedef struct {
    uint32_t led_count;  /**< count of leds to control */
    gpio_t pin;   /**< pin for data out */
    color_rgb_t leds[100];
} ws2812b_stripe_t;

/**
 * @brief Initialize a Neopixel by assigning it a gpio pin
 *
 * @param[out] dev         struct describing the neopixel led stripe
 * @param[in] leds         led count
 * @param[in] data         gpio pin for data out
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
int ws2812b_init(ws2812b_stripe_t *dev, uint32_t led_count, gpio_t data_out);

int ws2812b_set_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t led);

int ws2812b_get_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t *led);

void ws2812b_show(ws2812b_stripe_t *dev);

#ifdef __cplusplus
}
#endif

#endif // WS2812B_DRIVER

