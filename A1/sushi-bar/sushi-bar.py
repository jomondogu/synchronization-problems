"""
sushi bar problem in Python
uses threading_cleanup.py from http://greenteapress.com/semaphores/threading_cleanup.py
"""

import thread
from threading_cleanup import *
import time
import random
import psutil

crowd = 5
maximum = 5
eatsum = 0
waitsum = 0

class Shared:
    def __init__(self):
        self.eating = 0
        self.waiting = 0
        self.mutex = Semaphore(1)
        self.block = Semaphore(0)
        self.must_wait = False

def customer_code(shared):
    global waitsum
    global eatsum
    while True:
        shared.mutex.wait()
        if shared.must_wait:
            shared.waiting += 1
            shared.mutex.signal()
            #print("Customer is waiting for a spot.")
            waitsum += 1
            shared.block.wait()
        else:
            shared.eating += 1
            shared.must_wait = (shared.eating == maximum)
            shared.mutex.signal()

        #print("Customer is eating sushi!")
        #time.sleep(random.random())
        eatsum += 1

        shared.mutex.wait()
        shared.eating -= 1
        if shared.eating == 0:
            n = min(maximum, shared.waiting)
            shared.waiting -= n
            shared.eating += n
            shared.must_wait = (shared.eating == maximum)
            #print("Customer is leaving...")
            shared.block.signal(n)
        shared.mutex.signal()

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
        print "Total wait loops:", waitsum, "Total eats:", eatsum
        thread.interrupt_main()

Monitor().start()
time.sleep(1)

shared = Shared()
customers = [Thread(customer_code, shared) for i in range(crowd)]
for customer in customers: customer.join()
