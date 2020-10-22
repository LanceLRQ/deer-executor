# -*- coding: utf-8 -*-

while True:
    try:
        a, b = map(int, raw_input().split())
    except:
        break

    print(a + b)