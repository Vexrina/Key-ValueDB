import httpx
from typing import Final
import pytest
import requests
import json
BASEURL: Final = "http://127.0.0.1:8080/api/"
ENDPOINT: Final = {
    "db": "database",
    "table": "table",
    "key": "key"
}

DATABASE: Final = "database_new"
TABLE: Final = "table_database"

@pytest.fixture
def base_url():
    return BASEURL


@pytest.fixture
def endpoint():
    return ENDPOINT

@pytest.fixture
def database():
    return DATABASE

@pytest.fixture
def table():
    return TABLE


def test_create_database(base_url, endpoint, database):
    response = httpx.post(base_url + endpoint["db"], json={"db_name": database})
    assert response.status_code == 201, f"Unexpected status code: {response.status_code}"

def test_create_table(base_url, endpoint, database, table):
    response = httpx.post(base_url + endpoint["table"] + "?db_name="+database, json={"table_name": table})
    assert response.status_code == 201, f"Unexpected status code: {response.status_code}"

def test_put_table(base_url, endpoint, database, table):
    response = httpx.put(base_url + endpoint["table"] + "?db_name=" + database,
                         json={"table_name": table, "new_table_name": "table_database_new2"})
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"
    response = httpx.put(base_url + endpoint["table"] + "?db_name=" + database,
                         json={"table_name": "table_database_new2", "new_table_name": table})
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"

def test_get_table(base_url, endpoint, database, table):
    response = requests.get(base_url + endpoint["table"] + "?db_name=" + database, json={"table_name": table})
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"

def test_create_key(base_url, endpoint, database, table):
    json = {
        "key_name": "1",
        "value": {
            "Val": "table_database",
            "Ttl": "2025-01-01T15:30:30Z"
        }
    }
    response = httpx.post(base_url + endpoint["key"] + "?db_name=" + database+"&table_name="+table, json=json)
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"


def test_get_key(base_url, endpoint, database, table):
    response = requests.get(base_url + endpoint["key"] + "?db_name=" + database + "&table_name=" + table, json={"key_name": "1"})
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"
    assert response.json() == {
            "Val": "table_database",
            "Ttl": "2025-01-01T15:30:30Z"
        }

def test_put_key(base_url, endpoint, database, table):
    json = {
        "key_name": "1",
        "value": {
            "Val": "new_table_database",
            "Ttl": "2025-01-01T15:30:30Z"
        }
    }
    response = httpx.put(base_url + endpoint["key"] + "?db_name=" + database + "&table_name=" + table, json=json)
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"
    response = requests.get(base_url + endpoint["key"] + "?db_name=" + database + "&table_name=" + table,
                            json={"key_name": "1"})
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"
    assert response.json() == {
        "Val": "new_table_database",
        "Ttl": "2025-01-01T15:30:30Z"
    }

def test_delete_key(base_url, endpoint, database, table):
    response = requests.delete(base_url + endpoint["key"] + "?db_name=" + database + "&table_name=" + table, json={"key_name": "1"})
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}"

def test_delete_table(base_url, endpoint, database, table):
    response = requests.delete(base_url + endpoint["table"] + "?db_name=" + database,
                               json={"table_name":table})
    assert response.status_code == 202, f"Unexpected status code: {response.status_code}"

def test_delete_database(base_url, endpoint, database):
    response = requests.delete(base_url + endpoint["db"] , json={"db_name": database})
    assert response.status_code == 202, f"Unexpected status code: {response.status_code}"