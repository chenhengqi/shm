#pragma once

#include <sys/stat.h>
#include <sys/types.h>

int posix_create_shm(const char* name, int flag, mode_t mode, off_t length, void** addr);
int posix_destroy_shm(const char* name, void *addr, size_t length);
