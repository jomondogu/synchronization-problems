"""
dining savages problem in Python
uses threading_cleanup.py from http://greenteapress.com/semaphores/threading_cleanup.py
"""

import thread
from threading_cleanup import *
import time
import random
import psutil

M = 50

class Shared:
    def __init__(self):
        self.servings = 0
        self.mutex = Semaphore(1)
        self.emptyPot = Semaphore(0)
        self.fullPot = Semaphore(0)

def cook_code(shared):
    while True:
        shared.emptyPot.wait()

        print "Cooking : the OTHER other white meat!"
        time.sleep(random.random())
        shared.servings = M

        shared.fullPot.signal()

def savages_code(shared):
    while True:
        shared.mutex.wait()
        if shared.servings == 0:
            shared.emptyPot.signal()
            shared.fullPot.wait()
        shared.servings -= 1
        shared.mutex.signal()

        print "Eating : yum yum yum!"
        time.sleep(random.random())

class Monitor(Thread):
    def run(self):
        for i in range(10):
            print psutil.cpu_times()
            print "CPU usage:", psutil.cpu_percent()
            print psutil.virtual_memory()
            time.sleep(5)
        thread.interrupt_main()

Monitor().start()
time.sleep(1)

shared = Shared()
cook = Thread(cook_code, shared)
savages = [Thread(savages_code, shared) for i in range(30)]
for savage in savages: savage.join()
