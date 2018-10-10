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
cooksum = 0
eatsum = 0

class Shared:
    def __init__(self):
        self.servings = 0
        self.mutex = Semaphore(1)
        self.emptyPot = Semaphore(0)
        self.fullPot = Semaphore(0)

def cook_code(shared):
    global cooksum
    while True:
        shared.emptyPot.wait()

        #print "Cooking : the OTHER other white meat!"
        cooksum += 1
        #time.sleep(random.random())
        shared.servings = M

        shared.fullPot.signal()

def savages_code(shared):
    global eatsum
    while True:
        shared.mutex.wait()
        if shared.servings == 0:
            shared.emptyPot.signal()
            shared.fullPot.wait()
        shared.servings -= 1
        shared.mutex.signal()

        #print "Eating : yum yum yum!"
        eatsum += 1
        #time.sleep(random.random())

class Monitor(threading.Thread):
    def run(self):
        iterations = 10
        CPUusage = []
        VMusage = []
        for i in range(iterations):
            CPU = psutil.cpu_percent(0.2, False)
            print "CPU usage:", CPU
            CPUusage.append(CPU)
            VM = psutil.virtual_memory()
            print VM
            VMusage.append(VM[2])
            time.sleep(1)
        CPUavg = sum(CPUusage)/float(len(CPUusage))
        print "Average CPU usage:", CPUavg
        VMavg = sum(VMusage)/float(len(VMusage))
        print "Average VM usage:", VMavg
        print "Total cooks:", cooksum, "Total eats:", eatsum
        thread.interrupt_main()

Monitor().start()
time.sleep(1)

shared = Shared()
cook = Thread(cook_code, shared)
savages = [Thread(savages_code, shared) for i in range(30)]
for savage in savages: savage.join()
