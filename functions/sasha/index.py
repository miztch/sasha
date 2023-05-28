import random
import re
import time
from datetime import datetime

import requests
from dateutil import tz
from dateutil.parser import parse
from selectolax.parser import HTMLParser
from utils import headers


def sleep():
    '''
    sleep for 1~10 secs (randomly)
    '''
    sec = random.randint(1, 10)
    time.sleep(sec)

# TODO: Cache in scrape_match()


def scrape_event(event_url_path):
    '''
    scrape event page of url_path
    '''
    url = 'https://www.vlr.gg{}'.format(event_url_path)
    resp = requests.get(url, headers=headers)
    html = HTMLParser(resp.text)

    event_id = event_url_path.split('/')[2]

    event_name = html.css_first('.wf-title').text().strip()
    event_name = event_name.replace('\t', '').replace('\n', '')

    country_flag = html.css_first(
        '.event-desc-item-value .flag').attributes['class']
    country_flag = country_flag.replace(' mod-', '_').replace('flag_', '')

    data = {
        'event_id': event_id,
        'event_name': event_name,
        'country_flag': country_flag
    }

    return data


def scrape_match(match_url_path):
    '''
    scrape match page of url_path
    '''
    url = 'https://www.vlr.gg{}'.format(match_url_path)
    resp = requests.get(url, headers=headers)
    html = HTMLParser(resp.text)

    match_id = match_url_path.split('/')[1]

    match_name = html.css_first('.match-header-event-series').text()
    match_name = match_name.replace('\t', '').replace('\n', '')

    start_time = html.css_first('.moment-tz-convert').attributes['data-utc-ts']
    with_timezone = ' '.join([start_time, 'EST'])

    tzinfo = {'EST': tz.gettz('America/New_York'),
              'CST': tz.gettz('America/Chicago')}
    start_time_est = parse(with_timezone, tzinfos=tzinfo)
    start_time_utc = start_time_est.astimezone(tz.gettz('Etc/GMT'))
    start_time_utc = datetime.strftime(start_time_utc, '%Y-%m-%dT%H:%M:%S%z')

    teams = html.css('.wf-title-med')
    teams = [t.text().replace('\t', '').replace('\n', '') for t in teams]

    best_of = html.css('.match-header-vs-note')[-1].text()
    best_of = best_of.replace('Bo', '').replace('\t', '').replace('\n', '')

    event_url_path = html.css_first('a.match-header-event').attributes['href']
    event_info = scrape_event(event_url_path)

    data = {
        'match_id': match_id,
        'event_name': event_info['event_name'],
        'event_country_flag': event_info['country_flag'],
        'start_time': start_time_utc,
        'best_of': best_of,
        'match_name': match_name,
        'teams': teams
    }
    return data


def scrape_matches(page: str = 1):
    '''
    scrape /matches page
    '''
    url = 'https://www.vlr.gg/matches?page={}'.format(page)
    resp = requests.get(url, headers=headers)
    html = HTMLParser(resp.text)

    matches = []

    for item in html.css('a.wf-module-item'):
        match_url_path = item.attributes['href']

        sleep()
        match_detail = scrape_match(match_url_path)

        matches.append(
            {
                'id': match_detail['match_id'],
                'eventName': match_detail['event_name'],
                'eventCountryFlag': match_detail['event_country_flag'],
                'startTime': match_detail['start_time'],
                'bestOf': match_detail['best_of'],
                'matchName': match_detail['match_name'],
                'teams': [{'title': team} for team in match_detail['teams']],
                'match_page': match_url_path
            }
        )

    return matches


def lambda_handler(event, context):

    matches = scrape_matches()
    return matches
