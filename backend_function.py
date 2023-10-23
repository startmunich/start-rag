import requests

def extract_text_from_notion_page(token):
    # Ask the user to input the Notion page URL
    url = input("Enter the Notion page URL: ")

    # Extract page ID from the given Notion page URL
    page_id = url.rsplit('-', 1)[-1]
    print("page_id: ", page_id)

    headers = {
        "Authorization": "Bearer " + token,
        'Notion-Version': '2022-06-28',
    }

    readUrl = f"https://api.notion.com/v1/blocks/{page_id}/children"

    res = requests.get(readUrl, headers=headers)
    data = res.json()

    accepted_types = ['paragraph', 'heading_1', 'heading_2', 'heading_3', 'bulleted_list_item', 'numbered_list_item', 'to_do', 'toggle', 'callout', 'quote', 'code']

    page_content = ""

    def extract_text_from_block(block, content_string):
        if block['type'] in accepted_types:
            text_key = block['type'].lower()
            block_text = ''
            if text_key in block and 'rich_text' in block[text_key]:
                for text_segment in block[text_key]['rich_text']:
                    if 'text' in text_segment and 'content' in text_segment['text']:
                        block_text += text_segment['text']['content'] + ' '  # Add a space for separation
            content_string += block_text + '\n'

            # Check if the block has children indicating nested content
            if 'has_children' in block and block['has_children']:
                children_url = f"https://api.notion.com/v1/blocks/{block['id']}/children"
                res = requests.get(children_url, headers=headers)
                children_data = res.json()
                for nested_block in children_data['results']:
                    content_string = extract_text_from_block(nested_block, content_string)
        return content_string

    for block in data['results']:
        page_content = extract_text_from_block(block, page_content)

    return page_content

token = 'secret_c6vmUvoXD96zIBFolSK0f4wNgu6j1n1MFxZZMDQitGf'
page_content = extract_text_from_notion_page(token)
print("Page Content:\n", page_content)