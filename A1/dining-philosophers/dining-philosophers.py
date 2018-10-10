import thread
import threading
import random
import time
import psutil

# Dining philosophers, 5 Phillies with 5 forks. Must have two forks to eat.
#
# Deadlock is avoided by never waiting for a fork while holding a fork (locked)
# Procedure is to do block while waiting to get first fork, and a nonblocking
# acquire of second fork.  If failed to get second fork, release first fork,
# swap which fork is first and which is second and retry until getting both.
#
# See discussion page note about 'live lock'.
# from https://rosettacode.org/wiki/Dining_philosophers#Python

dinesum = 0
thinksum = 0

class Philosopher(threading.Thread):

    running = True

    def __init__(self, xname, forkOnLeft, forkOnRight):
        threading.Thread.__init__(self)
        self.name = xname
        self.forkOnLeft = forkOnLeft
        self.forkOnRight = forkOnRight

    def run(self):
        global thinksum
        while(self.running):
            #  Philosopher is thinking (but really is sleeping).
            thinksum += 1
            #time.sleep(1)
            #print '%s is hungry.' % self.name
            self.dine()

    def dine(self):
        fork1, fork2 = self.forkOnLeft, self.forkOnRight

        while self.running:
            fork1.acquire(True)
            locked = fork2.acquire(False)
            if locked: break
            fork1.release()
            #print '%s swaps forks' % self.name
            fork1, fork2 = fork2, fork1
        else:
            return

        self.dining()
        fork2.release()
        fork1.release()

    def dining(self):
        global dinesum
        #print '%s starts eating '% self.name
        dinesum += 1
        #time.sleep(1)
        #print '%s finishes eating and leaves to think.' % self.name

def DiningPhilosophers():
    forks = [threading.Lock() for n in range(5)]
    philosopherNames = ('Aristotle','Kant','Buddha','Marx', 'Russell')

    philosophers= [Philosopher(philosopherNames[i], forks[i%5], forks[(i+1)%5]) \
            for i in range(5)]

    random.seed(507129)
    Philosopher.running = True
    for p in philosophers: p.start()
    #time.sleep(100)
    #Philosopher.running = False
    #print ("Now we're finishing.")

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
        print "Total dines:", dinesum, "Total thinks:", thinksum
        thread.interrupt_main()

Monitor().start()
time.sleep(1)

DiningPhilosophers()
