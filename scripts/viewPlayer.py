import requests
from bs4 import BeautifulSoup

link = "https://www.pro-football-reference.com/players/B/BradTo00/gamelog/2021/advanced/"


resp = requests.get(link)
soup = BeautifulSoup(resp.text, "html.parser")
final = [td.get("data-stat") for td in soup.find("table").findAll("tr")[2].findAll("td")]
print(final)
print(len(final))