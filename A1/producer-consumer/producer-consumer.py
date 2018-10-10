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
psum = 0
csum = 0

class ProducerThread(Thread):
    def run(self):
        nums = range(5)
        global queue
        global psum
        while True:
            condition.acquire()
            if len(queue) == MAX_NUM:
                #print "Queue full, producer is waiting"
                condition.wait()
                #print "Space in queue, Consumer notified the producer"
            num = random.choice(nums)
            queue.append(num)
            #print "Produced", num
            psum += 1
            condition.notify()
            condition.release()

class ConsumerThread(Thread):
    def run(self):
        global queue
        global csum
        while True:
            condition.acquire()
            if not queue:
                #print "Nothing in queue, consumer is waiting"
                condition.wait()
                #print "Producer added something to queue and notified the consumer"
            num = queue.pop(0)
            #print "Consumed", num
            csum += 1
            condition.notify()
            condition.release()

class Monitor(Thread):
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
        print "Total productions:", psum, "Total consumptions:", csum
        thread.interrupt_main()

Monitor().start()
time.sleep(1)

ProducerThread().start()
ConsumerThread().start()
