/* correcthorse - a passphrase generator inspired by http://xkcd.com/936/
 *
 * Copyright (c) 2012-2013 Robin Martinjak
 *
 * This work is free. You can redistribute it and/or modify it under the
 * terms of the Do What The Fuck You Want To Public License, Version 2,
 * as published by Sam Hocevar. See the COPYING file for more details.
 */

#ifndef WLIST_H
#define WLIST_H

#include <stdlib.h>

struct wlist
{
    struct word *head;
    struct word *tail;
    size_t len;
};

struct word
{
    char *word;
    struct word *next;
};

struct word *word_init(const char *word);
void word_free(struct word *w);

struct wlist *wlist_init(void);
void wlist_free(struct wlist *wl);
size_t wlist_len(struct wlist *wl);

int wlist_add(struct wlist *wl, const char *word);
char *wlist_get(struct wlist *wl, size_t n);

struct wlist *wlist_read(const char *name);

#endif
