import requests
secret = open('secret', 'r').read()
headers = {'Authorization': 'Bearer ' + secret}

def upload(file_name, path):
    url = 'http://localhost:8080/upload'
    files = {'file': (file_name, open(path, 'rb'))}

    return requests.post(url, files=files, headers=headers)

def create_blog(title, page_hash, file_hashes = []):
    url = 'http://localhost:8080/create_blog'
    data = {'title': title, 'page_hash': page_hash, 'file_hashes': file_hashes}

    return requests.post(url, json=data, headers=headers)

response = create_blog("My Favourite Dish", "tjV6UjSXgLwReFFAq8atsDcJhwVPphK7do9kNUCJyOE=")

if response.status_code == 200:
    print('Request succeeded!')
    print(response.text)
else:
    print('Request failed with status code', response.status_code)