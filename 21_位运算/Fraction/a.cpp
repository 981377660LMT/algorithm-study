/* testlib.c */
double f(int n, double *x, void *user_data) {
  double c = *(double *)user_data;
  return c + x[0] - x[1] * x[2]; /* corresponds to c + x - y * z */
}