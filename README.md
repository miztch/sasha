# sasha

- scraper for [vlr.gg](https://www.vlr.gg/) upcoming matches
  - save matches information to DynamoDB
    - this is to be the source data of [dima](https://github.com/miztch/dima), API for upcoming matches
    - alternatively you can run locally with standard output
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
  "pagePath": "/213198/fnatic-vs-team-liquid-champions-tour-2023-emea-league-gf",
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

## Provisioning as Lambda Function (with State Machine)

You can use [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html) to provision this application.

```bash
cd sasha/
sam build
sam deploy --guided --capabilities CAPABILITY_IAM
```

## Running Locally

You can run locally with passing `--page` argument to scrape matches in a specific page.

```bash
go run *.go --page 1
```

Output is json formatted.

```json
[{"id":399763,"matchName":"Group Stage: Round 6","startDate":"2024-09-15","startTime":"2024-09-15T19:00:00+0000","bestOf":3,"teams":[{"title":"RETA Esports"},{"title":"Galorys"}],"pagePath":"/399763/reta-esports-vs-galorys-champions-tour-2024-americas-ascension-r6","eventName":"Champions Tour 2024 Americas: Ascension","eventCountryFlag":"mx"},...]
```
