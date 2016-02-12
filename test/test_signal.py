
from unittest import TestCase
import sys, os

sys.path.append(os.pardir)
from art_dl.signal import Signal, Slot, SignalRunningException

class SignallingTest(TestCase):

    def setUp(self):
        pass

    def test_add_func_to_signal_via_slot_contents(self):
        sig = Signal()
        def go():
            x = 1 + 2 # Do nothing
        slot = sig.add(go)
        self.assertEqual(slot.callback, go)

    def test_signal_sending_messages(self):
        sig = Signal()
        count = 0
        def go():
            nonlocal count
            count += 1
        slot = sig.add(go)
        [sig.send() for i in range(3)]
        self.assertEqual(count, 3)

    def test_signal_argument_passing(self):
        sig = Signal()
        count = 3
        def go(val):
            nonlocal count
            count += val
        slot = sig.add(go)
        sig.send(7)
        self.assertEqual(count, 10)

    def test_signal_args_and_kwargs(self):
        sig = Signal()
        count = 11
        data = 'foo'
        def go(val, say='default_string'):
            nonlocal count, data
            count += val
            data = say
        slot = sig.add(go)
        sig.send(3, say='bar')
        self.assertEqual(count, 14)
        self.assertEqual(data, 'bar')

    def test_slot_removal(self):
        sig = Signal()
        count = 3
        def go(val):
            nonlocal count
            count += val
        slot = sig.add(go)
        sig.send(5) # Count is now 8
        sig.remove_slot(slot)
        sig.send(7) # If count is 15, then removal failed
        self.assertEqual(count, 8)

    def test_slot_detach(self):
        sig = Signal()
        count = 7
        def go(val):
            nonlocal count
            count += val
        slot = sig.add(go)
        sig.send(3) # Count is now 10
        slot.detach()
        sig.send(7) # If count is 17, then removal failed
        self.assertEqual(count, 10)

    def test_cannot_remove_slot_while_signalling(self):

        class BadRemove:
            def __init__(self, signal):
                self.slot = signal.add(self.do_remove)
            def do_remove(self):
                self.slot.detach()

        sig = Signal()
        bad_remover = BadRemove(sig)
        with self.assertRaises(SignalRunningException):
            sig.send()

    def test_remove_function_instance_directly(self):
        sig = Signal()
        count = 4
        def go(val):
            nonlocal count
            count += val
        slot = sig.add(go)
        sig.send(3) # count is now 7
        sig.remove(go)
        sig.send(6) # If count is 13, then removal failed
        self.assertEqual(count, 7)
