#include <stdint.h>
#include <stddef.h>
#include <stdio.h>

#include "periph/gpio.h"
#include "cpu.h"
#include "debug.h"
#include "xtimer.h"
#include "vendor/MKW22D5.h"

#include "ws2812b_driver.h"

#define ENABLE_DEBUG    (0)

#define asm __asm__

/*
 * Full transmission: 1.25us +/- 600ns
 * 400ns +/- 150ns HIGH
 * 850ns +/- 150ns LOW
 */
#define SEND_ZERO	GPIOA->PSOR = 0x04;\
					asm volatile(\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
					);\
					GPIOA->PCOR = 0x04;\
					asm volatile(\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
					);\
					
/*
 * Full transmission: 1.25us +/- 600ns
 * 800ns +/- 150ns HIGH
 * 450ns +/- 150ns LOW
 */
#define SEND_ONE	GPIOA->PSOR = 0x04;\
					asm volatile(\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
					);\
					GPIOA->PCOR = 0x04;\
					asm volatile(\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
						"NOP\n"\
					);\


/*
 * @brief Initialize a LED-Stripe by assigning it a gpio pin
 *
 * @param[out] dev         struct describing the neopixel led stripe
 * @param[in] leds         led count (just for further development, the LED
 * count is described, by the macro 'MAX_LED_NUM')
 * @param[in] data         gpio pin for data out
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
int ws2812b_init(ws2812b_stripe_t *dev, uint32_t led_count, gpio_t pin) {
	if (gpio_init(pin, GPIO_OUT) < 0) {
		DEBUG("ws2812b_init: failed to init gpio %d\n", pin);
		return -1;
	}

	gpio_clear(pin);

	dev->led_count = led_count;
	dev->pin = pin;

	return 0;
}

/*
 * @brief sets the color of an single LED in an WS2812B-Stripe
 *
 * @param[out] dev         struct describing the WS2812B led stripe
 * @param[in] index        position of the LED on the stripe
 * @param[in] led          new color-values for the LED
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
int ws2812b_set_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t led) {
	if(index >= dev->led_count) {
		DEBUG("ws2812b_set_color: index out of bound.\n");
		return -1;
	}

	dev->leds[index].r = led.r;
	dev->leds[index].g = led.g;
	dev->leds[index].b = led.b;

	return 0;
}

/*
 * @brief returns the color of an single LED in an WS2812B-Stripe
 *
 * @param[out] dev         struct describing the WS2812B led stripe
 * @param[in] index        position of the LED on the stripe
 * @param[in] led          container for the LED-color-value
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
// maybe a double pointer for the LED?
int ws2812b_get_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t *led) {
	if (index >= dev->led_count) {
		DEBUG("ws2812b_get_color: index ouf of bound.\n");
		return -1;
	}

	/* NOTE: For color order see neopixel_set_pixel_color */
	led->r = dev->leds[index].r;
	led->g = dev->leds[index].g;
	led->b = dev->leds[index].b;

	return 0;
}

/*
 * @brief sets the given Color-Values on the real LED-Stripe
 *
 * @param[out] dev         struct describing the WS2812B led stripe
 *
 * @return                  0 on success
 * @return                  <0 on error
 */
void ws2812b_show(ws2812b_stripe_t *dev) {
	int irq_state;
	uint32_t i;	// Timing!
	// gpio_t pin = dev->pin;
	color_rgb_t *pixel; // Timing!
		
	gpio_init(GPIO_PIN(PORT_A, 2), GPIO_OUT);
	gpio_clear(GPIO_PIN(PORT_A, 2));


	/* Disable IRQ's, we have hard time bounds here.
	 * Restore it afterwards */

	irq_state = irq_disable();
	
	/* Loop over all connected leds
	 * For the WS2812B only the HIGH Signal has to be distinctive,
	 * for LOW we don't need any NOP's.
	 * We could also loop over all three colors for each LED,
	 * but since we need hard timings for each 24bit we unroled the loop.
	 */
	for(i=0; i<dev->led_count; i++) {
		pixel = &(dev->leds[i]);
		
		// green
		if(pixel->g & 0x80) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->g & 0x40) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->g & 0x20) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->g & 0x10) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->g & 0x08) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->g & 0x04) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->g & 0x02) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->g & 0x01) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
				
		// red
		if(pixel->r & 0x80) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->r & 0x40) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->r & 0x20) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->r & 0x10) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->r & 0x08) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->r & 0x04) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->r & 0x02) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->r & 0x01) {
			SEND_ONE
		} else {
			SEND_ZERO
		}

		// blue
		if(pixel->b & 0x80) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->b & 0x40) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->b & 0x20) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->b & 0x10) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->b & 0x08) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->b & 0x04) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->b & 0x02) {
			SEND_ONE
		} else {
			SEND_ZERO
		}
		
		if(pixel->b & 0x01) {
			SEND_ONE
		} else {
			SEND_ZERO
		}

	}
	
	/* Restore IRQ State */
	irq_restore(irq_state);
	xtimer_usleep(50);
}

