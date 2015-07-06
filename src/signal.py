
class Slot:

    def __init__(self, signal, callback):
        self._signal = signal
        self._callback = callback

    @property
    def signal(self):
        return self._signal

    @property
    def callback(self):
        return self._callback

    def execute(self, *args, **kwargs):
        cb = self._callback
        if cb:
            cb(*args, **kwargs)

    def detach(self):
        self.signal.remove_slot(self)

    def destroy(self):
        self._signal = None
        self._callback = None

class Signal:

    def __init__(self):
        self._slots = []
        self._running = False

    @property
    def running(self):
        return self._running

    @property
    def slots(self):
        return self._slots

    def add(self, func):
        slot = Slot(self, func)
        self._slots.append(func)
        return slot

    def send(self, *args, **kwargs):
        for slot in self._slots:
            slot.execute(*args, **kwargs)

    def remove(self, func):
        raise NotImplementedError('Todo')

    def remove_slot(self, func):
        raise NotImplementedError('Todo')
