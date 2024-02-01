import json
import logging
import os
import random
import re
import time
from datetime import datetime

import boto3
import requests
from dateutil import tz
from dateutil.parser import parse
from selectolax.parser import HTMLParser
from utils import headers

logger = logging.getLogger()
logger.setLevel(logging.INFO)

dynamodb = boto3.resource("dynamodb")
table = dynamodb.Table(os.environ["VLR_MATCHES_TABLE"])

# vlr match events cache
vlr_events_cache = {}


def insert(table, matches):
    """
    put items into specified DynamoDB table.
    """
    with table.batch_writer() as batch:
        for match in matches:
            logger.info("put match info into the table: {}".format(match))
            batch.put_item({k: v for k, v in match.items()})


def sleep():
    """
    sleep for 1~10 secs (randomly)
    """
    sec = random.randint(1, 10)
    time.sleep(sec)


def get_event_from_cache(event_url_path):
    global vlr_events_cache

    result = ""

    if event_url_path in vlr_events_cache:
        result = vlr_events_cache[event_url_path]

    return result


def scrape_event(event_url_path):
    """
    scrape event page of url_path
    """
    global vlr_events_cache

    url = "https://www.vlr.gg{}".format(event_url_path)
    logger.info("get event info: {}".format(url))

    resp = requests.get(url, headers=headers)
    html = HTMLParser(resp.text)

    event_id = int(event_url_path.split("/")[2])

    event_name = html.css_first(".wf-title").text().strip()
    event_name = event_name.replace("\t", "").replace("\n", "")

    country_flag = html.css_first(".event-desc-item-value .flag").attributes["class"]
    country_flag = country_flag.replace(" mod-", "_").replace("flag_", "")

    data = {
        "event_id": event_id,
        "event_name": event_name,
        "country_flag": country_flag,
    }

    # caching
    vlr_events_cache[event_url_path] = data

    return data


def scrape_match(match_url_path):
    """
    scrape match page of url_path
    """
    global vlr_events_cache

    url = "https://www.vlr.gg{}".format(match_url_path)
    logger.info("get match info: {}".format(url))

    resp = requests.get(url, headers=headers)
    html = HTMLParser(resp.text)

    match_id = int(match_url_path.split("/")[1])

    match_name = html.css_first(".match-header-event-series").text()
    match_name = match_name.replace("\t", "").replace("\n", "")

    start_time = html.css_first(".moment-tz-convert").attributes["data-utc-ts"]
    with_timezone = " ".join([start_time, "EST"])

    tzinfo = {"EST": tz.gettz("America/New_York"), "CST": tz.gettz("America/Chicago")}
    start_time_est = parse(with_timezone, tzinfos=tzinfo)
    start_time_utc = start_time_est.astimezone(tz.gettz("Etc/GMT"))
    start_time_utc = datetime.strftime(start_time_utc, "%Y-%m-%dT%H:%M:%S%z")

    teams = html.css(".wf-title-med")
    teams = [t.text().replace("\t", "").replace("\n", "") for t in teams]

    best_of = html.css(".match-header-vs-note")[-1].text()
    best_of = best_of.replace("Bo", "").replace(" Maps", "")
    best_of = best_of.replace("\t", "").replace("\n", "")
    best_of = int(best_of)

    event_url_path = html.css_first("a.match-header-event").attributes["href"]

    if event_url_path in vlr_events_cache:
        logger.info("get event info from cache: {}".format(event_url_path))
        event_info = vlr_events_cache[event_url_path]
    else:
        logger.info("get event info from website: {}".format(event_url_path))
        event_info = scrape_event(event_url_path)

    data = {
        "match_id": match_id,
        "event_name": event_info["event_name"],
        "event_country_flag": event_info["country_flag"],
        "start_time": start_time_utc,
        "best_of": best_of,
        "match_name": match_name,
        "teams": teams,
    }
    return data


def scrape_matches(page: str = 1):
    """
    scrape /matches page
    """
    url = "https://www.vlr.gg/matches?page={}".format(page)
    logger.info("fetch matches list from: {}".format(url))

    resp = requests.get(url, headers=headers)
    html = HTMLParser(resp.text)

    matches = []

    for item in html.css("a.wf-module-item"):
        match_url_path = item.attributes["href"]

        sleep()
        match_detail = scrape_match(match_url_path)

        item = {
            "id": match_detail["match_id"],
            "eventName": match_detail["event_name"],
            "eventCountryFlag": match_detail["event_country_flag"],
            "startTime": match_detail["start_time"],
            "bestOf": match_detail["best_of"],
            "matchName": match_detail["match_name"],
            "teams": [{"title": team} for team in match_detail["teams"]],
            "pagePath": match_url_path,
        }
        logger.info("add match to the list: {}".format(item))
        matches.append(item)

    return matches


def lambda_handler(event, context):
    page = str(event["page"])

    matches = scrape_matches(page)
    match_list.extend(matches)

    insert(table, matches)

    return {"matches_count": len(matches)}
