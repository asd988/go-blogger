import requests
import os
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

def download(file_hash, store_to_path = "tmp"):
    url = 'http://localhost:8080/download'
    data = {'hash': file_hash}

    response = requests.post(url, json=data, headers=headers)
    file_name = response.headers['Content-Disposition'].split('=')[1]
    open(os.path.join(store_to_path, file_name), 'wb').write(response.content)

    return response

def get_file_hashes(blog_id):
    url = 'http://localhost:8080/get_file_hashes'
    data = {'blog_id': blog_id}

    return requests.post(url, json=data, headers=headers)

def download_blog(blod_id, path = "tmp"):
    response = get_file_hashes(blod_id)
    file_hashes = response.json().get("file_hashes")
    for file_hash in file_hashes:
        download(file_hash, path)

def upload_blog(title, path = "tmp"):
    files = os.listdir(path)
    # check if there is one and only md
    md_file = [file for file in files if file.endswith(".md")]
    if len(md_file) != 1:
        print("There should be one and only one md file in the folder")
        return
    md_file = md_file[0]

    # upload all files
    file_hashes = []
    for file in files:
        r = upload(file, os.path.join(path, file))
        if file == md_file:
            page_hash = r.json().get("hash")
        else:
            file_hashes.append(r.json().get("hash"))

    # create blog
    r = create_blog(title, page_hash, file_hashes)
    return r



response = upload_blog("My Favorite dishes")

if response.status_code == 200:
    print('Request succeeded!')
    print(response.text)
else:
    print('Request failed with status code', response.status_code)