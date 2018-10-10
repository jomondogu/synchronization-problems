"""
unisex bathroom problem in Python
based on code from https://www.greenteapress.com/
"""

from threading_cleanup import *
import time
import random

students = 10

class Shared:
    def __init__(self):
        self.heathens = 0
        self.prudes = 0
        self.status = 'neutral'
        self.mutex = Semaphore(1)
        self.heathenTurn = Semaphore(1)
        self.prudeTurn = Semaphore(1)
        self.heathenQueue = Semaphore(0)
        self.prudeQueue = Semaphore(0)

def heathen_code(shared):
    while True:
        shared.heathenTurn.wait()
        shared.heathenTurn.signal()

        shared.mutex.wait()
        shared.heathens = shared.heathens + 1

        if shared.status == 'neutral':
            shared.status = 'heathens rule'
            shared.mutex.signal()
        elif shared.status == 'prudes rule':
            if shared.heathens > shared.prudes:
                shared.status = 'transition to heathens'
                shared.prudeTurn.wait()
            shared.mutex.signal()
            shared.heathenQueue.wait()
        elif shared.status == 'transition to heathens':
            shared.mutex.signal()
            shared.heathenQueue.wait()
        else:
            shared.mutex.signal()

        print("Heathen crossing the field! PRUDES WATCH OUT")
        time.sleep(random.random())

        shared.mutex.wait()
        shared.heathens = shared.heathens - 1

        if shared.heathens == 0:
            if shared.status == 'transition to prudes':
                shared.prudeTurn.signal()
            if shared.prudes:
                shared.prudeQueue.signal(shared.prudes)
                shared.status = 'prudes rule'
            else:
                shared.status = 'neutral'

        if shared.status == 'heathens rule':
            if shared.prudes > shared.heathens:
                shared.status = 'transition to prudes'
                shared.heathenTurn.wait()

        shared.mutex.signal()

def prude_code(shared):
    while True:
        shared.prudeTurn.wait()
        shared.prudeTurn.signal()

        shared.mutex.wait()
        shared.prudes = shared.prudes + 1

        if shared.status == 'neutral':
            shared.status = 'prudes rule'
            shared.mutex.signal()
        elif shared.status == 'heathens rule':
            if shared.prudes > shared.heathens:
                shared.status = 'transition to prudes'
                shared.heathenTurn.wait()
            shared.mutex.signal()
            shared.prudeQueue.wait()
        elif shared.status == 'transition to prudes':
            shared.mutex.signal()
            shared.prudeQueue.wait()
        else:
            shared.mutex.signal()

        print("Prude crossing the field! HEATHENS WATCH OUT")
        time.sleep(1)

        shared.mutex.wait()
        shared.prudes = shared.prudes - 1

        if shared.prudes == 0:
            if shared.status == 'transition to heathens':
                shared.heathenTurn.signal()
            if shared.heathens:
                shared.heathenQueue.signal(shared.heathens)
                shared.status = 'heathens rule'
            else:
                shared.status = 'neutral'

        if shared.status == 'prudes rule':
            if shared.heathens > shared.prudes:
                shared.status = 'transition to heathens'
                shared.prudeTurn.wait()

        shared.mutex.signal()

shared = Shared()
heathen = Thread(heathen_code, shared)
prude = Thread(prude_code, shared)
heathen.join()
prude.join()
