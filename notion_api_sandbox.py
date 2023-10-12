import requests, json

token = 'secret_c6vmUvoXD96zIBFolSK0f4wNgu6j1n1MFxZZMDQitGf'

headers = {
    "Authorization": "Bearer " + token,
    'Notion-Version': '2022-06-28',
}

# page_id = "bdee982442aa49d99268bb47cfaa7626"
# page_id = "a3d06a6077154099a2c1a61e3e1dcba2"
page_id = "b97d4c05c37c4c1697c60d473933eee2"

readUrl = f"https://api.notion.com/v1/blocks/{page_id}/children"


res = requests.get(readUrl, headers=headers)
data = res.json()
print("res status code: ", res.status_code)

accepted_types = ['paragraph', 'heading_1', 'heading_2', 'heading_3', 'bulleted_list_item', 'numbered_list_item', 'to_do', 'toggle', 'callout', 'quote', 'code']

for block in data['results']:
    if block['type'] in accepted_types:
        print('One block: ', block)
        print('Block type: ', block['type'])
    print('\n\n')

# Note: Doesn't yet get all the nested content of e.g. bulleted lists

def extract_text_from_block(block):
    if block['type'] in accepted_types:
        text_key = block['type'].lower()  # Get the appropriate text key dynamically
        
        block_text = ''
        for text_segment in block[text_key]['rich_text']:
            if 'text' in text_segment and 'content' in text_segment['text']:
                block_text += text_segment['text']['content']
        
        print(f'{block["type"].capitalize()} content: ', block_text)
        
        # Recurse into nested blocks if available
        if 'children' in block:
            for nested_block in block['children']:
                extract_text_from_block(nested_block)

for block in data['results']:
    extract_text_from_block(block)


