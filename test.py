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

response = upload("go.png", "gopher.png")
png_hash = response.json().get("hash")
response = upload("dish.md", "test.md")
md_hash = response.json().get("hash")
response = create_blog("My Favourite Dish", md_hash, [png_hash])

if response.status_code == 200:
    print('Request succeeded!')
    print(response.text)
else:
    print('Request failed with status code', response.status_code)