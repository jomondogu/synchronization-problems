"""
sushi bar problem in Python
uses threading_cleanup.py from http://greenteapress.com/semaphores/threading_cleanup.py
"""

from threading_cleanup import *
import time
import random

crowd = 5
maximum = 5

class Shared:
    def __init__(self):
        self.eating = 0
        self.waiting = 0
        self.mutex = Semaphore(1)
        self.block = Semaphore(0)
        self.must_wait = False

def customer_code(shared):
    while True:
        shared.mutex.wait()
        if shared.must_wait:
            shared.waiting += 1
            shared.mutex.signal()
            print("Customer is waiting for a spot.")
            shared.block.wait()
        else:
            shared.eating += 1
            shared.must_wait = (shared.eating == maximum)
            shared.mutex.signal()

        print("Customer is eating sushi!")
        time.sleep(random.random())

        shared.mutex.wait()
        shared.eating -= 1
        if shared.eating == 0:
            n = min(maximum, shared.waiting)
            shared.waiting -= n
            shared.eating += n
            shared.must_wait = (shared.eating == maximum)
            print("Customer is leaving...")
            shared.block.signal(n)
        shared.mutex.signal()

shared = Shared()
customers = [Thread(customer_code, shared) for i in range(crowd)]
for customer in customers: customer.join()
