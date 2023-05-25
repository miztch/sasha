import re
from datetime import datetime

import requests
from selectolax.parser import HTMLParser

from utils import headers


def scrape_matches(page: str = 1):
    url = 'https://www.vlr.gg/matches?page={}'.format(page)
    resp = requests.get(url, headers=headers)
    html = HTMLParser(resp.text)
    status = resp.status_code

    results = []

    for item in html.css("a.wf-module-item"):
        url_path = item.attributes['href']
        match_id = url_path.split("/")[1]

        eta = item.css_first(".match-item-eta").text().strip()
        eta = eta.replace("\t", " ").replace("\n", " ").split()
        try:
            if eta[0] == "ago":
                eta = "Live"
            else:
                eta = eta[1] + " " + eta[2] + " from now"
        except:
            eta = eta[0]

        rounds = item.css_first(".match-item-event-series").text().strip()

        event_name = item.css_first(".match-item-event").text().strip()
        event_name = event_name.replace("\t", " ")
        event_name = event_name.strip().split("\n")[1]
        event_name = event_name.strip()

        try:
            team_array = item.css_first(
                "div.match-item-vs").css_first("div:nth-child(2)").text()
        except:
            team_array = "TBD"

        team_array = team_array.replace("\t", " ").replace("\n", " ")
        team_array = team_array.strip().split(
            "                                  "
        )

        team_home = "TBD"
        team_away = "TBD"

        if team_array is not None and len(team_array) > 1:
            team_home = team_array[0]
            team_away = team_array[4].strip()

        results.append(
            {
                "id": match_id,
                "team_home": team_home,
                "team_away": team_away,
                "time_until_match": eta,
                "match_name": rounds,
                "event_name": event_name,
                "match_page": url_path
            }
        )

    segments = {"status": status, "segments": results}
    data = {"data": segments}

    if status != 200:
        raise Exception("API response: {}".format(status))

    return data


def lambda_handler(event, context):

    data = scrape_matches()

    return data
