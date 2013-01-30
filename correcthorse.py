#!/usr/bin/env python3

""" correcthorse - a passphrase generator inspired by http://xkcd.com/936/

Copyright (c) 2012-2013 Robin Martinjak <rob@rmartinjak.de>

This work is free. You can redistribute it and/or modify it under the
terms of the Do What The Fuck You Want To Public License, Version 2,
as published by Sam Hocevar. See the COPYING file for more details.
"""

import os.path
import random

CHARS_MIN_DEFAULT = 12
WORDS_MIN_DEFAULT = 4
WORDLISTS_DEFAULT = ['english']

WORDLIST_DIR = '/usr/share/correcthorse'


def read_wordlist(name):
    if not os.path.exists(name):
        name = os.path.join(WORDLIST_DIR, name)

    words = []

    with open(name, 'r') as f:
        for line in f:
            words.append(line.strip())

    return words


def make_passphrases(
        count=1,
        chars_min=CHARS_MIN_DEFAULT,
        words_min=WORDS_MIN_DEFAULT,
        wordlists=None,
        userwords=None,
        camelcase=False,
        sep=''):

    if not wordlists:
        wordlists = WORDLISTS_DEFAULT[:]

    words = []
    chars = 0
    wlistwords = []

    # add user-specified words
    if userwords:
        words.extend(userwords)
        for u in userwords:
            chars += len(u)

    # load words from wordlists
    if len(words) < words_min or chars < chars_min:
        for w in wordlists:
            wlistwords.extend(read_wordlist(w))

    # generate `count` passphrases
    for i in range(count):
        w = words[:]
        c = chars

        while len(w) < words_min or c < chars_min:
            word = random.choice(wlistwords)
            w.append(word)
            c += len(word)

        if camelcase:
            w = [word.capitalize() for word in w]

        random.shuffle(w)

        yield sep.join(w)


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(
        description='Generate passphrase(s) consisting of random words.')

    parser.add_argument(
        '-c', '--chars',
        type=int,
        dest='chars_min',
        metavar='N',
        default=CHARS_MIN_DEFAULT,
        help='minimal number of characters')

    parser.add_argument(
        '-w', '--words',
        type=int,
        dest='words_min',
        metavar='N',
        default=WORDS_MIN_DEFAULT,
        help='minimal number of words')

    parser.add_argument(
        '-i', '--include',
        dest='userwords',
        action='append',
        metavar='WORD',
        default=None,
        help='include WORD in passphrase')

    parser.add_argument(
        '-l', '--list',
        dest='wordlists',
        action='append',
        metavar='LIST',
        default=None,
        help='use words from LIST')

    parser.add_argument(
        '-u', '--camelcase',
        action='store_true',
        help='print words in CamelCase')

    parser.add_argument(
        '-s', '--sep',
        default='',
        help='separate words with SEP')

    parser.add_argument('count', nargs='?', type=int, default=1)

    args = parser.parse_args()
    for p in make_passphrases(**vars(args)):
        print(p)
