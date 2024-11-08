import time
import requests
import schedule


def job():
    url = "http://api_gateway:3000/api/v1/indego-data-fetch-and-store-it-db"
    payload = {}
    headers = {
        'Authorization': 'Bearer secret_token_static'
    }

    response = requests.request("POST", url, headers=headers, data=payload)

    print(response.text)


schedule.every().hour.do(job)
# schedule.every(2).minutes.do(job)

while True:
    schedule.run_pending()
    time.sleep(1)
