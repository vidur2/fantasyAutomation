import requests
from bs4 import BeautifulSoup

res = requests.get("https://www.pro-football-reference.com/teams/oti/2021.htm")

soup = BeautifulSoup(res.content, "html.parser")

for table in soup.findAll("table"):
    dataStats = []
    for row in table.findAll("tr"):
        for td in row.findAll("td"):
            dataStats.append(td.get("data-stat"))
        
        if len(dataStats) > 0:
            print(dataStats)