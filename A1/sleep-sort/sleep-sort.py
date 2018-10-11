"""
    sleep sort in Python
    from https://gist.github.com/ssirai/4115014
"""

import thread
import sys
from time import sleep
from random import shuffle
from threading import Thread
from numpy import arange
import psutil

r = []
class worker(Thread):
    def __init__(self, t):
        Thread.__init__(self)
        self.t = t

    def run(self):
        global r
        sleep(self.t)
        self.t = "%0.2f" % self.t
        r.append(self.t)
        print self.t
        sys.stdout.flush()
        
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
            sleep(1)
        CPUavg = sum(CPUusage)/float(len(CPUusage))
        print "Average CPU usage:", CPUavg
        VMavg = sum(VMusage)/float(len(VMusage))
        print "Average VM usage:", VMavg
        #thread.interrupt_main()

Monitor().start()
sleep(1)

high = float(sys.argv[1])
step = float(sys.argv[2])
xs = [i for i in arange(0, high, step)]
shuffle(xs)

xs = map(worker, xs)
for x in xs:
    x.start()

for x in xs:
    x.join()

print "Sorted:", all(r[i] <= r[i+1] for i in xrange(len(r)-1))
