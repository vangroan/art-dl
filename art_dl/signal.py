
class SignalRunningException(Exception):
    pass


class CallbackNotFoundException(Exception):
    pass


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
        self._slots.append(slot)
        return slot

    def send(self, *args, **kwargs):
        self._running = True
        for slot in self._slots:
            slot.execute(*args, **kwargs)
        self._running = False

    def remove(self, func):
        filtered = list(filter(lambda el: el.callback == func, self._slots))

        if len(filtered) == 0:
            raise CallbackNotFoundException('Callback not found in signal')

        for slot in filtered:
            self.remove_slot(slot)

    def remove_slot(self, slot):
        if self.running:
            raise SignalRunningException(
                    'Slot cannot be detached while signalling')
        self._slots = list(filter(lambda el: el != slot, self._slots))
        slot.destroy()

    def __len__(self):
        return len(self._slots)
