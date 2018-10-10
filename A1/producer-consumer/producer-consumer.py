"""
producer-consumer problem in Python
from https://www.agiliq.com/blog/2013/10/producer-consumer-problem-in-python/
"""

from threading import Thread, Condition
import thread
import time
import random
import psutil

queue = []
MAX_NUM = 10
condition = Condition()

class ProducerThread(Thread):
    def run(self):
        nums = range(5)
        global queue
        while True:
            condition.acquire()
            if len(queue) == MAX_NUM:
                print "Queue full, producer is waiting"
                condition.wait()
                print "Space in queue, Consumer notified the producer"
            num = random.choice(nums)
            queue.append(num)
            print "Produced", num
            condition.notify()
            condition.release()
            time.sleep(random.random())


class ConsumerThread(Thread):
    def run(self):
        global queue
        while True:
            condition.acquire()
            if not queue:
                print "Nothing in queue, consumer is waiting"
                condition.wait()
                print "Producer added something to queue and notified the consumer"
            num = queue.pop(0)
            print "Consumed", num
            condition.notify()
            condition.release()
            time.sleep(random.random())

class Monitor(Thread):
    def run(self):
        for i in range(10):
            print "CPU usage:", psutil.cpu_percent()
            print psutil.virtual_memory()
            time.sleep(5)
        thread.interrupt_main()

Monitor().start()
time.sleep(1)

ProducerThread().start()
ConsumerThread().start()
