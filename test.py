import requests
secret = open('secret', 'r').read()
headers = {'Authorization': 'Bearer ' + secret}

def upload(file_name, path):
    url = 'http://localhost:8080/upload'
    files = {'file': (file_name, open(path, 'rb'))}

    return requests.post(url, files=files, headers=headers)

response = upload("hello.md", "test.md")

if response.status_code == 200:
    print('Request succeeded!')
    print(response.text)
else:
    print('Request failed with status code', response.status_code)