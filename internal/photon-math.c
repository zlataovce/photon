#include <photon-math.h>
#include <stdio.h>
#include <stdlib.h>

const double LDBL_EXP = pow(FLT_RADIX, LDBL_MANT_DIG);

photon_ldouble_t photon_ldouble_wrap(long double v) {
  photon_ldouble_t wrap = malloc(sizeof(struct photon_ldouble));
  if (wrap == NULL) {
    fprintf(stderr, "malloc failed\n");
    return NULL;
  }

  int exp;
  long double mant = frexpl(v, &exp);

  int exp_mul = v < 0.0 ? -1 : 1; // move mantissa signbit to exponent
  wrap->exp = exp * exp_mul;
  wrap->mant = mant * LDBL_EXP * exp_mul;
  return wrap;
}

photon_ldouble_t* photon_ldouble_wrap_arr(void* v, size_t n) {
  photon_ldouble_t* w_arr = malloc(n * sizeof(photon_ldouble_t));
  if (w_arr == NULL) {
    fprintf(stderr, "malloc failed\n");
    return NULL;
  }

  long double* w = v;
  for (int i = 0; i < n; i++) {
    w_arr[i] = photon_ldouble_wrap(w[i]);
  }

  return w_arr;
}

long double photon_ldouble_unwrap(photon_ldouble_t wrap) {
  int exp = wrap->exp;
  uint64_t mant = wrap->mant;

  int exp_mul = exp < 0 ? -1 : 1; // move exponent signbit to mantissa
  return ldexpl((((long double) mant) / LDBL_EXP) * exp_mul, exp * exp_mul);
}

void* photon_ldouble_unwrap_arr(photon_ldouble_t* v, size_t n) {
  long double* w_arr = malloc(n * sizeof(long double));
  if (w_arr == NULL) {
    fprintf(stderr, "malloc failed\n");
    return NULL;
  }

  for (int i = 0; i < n; i++) {
    w_arr[i] = photon_ldouble_unwrap(v[i]);
  }

  return w_arr;
}
