#include <stdint.h>
#include <float.h>
#include <math.h>

extern const double LDBL_EXP;

struct photon_ldouble {
  int exp;
  uint64_t mant;
};

typedef struct photon_ldouble* photon_ldouble_t;

extern photon_ldouble_t photon_ldouble_wrap(long double v);
extern photon_ldouble_t* photon_ldouble_wrap_arr(void* v, size_t n);
extern long double photon_ldouble_unwrap(photon_ldouble_t wrap);
extern void* photon_ldouble_unwrap_arr(photon_ldouble_t* v, size_t n);
