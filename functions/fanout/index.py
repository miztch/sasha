import json
import logging
import os

import boto3

logger = logging.getLogger()
logger.setLevel(logging.INFO)

sqs = boto3.client('sqs')
queue_url = os.environ['FANOUT_QUEUE_URL']


def publish_page_numbers():
    '''
    publish page numbers to get match information to the queue
    '''

    pages_to_scrape = int(os.environ['PAGES_TO_SCRAPE'])

    base_delay_seconds = int(os.environ['BASE_DELAY_SECONDS'])
    for page in range(1, pages_to_scrape+1):
        logger.info('request to fetch match list for the day: {}'.format(page))

        payload = {'page': page}
        message = json.dumps(payload)

        response = sqs.send_message(
            QueueUrl=queue_url,
            MessageBody=message,
            DelaySeconds=min(page*base_delay_seconds, 900)
        )

        logger.info('message sent. queue: {} response: {}'.format(
            queue_url, response))


def lambda_handler(event, context):
    publish_page_numbers()
