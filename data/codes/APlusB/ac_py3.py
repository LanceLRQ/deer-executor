# -*- coding: utf-8 -*-
# coding:utf-8

import sys
import codecs
sys.stdout = codecs.getwriter("utf-8")(sys.stdout.detach())

while True:
    try:
        a, b = map(int, input().split())
    except:
        break

    print(a + b)