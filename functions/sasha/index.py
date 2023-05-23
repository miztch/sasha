import re
from datetime import datetime

import requests
from selectolax.parser import HTMLParser

from utils import headers


def scrape_matches(page: str = 1):
    url = 'https://www.vlr.gg/matches?page={}'.format(page)
    resp = requests.get(url, headers=headers)
    html = selectolax.parser.HTMLParser(resp.text)
    status = resp.status_code

    results = []

    for item in html.css("a.wf-module-item"):
        url_path = item.attributes['href']

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

        tourney = item.css_first(".match-item-event").text().strip()
        tourney = tourney.replace("\t", " ")
        tourney = tourney.strip().split("\n")[1]
        tourney = tourney.strip()

        tourney_icon_url = item.css_first("img").attributes['src']
        tourney_icon_url = f"https:{tourney_icon_url}"

        flag_list = [flag_parent.attributes["class"].replace(
            " mod-", "_") for flag_parent in item.css('.flag')]
        flag1 = flag_list[0]
        flag2 = flag_list[1]

        try:
            team_array = item.css_first(
                "div.match-item-vs").css_first("div:nth-child(2)").text()
        except:
            team_array = "TBD"

        team_array = team_array.replace("\t", " ").replace("\n", " ")
        team_array = team_array.strip().split(
            "                                  "
        )

        team1 = "TBD"
        team2 = "TBD"

        if team_array is not None and len(team_array) > 1:
            team1 = team_array[0]
            team2 = team_array[4].strip()

        score1 = "-"
        score2 = "-"

        if team_array is not None and len(team_array) > 1:
            score1 = team_array[1].replace(" ", "").strip()
            score2 = team_array[-1].replace(" ", "").strip()

        results.append(
            {
                "team1": team1,
                "team2": team2,
                "flag1": flag1,
                "flag2": flag2,
                "score1": score1,
                "score2": score2,
                "time_until_match": eta,
                "round_info": rounds,
                "tournament_name": tourney,
                "match_page": url_path,
                "tournament_icon": tourney_icon_url,
            }
        )

    segments = {"status": status, "segments": results}
    data = {"data": segments}
    if data['status'] != 200:
        raise Exception("API response: {}".format(status))

    return data


def lambda_handler(event, context):

    data = scrape_matches(page)

    return data
