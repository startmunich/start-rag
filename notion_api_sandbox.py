import requests, json

token = 'secret_c6vmUvoXD96zIBFolSK0f4wNgu6j1n1MFxZZMDQitGf'

headers = {
    "Authorization": "Bearer " + token,
    'Notion-Version': '2022-06-28',
}

page_id = "bdee982442aa49d99268bb47cfaa7626"

readUrl = f"https://api.notion.com/v1/blocks/{page_id}/children"


res = requests.get(readUrl, headers=headers)
data = res.json()
print("res status code: ", res.status_code)

accepted_types = ['paragraph', 'heading_1', 'heading_2', 'heading_3', 'bulleted_list_item', 'numbered_list_item', 'callout']

for block in data['results']:
    if block['type'] in accepted_types:
        print('One block: ', block)
        print('Block type: ', block['type'])
    print('\n\n')


