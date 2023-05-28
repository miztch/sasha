# sasha

- scraper for [vlr.gg](https://www.vlr.gg/) upcoming matches
  - save matches information to DynamoDB
    - this is to be the source data of [dima](https://github.com/miztch/dima), API for upcoming matches
  - named after the real name of Valorant agent Sova 🦉🏹

## Data format: sample

```json
{
    "id": 213198,
    "startTime": "2023-05-28T15:00:00+0000",
    "bestOf": "5",
    "eventCountryFlag": "de",
    "eventName": "Champions Tour 2023: EMEA League",
    "matchName": "Playoffs: Grand Final",
    "match_page": "/213198/fnatic-vs-team-liquid-champions-tour-2023-emea-league-gf",
    "startDate": "2023-05-28",
    "teams": [
        {
            "title": "FNATIC"
        },
        {
           "title": "Team Liquid"
        }
    ]
}
```

## Provisioning

You can use [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) to provision this application.

```bash
cd sasha/
sam build
sam deploy --guided --capabilities CAPABILITY_IAM
```
