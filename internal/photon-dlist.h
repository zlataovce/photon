#include <mir-dlist.h>

// _name is the name of the accumulated struct
// _type is the accumulated type
// _item is the type being accumulated
// _selector is the code transforming _item to a DLIST of _type (item is the name of the variable)
#define DLIST_ACCUMULATOR(_name, _type, _item, _selector) \
    struct photon_##_name { \
      _type* items; \
      size_t length; \
    }; \
    typedef struct photon_##_name* photon_##_name##_t; \
    photon_##_name##_t photon_list_##_name(_item item) { \
      DLIST(_type) dlist = _selector; \
      int length = DLIST_LENGTH(_type, dlist); \
      _type* items = malloc(length * sizeof(_type)); \
      if (items == NULL) { \
          fprintf(stderr, "malloc failed\n"); \
          return NULL; \
      } \
      int index = 0; \
      _type item0; \
      for (item0 = DLIST_HEAD(_type, dlist); item0 != NULL; \
           item0 = DLIST_NEXT(_type, item0)) { \
          items[index++] = item0; \
      } \
      photon_##_name##_t result = malloc(sizeof(struct photon_##_name)); \
      if (result == NULL) { \
          fprintf(stderr, "malloc failed\n"); \
          return NULL; \
      } \
      result->items = items; \
      result->length = length; \
      return result; \
    }
