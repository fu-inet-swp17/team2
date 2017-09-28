#include <stdio.h>
#include <string.h>
#include <math.h>

/*
 * @brief                  Creates a string with a float value without using floating point arithmetics
 * 
 * @params[out]    buf     Buffer where the string representation of the floating point number is to be written
 * @params[in]     val     Integer representation of the number to be converted
 * @params[in]     scale   Decimal scale of the number represented by val
 */
void int_to_floatstring(char* buf, int16_t val, int8_t scale) {
    char valstring[6];
    sprintf(valstring, "%d", val);
    size_t val_len = strlen(valstring);
    if (scale < 0) {
        snprintf(buf, val_len+scale+1, valstring);
        strncat(buf, ".", 1);
        strncat(buf, &valstring[val_len+scale], val_len);
    } else if (scale == 0) {
        snprintf(buf, val_len, valstring);
    } else {
        snprintf(buf, val_len+1, valstring);
        int8_t i = 0;
        while (i<scale) {
            strncat(buf, "0", 1);
            i++;
        }
    }
    
}

/*
 * @brief                  Rounds a float represented by an integer value and a decimal scale to an integer
 * 
 * @params[in]     val     Integer representation of the number to be converted
 * @params[in]     scale   Decimal scale of the number represented by val
 * 
 * @return                 The rounded integer
 */
int32_t int_from_float(int16_t val, int8_t scale) {
    if (scale < 0) {
        int16_t val2 = val % (int16_t) pow(10, -scale);
        int16_t val1 = (val - val2) / (int16_t) pow(10, -scale);
        if (val2 >= (int16_t) pow(10, -scale) / 2) {
            val1 += 1;
        }
        return val1;
    } else if (scale == 0) {
        return val;
    } else {
        return val * (int16_t) pow(10, scale);
    }
}
