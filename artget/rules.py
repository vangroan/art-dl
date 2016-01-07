
from artget.rulematch import PatternRules
from artget.scrapers.deviantart import DeviantartScraper
from artget.scrapers.drawcrowd import DrawcrowdScraper


# TODO: Declarative syntax
def configure_rules(resolver):

    # Deviantart
    resolver.add_rule(r'(?P<username>[\w\d-]+)\.deviantart.com', DeviantartScraper.create_scraper, inject_context=True)

    # Drawcrowd
    resolver.add_rule(r'drawcrowd.com/(?P<username>[\w\d-]+)', DrawcrowdScraper.create_scraper, inject_context=True)