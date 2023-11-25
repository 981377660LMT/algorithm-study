
from distutils.core import setup, Extension
module = Extension(
    "multiset",
    sources=["multiset_wrapper.cpp"],
    extra_compile_args=["-O3", "-march=native"]
)
setup(
    name="MultiSetMethod",
    version="0.0.4",
    description="wrapper for C++ multiset",
    ext_modules=[module]
)
