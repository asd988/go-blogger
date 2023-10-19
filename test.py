import requests

url = 'http://localhost:8080/upload'
#data = {'key1': 'value1', 'key2': 'value2'}
files = {'file': ('example.txt', open('test.md', 'rb'))}

# read file "secret"
secret = open('secret', 'r').read()
headers = {'Authorization': 'Bearer ' + secret}

response = requests.post(url, files=files, headers=headers)

if response.status_code == 200:
    print('Request succeeded!')
    print(response.text)
else:
    print('Request failed with status code', response.status_code)