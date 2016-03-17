
from art_dl.scrapers.artstation import ArtstationScraper
from art_dl.scrapers.deviantart import DeviantartScraper
from art_dl.scrapers.drawcrowd import DrawcrowdScraper
from art_dl.rulematch import rule

rules = (
    rule(r'artstation\.com/artist/(?P<username>[\w\d-]+)/?', ArtstationScraper.create_scraper, inject_context=True),
    rule(r'(?P<username>[\w\d-]+)\.deviantart.com', DeviantartScraper.create_scraper, inject_context=True),
    rule(r'drawcrowd.com/(?P<username>[\w\d-]+)', DrawcrowdScraper.create_scraper, inject_context=True),
)
