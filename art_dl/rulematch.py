
from collections import OrderedDict
import re


class RuleException(Exception):
    pass


class RegexRule(object):

    def __init__(self, pattern, handler, inject_context):

        self._original_string = pattern
        self.pattern = re.compile(pattern)
        self.handler = handler
        self.inject_context = inject_context

    @staticmethod
    def _split_groups(regex_match):
        """Split groups in match into unnamed and named
        """
        # http://stackoverflow.com/questions/30293064/get-all-unnamed-groups-in-a-python-match-object

        # Keep ordering so unnamed args will be added to
        # the tuple in the order they appear in the pattern
        named = OrderedDict()
        unnamed = OrderedDict()
        all_groups = regex_match.groups()
        groupdict = regex_match.groupdict()

        # Index every named group by its span
        for k, v in groupdict.items():
            named[regex_match.span(k)] = v

        # Index every other group by its span, skipping groups with same
        # span as a named group
        for i, v in enumerate(all_groups):
            sp = regex_match.span(i + 1)
            if sp not in named:
                unnamed[sp] = v

        return tuple(unnamed.values()), groupdict

    def create_context(self, input_url, args, kwargs):
        return {
            'pattern': self._original_string,
            'input': input_url,
            'args': args,
            'kwargs': kwargs,
        }

    def execute(self, url, context_processor=None):
        matched = self.pattern.search(url)
        if matched:
            groups = matched.groups()
            args, kwargs = self._split_groups(matched)
            response = None

            if self.inject_context:
                context = self.create_context(url, args, kwargs)
                if context_processor:
                    context_processor(context)
                response = self.handler(context, *args, **kwargs)
            else:
                response = self.handler(*args, **kwargs)

            if response is None:
                response = True

            return response


# TODO: Rename class to RuleResolver
class PatternRules(object):

    def __init__(self):

        self._rules = []

    def add_rule(self, pattern, handler, inject_context=False):

        self._rules.append(RegexRule(
            pattern,
            handler,
            inject_context
        ))

    def add_rules(self, rules):
        for r in rules:
            if not isinstance(r, RegexRule):
                raise TypeError("Attempt to add rule of type %s" % type(r))
            self._rules.append(r)

    def dispatch(self, url, context_processor=None):
        for rule in self._rules:
            response = rule.execute(url, context_processor)
            if response:
                return response
        else:
            raise RuleException('No pattern matched "%s"' % url)

    def __len__(self):
        return len(self._rules)

    def __iter__(self):
        for rule in self._rules:
            yield rule


def rule(pattern, handler, inject_context=False):
    return RegexRule(pattern, handler, inject_context)
