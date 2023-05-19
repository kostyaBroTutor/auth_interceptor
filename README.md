# Auth Interceptor

This repository stores a implementation of the gRPC interceptor,
that check auth roles or permissions, and authenticate user. \
This repository created for the [habr post(add link)](example.com).


## Cloning repository

This repository contains submodules, so you need to clone it with `--recursive-submodules` flag:

```bash
git clone --recursive https://github.com/kostyaBroTutor/auth_interceptor.git
```

If u want to update submodules, use:

```bash
git submodule update --init --recursive
```
