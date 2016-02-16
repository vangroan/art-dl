import unittest
from unittest.mock import Mock
import sys, os
import re

sys.path.insert(0, os.path.join(os.pardir))
from art_dl.rulematch import RegexRule, PatternRules, RuleException


class RegexRuleTests(unittest.TestCase):
    def test_group_split(self):
        pattern = re.compile(r'(?P<title>[a-z]+)/([0-9]+)')
        args, kwargs = RegexRule._split_groups(pattern.search('foo/0042'))

        self.assertEqual('0042', args[0])
        self.assertEqual('foo', kwargs['title'])


class PatternRulesTests(unittest.TestCase):
    def test_adding_rules(self):
        rules = PatternRules()
        rules.add_rule(r'[a-z]', lambda: True)

        self.assertEqual(1, len(rules))

    def test_unhandled_exception(self):
        rules = PatternRules()
        rules.add_rule(r'[a-z]', lambda: True)

        with self.assertRaises(RuleException):
            rules.dispatch('123')

    def test_basic_dispatch(self):
        rules = PatternRules()
        f = Mock()
        rules.add_rule(r'http://my\.test\.com/', f)
        rules.dispatch('http://my.test.com/')

        f.assert_any_calls()

    def test_dispatch_with_args(self):
        rules = PatternRules()
        f = Mock()
        rules.add_rule(r'^http://my\.test\.com/([a-zA-Z]+)/([0-9]{4})/$', f)
        rules.dispatch('http://my.test.com/foo/0042/')

        f.assert_called_with('foo', '0042')

    def test_dispatch_with_kwrags(self):
        rules = PatternRules()
        f = Mock()
        rules.add_rule(r'^http://my\.test\.com/(?P<title>[a-zA-Z]+)/([0-9]{4})/$', f)
        rules.dispatch('http://my.test.com/foo/0042/')

        f.assert_called_with('0042', title='foo')

    def test_inject_context(self):
        rules = PatternRules()
        f = Mock()
        rules.add_rule(r'[a-z]+/([0-9]+)', f, inject_context=True)
        rules.dispatch('foo/0042')

        f.assert_called_once_with({
            'input': 'foo/0042',
            'pattern': r'[a-z]+/([0-9]+)',
            'args': ('0042',),
            'kwargs': {}
        }, '0042')

    def test_dispatch_context_processor(self):
        rules = PatternRules()
        value = None

        def processor(ctx):
            ctx['test_processor'] = 'foobar'

        def handler(ctx):
            nonlocal value
            value = ctx['test_processor']

        rules.add_rule(r'[a-z]+', handler, inject_context=True)
        rules.dispatch('foo', context_processor=processor)

        self.assertEqual(value, 'foobar')

    def test_dispatch_return(self):
        rules = PatternRules()

        def factory():
            return ('foo', 'bar')

        rules.add_rule(r'[a-z]+', factory)
        result = rules.dispatch('foo')

        self.assertEqual(('foo', 'bar'), result)


if __name__ == '__main__':
    unittest.main()
